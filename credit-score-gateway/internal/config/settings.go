package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrpcPort                string `envconfig:"GRPC_PORT"`
	HttpPort                string `envconfig:"HTTP_PORT"`
	TcpPort                 string `envconfig:"TCP_PORT"`
	ChampionBaseUrl         string `envconfig:"CHAMPION_BASE_URL"`
	ChampionTimeoutMillis   int    `envconfig:"CHAMPION_TIMEOUT_MILLIS"`
	ChallengerBaseUrl       string `envconfig:"CHALLENGER_BASE_URL"`
	ChallengerTimeoutMillis int    `envconfig:"CHALLENGER_TIMEOUT_MILLIS"`
}

func LoadConfig() Config {
	godotenv.Load(".env")
	cfg := Config{}
	if err := envconfig.Process("CSG", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Config: %+v", cfg)
	return cfg
}
