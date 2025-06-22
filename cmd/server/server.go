package server

import (
	"context"
	"github.com/ewik2k21/grpcSpotInstrumentService/config"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/handlers"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/interceptors"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/repositories"
	"github.com/ewik2k21/grpcSpotInstrumentService/internal/services"
	spot_instrument_v1 "github.com/ewik2k21/grpcSpotInstrumentService/pkg/spot_instrument_v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	wg := sync.WaitGroup{}

	//chain interceptors for logging request and more...
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor(),
			interceptors.LoggerRequestInterceptor(logger),
			interceptors.PrometheusInterceptor(),
			interceptors.UnaryPanicRecoveryInterceptor(logger),
		),
	)

	//init all layers
	spotInstrumentRepo := repositories.NewSpotInstrumentRepository(logger)
	spotInstrumentService := services.NewSpotInstrumentService(*spotInstrumentRepo, logger)
	spotInstrumentHandler := handlers.NewSpotInstrumentHandler(*spotInstrumentService)

	spot_instrument_v1.RegisterSpotInstrumentServiceServer(grpcServer, spotInstrumentHandler)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logger.Error("failed listen tcp server ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server listening at %s", slog.Any("port", lis.Addr()))

	//start grpc server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start tcp server")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed start grpc server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	//cfg for httpserver for metrics
	metricsServer := &http.Server{
		Addr: cfg.HttpPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				promhttp.Handler().ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	//start server for metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("metrics endpoint start on :", slog.String("port", cfg.HttpPort))
		http.Handle("/metrics", promhttp.Handler())
		if err := metricsServer.ListenAndServe(); err != nil {
			logger.Error("metrics endpoint failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("received shutdown signal, start graceful shutdown")

	//shutdown metrics server
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("metrics server shutdown failed", slog.String("error", err.Error()))
	}

	//shutdown grpc
	grpcServer.GracefulStop()
	logger.Info("server stopped")

	wg.Wait()
	logger.Info("all stopped")
}
