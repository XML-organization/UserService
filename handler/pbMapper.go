package handler

import (
	"log"
	"time"
	"user_service/model"

	pb "github.com/XML-organization/common/proto/user_service"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func mapUserFromCreateUserRequest(user *pb.CreateUserRequest) model.User {

	println(user.Email)

	userId, err := uuid.Parse(user.Id)
	if err != nil {
		log.Println(err)
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

func mapRating(r *pb.CreateRatingRequest) model.Rating {
	id := uuid.New()

	reterId, err := uuid.Parse(r.RaterId)
	if err != nil {
		log.Println(err)
	}

	userId, err := uuid.Parse(r.UserId)
	if err != nil {
		log.Println(err)
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, r.Date)
	if err != nil {
		log.Println(err)
	}

	return model.Rating{
		Id:           id,
		Rating:       int(r.Rating),
		RaterID:      reterId,
		RaterName:    r.RaterName,
		RaterSurname: r.RaterSurname,
		UserId:       userId,
		Date:         date,
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
		log.Println(err)
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

func mapUserToGetUserByEmailResponse(user *model.User) *pb.GetUserByEmailResponse {

	id := " |" + user.ID.String() + " |"
	println(id)

	return &pb.GetUserByEmailResponse{
		Id:      id,
		Email:   user.Email,
		Role:    pb.Role(user.Role),
		Name:    user.Name,
		Surname: user.Surname,
		Country: user.Country,
		City:    user.City,
		Street:  user.Street,
		Number:  user.Number,
	}
}

func mapUserFromDeleteUserRequest(id *pb.DeleteUserRequest) uuid.UUID {

	userId, err := uuid.Parse(id.Id)
	if err != nil {
		log.Println(err)
	}

	return userId
}
