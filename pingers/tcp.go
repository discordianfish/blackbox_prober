package pingers

import (
	"log"
	"net"
	"net/url"
	"time"
)

func init() {
	pingers["tcp"] = pingerTCP
}

func pingerTCP(url *url.URL, m Metrics) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", url.Host, *timeout)
	if err != nil {
		log.Printf("Couldn't connect to %s: %s", url.Host, err)
		m.Up.WithLabelValues(url.String()).Set(0)
		return
	}
	defer conn.Close()
	if url.Path != "" {
		conn.SetDeadline(time.Now().Add(*timeout))

		size, err := readSize(conn)
		if err != nil {
			log.Printf("Error reading from %s: %s", url, err)
			m.Up.WithLabelValues(url.String()).Set(0)
		}
		m.Size.WithLabelValues(url.String()).Set(float64(size))
	}
	m.Latency.WithLabelValues(url.String()).Observe(time.Since(start).Seconds())
	m.Up.WithLabelValues(url.String()).Set(1)
}
