package main

import (
	"strings"
	"io"
	"time"
	"flag"
	"log"
	"net/http"
	"net/http/httptrace"
	"github.com/montanaflynn/stats"
	"fmt"
)

var (
	httpClient *http.Client
)

const (
	MaxIdleConnections int = 10
	RequestTimeout int = 0
)

func init() {
	httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

func main() {
	url := flag.String("url", "test", "Url to performance profile")
	iterations := flag.Int("iterations", 100, "number of times to send requests")
	flag.Parse();

	log.Print(*url)

	// in Milliseconds
	ttFBs := make([]float64, *iterations)
	ttLBs := make([]float64, *iterations)
	bodySizes := make([]float64, *iterations)

	// do one request to establish connection (sunk cost)
	firstTrace, err := test(*url)
	if err != nil {
		log.Fatal(err)
	}

	// now that connection *should* be established, run tests
	for i := 0; i < *iterations; i++ {
		trace, err := test(*url)
		if err != nil {
			log.Fatal(err)
		}

		ttFBs[i] = trace.TTFB.Sub(trace.StartTime).Seconds() * 1000
		ttLBs[i] = trace.TTLB.Sub(trace.StartTime).Seconds() * 1000
		bodySizes[i] = float64(trace.BodySize)
		log.Printf("Iteration #%d", (i+1))
		log.Printf("\tConnection Reuse: %t\n", trace.ConnectReused)
		log.Printf("\tHTTP Status: %s\n", trace.HTTPStatus)
	        log.Printf("\tBody Size: %d KiB\n", trace.BodySize/1024)
	        log.Printf("\tttfb: %s\n", trace.TTFB.Sub(trace.StartTime))
	        log.Printf("\tttlb: %s\n", trace.TTLB.Sub(trace.StartTime))
	        log.Printf("\ttime creating connection: %s\n", trace.ConnectDone.Sub(trace.ConnectStart))
		log.Printf("\tQuic Supported: %s\n", trace.QUICSupport)

	}

	// calculate percentiles
	medianBodySize, _ := stats.Median(bodySizes)

	medianTTFB, _ := stats.Median(ttFBs)
	p95TTFB, _ := stats.Percentile(ttFBs, 95)
	p99TTFB, _ := stats.Percentile(ttFBs, 99)

	medianTTLB, _ := stats.Median(ttLBs)
	p95TTLB, _ := stats.Percentile(ttLBs, 95)
	p99TTLB, _ := stats.Percentile(ttLBs, 99)

	fmt.Println("Statistics:")
	fmt.Printf(
`
	URL: %s

	-- From frist response --
	BodySize: %.0f KiB
	CipherSuite: %s
	TLS Version: %s
	Quic: %s
	HTTP Version: %s
	HTTP Status: %s
	--

	Time to First Byte:
		Median: %.2f ms
		95th: %.2f ms
		99th: %.2f ms

	Time to Last Byte:
		Median: %.2f ms
		95th: %.2f ms
		99th: %.2f ms
`,
	*url,
	(medianBodySize / 1024),
	firstTrace.TLSCipherSuite,
	firstTrace.TLSVersion,
	firstTrace.QUICSupport,
	firstTrace.HTTPVersion,
	firstTrace.HTTPStatus,
	medianTTFB, p95TTFB, p99TTFB,
	medianTTLB, p95TTLB, p99TTLB)
}

func test(url string) (*Trace, error) {
	// Declare Buffer for storing response in
	buf := make([]byte, 1024*1024) // 1MiB

	req, _ := http.NewRequest("GET", url, nil)

	// Explicitly set Accept GZIP Encoding so client lib doesn't
	// de-compress paypload.
	// see: https://golang.org/pkg/net/http/#Header
	//
	// ...If the Transport requests gzip on
        // its own and gets a gzipped response, it's transparently
        // decoded in the Response.Body. However, if the *user
        // explicitly requested gzip* it is not automatically
        // uncompressed.
	req.Header.Add("Accept-Encoding", "gzip")

	trace := new(Trace)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace.ClientTrace()))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Close Connection when done so we can re-use the connection
	defer resp.Body.Close();

	// Read in data into dummy buffer.
	// This is better than reading entire object in RAM which would fail for substantially large files
	trace.BodySize = 0
	done := false
	for !done {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err

		}
		if n == 0 || err == io.EOF {
			done = true
		}
		trace.BodySize += n
	}
	trace.TTLB = time.Now()
	trace.QUICSupport = quicVersion(resp.Header.Get("alt-svc"))
	trace.HTTPVersion = resp.Proto
	trace.HTTPStatus = resp.Status
	return trace, nil
}

// Detects what QUIC versions supported in alt-svc header string
func quicVersion (altSvc string) string {
	if len(altSvc) == 0 {
		return ""
	}

	parts := strings.Split(altSvc, ";")
	services := strings.Split(parts[0], ",")

	loc := ""
	for _, service := range services {
		service = strings.TrimSpace(service)
		if strings.HasPrefix(service, "quic=") {
			loc = strings.Replace(service, "quic=\"", "", 1)
			loc = strings.Replace(loc, "\"", "", -1)
		}
	}

	versions := ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "v=\"") {
			versions = strings.Replace(part, "v=\"", "", 1)
			versions = strings.Replace(versions, "\"", "", -1)
		}
	}

	out := "No"
	if len(loc) != 0 {
		out = "Yes, at " + loc
	}
	if len(versions) != 0 {
		out += " with versions: " + versions
	}
	return out
}


