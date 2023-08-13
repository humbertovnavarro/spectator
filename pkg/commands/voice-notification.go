package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/humbertovnavarro/spectator/pkg/appcontext"
	"github.com/humbertovnavarro/spectator/pkg/models"
	"github.com/sirupsen/logrus"
)

var CommandVoiceNotificationRegisteration = TextCommand{
	Name:        "registervoice",
	Aliases:     []string{"vnr", "register"},
	Description: "Register voice notification",
	Handler: func(ctx *appcontext.Context, args []string, session *discordgo.Session, message *discordgo.MessageCreate) {
		db := ctx.DB

		userID := message.Author.ID

		subscriber := models.VoiceNotificationSubscriber{
			UserID: userID,
		}
		err := db.Create(&subscriber).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		_, err = session.ChannelMessageSend(message.ChannelID, "You have registered to receive voice channel notifications for all guilds I'm a part of.")
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

var CommandVoiceNotificationUnregisteration = TextCommand{
	Name:        "unregistervoice",
	Aliases:     []string{"vnur", "unregister"},
	Description: "Unregister voice notification",
	Handler: func(ctx *appcontext.Context, args []string, session *discordgo.Session, message *discordgo.MessageCreate) {
		db := ctx.DB

		userID := message.Author.ID

		err := db.Where("user_id = ?", userID).Delete(&models.VoiceNotificationSubscriber{}).Error
		if err != nil {
			logrus.Error(err)
			return
		}
		_, err = session.ChannelMessageSend(message.ChannelID, "You have unregistered to receive voice channel notifications for all guilds I'm a part of.")
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}
