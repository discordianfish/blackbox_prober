package pingers

import (
	"net/url"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"os/exec"
)

func TestICMP(t *testing.T) {

	_, err := exec.LookPath("ping")
	if err != nil {
		t.Log("Unable to test ICMP ping, as the ping executable is not in the path")
		return
	}

	u, err := url.Parse("icmp://localhost")
	if err != nil {
		t.Fatal(err)
	}
	pingerICMP(u, metrics)

	pb := &dto.Metric{}
	metrics.Up.WithLabelValues(u.String()).Write(pb)
	if expected, got := 1., pb.Gauge.GetValue(); expected != got {
		t.Fatalf("Expected: %f, Got: %f", expected, got)
	}

	metrics.Latency.WithLabelValues(u.String()).Write(pb)
	if pb.Gauge.GetValue() == 0 {
		t.Fatal("Expected non-zero value")
	}
}
