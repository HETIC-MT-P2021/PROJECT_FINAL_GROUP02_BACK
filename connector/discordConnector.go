package connector

import(
	"github.com/bwmarrin/discordgo"
)

// Discord Interface to mock Discordgo
type DiscordInterface interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
}