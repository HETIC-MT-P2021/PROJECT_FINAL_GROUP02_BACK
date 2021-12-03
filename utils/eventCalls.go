package utils

import (
	"errors"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func InsertEvent(event game.Event) (int, error) {
	query := `INSERT INTO event
			(event_model_id, category_id, name, description, is_always_active, was_activated)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING event_id;`

	var eventId int

	err := database.DB.QueryRow(query, event.EventModelId, event.CategoryId, event.Name, event.Description,
		event.IsAlwaysActive, event.WasActivated).Scan(&eventId)

	if err != nil {
		log.Println(err)
		return eventId, errors.New("Event could not be created")
	}

	return eventId, nil
}

func LinkEventTile(eventId, tileId int) error {
	query := `INSERT INTO link_event_tile
			(event_id, tile_id)
			VALUES ($1, $2);`

	_, err := database.DB.Exec(query, eventId, tileId)

	if err != nil {
		log.Println(err)
		return errors.New("Event could not be linked")
	}

	return nil
}

func UpdateEventWasActivated(eventId int) (error) {
	query := `UPDATE event 
	 SET was_activated = true
	 WHERE event_id=$1`

	_, err := database.DB.Exec(query, eventId)

	if err != nil {
		log.Println(err)
		return errors.New("Event could not be set activated")
	}

	return nil
}
