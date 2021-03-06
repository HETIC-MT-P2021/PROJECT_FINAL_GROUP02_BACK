package duels

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)

type error interface {
	Error() string
}

func DuelController(s *discordgo.Session, channelID string, involvedPlayers []string) {

	initialSetup := duelSetup(involvedPlayers[0], involvedPlayers[1])

	var err error
	s.ChannelMessageSend(channelID, "Rolling Initiative...")
	initialSetup.ActiveFighter, err = rollInitiative(initialSetup, s, channelID)
	if err != nil {
		s.ChannelMessageSend(channelID, "Error when determining which character would start")
		return
	}

	// Get duel Informations
	currentDuel, err := utils.GetActiveDuel()
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, initialSetup.ActiveFighter, currentDuel.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	currentDuelPlayers, err := utils.GetDuelPlayers(currentDuel.Id)
	if err != nil {
		fmt.Println(err.Error())
	}
		
	var playingChar int32

	// Get the corresponding character
	if(initialSetup.ActiveFighter == currentDuelPlayers.Challenger){
		playingChar = currentDuelPlayers.ChallengerChar.Int32
	}else{
		playingChar = currentDuelPlayers.ChallengedChar.Int32

	}


	FightOptionsInfo(s, channelID, initialSetup.ActiveFighter, playingChar)
}

// Load Duel Infos
func duelSetup(challenger string, challenged string) *game.DuelBattle {
	challengersArray := []string{challenger, challenged}
	initialSetup := game.DuelBattle{
		Challengers: challengersArray,
		IsOver:      false,
		Turn:        0,
	}

	return &initialSetup
}

// Do an initiative test (based on character Agility) to determine which character will play first
func rollInitiative(duelSetup *game.DuelBattle, s *discordgo.Session, channelID string) (string, error) {
	currentDuel, err := utils.GetActiveDuel()
	if err != nil {
		return "0", errors.New("Duel not found")
	}
	currentDuelPlayers, err := utils.GetDuelPlayers(currentDuel.Id)
	if err != nil {
		return "0", errors.New("Duel data not found")
	}

	challengerChar, err := utils.GetCharacterById(int(currentDuelPlayers.ChallengerChar.Int32))
	if err != nil {
		return "0", errors.New("Challenger character not found")
	}
	challengedChar, err := utils.GetCharacterById(int(currentDuelPlayers.ChallengedChar.Int32))
	if err != nil {
		return "0", errors.New("Challenged character not found")
	}

	challengerCharacterName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengerChar.Int32))
	challengedCharacterName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengedChar.Int32))

	challengerInitiative := challengerChar.Agility + (rand.Intn(9) + 1)
	challengedInitiative := challengedChar.Agility + (rand.Intn(9) + 1)
	s.ChannelMessageSend(channelID, fmt.Sprintln(challengerCharacterName, " Rolled ", challengerInitiative, " for it's initiative"))
	s.ChannelMessageSend(channelID, fmt.Sprintln(challengedCharacterName, " Rolled ", challengedInitiative, " for it's initiative"))

	if challengerInitiative > challengedInitiative {
		s.ChannelMessageSend(channelID, fmt.Sprintln(challengerCharacterName, " will play first"))
		return duelSetup.Challengers[0], nil
	} else if challengedInitiative > challengerInitiative {
		s.ChannelMessageSend(channelID, fmt.Sprintln(challengedCharacterName, " will play first"))
		return duelSetup.Challengers[1], nil
	} else {
		s.ChannelMessageSend(channelID, "Tie ! Choosing at random who will have the initiative...")
		if rand.Intn(10) < 5 {
			s.ChannelMessageSend(channelID, fmt.Sprintln(challengerCharacterName, " will play first"))
			return duelSetup.Challengers[0], nil
		} else {
			s.ChannelMessageSend(channelID, fmt.Sprintln(challengedCharacterName, " will play first"))
			return duelSetup.Challengers[1], nil
		}
	}
}
