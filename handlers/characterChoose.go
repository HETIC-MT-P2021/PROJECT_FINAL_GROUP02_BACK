package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/duels"
	"github.com/SteakBarbare/RPGBot/game"
	"github.com/bwmarrin/discordgo"
)

// Main function, will run till either all characters are selected, an error has occured, or the selection is aborted
func selectCharacterBase(s *discordgo.Session, channelID string, involvedPlayers []string, lastPlayer int) {
	if lastPlayer == 0 {
		s.ChannelMessageSend(channelID, "Choose a character")
		s.AddHandlerOnce(selectCharacter)
	} else {
		s.ChannelMessageSend(channelID, "All Players are ready, starting duel !")
		duels.DuelController(s, channelID, involvedPlayers)
	}
}

func selectCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "-quit" {
		s.ChannelMessageSend(m.ChannelID, "Aborting character selection")
	} else if m.Content == "-char Show" {
		s.AddHandlerOnce(selectCharacter)
	} else {
		authorId, err := strconv.ParseInt(m.Author.ID, 10, 64)

		if err != nil {
			panic(err)
		}

		// Get the selecting player and duel ID
		duelPreparation := database.DB.QueryRow(fmt.Sprintln("SELECT * FROM duelpreparation WHERE isReady=", 0))

		currentDuel := game.DuelPreparation{}

		err = duelPreparation.Scan(&currentDuel.Id, &currentDuel.SelectingPlayer, &currentDuel.IsReady, &currentDuel.IsOver, &currentDuel.Turn)
		fmt.Println(currentDuel.Id)

		if err != nil {
			fmt.Println(err.Error())
			s.ChannelMessageSend(m.ChannelID, "No Duel found, aborting duel preparation")
			return
		} else {
			// Check if the message author is the selecting player
			if m.Author.ID == currentDuel.SelectingPlayer {

				// Get the character from db
				var selectedCharacter int
				charRow := database.DB.QueryRow("SELECT character_id FROM character WHERE player_id=$1 AND name=$2;", authorId, m.Content)

				switch err = charRow.Scan(&selectedCharacter); err {
				case sql.ErrNoRows:
					
					s.ChannelMessageSend(m.ChannelID, "No character found, type -char show if you forgot about your characters name")
					s.AddHandlerOnce(selectCharacter)
					return
				case nil:
					// Select the opponents and their respective character informations
					duelPlayerRow := database.DB.QueryRow(fmt.Sprintln("SELECT challenger, challenged, challengedchar, challengerchar FROM duelplayers WHERE preparationid=", currentDuel.Id))
					duelPlayers := game.DuelPlayer{}

					switch err = duelPlayerRow.Scan(&duelPlayers.Challenger, &duelPlayers.Challenged, &duelPlayers.ChallengerChar, &duelPlayers.ChallengedChar); err {
					case sql.ErrNoRows:
						fmt.Println(duelPlayers)
						s.ChannelMessageSend(m.ChannelID, "No duel preparation found")
						s.AddHandlerOnce(selectCharacter)
						return
					case nil:

						// Insert the selected character in DB and return the controller function
						if currentDuel.SelectingPlayer == duelPlayers.Challenger {
							_, err = database.DB.Exec(`UPDATE duelplayers SET challengerchar=$1 WHERE preparationid=$2;`, selectedCharacter, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}

							_, err = database.DB.Exec(`UPDATE duelpreparation SET selectingPlayer=$1 WHERE id=$2;`, duelPlayers.Challenged, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}
							selectCharacterBase(s, m.ChannelID, []string{duelPlayers.Challenger, duelPlayers.Challenged}, 0)
						} else if currentDuel.SelectingPlayer == duelPlayers.Challenged {
							_, err = database.DB.Exec(`UPDATE duelplayers SET challengedchar=$1 WHERE preparationid=$2;`, selectedCharacter, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}

							_, err = database.DB.Exec(`UPDATE duelpreparation SET isReady=1 WHERE id=$1;`, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}

							selectCharacterBase(s, m.ChannelID, []string{duelPlayers.Challenger, duelPlayers.Challenged}, 1)
						} else {
							s.ChannelMessageSend(m.ChannelID, "Error loading character")
							s.AddHandlerOnce(selectCharacter)
							return
						}
					}

				}
				if err != nil {
					fmt.Println(selectedCharacter)
					fmt.Println(err.Error())
				}
			} else {
				// Repeat the handler if the user isn't the selecting player
				s.ChannelMessageSend(m.ChannelID, "It is not your turn to choose a character")
				s.AddHandlerOnce(selectCharacter)
			}
		}
	}
}
