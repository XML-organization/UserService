package handler

import (
	"user_service/model"

	pb "github.com/XML-organization/common/proto/user_service"

	"github.com/google/uuid"
)

func mapUserFromCreateUserRequest(user *pb.CreateUserRequest) model.User {

	println(user.Email)

	userId, err := uuid.Parse(user.Id)
	if err != nil {
		panic(err)
	}

	addressId, err := uuid.Parse(user.Address.Id)
	if err != nil {
		panic(err)
	}

	return model.User{
		ID:        userId,
		Email:     user.Email,
		Password:  user.Password,
		Role:      model.Role(user.Role.Number()),
		Name:      user.Name,
		Surname:   user.Surname,
		AddressID: addressId,
		Address:   mapAddress(user.Address),
	}
}

func mapPassword(passwordRequest *pb.ChangePasswordRequest) model.UserPassword {
	userId, err := uuid.Parse(passwordRequest.Id)
	if err != nil {
		panic(err)
	}

	return model.UserPassword{
		ID:          userId,
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

	addressId, err := uuid.Parse(user.Address.Id)
	if err != nil {
		panic(err)
	}

	return model.ChangeUserDTO{
		ID:        userId,
		Email:     user.Email,
		Password:  user.Password,
		Role:      model.Role(user.Role.Number()),
		Name:      user.Name,
		Surname:   user.Surname,
		AddressID: addressId,
		Address:   mapAddress(user.Address),
	}
}
