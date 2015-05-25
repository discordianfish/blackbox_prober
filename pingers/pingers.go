package pingers

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
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
	Latency *prometheus.HistogramVec
	Size    *prometheus.GaugeVec
	Expires *prometheus.GaugeVec
}

func readSize(r io.Reader) (int, error) {
	size := 0
	buf := make([]byte, bytes.MinRead) // Since we discard the buffer, alloc only once
	for {
		log.Print("reading now")
		n, err := r.Read(buf)
		log.Printf("-> finished")
		size += n
		if err != nil {
			if err == io.EOF {
				log.Printf("EOF")
				return size, nil
			}
			log.Printf("Other error")
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
