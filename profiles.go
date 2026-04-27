package architect

import (
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

var (
	// Chrome124 represents a Chrome 124 fingerprint on MacOS
	Chrome124 = Profile{
		ID:    "chrome_124_macos",
		Name:  "Chrome 124 MacOS",
		TLSID: utls.HelloChrome_120,
		HTTP2Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingMaxConcurrentStreams, Val: 1000},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxFrameSize, Val: 16384},
			{ID: http2.SettingHeaderTableSize, Val: 65536}, // Duplicate to match common fingerprint
		},
		InitialWinSize:    6291456,
		MaxHeaderListSize: 262144,
		UserAgent:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		TTL:               64,
		PseudoHeaderOrder: []string{":method", ":authority", ":scheme", ":path"},
		HeaderOrder: []string{
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
		},
	}

	// Safari17iOS represents a Safari 17 fingerprint on iOS
	Safari17iOS = Profile{
		ID:    "safari_17_ios",
		Name:  "Safari 17 on iOS",
		TLSID: utls.HelloSafari_16_0,
		HTTP2Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 4096},
			{ID: http2.SettingEnablePush, Val: 1},
			{ID: http2.SettingInitialWindowSize, Val: 2097152},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		},
		InitialWinSize:    2097152,
		MaxHeaderListSize: 262144,
		UserAgent:         "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
		TTL:               64,
		PseudoHeaderOrder: []string{":method", ":scheme", ":path", ":authority"},
		HeaderOrder: []string{
			"user-agent",
			"accept",
			"accept-language",
			"accept-encoding",
			"connection",
		},
	}
)
