package architect

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
	"golang.org/x/net/ipv4"
)

// Profile defines the fingerprint for both TLS and HTTP/2 layers.
type Profile struct {
	ID                string
	Name              string
	TLSID             utls.ClientHelloID
	HTTP2Settings     []http2.Setting
	InitialWinSize    uint32
	MaxHeaderListSize uint32
	UserAgent         string

	// Network Layer (TCP/IP Fingerprinting)
	TTL               int
	
	// ECH (Experimental)
	ECHConfig         []byte

	// Application Layer (Header Normalization)
	PseudoHeaderOrder []string
	HeaderOrder       []string
}

type profileKey struct{}

// WithProfile attaches a specific Profile to the context for the Transport to use.
func WithProfile(ctx context.Context, p Profile) context.Context {
	return context.WithValue(ctx, profileKey{}, p)
}

// Transport is the core of Architect's evasion logic.
type Transport struct {
	Manager            *ProfileManager
	DefaultProfile     Profile
	InsecureSkipVerify bool
	mu                 sync.Mutex
	transports         map[string]*http2.Transport
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" {
		return nil, fmt.Errorf("architect: only https is supported")
	}

	profile := t.DefaultProfile
	if p, ok := req.Context().Value(profileKey{}).(Profile); ok {
		profile = p
	}

	// Automatic header normalization on the request
	t.normalizeHeaders(req, profile)

	h2Tr := t.getTransportForProfile(profile)
	return h2Tr.RoundTrip(req)
}

func (t *Transport) normalizeHeaders(req *http.Request, p Profile) {
	// Set User-Agent from profile if not set
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", p.UserAgent)
	}
	
	// Note: Forced header ordering requires a custom HPACK encoder
	// or a proxy-based approach to bypass Go's default map sorting.
}

func (t *Transport) getTransportForProfile(p Profile) *http2.Transport {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.transports == nil {
		t.transports = make(map[string]*http2.Transport)
	}

	if tr, ok := t.transports[p.ID]; ok {
		return tr
	}

	tr := &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return t.dialTLSWithProfile(ctx, network, addr, p)
		},
		MaxHeaderListSize: p.MaxHeaderListSize,
	}
	t.transports[p.ID] = tr
	return tr
}

func (t *Transport) dialTLSWithProfile(ctx context.Context, network, addr string, p Profile) (net.Conn, error) {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	// TCP Normalization (Layer 3/4)
	if p.TTL > 0 {
		_ = ipv4.NewConn(conn).SetTTL(p.TTL)
	}

	host, _, _ := net.SplitHostPort(addr)
	config := &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: t.InsecureSkipVerify,
	}

	// Experimental ECH support
	if len(p.ECHConfig) > 0 {
		config.EncryptedClientHelloConfigList = p.ECHConfig
	}

	uConn := utls.UClient(conn, config, p.TLSID)
	if err := uConn.HandshakeContext(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}

	if uConn.ConnectionState().NegotiatedProtocol == "h2" {
		return &h2Conn{Conn: uConn, profile: p}, nil
	}
	return uConn, nil
}

// h2Conn wraps the connection to intercept the HTTP/2 preface and SETTINGS frame.
type h2Conn struct {
	net.Conn
	profile  Profile
	once     sync.Once
	preface  bool
}

func (c *h2Conn) Write(b []byte) (n int, err error) {
	c.once.Do(func() {
		if len(b) >= 24 && string(b[:24]) == "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n" {
			c.preface = true
		}
	})
	return c.Conn.Write(b)
}

// ProfileManager manages a collection of browser profiles.
type ProfileManager struct {
	profiles map[string]Profile
	mu       sync.RWMutex
}

func NewProfileManager() *ProfileManager {
	return &ProfileManager{profiles: make(map[string]Profile)}
}

func (pm *ProfileManager) Register(p Profile) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.profiles[p.ID] = p
}

// NewClient returns an http.Client configured with the Architect transport.
func NewClient(profile Profile) *http.Client {
	return &http.Client{
		Transport: &Transport{
			DefaultProfile: profile,
		},
	}
}
