package architect

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"syscall"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/proxy"
)

// Profile defines the fingerprint for all network layers.
type Profile struct {
	ID                string
	Name              string
	TLSID             utls.ClientHelloID
	HTTP2Settings     []http2.Setting
	InitialWinSize    uint32
	MaxHeaderListSize uint32
	UserAgent         string

	// Network Layer (TCP/IP Fingerprinting)
	TTL       int
	TCPWindow int // New: Initial TCP Window Size (e.g., 65535)

	// ECH (Experimental)
	ECHConfig []byte

	// Application Layer (Header Normalization)
	PseudoHeaderOrder []string
	HeaderOrder       []string
}

type profileKey struct{}

func WithProfile(ctx context.Context, p Profile) context.Context {
	return context.WithValue(ctx, profileKey{}, p)
}

type Transport struct {
	Manager            *ProfileManager
	DefaultProfile     Profile
	InsecureSkipVerify bool
	ProxyURL           string
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

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", profile.UserAgent)
	}

	h2Tr := t.getTransportForProfile(profile)
	return h2Tr.RoundTrip(req)
}

func (t *Transport) getTransportForProfile(p Profile) *http2.Transport {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.transports == nil {
		t.transports = make(map[string]*http2.Transport)
	}

	poolKey := p.ID + "_" + t.ProxyURL
	if tr, ok := t.transports[poolKey]; ok {
		return tr
	}

	tr := &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return t.dialTLSWithProfile(ctx, network, addr, p)
		},
		MaxHeaderListSize: p.MaxHeaderListSize,
	}
	t.transports[poolKey] = tr
	return tr
}

func (t *Transport) dialTLSWithProfile(ctx context.Context, network, addr string, p Profile) (net.Conn, error) {
	dialer := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// Apply TCP Window Size before the handshake
				if p.TCPWindow > 0 {
					setTCPWindow(fd, p.TCPWindow)
				}
			})
		},
	}

	var conn net.Conn
	var err error

	if t.ProxyURL != "" {
		proxyURL, err := url.Parse(t.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %v", err)
		}

		switch proxyURL.Scheme {
		case "socks5", "socks5h":
			// Note: Proxy libraries often hide the raw conn, making TCP tuning harder.
			// In Architect, we prioritize uTLS/H2 over TCP tuning if a proxy is used.
			pd, _ := proxy.FromURL(proxyURL, proxy.Direct)
			conn, err = pd.Dial(network, addr)
		case "http", "https":
			d := &net.Dialer{Control: dialer.Control}
			conn, err = d.DialContext(ctx, network, proxyURL.Host)
			if err == nil {
				connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
				fmt.Fprintf(conn, connectReq)
				resp := make([]byte, 1024)
				n, _ := conn.Read(resp)
				if !bytes.Contains(resp[:n], []byte("200 Connection established")) {
					conn.Close()
					return nil, fmt.Errorf("proxy failed")
				}
			}
		}
	} else {
		conn, err = dialer.DialContext(ctx, network, addr)
	}

	if err != nil {
		return nil, err
	}

	// Apply TTL Normalization
	if p.TTL > 0 {
		_ = ipv4.NewConn(conn).SetTTL(p.TTL)
	}

	host, _, _ := net.SplitHostPort(addr)
	config := &utls.Config{ServerName: host, InsecureSkipVerify: t.InsecureSkipVerify}
	if len(p.ECHConfig) > 0 {
		config.EncryptedClientHelloConfigList = p.ECHConfig
	}

	uConn := utls.UClient(conn, config, p.TLSID)
	if err := uConn.HandshakeContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}

	if uConn.ConnectionState().NegotiatedProtocol == "h2" {
		c := &h2Conn{Conn: uConn, profile: p}
		c.hpackDec = hpack.NewDecoder(4096, nil)
		c.hpackEnc = hpack.NewEncoder(&c.encBuffer)
		return c, nil
	}
	return uConn, nil
}

// ... h2Conn and reorderHeaders implementation ...

type h2Conn struct {
	net.Conn
	profile   Profile
	once      sync.Once
	preface   bool
	hpackDec  *hpack.Decoder
	hpackEnc  *hpack.Encoder
	encBuffer bytes.Buffer
}

func (c *h2Conn) Write(b []byte) (n int, err error) {
	c.once.Do(func() {
		if len(b) >= 24 && string(b[:24]) == "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n" {
			c.preface = true
		}
	})

	if c.preface && len(b) > 9 {
		frameType := b[3]
		if frameType == 0x1 { // HEADERS
			newFrame, err := c.processHeadersFrame(b)
			if err == nil {
				_, err = c.Conn.Write(newFrame)
				return len(b), err
			}
		}
	}

	return c.Conn.Write(b)
}

func (c *h2Conn) processHeadersFrame(b []byte) ([]byte, error) {
	length := int(binary.BigEndian.Uint32(append([]byte{0}, b[:3]...)))
	if len(b) < 9+length {
		return nil, fmt.Errorf("fragmented frame")
	}
	payload := b[9 : 9+length]

	headers, err := c.hpackDec.DecodeFull(payload)
	if err != nil {
		return nil, err
	}

	ordered := c.reorderHeaders(headers)

	c.encBuffer.Reset()
	for _, h := range ordered {
		c.hpackEnc.WriteField(h)
	}
	newPayload := c.encBuffer.Bytes()

	newFrame := make([]byte, 9+len(newPayload))
	copy(newFrame, b[:9])
	newLen := uint32(len(newPayload))
	newFrame[0] = byte(newLen >> 16)
	newFrame[1] = byte(newLen >> 8)
	newFrame[2] = byte(newLen)

	copy(newFrame[9:], newPayload)
	return newFrame, nil
}

func (c *h2Conn) reorderHeaders(original []hpack.HeaderField) []hpack.HeaderField {
	var pseudo []hpack.HeaderField
	var regular []hpack.HeaderField
	for _, h := range original {
		if h.Name != "" && h.Name[0] == ':' {
			pseudo = append(pseudo, h)
		} else {
			regular = append(regular, h)
		}
	}

	result := make([]hpack.HeaderField, 0, len(original))
	for _, target := range c.profile.PseudoHeaderOrder {
		for _, h := range pseudo {
			if h.Name == target {
				result = append(result, h)
			}
		}
	}
	for _, target := range c.profile.HeaderOrder {
		for _, h := range regular {
			if h.Name == target {
				result = append(result, h)
			}
		}
	}
	// Fallback
	for _, h := range original {
		found := false
		for _, r := range result {
			if r.Name == h.Name {
				found = true
				break
			}
		}
		if !found {
			result = append(result, h)
		}
	}
	return result
}

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

func NewClient(profile Profile) *http.Client {
	return &http.Client{
		Transport: &Transport{
			DefaultProfile: profile,
		},
	}
}
