package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password" gorm:"not null;type:string;default:null"`
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	AddressID uuid.UUID `json:"address_id"`
	Address   Address   `json:"address" gorm:"foreignKey:AddressID"`
}

type Role int

const (
	Host Role = iota
	Guest
	NK
)
