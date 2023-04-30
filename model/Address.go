package model

import (
	"github.com/google/uuid"
)

type Address struct {
	ID      uuid.UUID `json:"id"`
	Country string    `json:"country" gorm:"not null;type:string"`
	City    string    `json:"city" gorm:"not null;type:string"`
	Street  string    `json:"street" gorm:"not null;type:string"`
	Number  string    `json:"number" gorm:"not null;type:string"`
}
