package utils

import (
	"errors"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)
func InsertEntityInstance(entity game.EntityInstance) (int, error) {
	query := `INSERT INTO entity_instance
			(entity_model_id, precision, strength, endurance, agility, hitpoints, chosen_action_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING entity_instance_id;`

	var entityId int

	err := database.DB.QueryRow(query, entity.ModelId, entity.Precision, entity.Strength,
		 entity.Endurance, entity.Agility, entity.Hitpoints, entity.ChosenActionId).Scan(&entityId)

	if err != nil {
		log.Println(err)
		return entityId, errors.New("Entity could not be created")
	}

	return entityId, nil	
}

func LinkEntityTile(entityId, tileId int) error {
	query := `INSERT INTO link_entity_tile
			(entity_instance_id, tile_id)
			VALUES ($1, $2);`

	_, err := database.DB.Exec(query, entityId, tileId)

	if err != nil {
		log.Println(err)
		return errors.New("Entity could not be linked")
	}

	return nil
}