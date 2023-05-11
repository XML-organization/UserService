package handler

import (
	"context"

	"user_service/service"

	pb "github.com/XML-organization/common/proto/user_service"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	UserService *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
	}
}

func (userHandler *UserHandler) GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	println("/////////////////////////////")
	println("usao u metodu getuserbyemail")
	user, err := userHandler.UserService.UserRepo.FindByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	println(user.ID.String())

	return mapUserToGetUserByEmailResponse(&user), nil
}

func (userHandler *UserHandler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	println("////////")
	println(in.Email)
	println(in.Id)

	message, err := userHandler.UserService.UpdateUser(mapUserFromUpdateUserRequest(in))

	return &pb.UpdateUserResponse{
		Message: message.Message,
	}, err
}

func (userHandler *UserHandler) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	message, err := userHandler.UserService.ChangePassword(mapPassword(in))

	return &pb.ChangePasswordResponse{
		Message: message.Message,
	}, err
}

func (userHandler *UserHandler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	user := mapUserFromCreateUserRequest(in)

	message, err := userHandler.UserService.CreateUser(user)

	response := pb.CreateUserResponse{
		Message: message.Message,
	}

	return &response, err
}

func (userHandler *UserHandler) Print(ctx context.Context, in *pb.PrintRequest) (*pb.PrintResponse, error) {
	println("adasdasdasdas")

	println(in.Message)

	return &pb.PrintResponse{
		Message: in.Message,
	}, nil
}
