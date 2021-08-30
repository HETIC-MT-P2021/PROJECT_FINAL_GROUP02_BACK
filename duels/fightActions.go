package duels

import (
// 	"fmt"

// 	"github.com/SteakBarbare/RPGBot/utils"
 	"github.com/bwmarrin/discordgo"
)

func FightAttack(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	currentDuel, err := utils.GetActiveDuel()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	currentDuelPlayers, err := utils.GetDuelPlayers(currentDuel.Id)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	challengerChar, err := utils.GetBattleCharacterById(currentDuelPlayers.ChallengerChar)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	challengedChar, err := utils.GetBattleCharacterById(currentDuelPlayers.ChallengedChar)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	if m.Author != currentDuel.SelectingPlayer {
// 		s.ChannelMessageSend(m.ChannelID, "It is not your turn to play")
// 		s.AddHandlerOnce(FightTurnOptions)
// 	}
}
