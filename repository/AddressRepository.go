package repository

import (
	"user_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *AddressRepository) FindById(id uuid.UUID) (model.Address, error) {
	address := model.Address{}

	dbResult := repo.DatabaseConnection.First(&address, "id = ?", id)

	if dbResult != nil {
		return address, dbResult.Error
	}

	return address, nil
}
