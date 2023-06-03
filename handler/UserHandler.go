package handler

import (
	"context"
	"log"

	"user_service/service"

	autentificationServicePb "github.com/XML-organization/common/proto/autentification_service"
	pb "github.com/XML-organization/common/proto/user_service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	UserService                  *service.UserService
	RatingService                *service.RatingService
	AutentificationServiceClient *autentificationServicePb.AutentificationServiceClient
}

func NewUserHandler(service *service.UserService, ratingService *service.RatingService) *UserHandler {
	return &UserHandler{
		UserService:   service,
		RatingService: ratingService,
	}
}

func (userHandler *UserHandler) GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	user, err := userHandler.UserService.UserRepo.FindByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	println(user.ID.String())

	return mapUserToGetUserByEmailResponse(&user), nil
}

func (userHandler *UserHandler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
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

func (userHandler *UserHandler) CreateRating(ctx context.Context, in *pb.CreateRatingRequest) (*pb.CreateRatingResponse, error) {

	rating := mapRating(in)

	message, err := userHandler.RatingService.CreateRating(rating)

	response := pb.CreateRatingResponse{
		Message: message.Message,
	}

	return &response, err
}

// WasGuestRatedHost
func (userHandler *UserHandler) WasGuestRatedHost(ctx context.Context, in *pb.WasGuestRatedHostRequest) (*pb.WasGuestRatedHostResponse, error) {

	HostId, err := uuid.Parse(in.HostId)
	if err != nil {
		panic(err)
	}

	GuestId, err := uuid.Parse(in.GuestId)
	if err != nil {
		panic(err)
	}

	wasRated, err := userHandler.RatingService.WasGuestRatedHost(HostId, GuestId)

	response := pb.WasGuestRatedHostResponse{
		WasRated: wasRated,
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

func (userHandler *UserHandler) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteResponseMessage, error) {
	conn, err := grpc.Dial("autentification_service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	aService := autentificationServicePb.NewAutentificationServiceClient(conn)

	_, err1 := userHandler.UserService.DeleteUser(in.Id)
	_, err2 := aService.DeleteUser(context.TODO(), &autentificationServicePb.DeleteUserRequest{Email: in.Id})
	if err1 != nil {
		println(err1.Error())
		println(err2.Error())
	}
	response := pb.DeleteResponseMessage{
		Message: "ok",
	}

	return &response, err
}

func (ratingHandler *UserHandler) DeleteRating(ctx context.Context, in *pb.DeleteRatingRequest) (*pb.DeleteRatingResponse, error) {
	message, err := ratingHandler.RatingService.DeleteRating(in.HostId, in.GuestId)

	if err != nil {
		return &pb.DeleteRatingResponse{
			Message: "An error occured, please try again!",
		}, err
	}

	return &pb.DeleteRatingResponse{
		Message: message.Message,
	}, nil
}

func (ratingHandler *UserHandler) UpdateRating(ctx context.Context, in *pb.UpdateRatingRequest) (*pb.UpdateRatingResponse, error) {
	message, err := ratingHandler.RatingService.UpdateRating(in.HostId, in.GuestId, int(in.Rating))

	if err != nil {
		return &pb.UpdateRatingResponse{
			Message: "An error occured, please try again!",
		}, err
	}

	return &pb.UpdateRatingResponse{
		Message: message.Message,
	}, nil
}
