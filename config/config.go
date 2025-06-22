package config

import (
	"flag"
	"os"
)

type Config struct {
	Port string
}

func InitConfig() *Config {
	port := flag.String("port", "50052", "port to listen on")

	flag.Parse()

	cfg := &Config{
		Port: *port,
	}

	if *port == "8080" {
		if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
			cfg.Port = envPort
		}
	}

	return cfg

}
