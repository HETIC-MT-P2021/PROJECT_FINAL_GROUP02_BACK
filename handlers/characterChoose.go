package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/duels"
	"github.com/SteakBarbare/RPGBot/game"
	"github.com/SteakBarbare/RPGBot/utils"
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

func selectDunjeonCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "-quit" {
		s.ChannelMessageSend(m.ChannelID, "Aborting dungeon character selection")
	} else if m.Content == "-char Show" {
		s.AddHandlerOnce(selectDunjeonCharacter)
	} else {
		authorId, err := strconv.ParseInt(m.Author.ID, 10, 64)

		if err != nil {
			panic(err)
		}

		dungeon, err :=  utils.GetPlayerNotStartedDungeon(authorId)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "No Active dungeon creation found, aborting")
			return
		}

		var selectedCharacterId string

		selecteCharQuery := `SELECT character_id 
		 FROM character 
		 WHERE name=$1 
		 AND player_id=$2 
		 AND is_occupied=false 
		 AND is_alive=true`

		// Get the character from db
		charRow := database.DB.QueryRow(selecteCharQuery, m.Content, authorId)

		err = charRow.Scan(&selectedCharacterId)

		if err != nil {
			switch err {
				case sql.ErrNoRows:
					s.ChannelMessageSend(m.ChannelID, "Error, character not found or is Busy\n type -char Show if you forgot about your characters name")
					s.AddHandlerOnce(selectDunjeonCharacter)

					return
				default:
					utils.ErrorMessage("Bot error", "an error occured:" + err.Error())
					s.AddHandlerOnce(selectDunjeonCharacter)

					return
			}
		}

		character, err := utils.GetPlayerCharacterByName(authorId, m.Content)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No character found or is Busy\n type -char Show if you forgot about your characters name")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		s.ChannelMessageSend(m.ChannelID, "Character found, generating dungeon map !")

		dungeonTiles, playerPosX, playerPosY, err := utils.InitDungeonTiles(character.Id, dungeon)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Couldn't create dungeon, please retry")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		err = utils.UpdateDungeonCharacter(character.Id, dungeon.Id)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Couldn't create dungeon, please retry")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		displayDungeonString := utils.DungeonTilesToString(dungeonTiles, playerPosX, playerPosY)

		s.ChannelMessageSend(m.ChannelID, "SuccessFully generated dungeon map ! \n\n" + displayDungeonString + "\n\nID of the Dungeon :" + strconv.FormatInt(int64(dungeon.Id), 10))

		
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
		duelPreparation := database.DB.QueryRow(fmt.Sprintln("SELECT * FROM duelPreparation WHERE isReady=", 0))

		currentDuel := game.DuelPreparation{}

		err = duelPreparation.Scan(&currentDuel.Id, &currentDuel.SelectingPlayer, &currentDuel.IsReady, &currentDuel.IsOver, &currentDuel.Turn)

		if err != nil {
			fmt.Println(err.Error())
			s.ChannelMessageSend(m.ChannelID, "No Duel found, aborting duel preparation")
			return
		} else {
			// Check if the message author is the selecting player
			if m.Author.ID == currentDuel.SelectingPlayer {

				// Get the character from db
				charRow := database.DB.QueryRow("SELECT id FROM characters WHERE player_id=$1 AND name=$2;", authorId, m.Content)

				var selectedCharacter string
				switch err = charRow.Scan(&selectedCharacter); err {
				case sql.ErrNoRows:
					s.ChannelMessageSend(m.ChannelID, "No character found, type -char Show if you forgot about your characters name")
					s.AddHandlerOnce(selectCharacter)
					return
				case nil:
					// Select the opponents and their respective character informations
					duelPlayerRow := database.DB.QueryRow(fmt.Sprintln("SELECT challenger, challenged, challengerChar, challengedChar FROM duelPlayers WHERE preparationId=", currentDuel.Id))
					duelPlayers := game.DuelPlayer{}

					switch err = duelPlayerRow.Scan(&duelPlayers.Challenger, &duelPlayers.Challenged, &duelPlayers.ChallengerChar, &duelPlayers.ChallengedChar); err {
					case sql.ErrNoRows:
						s.ChannelMessageSend(m.ChannelID, "No character found, type -char Show if you forgot about your characters name")
						s.AddHandlerOnce(selectCharacter)
						return
					case nil:

						// Insert the selected character in DB and return the controller function
						if currentDuel.SelectingPlayer == duelPlayers.Challenger {
							_, err = database.DB.Exec(`UPDATE duelPlayers SET challengerChar=$1 WHERE preparationId=$2;`, selectedCharacter, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}

							_, err = database.DB.Exec(`UPDATE duelPreparation SET selectingPlayer=$1 WHERE id=$2;`, duelPlayers.Challenged, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}
							selectCharacterBase(s, m.ChannelID, []string{duelPlayers.Challenger, duelPlayers.Challenged}, 0)
						} else if currentDuel.SelectingPlayer == duelPlayers.Challenged {
							_, err = database.DB.Exec(`UPDATE duelPlayers SET challengedChar=$1 WHERE preparationId=$2;`, selectedCharacter, currentDuel.Id)
							if err != nil {
								fmt.Println(err.Error())
							}

							_, err = database.DB.Exec(`UPDATE duelPreparation SET isReady=1 WHERE id=$1;`, currentDuel.Id)
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
