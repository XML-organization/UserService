package handler

import (
	"user_service/model"

	pb "github.com/XML-organization/common/proto/user_service"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func mapUserFromCreateUserRequest(user *pb.CreateUserRequest) model.User {

	println(user.Email)

	userId, err := uuid.Parse(user.Id)
	if err != nil {
		panic(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	return model.User{
		ID:       userId,
		Email:    user.Email,
		Password: password,
		Role:     model.Role(user.Role.Number()),
		Name:     user.Name,
		Surname:  user.Surname,
		Country:  user.Country,
		City:     user.City,
		Street:   user.Street,
		Number:   user.Number,
	}
}

func mapPassword(passwordRequest *pb.ChangePasswordRequest) model.UserPassword {

	return model.UserPassword{
		Email:       passwordRequest.Email,
		NewPassword: passwordRequest.NewPassword,
		OldPassword: passwordRequest.OldPassword,
	}
}

func mapAddress(a *pb.Address) model.Address {
	return model.Address{
		Country: a.Country,
		City:    a.City,
		Street:  a.Street,
		Number:  a.Number,
	}
}

func mapUserFromUpdateUserRequest(user *pb.UpdateUserRequest) model.ChangeUserDTO {
	userId, err := uuid.Parse(user.Id)
	if err != nil {
		panic(err)
	}

	return model.ChangeUserDTO{
		ID:       userId,
		Email:    user.Email,
		Password: user.Password,
		Role:     model.Role(user.Role.Number()),
		Name:     user.Name,
		Surname:  user.Surname,
		Country:  user.Country,
		City:     user.City,
		Street:   user.Street,
		Number:   user.Number,
	}
}
