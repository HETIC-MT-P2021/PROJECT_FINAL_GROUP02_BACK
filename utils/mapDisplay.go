package utils

import (
	"fmt"
	"strconv"

	"github.com/SteakBarbare/RPGBot/game"
)

func handleTileToDisplayString(tile game.DungeonTile) string {
	playerD := "P"
	wallD:= "#"
	unknownD := " ? "
	exitD := "E"
	emptyD := "  "

	
	if tile.IsImpassable {
		return wallD
	}

	if !tile.IsDiscovered{
		return unknownD
	}

	if tile.IsExit {
		return exitD
	}

	if len(tile.Characters) > 0 {
		return playerD
	}

	if len(tile.Entities) > 0 || len(tile.Events) > 0 {
		return strconv.Itoa(len(tile.Entities) + len(tile.Events))
	}

	return emptyD
}

func DungeonTilesToString(dungeonTiles []game.DungeonTile) string {
	var dungeonDisplay string
	var dungeonMap [5][5]string

	for i := 0; i < len(dungeonTiles); i++ {
		dungeonTile := dungeonTiles[i]
		
		displayStyle := handleTileToDisplayString(dungeonTile)

		dungeonMap[dungeonTile.Y][dungeonTile.X] = displayStyle
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			dungeonDisplay = dungeonDisplay + "[" + dungeonMap[y][x] + "]"
		}

		dungeonDisplay += "\n"
	}

	return dungeonDisplay
}

func GenerateTileInfoDisplay(tile game.DungeonTile) string {
	var characterNames string
	var entityNames string
	var eventDisplay string
	var exitString string

	if tile.IsDiscovered {
		if len(tile.Characters) > 0 {
			for _, character := range tile.Characters {
				characterNames += character.Name + ", "
			}
		} else {
			characterNames = "no one here"
		}
	
		if len(tile.Entities) > 0 {
			for _, entity := range tile.Entities {
				entityNames += entity.Name + ", "
			}
		} else {
			entityNames = "no entity here"
		}

		if len(tile.Events) > 0 {
			for _, event := range tile.Events {
				wasActivatedDisplay := "Not yet activated"

				if event.WasActivated {
					wasActivatedDisplay = "Was activated"
				}

				isAlwaysActiveDisplay := "Won't reactivate after discovery"

				if event.IsAlwaysActive {
					isAlwaysActiveDisplay = "Will reactivate each time"
				}

				eventDisplay += "\n -" + event.Name + ", " + event.Description + ", " + wasActivatedDisplay + ", " + isAlwaysActiveDisplay + "\n"
			}
		} else {
			eventDisplay = "No events here"
		}

		if tile.IsExit {
			exitString = "Yes !"
		} else {
			exitString = "No !"
		}
	} else {
		entityNames, characterNames, eventDisplay, exitString = "Unknown", "Unknown", "Unknown", "Unknown"
	}
	
	return fmt.Sprintln(
		"**ID:** ", strconv.Itoa(tile.Id),
		"**X:** ", strconv.Itoa(tile.X),
		"**Y:** ", strconv.Itoa(tile.Y),
		"\n**Characters present: ** ", characterNames,
		"\n**Entities present: ** ", entityNames,
		"\n**Events: ** ", eventDisplay,
		"\n**Is tile explored: ** ", strconv.FormatBool(tile.IsDiscovered),
		"\n**Is this a wall: ** ", strconv.FormatBool(tile.IsImpassable),
		"\n**Is this the Exit: ** ", exitString,
	)
}