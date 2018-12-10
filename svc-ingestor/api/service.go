package api

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
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
	"github.com/luigi-riefolo/nlp/svc-ingestor/pb"
	storerpb "github.com/luigi-riefolo/nlp/svc-storer/pb"
)

// TODO: test if connection is fine when storer not up
// - normalise phone and email
// - stop current processes if signal caught

const (
	gracePeriod        = 3 * time.Second
	phoneRegion        = "GB"
	storeClientTimeout = 5 * time.Second
)

var (
	log = liblog.NewLogger()

	// pool for storer requests
	pool = &sync.Pool{
		New: func() interface{} {
			return &storerpb.StoreEntryRequest{}
		},
	}
)

// IngestorService represents the ingestor service.
type IngestorService struct {
	// config values for the service
	config Config

	server *grpc.Server

	// gRPC client for the Storer service
	storerClient storerpb.StorerClient
}

// Config represents the IngestorService configuration
// that is set via environment variables.
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`

	// version is injectected via ldflags
	Version string `ignored:"true"`

	Environment    string `default:"dev"`
	Port           int    `envconfig:"SERVICE_GRPC_PORT" default:"10080"`
	PrometheusPort int    `envconfig:"PROMETHEUS_PORT" default:"10060"`

	// the local file containing the list of entries to be parsed
	DataFile string `envconfig:"INGESTOR_DATA_FILE" required:"true"`

	StorerServiceHost string `envconfig:"STORER_SERVICE_GRPC_HOST" default:"0.0.0.0"`
	StorerServicePort int    `envconfig:"STORER_SERVICE_GRPC_PORT" default:"10080"`
}

// NewIngestorService returns a ready-to-use IngestorService.
func NewIngestorService(config Config) (pb.IngestorServer, error) {

	svc := &IngestorService{
		config: config,
	}

	log.Infof("%s service starting", config.ServiceName)
	log.Debugf("config: %#v", config)

	return svc, nil
}

// ConnectStorerClient initialise the connection with the Storer service gRPC client.
// NOTE: the client connection is only initialised and not tested for
// reliability and/or connectivity.
func (i *IngestorService) ConnectStorerClient(ctx context.Context) error {

	target := fmt.Sprintf("%s:%d",
		i.config.StorerServiceHost,
		i.config.StorerServicePort)

	// TODO: add interceptors
	conn, err := grpc.DialContext(
		ctx,
		target,
		// unsecured connection are feasible
		// only over trusted/private networks
		grpc.WithInsecure(),

		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
			grpc_logrus.UnaryClientInterceptor(log),
		)),
	)
	if err != nil {
		return errors.Wrap(err, "could not connect to Storer service")
	}

	i.storerClient = storerpb.NewStorerClient(conn)
	log.Debugf("connected to Storer gRPC server on '%s'", target)

	return nil
}

// Start starts listening for incoming requests.
func (i *IngestorService) Start(ctx context.Context) error {

	svcAddress := net.JoinHostPort("0.0.0.0", strconv.Itoa(i.config.Port))

	lis, err := net.Listen("tcp", svcAddress)
	if err != nil {
		return errors.Wrap(err, "could not set up listener")
	}

	log.Printf("%s service listening on %v", i.config.ServiceName, svcAddress)

	// make sure that log statements internal to gRPC
	// library are logged using the logrus Logger as well
	grpc_logrus.ReplaceGrpcLogger(log)

	i.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	pb.RegisterIngestorServer(i.server, i)

	metrics.StartPrometheus(i.server, i.config.PrometheusPort)

	server.HandleSignals(
		func() {
			i.Stop(ctx)
		})

	// start processing the data file
	// NOTE: this action needs to be triggered via a specific endpoint.
	// The developer decided to initiate the process here due to lack of
	// specifications in the requirements.
	if err := i.ingest(); err != nil {
		log.WithError(err).Error("could not ingest data file")
	}

	if err := i.server.Serve(lis); err != nil {
		log = log.WithError(err)
	}
	log.Info("server stopped")

	return nil
}

// Stop gracefully shuts down the server.
func (i *IngestorService) Stop(ctx context.Context) {
	if err := metrics.StopPrometheus(ctx); err != nil {
		log = log.WithError(err)
	}

	i.server.GracefulStop()
	time.Sleep(gracePeriod)
}
