package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/route"

	gkm "github.com/go-kit/kit/metrics"
	grpc_proxy "github.com/mwitkow/grpc-proxy/proxy"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	server *grpc.Server
}

func (s *gRPCServer) Close() error {
	s.server.Stop()
	return nil
}

func (s *gRPCServer) Shutdown(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *gRPCServer) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

func GetGRPCDirector(tlscfg *tls.Config, cfg *config.Config) func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {

	connectionPool := newGrpcConnectionPool(tlscfg, cfg)

	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return ctx, nil, fmt.Errorf("error extracting metadata from request")
		}

		outCtx := metadata.NewOutgoingContext(ctx, md.Copy())

		target, _ := ctx.Value(targetKey{}).(*route.Target)

		if target == nil {
			log.Println("[WARN] grpc: no route for ", fullMethodName)
			return outCtx, nil, fmt.Errorf("no route found")
		}

		conn, err := connectionPool.Get(outCtx, target)

		return outCtx, conn, err
	}

}

type GrpcProxyInterceptor struct {
	Config       *config.Config
	StatsHandler *GrpcStatsHandler
	GlobCache    *route.GlobCache
}

type targetKey struct{}

type proxyStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (p proxyStream) Context() context.Context {
	return p.ctx
}

func makeGRPCTargetKey(t *route.Target) string {
	return t.URL.String()
}

func (g GrpcProxyInterceptor) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()

	target, err := g.lookup(ctx, info.FullMethod)

	if err != nil {
		log.Println("[ERROR] grpc: error looking up route", err)
		return status.Error(codes.Internal, "internal error")
	}

	if target == nil {
		g.StatsHandler.NoRoute.Add(1)
		log.Println("[WARN] grpc: no route found for", info.FullMethod)
		return status.Error(codes.NotFound, "no route found")
	}

	ctx = context.WithValue(ctx, targetKey{}, target)

	proxyStream := proxyStream{
		ServerStream: stream,
		ctx:          ctx,
	}

	start := time.Now()

	err = handler(srv, proxyStream)

	end := time.Now()
	dur := end.Sub(start)

	target.Timer.Observe(dur.Seconds())

	return err
}

func (g GrpcProxyInterceptor) lookup(ctx context.Context, fullMethodName string) (*route.Target, error) {
	pick := route.Picker[g.Config.Proxy.Strategy]
	match := route.Matcher[g.Config.Proxy.Matcher]

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("error extracting metadata from request")
	}

	reqUrl, err := url.ParseRequestURI(fullMethodName)

	if err != nil {
		log.Print("[WARN] Error parsing grpc request url ", fullMethodName)
		return nil, fmt.Errorf("error parsing request url")
	}

	headers := http.Header{}

	for k, v := range md {
		for _, h := range v {
			headers.Add(k, h)
		}
	}

	//grpc client can specify a destination host in metadata
	dstHostSpecifiedByGRPCClient := g.getDestinationHostFromMetadata(md)
	//todo: better a configuration flag is required to disable/enable this function, and make it disabled by default configuration

	req := &http.Request{
		Host:   dstHostSpecifiedByGRPCClient,
		URL:    reqUrl,
		Header: headers,
	}

	return route.GetTable().Lookup(req, pick, match, g.GlobCache, g.Config.GlobMatchingDisabled), nil
}

// grpc client can specify a destination host in metadata by key 'dsthost', e.g. dsthost=betatest
// the backend service(s) tags should be urlprefix-betatest/grpcpackage.servicename proto=grpc
// the 'betatest' will be parsed as 'host' and '/grpcpackage.servicename' is the 'path',
// a route record will be setup in route Table, t['betatest']
// the dstHost is extracted from context's metadata of grpc client, that will trigger t[dstHost] is used.
// if t[dstHost] not exists, fallback to t[""] is used
// dstHost will be "" as before if not specified by grpc client side.
func (g GrpcProxyInterceptor) getDestinationHostFromMetadata(md metadata.MD) (dstHost string) {
	dstHost = ""
	hosts := md["dsthost"]
	if len(hosts) == 1 {
		dstHost = hosts[0]
	}
	return
}

type GrpcStatsHandler struct {
	Connect gkm.Counter
	Request gkm.Histogram
	NoRoute gkm.Counter
	Status  gkm.Histogram
}

