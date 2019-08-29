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

	// new connection's Dial begins
	ConnectStart time.Time

	// new connection established
	ConnectDone time.Time

	// Got first byte of response
	TTFB time.Time

	// Got last byte of response
	TTLB time.Time

	// Size of Body in Bytes
	BodySize int
}

func (trace *Trace) ClientTrace() (*httptrace.ClientTrace) {
	ct := new(httptrace.ClientTrace)
	ct.GetConn = func(hostPort string) {
		trace.StartTime = time.Now()
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
