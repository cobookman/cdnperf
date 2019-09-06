package main

import (
	"testing"
)

func TestTcpRtt(t *testing.T) {
	if _, err := TcpRtt("https://google.com", 5); err != nil {
		t.Error(err)
	}
	if _, err := TcpRtt("http://google.com?awdf=s&t=234", 5); err != nil {
		t.Error(err)
	}
}
