package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/luigi-riefolo/90poe/lib/log"
	storer "github.com/luigi-riefolo/90poe/svc-storer/api"
)

var (
	// version is injectect via ldflags
	version = "undefined"
)

func main() {

	log := log.NewLogger()

	var config storer.Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal("could not load the environment variables: ", err)
	}
	config.Version = version

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	svc, err := storer.NewStorerService(config)
	if err != nil {
		log.WithError(err).Fatalf("could not create %s service", config.ServiceName)
	}

	if err := svc.(*storer.StorerService).Start(ctx); err != nil {
		log.WithError(err).Fatal("service failure")
	}
}
