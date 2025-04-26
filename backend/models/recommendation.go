package models

import "time"

type Recommendation struct {
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`     // Include company name for better context
	Score      float64   `json:"score"`       // A score indicating the strength of the recommendation
	Reason     string    `json:"reason"`      // Brief explanation for the score
	LastUpdate time.Time `json:"last_update"` // Timestamp of the latest signal contributing to the score
}