package cert

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

// VaultPKISource implements a certificate source which issues TLS certificates
// on-demand using a Vault PKI backend. Client authorization certificates are
// loaded from a generic backend (same as in VaultSource). The Vault token
// should be set through the VAULT_TOKEN environment variable.
//
// The TLS certificates are re-issued automatically before they expire.
type VaultPKISource struct {
	Client       *vaultClient
	CertPath     string
	ClientCAPath string
	CAUpgradeCN  string

	// Re-issue certificates this long before they expire. Cannot be less then
	// one hour.
	Refresh time.Duration

	certsCh chan []tls.Certificate

	mu    sync.Mutex
	certs map[string]tls.Certificate // issued certs
}

func NewVaultPKISource() *VaultPKISource {
	return &VaultPKISource{
		certs:   make(map[string]tls.Certificate, 0),
		certsCh: make(chan []tls.Certificate, 1),
	}
}

func (s *VaultPKISource) LoadClientCAs() (*x509.CertPool, error) {
	return (&VaultSource{
		Client:       s.Client,
		ClientCAPath: s.ClientCAPath,
		CAUpgradeCN:  s.CAUpgradeCN,
	}).LoadClientCAs()
}

func (s *VaultPKISource) Certificates() chan []tls.Certificate {
	return s.certsCh
}

func (s *VaultPKISource) Issue(commonName string) (*tls.Certificate, error) {
	c, err := s.Client.Get()
	if err != nil {
		return nil, fmt.Errorf("vault: client: %s", err)
	}

	resp, err := c.Logical().Write(s.CertPath, map[string]interface{}{
		"common_name": commonName,
	})
	if err != nil {
		fmt.Printf("Issue: %v\n", err)
		return nil, fmt.Errorf("vault: issue: %s", err)
	}

	b, _ := json.Marshal(resp.Data)
	var data struct {
		// Try load from regular vault-pki
		PrivateKey  string   `json:"private_key"`
		Certificate string   `json:"certificate"`
		CAChain     []string `json:"ca_chain"`
		// Try to load data from acme-vault
		Certificate2 string `json:"cert"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("vault: issue: %s", err)
	}

	if data.PrivateKey == "" {
		return nil, fmt.Errorf("vault: issue: missing private key")
	}

	if data.Certificate == "" && data.Certificate2 == "" {
		return nil, fmt.Errorf("vault: issue: missing certificate")
	}

	key := []byte(data.PrivateKey)
	fullChain := []byte(data.Certificate)
	if data.Certificate2 != "" {
		fullChain = append(fullChain, data.Certificate2...)
	}
	for _, c := range data.CAChain {
		fullChain = append(fullChain, '\n')
		fullChain = append(fullChain, []byte(c)...)
	}

	cert, err := tls.X509KeyPair(fullChain, key)
	if err != nil {
		return nil, fmt.Errorf("vault: issue: %s", err)
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		// Should never happen because x509.ParseCertificate did this
		// successfully already, but threw the result away.
		return nil, fmt.Errorf("vault: issue: %s", err)
	}

	refresh := s.Refresh
	if refresh < time.Hour {
		refresh = time.Hour
	}

	expires := x509Cert.NotAfter
	certTTL := time.Until(expires) - refresh
	time.AfterFunc(certTTL, func() {
		_, err := s.Issue(commonName)
		if err != nil {
			log.Printf("[ERROR] cert: vault: Failed to re-issue cert for %s: %s", commonName, err)
			// TODO: Now what? Retry? Do nothing?
			return
		}
	})

	s.mu.Lock()
	s.certs[commonName] = cert
	allCerts := make([]tls.Certificate, 0, len(s.certs))
	for _, c := range s.certs {
		allCerts = append(allCerts, c)
	}
	s.mu.Unlock()

	go func() { s.certsCh <- allCerts }()
	log.Printf("[INFO] cert: vault: issued cert for %s; serial = %s", commonName, s.formatSerial(x509Cert.SerialNumber))

	return &cert, nil
}

func (*VaultPKISource) formatSerial(sn *big.Int) string {
	var buf bytes.Buffer
	for _, b := range sn.Bytes() {
		if buf.Len() > 0 {
			buf.WriteByte('-')
		}
		fmt.Fprintf(&buf, "%02x", b)
	}
	return buf.String()
}
