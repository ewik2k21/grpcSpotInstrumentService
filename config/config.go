package config

import (
	"flag"
	"os"
)

type Config struct {
	GRPCPort string
	HttpPort string
}

func InitConfig() *Config {
	grpcPort := flag.String("grpcPort", ":50052", "grpcPort to listen on")
	httpPort := flag.String("httpPort", ":2112", "httpPort for metrics")

	flag.Parse()

	cfg := &Config{
		GRPCPort: *grpcPort,
		HttpPort: *httpPort,
	}

	if *grpcPort == ":50052" {
		if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
			cfg.GRPCPort = envPort
		}
	}

	if *httpPort == ":2112" {
		if envPort := os.Getenv("HTTP_PORT"); envPort != "" {
			cfg.HttpPort = envPort
		}
	}
	return cfg

}
