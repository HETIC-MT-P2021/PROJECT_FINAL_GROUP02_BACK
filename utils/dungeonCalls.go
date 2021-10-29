package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetPlayerReadyDungeon(dungeonId, playerId int64) (*game.Dungeon, error) {
	query := `SELECT dungeon_id, created_at, created_by, selected_character_id, has_started, has_ended, is_paused FROM dungeon 
		WHERE dungeon_id=$1
		AND created_by=$2
		AND selected_character_id IS NOT NULL
		AND (is_paused=true OR has_started=false)`;

	// Get the character from db
	dungeonRow := database.DB.QueryRow(query, dungeonId, playerId)

	playableDungeon := game.Dungeon{}
	
	err := dungeonRow.Scan(&playableDungeon.Id, &playableDungeon.CreatedAt, &playableDungeon.CreatedBy, &playableDungeon.SelectedCharacterId, &playableDungeon.HasStarted, &playableDungeon.HasEnded, &playableDungeon.IsPaused)

	if err != nil {
		fmt.Println(err.Error())
		return &playableDungeon, errors.New("Dungeon not found")
	} else {
		return &playableDungeon, nil
	}
}

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

func unlinkCharacterToTile(tileId, characterId int) error {
	query := `DELETE FROM link_character_tile
	 WHERE tile_id=$1 
	 AND character_id=$2`

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

	query := `SELECT dungeon_id, created_at, created_by, has_started, has_ended, is_paused, selected_character_id FROM dungeon 
		WHERE created_by=$1`;
		
	rows, err := database.DB.Query(query, playerId)

	if err != nil {
		return playerDungeons, err
	}
	
	for rows.Next() {
		var dungeon game.Dungeon

		if err = rows.Scan(&dungeon.Id, &dungeon.CreatedAt, &dungeon.CreatedBy,
			&dungeon.HasStarted, &dungeon.HasEnded, &dungeon.IsPaused, &dungeon.SelectedCharacterId); err != nil {
			return playerDungeons, err
		}
		

		playerDungeons = append(playerDungeons, dungeon)
	}

	return playerDungeons, nil
}

func UpdateDungeonHasStarted(dungeonId int) error {
	query := `UPDATE dungeon 
	 SET has_started = true
	 WHERE dungeon_id=$1`

	_, err := database.DB.Exec(query, dungeonId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon could not be updated")
	}

	return nil
}

func UpdateDungeonIsPaused(dungeonId int, isPaused bool) error {
	query := `UPDATE dungeon 
	 SET is_paused = $1
	 WHERE dungeon_id=$2`

	_, err := database.DB.Exec(query, isPaused, dungeonId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon could not be updated")
	}

	return nil
}


func FetchDungeonTiles(dungeonId int) ([]game.DungeonTile, error) {
	var dungeonTiles []game.DungeonTile

	query := `SELECT 
	 tile_id, dungeon_id, x, y, is_discovered, is_exit, is_impassable
	 FROM dungeon_tile 
   	 WHERE dungeon_id=$1`;
		
	rows, err := database.DB.Query(query, dungeonId)

	if err != nil {
		return dungeonTiles, err
	}
	
	for rows.Next() {
		var dungeonTile game.DungeonTile

		if err = rows.Scan(&dungeonTile.Id, &dungeonTile.DungeonId, &dungeonTile.X,
			&dungeonTile.Y, &dungeonTile.IsDiscovered, &dungeonTile.IsExit, &dungeonTile.IsImpassable); err != nil {
			return dungeonTiles, err
		}
		

		dungeonTiles = append(dungeonTiles, dungeonTile)
	}

	return dungeonTiles, nil
}

