package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/humbertovnavarro/spectator/pkg/appcontext"
)

type TextCommand struct {
	Name        string
	Aliases     []string
	Description string
	Handler     func(*appcontext.Context, []string, *discordgo.Session, *discordgo.MessageCreate)
}
