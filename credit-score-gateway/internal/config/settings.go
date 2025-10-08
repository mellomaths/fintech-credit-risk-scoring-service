package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GrpcPort           string `envconfig:"GRPC_PORT"`
	HttpPort           string `envconfig:"HTTP_PORT"`
	TcpPort            string `envconfig:"TCP_PORT"`
	NykCsServiceName   string `envconfig:"NYKCS_SERVICE_NAME"`
	NykCsBaseUrl       string `envconfig:"NYKCS_BASE_URL"`
	NykCsTimeoutMillis int    `envconfig:"NYKCS_TIMEOUT_MILLIS"`
	BknCsServiceName   string `envconfig:"BKNCS_SERVICE_NAME"`
	BknCsBaseUrl       string `envconfig:"BKNCS_BASE_URL"`
	BknCsTimeoutMillis int    `envconfig:"BKNCS_TIMEOUT_MILLIS"`
}

type CreditScoreServiceConfig struct {
	ServiceName   string
	BaseUrl       string
	TimeoutMillis int
}

func LoadConfig() Config {
	godotenv.Load(".env")
	cfg := Config{}
	if err := envconfig.Process("CSG", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func (c *Config) GetCreditScoreServiceConfig(service string) CreditScoreServiceConfig {
	switch service {
	case "NYKCS":
		return CreditScoreServiceConfig{
			ServiceName:   c.NykCsServiceName,
			BaseUrl:       c.NykCsBaseUrl,
			TimeoutMillis: c.NykCsTimeoutMillis,
		}
	case "BKNCS":
		return CreditScoreServiceConfig{
			ServiceName:   c.BknCsServiceName,
			BaseUrl:       c.BknCsBaseUrl,
			TimeoutMillis: c.BknCsTimeoutMillis,
		}
	default:
		return CreditScoreServiceConfig{}
	}
}
