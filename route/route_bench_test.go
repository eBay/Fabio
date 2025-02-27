package route

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

var (
	b5Routes   Table
	b10Routes  Table
	b100Routes Table
	b500Routes Table

	once sync.Once
)

// initRoutes is used for lazy one time initialization of the test data for
// the parallel benchmarks via once
func initRoutes() {
	b5Routes = makeRoutes(1, 5, 1, 6)
	b10Routes = makeRoutes(1, 5, 2, 6)
	b100Routes = makeRoutes(10, 5, 2, 24)
	b500Routes = makeRoutes(10, 10, 5, 24)
}

func BenchmarkPrefixMatcherRndPicker5Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b5Routes, prefixMatcher, rndPicker, b) })
}

func BenchmarkPrefixMatcherRRPicker5Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b5Routes, prefixMatcher, rrPicker, b) })
}

func BenchmarkPrefixMatcherRndPicker10Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b10Routes, prefixMatcher, rndPicker, b) })
}

func BenchmarkPrefixMatcherRRPicker10Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b10Routes, prefixMatcher, rrPicker, b) })
}

func BenchmarkPrefixMatcherRndPicker100Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b100Routes, prefixMatcher, rndPicker, b) })
}

func BenchmarkPrefixMatcherRRPicker100Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b100Routes, prefixMatcher, rrPicker, b) })
}

func BenchmarkPrefixMatcherRndPicker500Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b500Routes, prefixMatcher, rndPicker, b) })
}

func BenchmarkPrefixMatcherRRPicker500Routes(b *testing.B) {
	once.Do(initRoutes)
	b.ResetTimer()
	b.SetParallelism(3)
	b.RunParallel(func(b *testing.PB) { benchmarkGet(b500Routes, prefixMatcher, rrPicker, b) })
}

// makeRoutes builds a set of routes for a set of domains
// and target urls. For each domain all paths up to depth
// are constructed and all host/path combinations have the
// same target URLs. The number of generated routes is
// domains * paths * depth.
func makeRoutes(domains, paths, depth, urls int) Table {
	s := ""
	for i := 0; i < domains; i++ {
		prefix := fmt.Sprintf("www.host-%d.com/", i)
		for j := 0; j < paths; j++ {
			for k := 0; k < depth; k++ {
				prefix += fmt.Sprintf("path-%d/", k)
				for l := 0; l < urls; l++ {
					s += fmt.Sprintf("route add svc %s http://host:12345/\n", prefix)
				}
			}
		}
	}

	t, err := NewTable(bytes.NewBufferString(s))
	if err != nil {
		panic(err)
	}
	return t
}

// makeRequests builds a list of http.Request objects with an
// additional path for benchmarking.
func makeRequests(t Table) []*http.Request {
	reqs := []*http.Request{}
	for host, hr := range t {
		for _, r := range hr {
			req := &http.Request{Host: host, RequestURI: r.Path + "/some/additional/path"}
			reqs = append(reqs, req)
		}
	}
	return reqs
}

// benchmarkGet runs the benchmark on the Table.Lookup() function with the
// given matcher and picker functions.
func benchmarkGet(t Table, match matcher, pick picker, pb *testing.PB) {
	reqs := makeRequests(t)
	k, n := len(reqs), 0
	//Glob Matching True
	for pb.Next() {
		t.Lookup(reqs[n%k], pick, match, globCache, globEnabled)
		n++
	}
}
