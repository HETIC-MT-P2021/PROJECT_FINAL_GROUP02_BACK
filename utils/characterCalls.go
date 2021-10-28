package utils

import (
	"database/sql"
	"errors"
	"log"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetCharacterById(id int) (*game.Character, error) {
	query := `SELECT 
	 character_id,
	 precision, strength, endurance, agility, hitpoints,
	 precision_max, strength_max, endurance_max, agility_max, hitpoints_max,
	 is_occupied, is_alive 
	 FROM character WHERE character_id=$1;`

	// Get the character from db
	charRow := database.DB.QueryRow(query, id)

	selectedCharacter := game.Character{}

	switch err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Name, &selectedCharacter.PlayerId, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.PrecisionMax, &selectedCharacter.StrengthMax, &selectedCharacter.EnduranceMax, &selectedCharacter.AgilityMax, &selectedCharacter.HitpointsMax, &selectedCharacter.IsAlive, &selectedCharacter.IsOccupied); err {
		case sql.ErrNoRows:
			return &selectedCharacter, errors.New("Character not found")
		case nil:
	}

	return &selectedCharacter, nil
}

func GetPlayerCharacterByName(playerId int64, characterName string) (game.Character, error){
	query := `SELECT 
	 character_id,
	 precision, strength, endurance, agility, hitpoints,
	 precision_max, strength_max, endurance_max, agility_max, hitpoints_max,
	 is_occupied, is_alive
	 FROM character 
	 WHERE player_id=$1 
	 AND name=$2 
	 AND is_occupied=false
	 AND is_alive=true;`

	charRow := database.DB.QueryRow(query, playerId, characterName)

	var selectedCharacter game.Character

	err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Endurance, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.PrecisionMax, &selectedCharacter.StrengthMax, &selectedCharacter.EnduranceMax, &selectedCharacter.AgilityMax, &selectedCharacter.HitpointsMax, &selectedCharacter.IsAlive, &selectedCharacter.IsOccupied);

	if err != nil{
		return selectedCharacter, err
	}

	selectedCharacter.Name = characterName
	selectedCharacter.PlayerId = playerId

	return selectedCharacter, nil
}

func CreateCharacter(character game.Character)(int, error){
	query := `INSERT INTO character
			 (player_id, name,
		 	 precision, strength, endurance, agility, hitpoints, 
			 precision_max, strength_max, endurance_max, agility_max, hitpoints_max)
			 VALUES (
				 $1, $2, 
				 $3, $4, $5, $6, $7, 
				 $8, $9, $10, $11, $12)
			 RETURNING character_id`

	var characterId int

	err := database.DB.QueryRow(query, character.PlayerId, character.Name, 
		character.Precision, character.Strength, character.Endurance, character.Agility, character.Hitpoints,
		character.Precision, character.Strength, character.Endurance, character.Agility, character.Hitpoints).Scan(&characterId)

	if err != nil {
		log.Println(err)
		return characterId, errors.New("Character could not be created")
	}

	return characterId, nil
}

func FindCharNameWithId(characterId int) (string, error) {
	query := `SELECT name FROM character 
	 WHERE character_id = 
		(SELECT character_id FROM character
		WHERE character_id=$1);`  

	charRow := database.DB.QueryRow(query, characterId)

	var characterName string

	err := charRow.Scan(&characterName);

	if err != nil {
		log.Println(err)
		return characterName, errors.New("Character name counldn't be found")
	}

	return characterName, nil
}