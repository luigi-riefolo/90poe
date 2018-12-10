package api

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	liblog "github.com/luigi-riefolo/nlp/lib/log"
	"github.com/luigi-riefolo/nlp/lib/metrics"
	"github.com/luigi-riefolo/nlp/lib/server"
	"github.com/luigi-riefolo/nlp/svc-storer/pb"
)

// TODO: test if connection is fine when storer not up

const (
	gracePeriod = 3 * time.Second
)

var (
	log = liblog.NewLogger()
)

// StorerService represents the ingestor service.
type StorerService struct {
	// config values for the service
	config Config

	server *grpc.Server

	// contains the list of user entries by ID
	userEntries map[string]*pb.Entry
}

// Config represents the StorerService configuration
// that is set via environment variables.
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`

	// version is injectected via ldflags
	Version string `ignored:"true"`

	Environment    string `default:"dev"`
	Port           int    `envconfig:"SERVICE_GRPC_PORT" default:"10080"`
	PrometheusPort int    `envconfig:"PROMETHEUS_PORT" default:"10060"`
}

// NewStorerService returns a ready-to-use StorerService.
func NewStorerService(config Config) (pb.StorerServer, error) {

	svc := &StorerService{
		config:      config,
		userEntries: map[string]*pb.Entry{},
	}

	log.Infof("%s service starting", config.ServiceName)
	log.Debugf("config: %#v", config)

	return svc, nil
}

// Start starts listening for incoming requests.
func (s *StorerService) Start(ctx context.Context) error {

	svcAddress := net.JoinHostPort("0.0.0.0", strconv.Itoa(s.config.Port))

	log.Printf("%s service listening on %v", s.config.ServiceName, svcAddress)

	lis, err := net.Listen("tcp", svcAddress)
	if err != nil {
		return errors.Wrap(err, "could not set up listener")
	}

	// make sure that log statements internal to gRPC
	// library are logged using the logrus Logger as well
	grpc_logrus.ReplaceGrpcLogger(log)

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	pb.RegisterStorerServer(s.server, s)

	metrics.StartPrometheus(s.server, s.config.PrometheusPort)

	server.HandleSignals(
		func() {
			s.Stop(ctx)
		})

	if err := s.server.Serve(lis); err != nil {
		log = log.WithError(err)
	}
	log.Info("server stopped")

	return nil
}

// Stop gracefully shuts down the server.
func (s *StorerService) Stop(ctx context.Context) {
	if err := metrics.StopPrometheus(ctx); err != nil {
		log = log.WithError(err)
	}

	s.server.GracefulStop()
	time.Sleep(gracePeriod)
}
