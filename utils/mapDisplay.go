package utils

import "github.com/SteakBarbare/RPGBot/game"

var DiplayMessage = []string{"[ ]", "[#]", "[P]", "[E]"}

func DungeonTilesToString(dungeonTiles []game.DungeonTile, playerPosX, playerPosY int) string {
	var dungeonDisplay string
	var dungeonMap [5][5]int

	for i := 0; i < len(dungeonTiles); i++ {
		dungeonTile := dungeonTiles[i]

		displayInt := 0

		if dungeonTile.IsExit {
			displayInt = 3
		}
		if dungeonTile.X == playerPosX && playerPosY == dungeonTile.Y {
			displayInt = 2
		}
		if dungeonTile.IsImpassable {
			displayInt = 1
		}


		dungeonMap[dungeonTile.X][dungeonTile.Y] = displayInt
	}

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			displayInt := dungeonMap[x][y]

			displayString := DiplayMessage[displayInt]

			dungeonDisplay = dungeonDisplay + displayString
		}

		dungeonDisplay += "\n"
	}

	return dungeonDisplay
}