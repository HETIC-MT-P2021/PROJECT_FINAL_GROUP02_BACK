package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/SteakBarbare/RPGBot/game"
	"github.com/bwmarrin/discordgo"
)

func findRandomNonBlockerPosInPattern(tilePattern [][]int, exitPosX, exitPosY int)(int, int){
	var posx int;
	var posy int;

	posx = rand.Intn(5)
	posy = rand.Intn(5)

	for tilePattern[posy][posx] != 0 && posx != exitPosX && posy != exitPosY {
		posx = rand.Intn(5)
		posy = rand.Intn(5)
	}

	return posx, posy
}

func InitDungeonTiles(characterId int, dungeon *game.Dungeon) ([]game.DungeonTile, error){
	basePattern := [][]int{{0,0,0,1,0}, {0,1,0,1,0}, {0,0,0,0,0}, {0,0,1,0,0}, {0,0,1,0,0}}

	exitPostX, exitPosY := findRandomNonBlockerPosInPattern(basePattern, -1, -1)
	playerPosX, playerPosY := findRandomNonBlockerPosInPattern(basePattern, exitPostX, exitPosY)

	var dungeonTiles []game.DungeonTile

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			var tile game.DungeonTile

			tile.DungeonId = dungeon.Id

			tile.X = x
			tile.Y = y

			if basePattern[y][x] == 1 {
				tile.IsImpassable = true
			} else {
				tile.IsImpassable = false
			}

			if exitPostX == x && exitPosY == y {
				tile.IsExit =true;
			} else {
				tile.IsExit = false
			}

			tileId, tileCreationErr := createDungeonTile(tile)

			if tileCreationErr != nil {
				log.Println(tileCreationErr)
				return dungeonTiles, errors.New("Dungeon tiles couln't be created")
			}

			if playerPosX == x && playerPosY == y {
				var char game.Character

				linkCharacterToTile(tileId, characterId)

				updateTileIsDiscovered(tileId)

				tile.IsDiscovered = true
				tile.Characters = append(tile.Characters, char)
			}

			dungeonTiles = append(dungeonTiles, tile)
		}
	}

	return dungeonTiles, nil
}

func showCharacterName (hasCharacter bool, characterName string) (string){
	if !hasCharacter {
		return "None"
	} else {
		return characterName
	}
}

func DisplayDungeonList (s *discordgo.Session, m *discordgo.MessageCreate, dungeons []game.Dungeon) (error) {
	s.ChannelMessageSend(m.ChannelID, "Here's the list of your dungeons :")
	var err error

	for _, dungeon := range dungeons {
		var characterName string
		
		hasDungeonCharacter := dungeon.SelectedCharacterId.Valid

		if hasDungeonCharacter {
			characterName, _ = FindCharNameWithId(int(dungeon.SelectedCharacterId.Int32))
		}

		// Show the new character stats & name
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintln("Here's one of your dungeon"),
			Description: fmt.Sprintln(
				"**ID:** ", strconv.Itoa(dungeon.Id),
				"\n**Selected Character:** ", showCharacterName(hasDungeonCharacter, characterName),
				"\n**Created At:** ", strconv.Itoa(int(dungeon.CreatedAt.Month())) + "/" + strconv.Itoa(dungeon.CreatedAt.Day()),
				"\n**Has dungeon started:** ", strconv.FormatBool(dungeon.HasStarted),
				"\n**Has dungeon Ended:** ", strconv.FormatBool(dungeon.HasEnded),
				"\n**Is dungeon paused:** ", strconv.FormatBool(dungeon.IsPaused),
			),
			Color: 0x00ff99,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Player: " + m.Author.ID,
			},
		})

	}

	return err
}

func GetCloseTiles(centerTile game.DungeonTile, dungeonTiles[]game.DungeonTile)([]game.DungeonTile){
	var closeByTiles []game.DungeonTile

	for _, tile := range dungeonTiles {
		if tile.Id == centerTile.Id ||
		 (tile.X == centerTile.X && tile.Y == centerTile.Y + 1) ||
		 (tile.X == centerTile.X && tile.Y == centerTile.Y - 1)||
		 (tile.X == centerTile.X + 1 && tile.Y == centerTile.Y)||
		 (tile.X == centerTile.X - 1 && tile.Y == centerTile.Y) {
			closeByTiles = append(closeByTiles, tile)
		} 
	}

	return closeByTiles
}

func InsertEventAndLinkToTiles(event game.Event, tiles []game.DungeonTile)error {
	var buffEffectId int

	buffEffectId, err := InsertEvent(event)

	if err !=nil {
		return err
	}

	for _, tile := range tiles {
		if !tile.IsExit || !tile.IsImpassable {
			err = LinkEventTile(buffEffectId, tile.Id)

			if err !=nil {
				return err
			}
		}
	}

	return nil
}