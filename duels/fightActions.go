package duels

import (
	"fmt"
	"math/rand"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/utils"

	"github.com/bwmarrin/discordgo"
)

func FightAttack(s *discordgo.Session, channelID string) {
	currentDuel, err := utils.GetActiveDuel()
	if err != nil {
		fmt.Println(err.Error())
	}
	currentDuelPlayers, err := utils.GetDuelPlayers(currentDuel.Id)
	if err != nil {
		fmt.Println(err.Error())
	}

	challengerChar, err := utils.GetCharacterById(int(currentDuelPlayers.ChallengerChar.Int32))
	if err != nil {
		fmt.Println(err.Error())
	}
	challengedChar, err := utils.GetCharacterById(int(currentDuelPlayers.ChallengedChar.Int32))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(challengerChar)	
	fmt.Println(challengerChar.Precision)	
	fmt.Println(challengedChar)	
	fmt.Println(challengedChar.Precision)	

	if(currentDuel.SelectingPlayer == currentDuelPlayers.Challenger){
		defender := currentDuelPlayers.Challenged
		attackerChar := challengerChar
		defenderChar := challengedChar
		attackerName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengerChar.Int32))
		defenderName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengedChar.Int32))
		// Check if the attacker can hit is attack
		attackerPrecision := (rand.Intn(99) + 1)
		s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Rolled ", attackerPrecision, " for it's precision"))
		s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Current precision is", attackerChar.Precision))
		// Apply damage if a hit is scored
		if(attackerPrecision <= challengedChar.Precision){
			damageDealt := ((rand.Intn(9) + 1) + (attackerChar.Strength / 10)) - (defenderChar.Endurance/10)
			remainingHitpoints := (defenderChar.Hitpoints - damageDealt)
			defenderChar.Hitpoints -= damageDealt
			s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Manage to hit ", defenderName, " and deals ", damageDealt, " damage."))
			s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " Has ", remainingHitpoints, " hitpoints remaining"))
			// Check if the defenderChar is dead
			if(defenderChar.Hitpoints <= 0){
				s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, "Slays", defenderName, " and win the duel !"))
				_, err = database.DB.Exec(`UPDATE character SET is_alive=false WHERE character_id=$1;`, defenderChar)
				if err != nil {
					fmt.Println(err.Error())
				}
			// Initiate new turn and update the db if the blow wasn't fatal
			}else{
				_, err = database.DB.Exec(`UPDATE character SET hitpoints=$1 WHERE character_id=$2;`, remainingHitpoints, defenderChar)
				if err != nil {
					fmt.Println(err.Error())
				}
				_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, defender, currentDuel.Id)
				if err != nil {
					fmt.Println(err.Error())
				}

				s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " It is now your turn to play"))
				s.AddHandlerOnce(FightTurnOptions)
			}
		// Just begin a new turn if the attacker has missed his attack	
		}else{
			s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, "miss", defenderName))
			_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, defender, currentDuel.Id)
			if err != nil {
				fmt.Println(err.Error())
			}

			s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " It is now your turn to play"))
			s.AddHandlerOnce(FightTurnOptions)
		}
	}else{
		defender := currentDuelPlayers.Challenger
		attackerChar := challengedChar
		defenderChar := challengerChar
		attackerName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengedChar.Int32))
		defenderName, _ := utils.FindCharNameWithId(int(currentDuelPlayers.ChallengerChar.Int32))
		// Check if the attacker can hit is attack
		attackerPrecision := (rand.Intn(99) + 1)
		s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Rolled ", attackerPrecision, " for it's precision"))
		s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Current precision is", attackerChar.Precision))
		// Apply damage if a hit is scored
		if(attackerPrecision <= challengedChar.Precision){
			damageDealt := ((rand.Intn(9) + 1) + (attackerChar.Strength / 10)) - (defenderChar.Endurance/10)
			remainingHitpoints := (defenderChar.Hitpoints - damageDealt)
			defenderChar.Hitpoints -= damageDealt
			s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, " Manage to hit ", defenderName, " and deals ", damageDealt, " damage."))
			s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " Has ", remainingHitpoints, " hitpoints remaining"))
			// Check if the defenderChar is dead
			if(defenderChar.Hitpoints <= 0){
				s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, "Slays", defenderName, " and win the duel !"))
				_, err = database.DB.Exec(`UPDATE character SET is_alive=false WHERE character_id=$1;`, defenderChar)
				if err != nil {
					fmt.Println(err.Error())
				}
			// Initiate new turn and update the db if the blow wasn't fatal
			}else{
				_, err = database.DB.Exec(`UPDATE character SET hitpoints=$1 WHERE character_id=$2;`, remainingHitpoints, defenderChar)
				if err != nil {
					fmt.Println(err.Error())
				}
				_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, defender, currentDuel.Id)
				if err != nil {
					fmt.Println(err.Error())
				}

				s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " It is now your turn to play"))
				s.AddHandlerOnce(FightTurnOptions)
			}
		// Just begin a new turn if the attacker has missed his attack	
		}else{
			s.ChannelMessageSend(channelID, fmt.Sprintln(attackerName, "miss", defenderName))
			_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, defender, currentDuel.Id)
			if err != nil {
				fmt.Println(err.Error())
			}

			s.ChannelMessageSend(channelID, fmt.Sprintln(defenderName, " It is now your turn to play"))
			s.AddHandlerOnce(FightTurnOptions)
		}
	}

	

}
