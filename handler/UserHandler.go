package handler

import (
	"encoding/json"
	"net/http"

	"user_service/dto"

	"user_service/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func (userHandler *UserHandler) UpdateUser(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var user dto.ChangeUserDTO
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}
	message := userHandler.UserService.UpdateUser(user)
	json.NewEncoder(writer).Encode(message)
}

func (userHandler *UserHandler) ChangePassword(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var userPasswords dto.UserPassword
	err := decoder.Decode(&userPasswords)
	if err != nil {
		panic(err)
	}
	message := userHandler.UserService.ChangePassword(userPasswords)
	json.NewEncoder(writer).Encode(message)
}
