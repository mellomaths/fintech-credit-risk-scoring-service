package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/domain"
	"github.com/mellomaths/fintech-credit-risk-scoring-service/credit-score-gateway/internal/models"
)

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func postScore(c *gin.Context) {
	var request models.CreditScoreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"score":    0,
			"decision": "UNDETERMINED",
			"error": gin.H{
				"error_code":    "INVALID_REQUEST",
				"error_message": "Invalid request body",
			},
		})
		return
	}

	response := domain.ProcessCreditScoreRequest(request)
	c.JSON(http.StatusOK, gin.H{
		"score":    response.Score,
		"decision": response.Decision,
		"error": gin.H{
			"error_code":    response.Error.ErrorCode,
			"error_message": response.Error.ErrorMessage,
		},
	})
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func StartHttpServer(port string) {
	router := gin.Default()
	router.Use(JSONMiddleware())
	router.GET("/health", getHealth)
	router.POST("/score", postScore)

	log.Println("HTTP server listening on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
