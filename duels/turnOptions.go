package duels

import (
	"fmt"

	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)

func FightOptionsInfo(s *discordgo.Session, channedId string) {
	currentDuel, err := utils.GetActiveDuel()
	if err != nil {
		fmt.Println(err.Error())
	}
	activePlayer, err := utils.GetCharacterById(currentDuel.SelectingPlayer)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = s.ChannelMessageSendEmbed(channedId, &discordgo.MessageEmbed{
		Title: fmt.Sprintln(activePlayer.Name, ", it is your turn to play"),
		Description: fmt.Sprintln("```--fight Attack will do a Weaponskill test in order to hit your opponent",
			"\n-fight Dodge will do a Weaponskill test divided by 2 in order to hit, but will allow you to do a Dodge test if you get hit this turn",
			"\n-fight Flee will do an Agility test to flee this battle, if failed, you opponent get a free attack on you"),
		Color: 0x0099ff,
		// Footer: &discordgo.MessageEmbedFooter{
		// 	Text: "generalDuelInvite:" + m.Author.ID,
		// },
	})

	if err != nil {
		s.ChannelMessageSend(channedId, utils.ErrorMessage("Bot error", "Error Showing Fight Options."))
	}

	s.AddHandlerOnce(FightTurnOptions)
}

func FightTurnOptions(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "-fight attack":
		s.AddHandlerOnce(FightAttack)
	case "-fight dodge":
	case "-fight flee":
	}
}

func MapTurnOptions(s *discordgo.Session, m *discordgo.MessageCreate) {
}
