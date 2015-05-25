package pingers

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	insecure = flag.Bool("ping.insecure", false, "Disable validation of server certificate for https.")

	expireTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "cert_expire_timestamp",
		Help:      "Certificate expiry date in seconds since epoch.",
	}, []string{"url"})
)

func init() {
	prometheus.MustRegister(expireTimestamp)
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

	m.Latency.WithLabelValues(url.String()).Set(time.Since(start).Seconds())
	m.Size.WithLabelValues(url.String()).Set(float64(size))
	m.Up.WithLabelValues(url.String()).Set(1)

	if resp.TLS != nil {
		var expires time.Time
		if *insecure { // If insecure, we check the unverified certs
			expires = resp.TLS.PeerCertificates[0].NotAfter
		} else {
			expires = resp.TLS.VerifiedChains[0][0].NotAfter
		}
		expireTimestamp.WithLabelValues(url.String()).Set(float64(expires.Unix()))
	}
}
