package service

import (
	"user_service/model"
	"user_service/repository"

	"github.com/google/uuid"
)

type RatingService struct {
	RatingRepo *repository.RatingRepository
}

func NewRatingService(repo *repository.RatingRepository) *RatingService {
	return &RatingService{
		RatingRepo: repo,
	}
}

func (service *RatingService) CreateRating(rating model.Rating) (model.RequestMessage, error) {

	message, err := service.RatingRepo.CreateRating(rating)

	if err != nil {
		response := model.RequestMessage{
			Message: "An error occured, please try again!",
		}

		return response, err
	}

	response := model.RequestMessage{
		Message: message.Message,
	}

	return response, nil
}

func (service *RatingService) WasGuestRatedHost(hostID uuid.UUID, guestID uuid.UUID) (bool, error) {

	wasRated, err := service.RatingRepo.WasGuestRatedHost(hostID, guestID)

	if err != nil {
		return false, err
	}

	return wasRated, nil
}

func (service *RatingService) DeleteRating(hostID string, guestID string) (model.RequestMessage, error) {

	message, err := service.RatingRepo.DeleteRating(hostID, guestID)

	if err != nil {
		return message, err
	}

	return message, nil
}

func (service *RatingService) UpdateRating(hostID string, guestID string, rating int) (model.RequestMessage, error) {

	message, err := service.RatingRepo.UpdateRating(hostID, guestID, rating)

	if err != nil {
		return message, err
	}

	return message, nil
}
