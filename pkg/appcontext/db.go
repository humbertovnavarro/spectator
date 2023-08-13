package appcontext

import (
	"github.com/humbertovnavarro/spectator/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("userdata.db"), &gorm.Config{})
	db.AutoMigrate(&models.VoiceNotificationSubscriber{})
	if err != nil {
		panic(err)
	}
	return db
}
