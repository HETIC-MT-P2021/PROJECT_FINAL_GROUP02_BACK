package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetCharacterById(id string) (*game.PlayerChar, error) {
	// Get the character from db
	charRow := database.DB.QueryRow("SELECT * FROM characters WHERE id=$1", id)

	selectedCharacter := game.PlayerChar{}
	switch err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Name, &selectedCharacter.Player, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Toughness, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.IsCharAlive); err {
	case sql.ErrNoRows:
		return &selectedCharacter, errors.New("Character not found")
	case nil:
	}

	return &selectedCharacter, nil
}

func GetBattleCharacterById(id string) (*game.CharBattle, error) {
	// Get the character from db
	charRow := database.DB.QueryRow("SELECT * FROM battlechars WHERE id=$1;", id)

	selectedCharacter := game.CharBattle{}
	err := charRow.Scan(&selectedCharacter.Id, &selectedCharacter.Name, &selectedCharacter.Player, &selectedCharacter.Precision, &selectedCharacter.Strength, &selectedCharacter.Toughness, &selectedCharacter.Agility, &selectedCharacter.Hitpoints, &selectedCharacter.IsFighting, &selectedCharacter.IsDodging, &selectedCharacter.IsFleeing)
	if err != nil {
		fmt.Println(err.Error())
		return &selectedCharacter, errors.New("Character not found")
	} else {
		return &selectedCharacter, nil
	}
}
