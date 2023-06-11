package handler

import (
	"log"
	"strings"
	"user_service/model"

	events "github.com/XML-organization/common/saga/create_user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func mapSagaUserToUser(u *events.User) *model.User {

	println(u.Id)
	idString := strings.Split(u.Id, " |")[1]
	log.Println("OVDJE MAPIRAM STRING ID U UUID")
	log.Println("ID: " + idString)
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println(err)
		return nil
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

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
