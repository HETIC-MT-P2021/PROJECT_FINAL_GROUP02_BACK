package duels

import (
	"fmt"

	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)

func FightOptionsInfo(s *discordgo.Session, channedId string, activePlayer string, activeCharacter int32) {

	activeCharacterName, _ := utils.FindCharNameWithId(int(activeCharacter))

	_, err := s.ChannelMessageSendEmbed(channedId, &discordgo.MessageEmbed{
		Title: fmt.Sprintln(activeCharacterName, ", it is your turn to play"),
		Description: fmt.Sprintln("```-fight Attack will do a Precision test in order to hit your opponent",
			"\n-fight Dodge will do a Precision test divided by 2 in order to hit, but will allow you to do an Agility test if you get hit this turn",
			"\n-fight Flee will do an Agility test to flee this battle, if failed, you opponent get a free attack on you```"),
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

	// Get duel Informations
	currentDuel, err := utils.GetActiveDuel()
	if err != nil {
		fmt.Println(err.Error())
	}

	if(m.Author.ID != currentDuel.SelectingPlayer){
		s.ChannelMessageSend(m.ChannelID, "It is not your turn to play")
		s.AddHandlerOnce(FightTurnOptions)
	}else{
		switch m.Content {
			case "-fight Attack":
				FightAttack(s, m.ChannelID)
			case "-fight Dodge":
			case "-fight Flee":
		}
	}
}

func MapTurnOptions(s *discordgo.Session, m *discordgo.MessageCreate) {
}
