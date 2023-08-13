package models

import "gorm.io/gorm"

type VoiceNotificationSubscriber struct {
	gorm.Model
	UserID string `gorm:"unique"`
}
