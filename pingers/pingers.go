package pingers

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Namespace for prober
const Namespace = "blackbox"

var (
	timeout = flag.Duration("ping.timeout", 10*time.Second, "Timeout for requests.")

	// ErrUnsupportedScheme will returned if no pinger function exists for given scheme
	ErrUnsupportedScheme = errors.New("Scheme not supported")

	pingers = make(map[string]func(url *url.URL, m Metrics))
)

// Metrics holds the prometheus metrics common among all pingers
type Metrics struct {
	Up      *prometheus.GaugeVec
	Latency *prometheus.GaugeVec
	Size    *prometheus.GaugeVec
	Expires *prometheus.GaugeVec
}

func readSize(r io.Reader) (int, error) {
	size := 0
	buf := make([]byte, bytes.MinRead) // Since we discard the buffer, alloc only once
	for {
		n, err := r.Read(buf)
		size += n
		if err != nil {
			if err == io.EOF {
				return size, nil
			}
			return size, err
		}
	}
}

// CanHandle return true if url is supported.
func CanHandle(url *url.URL) bool {
	_, ok := pingers[url.Scheme]
	return ok
}

// Ping executes the matching pinger function for the url.
// If no pinger function can be found, it return ErrUnsupportedScheme.
func Ping(url *url.URL, m Metrics) error {
	pinger, ok := pingers[url.Scheme]
	if !ok {
		return ErrUnsupportedScheme
	}
	pinger(url, m)
	return nil
}
