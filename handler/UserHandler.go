package handler

import (
	"context"
	"log"

	"user_service/service"

	autentificationServicePb "github.com/XML-organization/common/proto/autentification_service"
	bookingServicePb "github.com/XML-organization/common/proto/booking_service"
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
	BookingServiceClient         *bookingServicePb.BookingServiceClient
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
		log.Println(err)
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

func (ratingHandler *UserHandler) IsExceptional(ctx context.Context, in *pb.IsExceptionalRequest) (*pb.IsExceptionalResponse, error) {
	ratings, err := ratingHandler.RatingService.GetHostRatings(in.UserId)
	println("HostID: ", in.UserId)
	log.Println("HostID: ", in.UserId)

	//Provjera prosjecne ocjene host-a
	var ratingSum float64 = 0
	var i float64 = 0
	var threshold float64 = 4.7
	for _, rating := range ratings {
		println("Trenutna rating suma: ", ratingSum)
		ratingSum = ratingSum + float64(rating.Rating)
		i = i + 1
	}
	var rating float64 = ratingSum / i
	println("Rating suma: ", rating)
	println("Rating za ovog hosta je: ")
	if rating <= threshold {
		response := pb.IsExceptionalResponse{
			IsExceptional: true,
		}
		println("Rating je nizi od 4.7!")
		return &response, err
	}

	println("Rating je visi od 4.7!")
	//Provjeravam broj i duzinu rezervacija
	conn, err := grpc.Dial("booking-service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	bookingService := bookingServicePb.NewBookingServiceClient(conn)
	isExceptional, err := bookingService.IsExceptional(context.TODO(), &bookingServicePb.IsExceptionalRequest{UserId: in.UserId})
	if err != nil {
		log.Println(err)
		println(err.Error())
		return nil, err
	}

	if isExceptional.IsExceptional == false {
		response := pb.IsExceptionalResponse{
			IsExceptional: false,
		}
		return &response, err
	} else {
		response := pb.IsExceptionalResponse{
			IsExceptional: true,
		}
		return &response, err
	}
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
		log.Println(err)
	}

	GuestId, err := uuid.Parse(in.GuestId)
	if err != nil {
		log.Println(err)
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
		log.Println(err1)
		log.Println(err2)
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
		log.Println(err)
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
		log.Println(err)
		return &pb.UpdateRatingResponse{
			Message: "An error occured, please try again!",
		}, err
	}

	return &pb.UpdateRatingResponse{
		Message: message.Message,
	}, nil
}

func (ratingHandler *UserHandler) GetHostRatings(ctx context.Context, in *pb.GetHostRatingsRequest) (*pb.GetHostRatingsResponse, error) {
	ratings, err := ratingHandler.RatingService.GetHostRatings(in.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Konvertiranje rezultata ocjena u odgovarajući format protobuf-a
	var pbRatings []*pb.Rating
	for _, rating := range ratings {
		pbRating := &pb.Rating{
			Id:           rating.Id.String(),
			Rating:       int32(rating.Rating),
			Date:         rating.Date.Format("2006-01-02"),
			RaterId:      rating.RaterID.String(),
			RaterName:    rating.RaterName,
			RaterSurname: rating.RaterSurname,
			UserId:       rating.UserId.String(),
		}
		pbRatings = append(pbRatings, pbRating)
	}

	response := &pb.GetHostRatingsResponse{
		Ratings: pbRatings,
	}

	return response, nil
}
