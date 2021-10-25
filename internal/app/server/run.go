package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Run starts accepting new connections and waits for signals from the OS to shut down gracefully.
func (s Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	errWg := errgroup.Group{}

	errWg.Go(func() error {
		lis, err := net.Listen("tcp", s.grpcAddr)
		if err != nil {
			return err
		}

		return s.grpcServer.Serve(lis)
	})

	errWg.Go(func() error {
		l, err := net.Listen("tcp", s.httpAddr)
		if err != nil {
			return err
		}

		return s.httpServer.Serve(l)
	})

	errWg.Go(func() error {
		swaggerAddr := s.httpAddr + swaggerUIPrefix
		log.Info().Msgf("App started. HTTP: %s, Swagger UI: %s, gRPC: %s", s.httpAddr, swaggerAddr, s.grpcAddr)
		return nil
	})

	errWg.Go(func() error {
		shutdownCh := make(chan os.Signal, 1)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		sig := <-shutdownCh

		s.stop(ctx)
		cancel()

		log.Fatal().Msgf("exit reason: %s", sig)

		return nil
	})

	return errWg.Wait()
}

// Stop the gRPC and HTTP servers.
func (s Server) stop(ctx context.Context) {
	s.grpcServer.Stop()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Err(err).Msg("error shutting down the http server")
	}
}
