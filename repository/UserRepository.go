package repository

import (
	"user_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	err := db.AutoMigrate(&model.User{}, &model.Address{})
	if err != nil {
		return nil
	}

	return &UserRepository{
		DatabaseConnection: db,
	}
}

/*sqlStatementUser := `
UPDATE users
SET name = $2, surname = $3, address_id = $4, email = $5
WHERE id = $1;`*/

func (repo *UserRepository) UpdateUser(user model.ChangeUserDTO) error {
	//dbResult := repo.DatabaseConnection.Update(user)

	sqlStatementUser := `
		UPDATE users
		SET name = $2, surname = $3, email = $4, country = $5, city = $6, street = $7, number = $8
		WHERE id = $1;`

	dbResult := repo.DatabaseConnection.Exec(sqlStatementUser, user.ID, user.Name, user.Surname, user.Email, user.Country, user.City, user.Street, user.Number)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}

func (repo *UserRepository) ChangePassword(userPasswords model.UserPassword) (model.RequestMessage, error) {

	sqlStatementUser := `
		UPDATE users
		SET password = $2
		WHERE email = $1;`

	dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, userPasswords.Email, userPasswords.NewPassword)

	if dbResult1.Error != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message, dbResult1.Error
	}

	message := model.RequestMessage{
		Message: "Success!",
	}
	return message, nil
}

func (repo *UserRepository) ChangeEmail(emails model.UpdateEmailDTO) error {

	sqlStatementUser := `
		UPDATE users
		SET email = $2
		WHERE email = $1;`

	dbResult1 := repo.DatabaseConnection.Exec(sqlStatementUser, emails.OldEmail, emails.NewEmail)

	if dbResult1.Error != nil {
		return dbResult1.Error
	}
	return nil
}

func (repo *UserRepository) FindById(id uuid.UUID) (model.User, error) {
	user := model.User{}

	dbResult := repo.DatabaseConnection.First(&user, "id = ?", id)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (model.User, error) {
	user := model.User{}

	dbResult := repo.DatabaseConnection.First(&user, "email = ?", email)

	if dbResult != nil {
		return user, dbResult.Error
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user model.User) (model.RequestMessage, error) {

	dbResult := repo.DatabaseConnection.Save(&user)

	if dbResult.Error != nil {
		return model.RequestMessage{
			Message: "An error occured, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Success!",
	}, nil
}

func (repo *UserRepository) DeleteUserById(id uuid.UUID) (model.RequestMessage, error) {
	dbResult := repo.DatabaseConnection.Delete(&model.User{}, id)
	if dbResult.Error != nil {
		return model.RequestMessage{
			Message: "An error occured, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Success!",
	}, nil
}
