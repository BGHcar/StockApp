// model/stock.go
// File: backend/models/stock.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Stock representa un registro de stock en la base de datos
type Stock struct {
	gorm.Model           // Añade campos ID, CreatedAt, UpdatedAt, DeletedAt
	Ticker     string    `json:"ticker" gorm:"index;not null;type:text"`
	Company    string    `json:"company" gorm:"type:text"`
	TargetFrom string    `json:"target_from" gorm:"column:target_from;type:text"`
	TargetTo   string    `json:"target_to" gorm:"column:target_to;type:text"`
	Action     string    `json:"action" gorm:"index;type:text"`
	Brokerage  string    `json:"brokerage" gorm:"index;type:text"`
	RatingFrom string    `json:"rating_from" gorm:"column:rating_from;type:text"`
	RatingTo   string    `json:"rating_to" gorm:"column:rating_to;type:text"`
	Time       time.Time `json:"time" gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (Stock) TableName() string {
	return "stocks"
}
