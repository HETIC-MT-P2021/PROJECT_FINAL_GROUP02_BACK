package utils

import (
	"strconv"

	"github.com/SteakBarbare/RPGBot/game"
)

func handleTileToDisplayString(tile game.DungeonTile) string {
	var tileDisplayString string

	playerD := "P"
	wallD:= "####"
	unknownD := "?????"
	exitD := "EXT"
	eventD := "EV"
	monsterD := "M"
	spaceD := "  "
	halfSpaceD := " "
	emptyD := "            "

	
	if tile.IsImpassable {
		return wallD
	}

	if !tile.IsDiscovered{
		return unknownD
	}

	if tile.IsExit {
		tileDisplayString += exitD
	}

	if len(tile.Entities) > 0 {
		tileDisplayString += strconv.Itoa(len(tile.Entities)) + monsterD
	}

	if len(tile.Events) > 0 {
		tileDisplayString += strconv.Itoa(len(tile.Events)) + eventD
	}

	if len(tile.Characters) > 0 {
		tileDisplayString += strconv.Itoa(len(tile.Characters)) + playerD
	}

	if len(tileDisplayString) == 0 {
		return emptyD
	}

	if len(tileDisplayString) < 4 {
		return spaceD + tileDisplayString + spaceD
	}

	if len(tileDisplayString) == 4 {
		return halfSpaceD + tileDisplayString + halfSpaceD
	}

	return tileDisplayString
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