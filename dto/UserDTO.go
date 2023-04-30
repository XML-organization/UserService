package dto

import (
	"user_service/model"

	"github.com/google/uuid"
)

type UserPassword struct {
	ID          uuid.UUID `json:"id"`
	NewPassword string    `json:"new_password"`
	OldPassword string    `json:"old_password"`
}

type ChangeUserDTO struct {
	ID        uuid.UUID     `json:"id"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	Role      model.Role    `json:"role"`
	Name      string        `json:"name"`
	Surname   string        `json:"surname"`
	AddressID uuid.UUID     `json:"address_id"`
	Address   model.Address `json:"address" gorm:"foreignKey:AddressID"`
}
