package metrics

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	liblog "github.com/luigi-riefolo/nlp/lib/log"
)

var (
	log = liblog.NewLogger().
		WithField("component", "metrics")
	server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
)

// StartPrometheus starts the Prometheus server.
func StartPrometheus(srv *grpc.Server, port int) {

	// after all your registrations, make sure all
	// of the Prometheus metrics are initialized.
	grpc_prometheus.Register(srv)

	mux := http.NewServeMux()

	// register Prometheus metrics handler.
	mux.Handle("/metrics", promhttp.Handler())

	server.Addr = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	server.Handler = mux

	log.Debugf("prometheus serving metrics on %s", server.Addr)

	go func() {
		server.ListenAndServe()
	}()
}

// StopPrometheus gracefully stops the Prometheus server.
func StopPrometheus(ctx context.Context) error {
	if err := server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "could not stop Prometheus server")
	}
	log.Debug("Prometheus server stopped")

	return nil
}
