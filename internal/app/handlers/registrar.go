package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	sink "github.com/go-sink/sink/pkg/sink/v1"
)

// Registrar is registrar of gRPC and gRPC-Gateway handlers.
type Registrar struct {
	sinkServer sink.SinkServiceServer
}

// NewRegistrar returns new registrar instance.
func NewRegistrar(sink sink.SinkServiceServer) Registrar {
	return Registrar{
		sinkServer: sink,
	}
}

// RegisterHandlers registers all handlers for grpcServer and runtime.ServeMux.
// Moved separately to avoid mixing business logic and infrastructure.
func (r Registrar) RegisterHandlers(ctx context.Context, grpcServer *grpc.Server, wgMux *runtime.ServeMux) error {
	// register gRPC
	sink.RegisterSinkServiceServer(grpcServer, r.sinkServer)
	// register gRPC Gateway
	err := sink.RegisterSinkServiceHandlerServer(ctx, wgMux, r.sinkServer)
	if err != nil {
		return err
	}

	return nil
}
