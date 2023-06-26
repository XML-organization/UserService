package model

import (
	"github.com/google/uuid"
)

type UsersNotification struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `gorm:"column:id_user" json:"id_user"`
	RequestCreated      bool      `json:"requestCreated"`
	ReservationCanceled bool      `json:"reservationCanceled"`
	HostGraded          bool      `json:"hostGraded"`
	AccommodationGraded bool      `json:"accommodationGraded"`
	StatusChange        bool      `json:"statusChange"`
	ReservationReply    bool      `json:"reservationReply"`
}
