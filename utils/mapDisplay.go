package utils

import "github.com/SteakBarbare/RPGBot/game"

var DiplayMessage = []string{"[?]", "[#]", "[P]", "[E]", "[  ]"}

func DungeonTilesToString(dungeonTiles []game.DungeonTile) string {
	var dungeonDisplay string
	var dungeonMap [5][5]int

	for i := 0; i < len(dungeonTiles); i++ {
		dungeonTile := dungeonTiles[i]

		displayInt := 0

		if dungeonTile.IsDiscovered {
			displayInt = 4
		}
		if dungeonTile.IsExit {
			displayInt = 3
		}
		if dungeonTile.IsImpassable {
			displayInt = 1
		}
		if len(dungeonTile.Characters) > 0 {
			displayInt = 2
		}

		dungeonMap[dungeonTile.Y][dungeonTile.X] = displayInt
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			displayInt := dungeonMap[y][x]

			displayString := DiplayMessage[displayInt]

			dungeonDisplay = dungeonDisplay + displayString
		}

		dungeonDisplay += "\n"
	}

	return dungeonDisplay
}