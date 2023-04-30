package repository

import (
	"user_service/dto"
	"user_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
	AddressRepository  *AddressRepository
}

func (repo *UserRepository) UpdateUser(user dto.ChangeUserDTO) error {
	//dbResult := repo.DatabaseConnection.Update(user)

	loggedUser, err := repo.FindById(user.ID)
	var userOldAddress model.Address
	err1 := repo.DatabaseConnection.Where("id = ?", loggedUser.AddressID).First(&userOldAddress)

	if err != nil {
		return err
	}

	if err1.Error != nil {
		return err1.Error
	}

	sqlStatementUser := `
		UPDATE users
		SET name = $2, surname = $3, address_id = $4, email = $5
		WHERE id = $1;`

	sqlStatementUser1 := `
		UPDATE users
		SET name = $2, surname = $3, email = $4
		WHERE id = $1;`

	id := uuid.New()
	address := model.Address{
		ID:      id,
		Country: user.Address.Country,
		City:    user.Address.City,
		Street:  user.Address.Street,
		Number:  user.Address.Number,
	}

	if userOldAddress.Country == address.Country && userOldAddress.City == address.City && userOldAddress.Street == address.Street && userOldAddress.Number == address.Number {
		dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser1, user.ID, user.Name, user.Surname, user.Email)
		if dbResult1.Error != nil {
			return dbResult1.Error
		}
	} else {
		dbResult2 := repo.DatabaseConnection.Save(address)
		dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, user.ID, user.Name, user.Surname, id, user.Email)

		if dbResult1.Error != nil {
			return dbResult1.Error
		}
		if dbResult2.Error != nil {
			return dbResult2.Error
		}
	}
	return nil
}

func (repo *UserRepository) ChangePassword(userPasswords dto.UserPassword) model.RequestMessage {

	sqlStatementUser := `
		UPDATE users
		SET password = $2
		WHERE id = $1;`

	dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, userPasswords.ID, userPasswords.NewPassword)

	if dbResult1.Error != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message
	}

	message := model.RequestMessage{
		Message: "Success!",
	}
	return message
}

func (repo *UserRepository) FindById(id uuid.UUID) (model.User, error) {
	user := model.User{}

	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user model.User) model.RequestMessage {
	dbResult := repo.DatabaseConnection.Save(user)

	if dbResult.Error != nil {
		return model.RequestMessage{
			Message: "An error occured, please try again!",
		}
	}

	return model.RequestMessage{
		Message: "Success!",
	}
}
