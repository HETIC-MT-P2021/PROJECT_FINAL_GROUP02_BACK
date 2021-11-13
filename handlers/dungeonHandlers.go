package handlers

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)


func selectDunjeonCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "-quit" {
		s.ChannelMessageSend(m.ChannelID, "Aborting dungeon character selection")
		return
	} else if m.Content == "-char Show" {
		s.AddHandlerOnce(selectDunjeonCharacter)
		return
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
		 AND is_alive=true;`

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
					s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "an error occured:" + err.Error()))
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

		dungeonTiles, err := utils.InitDungeonTiles(character.Id, dungeon)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Couldn't create dungeon, please retry")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		err = utils.UpdateDungeonCharacter(character.Id, dungeon.Id)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Couldn't create dungeon, please retry")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		err = utils.LinkCharacterDungeon(dungeon.Id, character.Id)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Couldn't create dungeon please retry")
			s.AddHandlerOnce(selectDunjeonCharacter)

			return
		}

		displayDungeonString := utils.DungeonTilesToString(dungeonTiles)

		s.ChannelMessageSend(m.ChannelID, "SuccessFully generated dungeon map ! \n\n" + displayDungeonString + "\n\nID of the Dungeon :" + strconv.FormatInt(int64(dungeon.Id), 10))		

		return
	}
}

func selectDunjeonToPlay(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "-quit" {
		s.ChannelMessageSend(m.ChannelID, "Aborting dungeon selection")

		return
	} else if m.Content == "-char Show" {
		s.AddHandlerOnce(selectDunjeonToPlay)

		return
	} else if m.Content == "-dungeon list" {
		s.AddHandlerOnce(selectDunjeonToPlay)

		return
	} else {
		authorId, err := strconv.ParseInt(m.Author.ID, 10, 64)

		if err != nil {
			panic(err)
		}
		parsedDungeonId, err := strconv.ParseInt(m.Content, 10, 64)
		
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Message error", "Enter a correct ID"))
			s.AddHandlerOnce(selectDunjeonToPlay)

			return
		}

		currentDugeon, err := utils.GetPlayerCurrentStartedDungeon(authorId)

		if err == nil {
			if currentDugeon.Id == int(parsedDungeonId) {
				s.ChannelMessageSend(m.ChannelID, "You are already on a dungeon adventure with dungeon id:" + strconv.FormatInt(int64(currentDugeon.Id), 10) + "\n You can move")
				s.AddHandlerOnce(dungeonTileMove)

				return
			} else {
				utils.UpdateDungeonIsPaused(currentDugeon.Id, true)
				s.AddHandlerOnce(selectDunjeonToPlay)
				s.ChannelMessageSend(m.ChannelID, "You are already on a dungeon adventure with dungeon id:" + strconv.FormatInt(int64(currentDugeon.Id), 10) + "\n It is now paused, you can enter a new dungeon id to play")

				return
			}
		}

		dungeon, err := utils.GetPlayerReadyDungeon(parsedDungeonId, authorId)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Selected dungeon not found,\n -dungeon list to list yours")
			return
		}

		if !dungeon.HasStarted {
			err = utils.UpdateDungeonHasStarted(dungeon.Id)
	
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, "Selected dungeon coulnd't be started")
				return
			}
		}

		if dungeon.IsPaused {
			err = utils.UpdateDungeonIsPaused(dungeon.Id, false)

			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, "Selected dungeon coulnd't be un paused")
				return
			}
		}

		dungeonTiles, err := utils.GetFullDungeonTiles(dungeon.Id)

		dungeonString := utils.DungeonTilesToString(dungeonTiles)

		instructionString := `
		You can now select where you want to go with:
		 -dungeon move [left, right, top, bot]

		You can also pause the exploration with:
		 -dungeon pause or -quit
		 
		On the exit tile, you will be able to leave with:
		 -dungeon exit`

		s.ChannelMessageSend(m.ChannelID, "Successfully found the dungeon, here's the current map of it ! \n\n"+ dungeonString + instructionString)

		s.AddHandlerOnce(dungeonTileMove)

		return
	}
}

func dungeonTileMove(s *discordgo.Session, m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID {
		return
	}

	authorId, err := strconv.ParseInt(m.Author.ID, 10, 64)

	if err != nil {
		panic(err)
	}

	if m.Content == "-quit" || m.Content == "-dungeon pause"{
		dungeon, err := utils.GetPlayerCurrentStartedDungeon(authorId)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot Error", "Couldn't find current dungeon"))

			return
		}

		err = utils.UpdateDungeonIsPaused(dungeon.Id, true)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot Error", "Couldn't stop current dungeon"))

			return
		}

		s.ChannelMessageSend(m.ChannelID, "Pausing dungeon, you can restart again with -dungeon start")

		return
	} else if m.Content == "-dungeon exit"{
		dungeon, err := utils.GetPlayerCurrentStartedDungeon(authorId)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot Error", "Couldn't find current dungeon"))

			return
		}

		characterId, err := utils.GetPlayerDungeonLinkedCharacter(dungeon.Id, authorId)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot Error", "Couldn't find your character"))

			return
		}

		characterTile, err := utils.GetCharacterTile(characterId, dungeon.Id)

		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot Error", "Couldn't find your character's tile"))

			return
		}

		if characterTile.IsExit == true {
			utils.EndDungeon(dungeon.Id, characterId)
			s.ChannelMessageSend(m.ChannelID, "You successfully left the dungeon, this character is now available")

			return
		} else {
			s.ChannelMessageSend(m.ChannelID, "You are not on a exit Tile")
			s.AddHandlerOnce(dungeonTileMove)

			return
		}
	} else {
		messageSplit := strings.Split(m.Content, " ")

		if messageSplit[0] != "-dungeon" || messageSplit[1] != "move" {
			s.AddHandlerOnce(dungeonTileMove)

			return
		}

		direction := messageSplit[2]

		newDungeonTiles, wasTileAlreadyDiscovered, dungeonId, err := utils.HandleTileMove(direction, authorId)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Couldn't move there ! \n:"+ err.Error())

			s.AddHandlerOnce(dungeonTileMove)

			return
		}

		newMapString := utils.DungeonTilesToString(newDungeonTiles)

		if !wasTileAlreadyDiscovered {
			s.ChannelMessageSend(m.ChannelID, "You arrive in a new room! \n\n" + newMapString)
		} else {
			s.ChannelMessageSend(m.ChannelID, "You have already been there, \n\n" + newMapString)
		}
		
		err, updatedCharacter := utils.HandleNewTileEvents(newDungeonTiles, s, m, authorId)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "an error occured:" + err.Error()))
		}

		if !updatedCharacter.IsAlive{
			err = utils.EndDungeon(dungeonId, updatedCharacter.Id)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "an error occured:" + err.Error()))
			}

			s.ChannelMessageSend(m.ChannelID, "This dungeon is over.")

			return
		}

		s.ChannelMessageSend(m.ChannelID, "\n\n What's your next move ?")

		s.AddHandlerOnce(dungeonTileMove)

		return
	}
}
