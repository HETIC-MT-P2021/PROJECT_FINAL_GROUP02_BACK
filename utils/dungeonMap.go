package utils

import (
	"errors"
	"log"
	"math/rand"

	"github.com/SteakBarbare/RPGBot/game"
)

func findRandomNonBlockerPosInPattern(tilePattern [][]int, exitPosX, exitPosY int)(int, int){
	var posx int;
	var posy int;

	posx = rand.Intn(4)
	posy = rand.Intn(4)

	for tilePattern[posx][posy] != 0 && posx != exitPosX && posy != exitPosY {
		posx = rand.Intn(4)
		posy = rand.Intn(4)
	}

	return posx, posy
}

func InitDungeonTiles(characterInstanceId int, dungeon *game.Dungeon) ([]game.DungeonTile, int, int, error){
	basePattern := [][]int{{0,0,0,1,0}, {0,1,0,1,0}, {0,0,0,0,0}, {0,0,1,0,0}, {0,0,1,0,0}}

	exitPostX, exitPosY := findRandomNonBlockerPosInPattern(basePattern, 5, 5)
	playerPosX, playerPosY := findRandomNonBlockerPosInPattern(basePattern, exitPostX, exitPosY)

	var dungeonTiles []game.DungeonTile

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			var tile game.DungeonTile

			tile.DungeonId = dungeon.Id

			tile.X = x
			tile.Y = y

			if basePattern[x][y] == 1 {
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
				return dungeonTiles, playerPosX, playerPosY, errors.New("Dungeon tiles couln't be created")
			}

			if playerPosX == x && playerPosY == y {
				linkCharacterInstanceToTile(tileId, characterInstanceId)
			}

			dungeonTiles = append(dungeonTiles, tile)
		}
	}

	return dungeonTiles, playerPosX, playerPosY, nil

	// paterne de creation [[0,0,0,1,0], [0,1,0,1,0], [0,0,0,0,0], [0,0,1,0,0], [0,0,1,0,0]]
	// rdm 0-> 4 2 fois pour avoir exit x/y, while [x][y] = 1, refaire
	// rdm 0-> 4 2 fois pour avoir playerPos x/y, while [x][y] = 1, refaire
	// pour chaque row
	// pour chaque case
	// if 0 => tile with impassable =false else true
	// if exit set exit
	// if playerPos add init playerModelInstance add id
	// return array of tiles and playerPos
}