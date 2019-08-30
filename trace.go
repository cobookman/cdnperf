package main

import (
	"time"
	"net/http/httptrace"
	"crypto/tls"
)

type Trace struct {
	// Time request was started
	StartTime time.Time

	// DNS lookup begins
	DNSStart time.Time

	// DNS lookup ends
	DNSDone time.Time

	// TLS handshake is started
	TLSHandshakeStart time.Time

	// TLS handshake complete
	TLSHandshakeDone time.Time

	// TLS Version
	TLSVersion string

	// TLS Cipher Suite
	TLSCipherSuite string

	// new connection's Dial begins
	ConnectStart time.Time

	// new connection established
	ConnectDone time.Time

	// False when a new connection was established, and a connection
	// was not reused
	ConnectReused bool

	// Got first byte of response
	TTFB time.Time

	// Got last byte of response
	TTLB time.Time

	// Size of Body in Bytes
	BodySize int

	// Set if endpoint supports QUIC, with list of versions
	QUICSupport string

	// HTTP Version used
	HTTPVersion string

	// Status Code
	HTTPStatus string
}

func (trace *Trace) ClientTrace() (*httptrace.ClientTrace) {
	ct := new(httptrace.ClientTrace)
	ct.GetConn = func(hostPort string) {
		trace.StartTime = time.Now()
	}
	ct.GotConn = func(gci httptrace.GotConnInfo) {
		trace.ConnectReused = gci.Reused
	}
	ct.DNSStart = func(dsi httptrace.DNSStartInfo) {
		trace.DNSStart = time.Now()
	}
	ct.DNSDone = func(ddi httptrace.DNSDoneInfo) {
		trace.DNSDone = time.Now()
	}
	ct.TLSHandshakeStart = func() {
		trace.TLSHandshakeStart = time.Now()
	}
	ct.TLSHandshakeDone = func(cs tls.ConnectionState, err error) {
		trace.TLSHandshakeDone = time.Now()

		switch cs.Version {
		case tls.VersionSSL30:
			trace.TLSVersion = "SSL/3.0"
		case tls.VersionTLS10:
			trace.TLSVersion = "TLS/1.0"
		case tls.VersionTLS11:
			trace.TLSVersion = "TLS/1.1"
		case tls.VersionTLS12:
			trace.TLSVersion = "TLS/1.2"
		case tls.VersionTLS13:
			trace.TLSVersion = "TLS/1.3"
		default:
			trace.TLSVersion = "Unknown"
		}

		switch cs.CipherSuite {
		// TLS 1.0 - 1.2 cipher suites.
		case tls.TLS_RSA_WITH_RC4_128_SHA:
			trace.TLSCipherSuite = "TLS_RSA_WITH_RC4_128_SHA"
		case tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:
			trace.TLSCipherSuite = "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
		case tls.TLS_RSA_WITH_AES_128_CBC_SHA:
			trace.TLSCipherSuite = "TLS_RSA_WITH_AES_128_CBC_SHA"
		case tls.TLS_RSA_WITH_AES_256_CBC_SHA:
			trace.TLSCipherSuite = "TLS_RSA_WITH_AES_256_CBC_SHA"
		case tls.TLS_RSA_WITH_AES_128_CBC_SHA256:
			trace.TLSCipherSuite = "TLS_RSA_WITH_AES_128_CBC_SHA256"
		case tls.TLS_RSA_WITH_AES_128_GCM_SHA256:
			trace.TLSCipherSuite = "TLS_RSA_WITH_AES_128_GCM_SHA256"
		case tls.TLS_RSA_WITH_AES_256_GCM_SHA384:
			trace.TLSCipherSuite = "TLS_RSA_WITH_AES_256_GCM_SHA384"
		case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"
		case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
		case tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
		case tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_RC4_128_SHA"
		case tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA"
		case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
		case tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
		case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
		case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
		case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
		case tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
		case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
		case tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
		case tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:
			trace.TLSCipherSuite = "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"
		case tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:
			trace.TLSCipherSuite = "TLS_ECDHE_ECDSA_WITH_CACHA20_POLY1305"

		// TLS 1.3 cipher suites.
		case tls.TLS_AES_128_GCM_SHA256:
			trace.TLSCipherSuite = "TLS_AES_128_GCM_SHA256"
		case tls.TLS_AES_256_GCM_SHA384:
			trace.TLSCipherSuite = "TLS_AES_256_GCM_SHA384"
		case tls.TLS_CHACHA20_POLY1305_SHA256:
			trace.TLSCipherSuite = "TLS_CACHA20_POLY1305_SHA256"

		// TLS_FALLBACK_SCSV isn't a standard cipher suite but an indicator
		// that the client is doing version fallback. See RFC 7507.
		case tls.TLS_FALLBACK_SCSV:
			trace.TLSCipherSuite = "TLS_FALLBACK_SCSV"
		}
	}
	ct.ConnectStart = func(network, addr string) {
		trace.ConnectStart = time.Now()
	}
	ct.ConnectDone = func(network, addr string, err error) {
		trace.ConnectDone = time.Now()
	}
	ct.GotFirstResponseByte = func() {
		trace.TTFB = time.Now()
	}
	return ct
}
