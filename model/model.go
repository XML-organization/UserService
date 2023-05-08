package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password []byte    `json:"password" gorm:"not null;type:string;default:null"`
	Role     Role      `json:"role"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Country  string    `json:"country" gorm:"not null;type:string"`
	City     string    `json:"city" gorm:"not null;type:string"`
	Street   string    `json:"street" gorm:"not null;type:string"`
	Number   string    `json:"number" gorm:"not null;type:string"`
}

type Role int

const (
	Host Role = iota
	Guest
	NK
)

type Address struct {
	ID      uuid.UUID `json:"id"`
	Country string    `json:"country" gorm:"not null;type:string"`
	City    string    `json:"city" gorm:"not null;type:string"`
	Street  string    `json:"street" gorm:"not null;type:string"`
	Number  string    `json:"number" gorm:"not null;type:string"`
}

type RequestMessage struct {
	Message string `json:"message"`
}

type UserPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type ChangeUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Country  string    `json:"country" gorm:"not null;type:string"`
	City     string    `json:"city" gorm:"not null;type:string"`
	Street   string    `json:"street" gorm:"not null;type:string"`
	Number   string    `json:"number" gorm:"not null;type:string"`
}