type connCtxKey struct{}
type rpcCtxKey struct{}

func (h *GrpcStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return context.WithValue(ctx, connCtxKey{}, info)
}

func (h *GrpcStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return context.WithValue(ctx, rpcCtxKey{}, info)
}

func (h *GrpcStatsHandler) HandleRPC(ctx context.Context, rpc stats.RPCStats) {

	rpcStats, _ := rpc.(*stats.End)

	if rpcStats == nil {
		return
	}

	dur := rpcStats.EndTime.Sub(rpcStats.BeginTime)

	h.Request.Observe(dur.Seconds())

	s, _ := status.FromError(rpcStats.Error)

	h.Status.With("code", s.Code().String()).Observe(dur.Seconds())
}

// HandleConn processes the Conn stats.
func (h *GrpcStatsHandler) HandleConn(ctx context.Context, conn stats.ConnStats) {
	connBegin, _ := conn.(*stats.ConnBegin)

	if connBegin != nil {
		h.Connect.Add(1)
	}
}

type grpcConnectionPool struct {
	connections     map[string]*grpc.ClientConn
	lock            sync.RWMutex
	cleanupInterval time.Duration
	tlscfg          *tls.Config
	cfg             *config.Config
}

func newGrpcConnectionPool(tlscfg *tls.Config, cfg *config.Config) *grpcConnectionPool {
	cp := &grpcConnectionPool{
		connections:     make(map[string]*grpc.ClientConn),
		lock:            sync.RWMutex{},
		cleanupInterval: time.Second * 5,
		tlscfg:          tlscfg,
		cfg:             cfg,
	}

	go cp.cleanup()

	return cp
}

func (p *grpcConnectionPool) Get(ctx context.Context, target *route.Target) (*grpc.ClientConn, error) {
	p.lock.RLock()
	conn := p.connections[makeGRPCTargetKey(target)]
	p.lock.RUnlock()

	if conn != nil && conn.GetState() != connectivity.Shutdown {
		return conn, nil
	}

	return p.newConnection(ctx, target)
}

func (p *grpcConnectionPool) newConnection(ctx context.Context, target *route.Target) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.CallCustomCodec(grpc_proxy.Codec()), grpc.MaxCallRecvMsgSize(p.cfg.Proxy.GRPCMaxRxMsgSize)),
	}

	if target.URL.Scheme == "grpcs" && p.tlscfg != nil {
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				ClientCAs:          p.tlscfg.ClientCAs,
				InsecureSkipVerify: target.TLSSkipVerify,
				// as per the http/2 spec, the host header isn't required, so if your
				// target service doesn't have IP SANs in it's certificate
				// then you will need to override the servername
				ServerName: target.Opts["grpcservername"],
			})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.DialContext(ctx, target.URL.Host, opts...)

	if err == nil {
		p.Set(target, conn)
	}

	return conn, err
}

func (p *grpcConnectionPool) Set(target *route.Target, conn *grpc.ClientConn) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.connections[makeGRPCTargetKey(target)] = conn
}

func (p *grpcConnectionPool) cleanup() {
	for {
		p.lock.Lock()
		table := route.GetTable()
		for tKey, cs := range p.connections {
			state := cs.GetState()
			if state == connectivity.Shutdown {
				delete(p.connections, tKey)
				continue
			}

			if !hasTarget(tKey, table) {
				log.Println("[DEBUG] grpc: cleaning up connection to", tKey)
				go func(cs *grpc.ClientConn, state connectivity.State) {
					ctx, cancel := context.WithTimeout(context.Background(), p.cfg.Proxy.GRPCGShutdownTimeout)
					defer cancel()
					// wait for state to change, or timeout, before closing, in case it's still handling traffic.
					cs.WaitForStateChange(ctx, state)
					cs.Close()
				}(cs, state)
				delete(p.connections, tKey)
			}
		}
		p.lock.Unlock()
		time.Sleep(p.cleanupInterval)
	}
}

func hasTarget(tKey string, table route.Table) bool {
	for _, routes := range table {
		for _, r := range routes {
			for _, t := range r.Targets {
				if tKey == makeGRPCTargetKey(t) {
					return true
				}
			}
		}
	}
	return false
}
