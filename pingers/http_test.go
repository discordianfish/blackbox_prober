package pingers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

const payload = "hello"

func TestHTTP(t *testing.T) {
	serverHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payload)
	}))
	serverHTTPS := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payload)
	}))

	*insecure = true
	for _, server := range []*httptest.Server{serverHTTP, serverHTTPS} {
		u, err := url.Parse(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		pingerHTTP(u, metrics)

		pb := &dto.Metric{}
		metrics.Up.WithLabelValues(server.URL).Write(pb)
		if expected, got := 1., pb.Gauge.GetValue(); expected != got {
			t.Fatalf("Expected: %f, Got: %f", expected, got)
		}

		metrics.Size.WithLabelValues(server.URL).Write(pb)
		if expected, got := float64(len(payload)), pb.Gauge.GetValue(); expected != got {
			t.Fatalf("Expected: %f, Got: %f", expected, got)
		}

		metrics.Latency.WithLabelValues(server.URL).Write(pb)
		if pb.Gauge.GetValue() == 0 {
			t.Fatal("Expected non-zero value")
		}

		statusCode.WithLabelValues(server.URL).Write(pb)
		if expected, got := float64(200), pb.Gauge.GetValue(); expected != got {
			t.Fatalf("Expected: %f, Got: %f", expected, got)
		}

	}
}

func TestHTTPSExpire(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payload)
	}))

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	pingerHTTP(u, metrics)

	pb := &dto.Metric{}
	expireTimestamp.WithLabelValues(server.URL).Write(pb)
	if expected, got := 3.6e+09, pb.Gauge.GetValue(); expected != got {
		t.Fatalf("Expected: %f, Got: %f", expected, got)
	}
}