func GetTileEntities(tileId int) ([]game.EntityInstance, error) {
	var entities []game.EntityInstance

	query := `SELECT 
	 l.entity_instance_id,
	 e.entity_model_id, e.precision, e.strength, e.endurance, e.agility, e.hitpoints, e.chosen_action_id
	 FROM link_entity_tile l 
	 INNER JOIN entity_instance e ON l.entity_instance_id = e.entity_instance_id
   	 WHERE tile_id=$1`;
		
	rows, err := database.DB.Query(query, tileId)

	if err != nil {
		return entities, err
	}
	
	for rows.Next() {
		var entity game.EntityInstance

		if err = rows.Scan(&entity.Id, &entity.ModelId, &entity.Precision,
			&entity.Strength, &entity.Endurance, &entity.Agility, &entity.Hitpoints, entity.ChosenActionId); err != nil {
			return entities, err
		}
		

		entities = append(entities, entity)
	}

	return entities, nil
}

func GetTileCharacter(tileId int) ([]game.Character, error) {
	var characters []game.Character

	query := `SELECT 
	 l.character_id,
	 c.name, c.player_id,
	 c.precision, c.strength, c.endurance, c.agility, c.hitpoints,
	 c.precision_max, c.strength_max, c.endurance_max, c.agility_max, c.hitpoints_max,
	 c.is_occupied, c.is_alive, c.chosen_action_id
	 FROM link_character_tile l 
	 INNER JOIN character c ON l.character_id = c.character_id
   	 WHERE tile_id=$1`;
		
	rows, err := database.DB.Query(query, tileId)

	if err != nil {
		return characters, err
	}
	
	for rows.Next() {
		var character game.Character

		if err = rows.Scan(&character.Id, &character.Name, &character.PlayerId,
			&character.Precision, &character.Strength, &character.Endurance, &character.Agility, &character.Hitpoints,
			&character.PrecisionMax, &character.StrengthMax, &character.EnduranceMax, &character.AgilityMax, &character.HitpointsMax,
			&character.IsOccupied, &character.IsAlive, &character.ChosenActionId); err != nil {
			return characters, err
		}
		

		characters = append(characters, character)
	}

	return characters, nil
}

func GetTileEvents(tileId int) ([]game.EventModel, error) {
	var events []game.EventModel

	query := `SELECT 
	 l.event_id, e.event_name
	 FROM link_event_tile l 
	 INNER JOIN event_model e ON l.event_id = e.event_id
   	 WHERE tile_id=$1`;
		
	rows, err := database.DB.Query(query, tileId)

	if err != nil {
		return events, err
	}
	
	for rows.Next() {
		var event game.EventModel

		if err = rows.Scan(&event.Id, &event.Name); err != nil {
			return events, err
		}
		

		events = append(events, event)
	}

	return events, nil
}

func GetFullDungeonTiles(dungeonId int)([]game.DungeonTile, error) {
	dungeonTiles, err := FetchDungeonTiles(dungeonId)
	
	if err != nil {
		return dungeonTiles, err
	}

	for i := 0; i < len(dungeonTiles); i++ {
		tile := &dungeonTiles[i]

		entities, err := GetTileEntities(tile.Id)

		if err != nil {
			return dungeonTiles, err
		}

		characters, err := GetTileCharacter(tile.Id)

		if err != nil {
			return dungeonTiles, err
		}

		events, err := GetTileEvents(tile.Id)

		if err != nil {
			return dungeonTiles, err
		}

		(*tile).Entities = entities
		(*tile).Characters = characters
		(*tile).Events = events
	}	

	return dungeonTiles, nil
}

func GetPlayerCurrentStartedDungeon(playerId int64) (game.Dungeon, error) {
	query := `SELECT dungeon_id, created_at, selected_character_id FROM dungeon 
		WHERE created_by=$1
		AND has_started=true
		AND has_ended=false
		AND is_paused=false`;

	// Get the character from db
	dungeonRow := database.DB.QueryRow(query, playerId)

	dungeon := game.Dungeon{}
	
	err := dungeonRow.Scan(&dungeon.Id, &dungeon.CreatedAt, &dungeon.SelectedCharacterId)
	if err != nil {
		fmt.Println(err.Error())
		return dungeon, errors.New("Dungeon not found")
	} 

	dungeon.CreatedBy = playerId
	dungeon.HasEnded = false
	dungeon.HasStarted = true
	dungeon.IsPaused = false

	return dungeon, nil
}

