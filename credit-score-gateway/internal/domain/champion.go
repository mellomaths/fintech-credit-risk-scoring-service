package domain

import (
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/config"
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/models"
)

func GetChampionCreditScore(request models.CreditScoreRequest) string {
	cfg := config.LoadConfig()
	return cfg.NykCsServiceName
}
