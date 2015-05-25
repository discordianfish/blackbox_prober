package pingers

import (
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func init() {
	pingers["icmp"] = pingerICMP
}

func pingerICMP(url *url.URL, m Metrics) {
	hostPort := strings.Split(url.Host, ":")
	start := time.Now()
	err := exec.Command("ping", "-n", "-c", "1", "-W", strconv.Itoa(int(timeout.Seconds())), hostPort[0]).Run()
	if err != nil {
		log.Printf("Couldn't ping %s: %s", url, err)
		m.Up.WithLabelValues(url.String()).Set(0)
		return
	}
	m.Latency.WithLabelValues(url.String()).Set(time.Since(start).Seconds())
	m.Up.WithLabelValues(url.String()).Set(1)
}
