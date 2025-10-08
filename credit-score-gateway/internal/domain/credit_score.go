package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/config"
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/models"
)

func GetCreditScore(request models.CreditScoreRequest, service string) (models.CreditScoreResponse, error) {
	cfg := config.LoadConfig()
	baseUrl := cfg.ChampionBaseUrl
	timeoutMillis := cfg.ChampionTimeoutMillis
	if service == "CHALLENGER" {
		baseUrl = cfg.ChallengerBaseUrl
		timeoutMillis = cfg.ChallengerTimeoutMillis
	}
	httpClient := &http.Client{
		Timeout: time.Duration(timeoutMillis) * time.Millisecond,
	}

	url := baseUrl + "/score"
	response, err := httpClient.Get(url)
	if err != nil {
		return models.CreditScoreResponse{}, fmt.Errorf("failed to get champion score: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return models.CreditScoreResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return models.CreditScoreResponse{}, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(body))
	}
	log.Printf("Response body: %s", string(body))

	var scoreResponse models.CreditScoreResponse
	err = json.Unmarshal(body, &scoreResponse)
	if err != nil {
		return models.CreditScoreResponse{}, fmt.Errorf("failed to unmarshal response body (raw: %s): %w", string(body), err)
	}
	return scoreResponse, nil
}
