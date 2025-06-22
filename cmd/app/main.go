package main

import (
	"context"
	"github.com/ewik2k21/grpcSpotInstrumentService/cmd/server"
	config "github.com/ewik2k21/grpcSpotInstrumentService/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.InitConfig() // init config

	ctx := context.Background() // create ctx

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	server.Execute(ctx, cfg, logger) // execute grpc server
}
