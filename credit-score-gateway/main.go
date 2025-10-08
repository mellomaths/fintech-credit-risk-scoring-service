package main

import (
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/config"
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	go server.StartGrpcServer(cfg.GrpcPort)
	go server.StartTcpServer(cfg.TcpPort)
	server.StartHttpServer(cfg.HttpPort)
}
