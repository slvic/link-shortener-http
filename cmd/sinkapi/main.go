package main

import (
	"context"
	_ "embed"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-sink/sink/internal/app/config"
	"github.com/go-sink/sink/internal/app/server"
)

var configPath = flag.String("config", "configs/app.example.hcl", "application config")

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	cfgBytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config file")
	}
	cfg, err := config.Parse(cfgBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse config")
	}

	app, err := server.InitApp(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init app")
	}

	if err = app.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("could not run app")
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh)

	sig := <-shutdownCh

	cancel()

	log.Fatal().Msgf("exit reason: %s\n", sig)
}
