package model

import (
	"time"

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

type Rating struct {
	Id           uuid.UUID `json:"id"`
	Rating       int       `json:"rating"`
	Date         time.Time `json:"date"`
	RaterID      uuid.UUID `json:"rater_id"`
	RaterName    string    `json:"rater_name"`
	RaterSurname string    `json:"rater_surname"`
	UserId       uuid.UUID `json:"user_id"`
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

type UpdateEmailDTO struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}
