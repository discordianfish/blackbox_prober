package pingers

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	insecure = flag.Bool("ping.insecure", false, "Disable validation of server certificate for https.")
)

func init() {
	pingers["http"] = pingerHTTP
	pingers["https"] = pingerHTTP
}

func pingerHTTP(url *url.URL, m Metrics) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: *insecure},
			DisableKeepAlives: true,
		},
		Timeout: *timeout,
	}
	start := time.Now()
	resp, err := client.Get(url.String())
	if err != nil {
		log.Printf("Couldn't get %s: %s", url, err)
		m.Up.WithLabelValues(url.String()).Set(0)
		return
	}
	defer resp.Body.Close()
	size, err := readSize(resp.Body)
	if err != nil {
		log.Printf("Couldn't read from %s: %s", url, err)
	}

	m.Latency.WithLabelValues(url.String()).Observe(time.Since(start).Seconds())
	m.Size.WithLabelValues(url.String()).Set(float64(size))
	m.Up.WithLabelValues(url.String()).Set(1)
}
