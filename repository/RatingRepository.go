package repository

import (
	"errors"
	"log"
	"time"
	"user_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingRepository struct {
	DatabaseConnection *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	err := db.AutoMigrate(&model.Rating{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return &RatingRepository{
		DatabaseConnection: db,
	}
}

func (repo *RatingRepository) CreateRating(rating model.Rating) (model.RequestMessage, error) {

	dbResult := repo.DatabaseConnection.Save(&rating)

	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occured, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Success!",
	}, nil
}

func (repo *RatingRepository) WasGuestRatedHost(hostID uuid.UUID, guestID uuid.UUID) (bool, error) {
	var rating model.Rating
	err := repo.DatabaseConnection.Where("user_id = ? AND rater_id = ?", hostID, guestID).First(&rating).Error
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Ocena nije pronađena
		}
		return false, err // Greška pri izvršavanju upita
	}
	return true, nil // Ocena je pronađena
}

func (repo *RatingRepository) DeleteRating(hostId string, guestId string) (model.RequestMessage, error) {
	// Pretraživanje ocjene na temelju hostId-a i guestId-a
	rating := model.Rating{}
	err := repo.DatabaseConnection.Where("user_id = ? AND rater_id = ?", hostId, guestId).First(&rating).Error
	if err != nil {
		log.Println(err)
		return model.RequestMessage{
			Message: "Rating not found!",
		}, err
	}

	// Brisanje ocjene iz baze podataka
	dbResult := repo.DatabaseConnection.Delete(&rating)
	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Success!",
	}, nil
}

func (repo *RatingRepository) UpdateRating(hostId string, guestId string, newRating int) (model.RequestMessage, error) {
	// Pretraživanje ocjene na temelju hostId-a i guestId-a
	rating := model.Rating{}
	err := repo.DatabaseConnection.Where("user_id = ? AND rater_id = ?", hostId, guestId).First(&rating).Error
	if err != nil {
		log.Println(err)
		return model.RequestMessage{
			Message: "Rating not found!",
		}, err
	}

	// Ažuriranje ocjene i datuma
	rating.Rating = newRating
	rating.Date = time.Now() // Postavljanje trenutnog datuma

	dbResult := repo.DatabaseConnection.Save(&rating)
	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Rating updated successfully!",
	}, nil
}
func (repo *RatingRepository) GetHostRatings(hostId string) ([]model.Rating, error) {
	var ratings []model.Rating
	err := repo.DatabaseConnection.Where("user_id = ?", hostId).Find(&ratings).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ratings, nil
}
