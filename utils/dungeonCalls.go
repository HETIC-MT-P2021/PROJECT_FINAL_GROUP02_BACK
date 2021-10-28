package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetPlayerNotStartedDungeon(player_id int64) (*game.Dungeon, error) {
	query := `SELECT dungeon_id, created_at, created_by, has_started FROM dungeon 
		WHERE created_by=$1 
		AND selected_character_id IS NULL`;

	// Get the character from db
	dungeonRow := database.DB.QueryRow(query, player_id)

	notInitiatedDungeon := game.Dungeon{}
	
	err := dungeonRow.Scan(&notInitiatedDungeon.Id, &notInitiatedDungeon.CreatedAt, &notInitiatedDungeon.CreatedBy, &notInitiatedDungeon.HasStarted)

	if err != nil {
		fmt.Println(err.Error())
		return &notInitiatedDungeon, errors.New("Dungeon not found")
	} else {
		return &notInitiatedDungeon, nil
	}
}


func SaveInitDungeon(playerId int64) (error) {
	query := `INSERT INTO dungeon
	 	(created_by)
	 	VALUES ($1)`

	_, err := database.DB.Exec(query,
			playerId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon not created")
	}

	return nil
}

func linkCharacterToTile(tileId, characterId int) error {
	query := `INSERT INTO link_character_tile
	 	(tile_id, character_id)
	 	VALUES ($1, $2)`

	_, err := database.DB.Exec(query, tileId, characterId)

	if err != nil {
		log.Println(err)
		return errors.New("Link character tile could not be created")
	}

	return nil
}

func createDungeonTile(tile game.DungeonTile) (int, error) {
	query := `INSERT INTO dungeon_tile
	 	(dungeon_id, x, y, is_exit, is_impassable)
	 	VALUES ($1, $2, $3, $4, $5) RETURNING tile_id`

	var tileId int

	err := database.DB.QueryRow(query, tile.DungeonId, tile.X, tile.Y, tile.IsExit, tile.IsImpassable).Scan(&tileId)

	if err != nil {
		log.Println(err)
		return tileId, errors.New("DungeonTile could not be created")
	}

	return tileId, nil
}

func UpdateDungeonCharacter(characterId, dungeonId int) error {
	query := `UPDATE dungeon 
	 SET selected_character_id = $1
	 WHERE dungeon_id=$2`

	_, err := database.DB.Exec(query, characterId, dungeonId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon could not be updated")
	}

	updateCharQuery := `UPDATE character
	 SET is_occupied = true
	 WHERE character_id=$1`

	_, err = database.DB.Exec(updateCharQuery, characterId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon could not be updated")
	}

	return nil
}

func GetPlayerDungeons(playerId int64) ([]game.Dungeon, error){
	var playerDungeons []game.Dungeon

	query := `SELECT dungeon_id, created_at, created_by, has_started, has_ended, selected_character_id FROM dungeon 
		WHERE created_by=$1`;
		
	rows, err := database.DB.Query(query, playerId)

	if err != nil {
		return playerDungeons, err
	}
	
	for rows.Next() {
		var dungeon game.Dungeon

		if err = rows.Scan(&dungeon.Id, &dungeon.CreatedAt, &dungeon.CreatedBy,
			&dungeon.HasStarted, &dungeon.HasEnded, &dungeon.SelectedCharacterId); err != nil {
			return playerDungeons, err
		}
		

		playerDungeons = append(playerDungeons, dungeon)
	}

	return playerDungeons, nil
}
