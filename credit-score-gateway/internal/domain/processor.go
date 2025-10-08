package domain

import (
	"log"

	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/models"
)

func ProcessCreditScoreRequest(request models.CreditScoreRequest) models.CreditScoreResponse {
	championScore, err := GetCreditScore(request, "CHAMPION")
	if err != nil {
		championScore = models.CreditScoreResponse{
			Score:    0,
			Decision: "UNDETERMINED",
			Error: &models.Error{
				ErrorCode:    "INTERNAL_ERROR",
				ErrorMessage: err.Error(),
			},
		}
	}

	challengerScore, err := GetCreditScore(request, "CHALLENGER")
	if err != nil {
		challengerScore = models.CreditScoreResponse{
			Score:    0,
			Decision: "UNDETERMINED",
			Error: &models.Error{
				ErrorCode:    "INTERNAL_ERROR",
				ErrorMessage: err.Error(),
			},
		}
	}
	log.Printf("Champion Score: %+v", championScore)
	log.Printf("Challenger Score: %+v", challengerScore)

	return championScore
}
