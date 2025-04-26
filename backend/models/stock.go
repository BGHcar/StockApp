package models

import "time"

type Stock struct {
	Ticker string `gorm:"primaryKey"`
	TargetFrom float64
	TargetTo float64
	Company string 
	Action string 
	Brokerage string 
	RatingFrom string 
	RatingTo string 
	Time time.Time `gorm:"primaryKey"`
}