package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/luigi-riefolo/nlp/lib/log"
	ingestor "github.com/luigi-riefolo/nlp/svc-ingestor/api"
)

var (
	// version is injectect via ldflags
	version = "undefined"
)

func main() {

	log := log.NewLogger()

	var config ingestor.Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal("could not load the environment variables: ", err)
	}
	config.Version = version

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// profiling handler
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	svc, err := ingestor.NewIngestorService(config)
	if err != nil {
		log.WithError(err).Fatal("could not create ingestor service")
	}

	ingestorSvc := svc.(*ingestor.IngestorService)

	if err := ingestorSvc.ConnectStorerClient(ctx); err != nil {
		log.WithError(err).Fatal("could not create client connection")
	}

	if err := ingestorSvc.Start(ctx); err != nil {
		log.WithError(err).Fatal("service failure")
	}
}
