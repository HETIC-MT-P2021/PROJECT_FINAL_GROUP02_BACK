package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetCharacterById(id string) (*game.CharacterModel, error) {
	// Get the character from db
	charRow := database.DB.QueryRow("SELECT * FROM character_model WHERE character_model_id=$1", id)

	selectedCharacter := game.CharacterModel{}
	switch err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Name, &selectedCharacter.PlayerId, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.IsAlive, &selectedCharacter.IsOccupied); err {
	case sql.ErrNoRows:
		return &selectedCharacter, errors.New("Character not found")
	case nil:
	}

	return &selectedCharacter, nil
}

func GetCharacterInstanceById(id string) (*game.CharacterInstance, error) {
	// Get the character from db
	charRow := database.DB.QueryRow("SELECT * FROM character_instance WHERE character_instance_id=$1;", id)

	selectedCharacter := game.CharacterInstance{}
	err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.ChosenActionId)
	if err != nil {
		fmt.Println(err.Error())
		return &selectedCharacter, errors.New("Character not found")
	} else {
		return &selectedCharacter, nil
	}
}

func GetCharacterInstanceByModelId(id int) (*game.CharacterInstance, error) {
	// Get the character from db
	charRow := database.DB.QueryRow("SELECT character_instance_id, precision, strength, endurance, agility, hitpoints FROM character_instance WHERE character_model_id=$1;", id)

	selectedCharacter := game.CharacterInstance{}
	err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints)
	if err != nil {
		fmt.Println(err.Error())
		return &selectedCharacter, errors.New("Character not found")
	} else {
		return &selectedCharacter, nil
	}
}

func GetPlayerCharacterModel(playerId int64, characterName string) (game.CharacterModel, error){
	query := `SELECT character_model_id, precision, strength, endurance, agility, hitpoints, is_occupied, is_alive  
	 FROM character_model
	 WHERE player_id=$1 
	 AND name=$2 
	 AND is_occupied=false
	 AND is_alive=true;`

	charRow := database.DB.QueryRow(query, playerId, characterName)

	var selectedCharacter game.CharacterModel

	err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.IsOccupied, &selectedCharacter.IsAlive);

	selectedCharacter.Name = characterName
	selectedCharacter.PlayerId = playerId

	if err !=nil{
		return selectedCharacter, err
	}

	return selectedCharacter, nil
}

func CreateCharacterInstance(characterModel game.CharacterModel)(int, error){
	query := `INSERT INTO character_instance
			(character_model_id, precision, strength, endurance, agility, hitpoints)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING character_instance_id`

	var characterInstanceId int

	err := database.DB.QueryRow(query, characterModel.Id, characterModel.Precision, characterModel.Strength, characterModel.Endurance, characterModel.Agility, characterModel.Hitpoints).Scan(&characterInstanceId)

	if err != nil {
		log.Println(err)
		return characterInstanceId, errors.New("Link character tile could not be created")
	}

	return characterInstanceId, nil
}