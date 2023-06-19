package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID               uuid.UUID          `json:"id"`
	Text             string             `json:"text"`
	NotificationTime time.Time          `json:"notificationTime" gorm:"not null"`
	UserID           uuid.UUID          `gorm:"column:id_user" json:"id_user"`
	Status           NotificationStatus `json:"status"`
}

/*func (not *Notification) BeforeCreate(scope *gorm.DB) error {
	not.ID = uuid.New()
	return nil
}*/