func updateTileIsDiscovered(tileId int) error {
	query := `UPDATE dungeon_tile 
	 SET is_discovered = true
	 WHERE tile_id=$1`

	_, err := database.DB.Exec(query, tileId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon tile could not be updated")
	}

	return nil
}

func DiscoverTile(tile game.DungeonTile)(game.DungeonTile, error){
	// TODO: generate event/entities
	err := updateTileIsDiscovered(tile.Id)

	if err != nil {
		return tile, err
	}

	tile.IsDiscovered = true

	return tile, nil
}

func UpdatePlayerTile(character game.Character, dungeonTiles []game.DungeonTile, oldTile, newTile game.DungeonTile) ([]game.DungeonTile, error){
	err := linkCharacterToTile(newTile.Id, character.Id)

	if err != nil {
		return dungeonTiles, err
	}

	err = unlinkCharacterToTile(oldTile.Id, character.Id)

	if err != nil {
		return dungeonTiles, err
	}

	for i := 0; i < len(dungeonTiles); i++ {
		tile := &dungeonTiles[i]

		if tile.Id == newTile.Id {
			(*tile).Characters = append(tile.Characters, character)

			if !tile.IsDiscovered{
				updatedTile, err := DiscoverTile(*tile)

				if err != nil {
					return dungeonTiles, err
				}
	
				(*tile) = updatedTile
			}
		}

		if tile.Id == oldTile.Id {
			var newCharacters []game.Character

			for _, oldCharacter := range tile.Characters {
				if oldCharacter.Id != character.Id {
					newCharacters = append(newCharacters, oldCharacter)
				}
			}

			(*tile).Characters = newCharacters
		}
	}

	return dungeonTiles, err
}

func LinkCharacterDungeon(dungeonId, character_id int) error {
	query := `INSERT INTO link_character_dungeon
	 	(dungeon_id, character_id)
	 	VALUES ($1, $2)`

	_, err := database.DB.Exec(query, dungeonId, character_id)

	if err != nil {
		log.Println(err)
		return errors.New("LinkDungeonCharacter could not be created")
	}

	return nil
}

func GetPlayerDungeonLinkedCharacter(dungeonId int, playerId int64) (int, error) {
	var characterId int

	query := `SELECT l.character_id
	 FROM link_character_dungeon l
	 INNER JOIN character c ON l.character_id=c.character_id
   	 WHERE l.dungeon_id=$1
	 AND c.player_id=$2`;
		
	row := database.DB.QueryRow(query, dungeonId, playerId)
	
	err := row.Scan(&characterId)

	if err != nil {
		fmt.Println(err.Error())
		return characterId, errors.New("Linked Character not found")
	} else {
		return characterId, nil
	}
}

func GetCharacterTile(characterId, dungeonId int) (game.DungeonTile, error) {
	query := `SELECT l.tile_id, t.x, t.y, t.is_discovered, t.is_exit, t.is_impassable
	 FROM link_character_tile l 
	 INNER JOIN dungeon_tile t ON t.tile_id = l.tile_id
	 WHERE t.dungeon_id =$1 
	 AND l.character_id=$2`;

	// Get the character from db
	tileRow := database.DB.QueryRow(query, dungeonId, characterId)

	characterTile := game.DungeonTile{}
	
	err := tileRow.Scan(&characterTile.Id, &characterTile.X, &characterTile.Y, &characterTile.IsDiscovered, &characterTile.IsExit, &characterTile.IsImpassable)

	if err != nil {
		fmt.Println(err.Error())
		return characterTile, errors.New("Dungeon not found")
	} else {
		return characterTile, nil
	}
}

func EndDungeon(dungeonId, characterId int) error {
	query := `UPDATE dungeon 
	 SET is_paused = true, has_ended = true
	 WHERE dungeon_id=$1`

	_, err := database.DB.Exec(query, dungeonId)

	if err != nil {
		log.Println(err)
		return errors.New("Dungeon could not be updated")
	}

	characterQuery := `UPDATE character 
	 SET is_occupied = false
	 WHERE character_id=$1`

	_, err = database.DB.Exec(characterQuery, characterId)

	if err != nil {
		log.Println(err)
		return errors.New("Character could not be updated")
	}

	return nil
}
