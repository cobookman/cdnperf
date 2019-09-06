package main

import (
	"log"
	"time"
	"net"
	"net/url"
)


func TcpRtt(uri string, iterations int) (*[]time.Duration, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	addr := u.Hostname() + ":80"
	results := make([]time.Duration, iterations)

	// do 1 connection before hand in-case DNS lookup required
	// disregarding any DNS lookup latency impact
	if _, err := net.Dial("tcp", addr); err != nil {
		return nil, err
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		_, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		results[i] = time.Since(start)
		log.Printf("TCP RTT [#%d] - %s - %d ms\n", i, addr, results[i].Milliseconds())
	}
	return &results, nil
}
