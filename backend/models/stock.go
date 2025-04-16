// model/stock.go
// File: backend/models/stock.go

package models

import "time"

// Stock representa un registro de stock en la base de datos
type Stock struct {
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

/*
   "ticker": "RMTI",
   "target_from": "$7.00",
   "target_to": "$3.00",
   "company": "Rockwell Medical",
   "action": "target lowered by",
   "brokerage": "HC Wainwright",
   "rating_from": "Buy",
   "rating_to": "Buy",
   "time": "2025-03-25T00:30:06.00066843Z"
*/
