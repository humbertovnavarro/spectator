package voicenotificationsubscriber

import (
	"github.com/humbertovnavarro/spectator/pkg/models"
	"gorm.io/gorm"
)

func FetchAll(db *gorm.DB) ([]*models.VoiceNotificationSubscriber, error) {
	var subscribers []*models.VoiceNotificationSubscriber
	err := db.Find(&subscribers).Error
	return subscribers, err
}
