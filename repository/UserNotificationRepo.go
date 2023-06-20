package repository

import (
	"log"
	"user_service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserNotificationRepository struct {
	DatabaseConnection *gorm.DB
}

func NewUserNotificationRepository(db *gorm.DB) *UserNotificationRepository {
	err := db.AutoMigrate(&model.UsersNotification{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return &UserNotificationRepository{
		DatabaseConnection: db,
	}
}

func (repo *UserNotificationRepository) GetForUserID(userID uuid.UUID) (model.UsersNotification, error) {
	notificationsSettings := model.UsersNotification{}
	result := repo.DatabaseConnection.Where("id_user = ?", userID).Find(&notificationsSettings)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return notificationsSettings, nil
}

func (repo *UserNotificationRepository) Save(notificationSettings model.UsersNotification) error {
	result := repo.DatabaseConnection.Save(&notificationSettings)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func (repo *UserNotificationRepository) FindByID(id uuid.UUID) (model.UsersNotification, error) {
	notificationsSettings := model.UsersNotification{}
	err := repo.DatabaseConnection.Where("id_user = ?", id).First(&notificationsSettings)
	if err != nil {
		println(err)
	}
	return notificationsSettings, err.Error
}

func (repo *UserNotificationRepository) UpdateNotificationSettings(userNot model.UsersNotification) (model.RequestMessage, error) {
	settings := model.UsersNotification{}
	err := repo.DatabaseConnection.Where("id_user = ?", userNot.UserID).First(&settings)
	if err != nil {
		log.Println(err)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}, err.Error
	}
	println("---------------------ovo je podesavanje za tog korisnika prije azuriranja:")
	println("id:", settings.ID.String())
	println("userid:", settings.UserID.String())
	println("RequestCreated:", settings.RequestCreated)
	println("ReservationCanceled:", settings.RequestCreated)
	println("HostGraded:", settings.RequestCreated)
	println("AccommodationGraded:", settings.RequestCreated)
	println("StatusChange:", settings.RequestCreated)
	println("ReservationReply:", settings.RequestCreated)

	// Ažuriranje
	settings.AccommodationGraded = userNot.AccommodationGraded
	settings.RequestCreated = userNot.RequestCreated
	settings.ReservationCanceled = userNot.ReservationCanceled
	settings.ReservationReply = userNot.ReservationReply
	settings.StatusChange = userNot.StatusChange
	settings.HostGraded = userNot.HostGraded

	println("---------------------ovo je podesavanje za tog korisnika prije azuriranja:")
	println("id:", settings.ID.String())
	println("userid:", settings.UserID.String())
	println("RequestCreated:", settings.RequestCreated)
	println("ReservationCanceled:", settings.RequestCreated)
	println("HostGraded:", settings.RequestCreated)
	println("AccommodationGraded:", settings.RequestCreated)
	println("StatusChange:", settings.RequestCreated)
	println("ReservationReply:", settings.RequestCreated)

	dbResult := repo.DatabaseConnection.Model(&settings).Updates(settings)
	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}, dbResult.Error
	}

	return model.RequestMessage{
		Message: "Notification settings updated successfully!",
	}, nil
}

func (repo *UserNotificationRepository) DeleteAll() error {
	result := repo.DatabaseConnection.Exec("DELETE FROM users_notifications")
	if result.Error != nil {
		return result.Error
	}
	return nil
}
