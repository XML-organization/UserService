package handler

import (
	"user_service/model"

	events "github.com/XML-organization/common/saga/create_user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func mapSagaUserToUser(u *events.User) *model.User {

	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	println("nalazim se u mapSagaUserToUser ispod metode za generisanje passworda")

	id, _ := uuid.Parse(u.ID)

	println("ovo je string reprezentacija uuid iz usera: " + u.ID)
	println("ovo je nakon parsiranja pa ponovo prevod u string: " + id.String())

	return &model.User{
		ID:       id,
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Password: password,
		Role:     model.Role(u.Role),
		Country:  u.Country,
		City:     u.City,
		Street:   u.Street,
		Number:   u.Number,
	}
}
