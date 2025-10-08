package models

type CreditScoreRequest struct {
	ApplicantID   string  `json:"applicant_id"`
	Income        float64 `json:"income"`
	LoanAmount    float64 `json:"loan_amount"`
	CreditHistory float64 `json:"credit_history"`
}

type CreditScoreResponse struct {
	Score    float64 `json:"score"`
	Decision string  `json:"decision"`
	Error    *Error  `json:"error,omitempty"`
}

type Error struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
