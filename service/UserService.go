package service

import (
	"context"

	"user_service/dto"
	"user_service/model"
	user "user_service/proto"
	"user_service/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	user.UnimplementedUserServiceServer
	UserRepo *repository.UserRepository
}

func (service *UserService) UpdateUser(user dto.ChangeUserDTO) model.RequestMessage {
	err := service.UserRepo.UpdateUser(user)
	if err != nil {
		return model.RequestMessage{Message: "An error occurred, please try again!"}
	}
	return model.RequestMessage{Message: "Success!"}
}

func (service *UserService) ChangePassword(userPasswords dto.UserPassword) model.RequestMessage {
	user, err := service.UserRepo.FindById(userPasswords.ID)

	if err != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message
	} else if err := bcrypt.CompareHashAndPassword(user.Password, []byte(userPasswords.OldPassword)); err != nil {
		message := model.RequestMessage{
			Message: "The old password is not correct!",
		}
		return message
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(userPasswords.NewPassword), 14)
	userPasswords.NewPassword = string(newPassword)

	return service.UserRepo.ChangePassword(userPasswords)
}

func (service *UserService) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error) {

	id, err := uuid.Parse(in.Id)

	if err != nil {
		panic(err)
	}

	addressId, err1 := uuid.Parse(in.Address.Id)

	if err1 != nil {
		panic(err1)
	}

	address := model.Address{
		ID:      addressId,
		Country: in.Address.Country,
		City:    in.Address.City,
		Street:  in.Address.Street,
		Number:  in.Address.Number,
	}

	user1 := model.User{
		ID:        id,
		Name:      in.Name,
		Surname:   in.Surname,
		Email:     in.Email,
		Password:  []byte(in.Password),
		AddressID: addressId,
		Address:   address,
	}

	response := &user.CreateUserResponse{
		Message: service.UserRepo.CreateUser(user1).Message,
	}

	return response, nil
}

/*func (service1 *UserService) NewGrpcServer() (*user.CreateUserResponse, error) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	user.RegisterUserServiceServer(grpcServer, &UserService{})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}*/
