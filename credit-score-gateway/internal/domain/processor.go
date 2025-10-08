package domain

import (
	"log"

	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/models"
)

func ProcessCreditScoreRequest(request models.CreditScoreRequest) models.CreditScoreResponse {
	// wg := sync.WaitGroup{}
	nykcsScore, err := GetCreditScore(request, "NYKCS")
	if err != nil {
		nykcsScore = models.CreditScoreResponse{
			Score:    0,
			Decision: "UNDETERMINED",
			Error: &models.Error{
				ErrorCode:    "INTERNAL_ERROR",
				ErrorMessage: err.Error(),
			},
		}
	}

	bkncsScore, err := GetCreditScore(request, "BKNCS")
	if err != nil {
		bkncsScore = models.CreditScoreResponse{
			Score:    0,
			Decision: "UNDETERMINED",
			Error: &models.Error{
				ErrorCode:    "INTERNAL_ERROR",
				ErrorMessage: err.Error(),
			},
		}
	}
	log.Printf("NYKCS Score: %+v", nykcsScore)
	log.Printf("BKNCS Score: %+v", bkncsScore)
	championServiceName := GetChampionCreditScore(request)
	log.Printf("Champion Service Name: %s", championServiceName)
	switch championServiceName {
	case "NYKCS":
		return nykcsScore
	case "BKNCS":
		return bkncsScore
	default:
		return models.CreditScoreResponse{
			Score:    0,
			Decision: "UNDETERMINED",
			Error: &models.Error{
				ErrorCode:    "INTERNAL_ERROR",
				ErrorMessage: "Unknown champion service name",
			},
		}
	}
}
