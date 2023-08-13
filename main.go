package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/humbertovnavarro/spectator/pkg/appcontext"
	"github.com/humbertovnavarro/spectator/pkg/commands"
	voicenotificationsubscriber "github.com/humbertovnavarro/spectator/pkg/voice-notification-subscriber"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := &appcontext.Context{
		DB: appcontext.NewDB(),
	}

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		panic(err)
	}

	discord.AddHandler(func(discord *discordgo.Session, conn *discordgo.VoiceStateUpdate) {

		var channel *discordgo.Channel

		var channelID string

		if conn.BeforeUpdate != nil && conn.BeforeUpdate.ChannelID != "" {
			channelID = conn.BeforeUpdate.ChannelID
		} else {
			channelID = conn.ChannelID
		}

		channel, err := discord.Channel(channelID)
		if err != nil {
			logrus.Error("Error feching channel: ", channel)
		}

		subscribers, err := voicenotificationsubscriber.FetchAll(ctx.DB)

		if err != nil {
			logrus.Error("Error fetching subscribers: ", err)
			return
		}

		for _, subscriber := range subscribers {
			logrus.Info("Sending dm to ", subscriber)
			user, err := discord.User(subscriber.UserID)
			if err != nil {
				logrus.Error("Error getting user: ", err)
				return
			}
			dmChannel, err := discord.UserChannelCreate(user.ID)
			if err != nil {
				logrus.Error("Error creating dm channel: ", err)
				return
			}
			event := getVoiceNotificationEvent(conn)
			_, err = discord.ChannelMessageSend(dmChannel.ID, createDMVoiceNotification(discord, user, channel, event))
			if err != nil {
				logrus.Error("Error sending message to user: ", err)
				return
			}
		}
	})

	prefix := "!"
	cmds := commands.NewTextCommandManager()

	cmds.Register(commands.CommandVoiceNotificationRegisteration)
	cmds.Register(commands.CommandVoiceNotificationUnregisteration)

	discord.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		if !strings.HasPrefix(message.Content, prefix) {
			return
		}

		args := strings.Split(message.Content, " ")
		if len(args) == 0 {
			return
		}
		args[0] = strings.TrimPrefix(args[0], prefix)

		cmd, ok := cmds.Get(args[0])
		if !ok {
			return
		}
		cmd.Handler(ctx, args[1:], discord, message)
	})

	err = discord.Open()

	if err != nil {
		panic(err)
	}

	logrus.Info("Bot is running!")

	<-make(chan struct{})
}

func GetEnv(s string) {
	panic("unimplemented")
}

func createChannelLink(channel *discordgo.Channel) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s", channel.GuildID, channel.ID)
}

func createDMVoiceNotification(discord *discordgo.Session, user *discordgo.User, channel *discordgo.Channel, event VoiceNotificationEvent) string {
	guild, err := discord.Guild(channel.GuildID)
	if err != nil {
		logrus.Error("Error fetching guild: ", err)
		return ""
	}
	out := fmt.Sprintf("%s in %s\n", channel.Name, guild.Name)
	switch event {
	case JoinVoiceChannel:
		return out + fmt.Sprintf("User %s has joined voice channel %s", user.Mention(), createChannelLink(channel))
	case LeftVoiceChannel:
		return out + fmt.Sprintf("User %s has left voice channel %s", user.Mention(), createChannelLink(channel))
	case Muted:
		return out + fmt.Sprintf("User %s has been muted in voice channel %s", user.Mention(), createChannelLink(channel))
	case Unmuted:
		return out + fmt.Sprintf("User %s has been unmuted in voice channel %s", user.Mention(), createChannelLink(channel))
	case Deafened:
		return out + fmt.Sprintf("User %s has been deafened in voice channel %s", user.Mention(), createChannelLink(channel))
	case Undeafened:
		return out + fmt.Sprintf("User %s has been undeafened in voice channel %s", user.Mention(), createChannelLink(channel))
	default:
		return ""
	}
}

type VoiceNotificationEvent string

const (
	LeftVoiceChannel VoiceNotificationEvent = "left"
	JoinVoiceChannel VoiceNotificationEvent = "join"
	Muted            VoiceNotificationEvent = "muted"
	Unmuted          VoiceNotificationEvent = "unmuted"
	Deafened         VoiceNotificationEvent = "deafened"
	Undeafened       VoiceNotificationEvent = "undeafened"
)

func getVoiceNotificationEvent(state *discordgo.VoiceStateUpdate) VoiceNotificationEvent {
	if state.BeforeUpdate == nil {
		return JoinVoiceChannel
	}

	if state.ChannelID == "" {
		return LeftVoiceChannel
	}

	if state.BeforeUpdate.Mute != state.Mute {
		if state.Mute {
			return Muted
		}
		return Unmuted
	}

	if state.BeforeUpdate.Deaf != state.Deaf {
		if state.Deaf {
			return Deafened
		}
		return Undeafened
	}

	return JoinVoiceChannel
}
