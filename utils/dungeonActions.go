package utils

import (
	"errors"

	"github.com/SteakBarbare/RPGBot/game"
)

func findPlayerPosInDungeon(playerId int64, dungeonTiles []game.DungeonTile) (game.DungeonTile, game.Character, bool) {
	var foundTile game.DungeonTile
	var foundCharacter game.Character
	isFound := false

	for _, tile := range dungeonTiles {
		if len(tile.Characters) > 0 && !isFound {
			for _, character := range tile.Characters {
				if character.PlayerId == playerId {
					foundTile = tile
					foundCharacter = character
					isFound = true
					break
				}
			}
		}
	}

	return foundTile, foundCharacter, isFound
}

func tranlasteDirectionToNewPos(direction string, posX, posY int) (int, int, error) {
	var newPosX, newPosY int

	switch direction {
	case "top", "up":
		newPosX = posX
		newPosY = posY - 1 
		break
	case "right":
		newPosX = posX + 1 
		newPosY = posY
		break
	case "bot", "bottom":
		newPosX = posX
		newPosY = posY + 1
		break
	case "left":
		newPosX = posX - 1
		newPosY = posY
		break
	default:
		return newPosX, newPosY, errors.New("Direction wasn't correct")
	}

	if newPosX < 0 || newPosY < 0 {
		return newPosX, newPosY, errors.New("Cannot go out of the map")
	}

	return newPosX, newPosY, nil
}

func checkTileIsGood(dungeonTiles []game.DungeonTile, newPosX, newPosY int) (game.DungeonTile, error) {
	var newTile game.DungeonTile
	found := false

	for _, tile := range dungeonTiles {
		if tile.X == newPosX && tile.Y == newPosY {
			newTile = tile
			found = true
			break
		}
	}

	if !found {
		return newTile, errors.New("This tile doesn't exist")

	}

	if newTile.IsImpassable {
		return newTile, errors.New("This tile is blocked")
	}

	return newTile, nil
}

func HandleTileMove(direction string, playerId int64) (string, error) {
	dungeon, err := GetPlayerCurrentStartedDungeon(playerId)

	if err != nil {
		return "", err
	}

	dungeonTiles, err := GetFullDungeonTiles(dungeon.Id)

	if err != nil {
		return "", err
	}

	oldTile, character, isFound := findPlayerPosInDungeon(playerId, dungeonTiles)

	if !isFound {
		return "", errors.New("Player not found in the dungeon")
	}

	newPosX, newPosY, err := tranlasteDirectionToNewPos(direction, oldTile.X, oldTile.Y)

	if err != nil {
		return "", err
	}

	newPlayerTile, err := checkTileIsGood(dungeonTiles, newPosX, newPosY)

	if err != nil {
		return "", err
	}

	newDungeonTiles, err := UpdatePlayerTile(character, dungeonTiles, oldTile, newPlayerTile)

	if err != nil {
		return "", err
	}

	newMapString := DungeonTilesToString(newDungeonTiles)

	return newMapString, nil
}	
