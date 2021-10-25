package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) initTransport(ctx context.Context, wgOpts runtime.ServeMuxOption) (err error) {
	gatewayMux := runtime.NewServeMux(wgOpts)
	s.grpcServer = grpc.NewServer()
	reflection.Register(s.grpcServer)

	if err = s.registrar.RegisterHandlers(ctx, s.grpcServer, gatewayMux); err != nil {
		return fmt.Errorf("error registering handlers: %s", err)
	}

	mux := http.NewServeMux()
	if err := serveSwaggerUI(mux); err != nil {
		return fmt.Errorf("error serving swagger-ui: %s", err)
	}
	mux.Handle("/", gatewayMux)

	s.httpServer = &http.Server{
		Handler:        cors(mux, s.allowedOrigins),
		ReadTimeout:    time.Second * 15, // @TODO: move to config
		WriteTimeout:   time.Second * 15,
		IdleTimeout:    time.Second * 15,
		MaxHeaderBytes: 2 << 11,
	}

	return nil
}
