package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreditScoreRequest struct {
	ApplicantID   string  `json:"applicant_id"`
	Income        float64 `json:"income"`
	LoanAmount    float64 `json:"loan_amount"`
	CreditHistory float64 `json:"credit_history"`
}

type CreditScoreResponse struct {
	Score    float64 `json:"score"`
	Decision string  `json:"decision"`
}

type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func calculateScore(request CreditScoreRequest) float64 {
	// TODO: Implement fake score calculation
	return request.Income / request.LoanAmount
}

func postScore(c *gin.Context) {
	var request CreditScoreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode:    "INVALID_REQUEST",
			ErrorMessage: "Invalid request body",
		})
		return
	}

	score := calculateScore(request)
	decision := "UNDETERMINED"
	if score < 0.5 {
		decision = "REJECTED"
	} else if score < 0.7 {
		decision = "REVIEW"
	} else {
		decision = "APPROVED"
	}
	response := CreditScoreResponse{
		Score:    score,
		Decision: decision,
	}
	c.JSON(http.StatusOK, response)
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	router := gin.Default()
	router.GET("/health", getHealth)
	router.POST("/score", postScore)
	router.Run(":8081")
}
