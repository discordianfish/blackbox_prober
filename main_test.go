package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"time"
)

func TestCollector(t *testing.T) {

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add a 2 second delay to response time
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "hello")
	}))

	targets := targets{}
	targets.Set(httpServer.URL)
	targets.Set(httpServer.URL)

	collector := NewPingCollector(targets)
	ch := make(chan prometheus.Metric)

	start_time := time.Now()

	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	for m := range ch {
		pb := &dto.Metric{}
		m.Write(pb)
	}

	if time.Since(start_time).Seconds() > 4.0 {
		t.Fatalf("Elapsed time was more than 4s, indicating non-concurrent checks.")
	}

}
