package server

import (
	"context"
	"github.com/ewik2k21/grpcSpotInstrumentService/config"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/handlers"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/repositories"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/services"
	spot_instrument_v1 "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	wg := sync.WaitGroup{}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(),
	)

	//init all layers
	spotInstrumentRepo := repositories.NewSpotInstrumentRepository(logger)
	spotInstrumentService := services.NewSpotInstrumentService(*spotInstrumentRepo, logger)
	spotInstrumentHandler := handlers.NewSpotInstrumentHandler(*spotInstrumentService)

	spot_instrument_v1.RegisterSpotInstrumentServiceServer(grpcServer, spotInstrumentHandler)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		logger.Error("failed listen tcp server ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server listening at %v", lis.Addr())

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start tcp server")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed start grpc server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("received shutdown signal, start graceful shutdown")

	grpcServer.GracefulStop()
	logger.Info("server stopped")

	wg.Wait()
	logger.Info("all stopped")
}
