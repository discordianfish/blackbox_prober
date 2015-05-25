package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var (
	upDesc  = prometheus.NewDesc("blackbox_up", "1 if url is reachable, 0 if not", nil, nil)
	expectedMetrics = []testMetric{
		testMetric{desc: upDesc, value: 1, url}
	}
}
)

func TestCollector(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	}))

	targets := targets{}
	targets.Set(httpServer.URL)

	collector := NewPingCollector(targets)
	ch := make(chan prometheus.Metric)

	go func() {
		e.Collect(ch)
		close(ch)
	}()

	for m := range ch {
		pb := &dto.Metric{}
		m.Write(pb)

	}
}
