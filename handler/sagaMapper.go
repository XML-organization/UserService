package handler

import (
	"user_service/model"

	events "github.com/XML-organization/common/saga/create_user"
	"golang.org/x/crypto/bcrypt"
)

func mapSagaUserToUser(u *events.User) *model.User {

	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	println("nalazim se u mapSagaUserToUser ispod metode za generisanje passworda")

	return &model.User{
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
