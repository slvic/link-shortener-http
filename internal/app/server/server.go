package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/go-sink/sink/internal/app/config"
	"github.com/go-sink/sink/internal/app/handlers"
	"github.com/go-sink/sink/internal/app/handlers/sinkapi"
	"github.com/go-sink/sink/internal/app/repository"
	"github.com/go-sink/sink/internal/cron"
	"github.com/go-sink/sink/internal/pkg/gateway"

	"github.com/go-sink/sink/internal/app/service"
	"github.com/go-sink/sink/internal/pkg/urlgen"
)

// Server contains application dependencies.
type Server struct {
	httpAddr, grpcAddr string
	allowedOrigins []string
	grpcServer         *grpc.Server
	httpServer         *http.Server
	registrar          handlers.Registrar
}

// InitApp initializes handlers and transport.
func InitApp(ctx context.Context, config config.Config) (*Server, error) {
	s := &Server{
		allowedOrigins: config.App.AllowedOrigins,
		httpAddr: config.App.HTTPAddr,
		grpcAddr: config.App.GRPCAddr,
	}

	// set up database
	dbcfg := config.Database
	db, err := setUpDb(dbcfg)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %s", err)
	}

	random := rand.New(rand.NewSource(time.Now().Unix())) //nolint:gosec
	linkRepo := repository.NewLinkRepository(db)
	linkStatusChecker := cron.NewLinkStatusChecker(linkRepo)

	// TODO: Add leader election
	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Every(10).Second().Do(linkStatusChecker.CheckLinks, ctx)
	if err != nil {
		return nil, err
	}
	scheduler.StartAsync()

	urlEncoder := service.NewGenerator(urlgen.NewRandomURLGenerator(random), linkRepo)

	sinkAPIHandler := sinkapi.New(urlEncoder)

	s.registrar = handlers.NewRegistrar(sinkAPIHandler)

	if err = s.initTransport(ctx, runtime.WithForwardResponseOption(gateway.ReplaceHeaders())); err != nil {
		return nil, fmt.Errorf("error initializing transport: %w", err)
	}

	return s, nil
}
