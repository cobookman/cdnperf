package main

import (
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

	for i := 0; i < *iterations; i++ {
		trace, err := test(*url)
		if err != nil {
			log.Fatal(err)
		}

		ttFBs[i] = trace.TTFB.Sub(trace.StartTime).Seconds() * 1000
		ttLBs[i] = trace.TTLB.Sub(trace.StartTime).Seconds() * 1000
		bodySizes[i] = float64(trace.BodySize)
	        log.Printf("\tBody Size: %d KiB\n", trace.BodySize/1024)
	        log.Printf("\tttfb: %s\n", trace.TTFB.Sub(trace.StartTime))
	        log.Printf("\tttlb: %s\n", trace.TTLB.Sub(trace.StartTime))
	        log.Printf("\ttime creating connection: %s\n", trace.ConnectDone.Sub(trace.ConnectStart))
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
	BodySize: %.0f KiB
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
	medianTTFB, p95TTFB, p99TTFB,
	medianTTLB, p95TTLB, p99TTLB)
}

func test(url string) (*Trace, error) {
	// Declare Buffer for storing response in
	buf := make([]byte, 1024*1024) // 1MiB

	req, _ := http.NewRequest("GET", url, nil)
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
	/** log.Printf("\tBody Size: %d KiB\n", size/1024)
	log.Printf("\tttfb: %s\n", trace.TTFB.Sub(trace.StartTime))
	log.Printf("\tttlb: %s\n", trace.TTLB.Sub(trace.StartTime))
	log.Printf("\ttime creating connection: %s\n", trace.ConnectDone.Sub(trace.ConnectStart)) **/
	return trace, nil
}


