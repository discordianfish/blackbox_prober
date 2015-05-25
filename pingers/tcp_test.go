package pingers

import (
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
)

func TestTCPConnect(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	// defer listener.Close()
	go serve(t, listener)

	u, err := url.Parse("tcp://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	pingerTCP(u, metrics)

	pb := &dto.Metric{}
	metrics.Up.WithLabelValues(u.String()).Write(pb)
	if expected, got := 1., pb.Gauge.GetValue(); expected != got {
		t.Fatalf("Expected: %f, Got: %f", expected, got)
	}

	metrics.Size.WithLabelValues(u.String()).Write(pb)
	if expected, got := 0., pb.Gauge.GetValue(); expected != got {
		t.Fatalf("Expected: %f, Got: %f", expected, got)
	}

	metrics.Latency.WithLabelValues(u.String()).Write(pb)
	if pb.Gauge.GetValue() == 0 {
		t.Fatal("Expected non-zero value")
	}

	// Again this time with reading until EOF
	u.Path = "/foo"
	pingerTCP(u, metrics)

	metrics.Size.WithLabelValues(u.String()).Write(pb)
	if expected, got := float64(len(payload)), pb.Gauge.GetValue(); expected != got {
		t.Fatalf("Expected: %f, Got: %f", expected, got)
	}
}

func serve(t *testing.T, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(100 * time.Millisecond)
		fmt.Fprint(conn, payload)
		conn.Close()
	}
}
