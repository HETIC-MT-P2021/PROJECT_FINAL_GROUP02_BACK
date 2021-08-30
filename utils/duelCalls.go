package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
)

func GetDuelPreparation() (*game.DuelPreparation, error) {
	// Get the selecting player and duel ID
	duelPreparation := database.DB.QueryRow(fmt.Sprintln("SELECT * FROM duelPreparation WHERE isReady=", 0))

	currentDuel := game.DuelPreparation{}
	switch err := duelPreparation.Scan(&currentDuel.Id, &currentDuel.SelectingPlayer, &currentDuel.IsReady); err {
	case sql.ErrNoRows:
		return &currentDuel, errors.New("No duel in preparation found")
	case nil:
		fmt.Println(duelPreparation)
	}
	return &currentDuel, nil
}

func GetActiveDuel() (*game.DuelPreparation, error) {
	// Get the selecting player and duel ID
	duelPreparation := database.DB.QueryRow("SELECT * FROM duelPreparation WHERE isReady=$1 AND isOver=$2;", 1, false)

	currentDuel := game.DuelPreparation{}
	switch err := duelPreparation.Scan(&currentDuel.Id, &currentDuel.SelectingPlayer, &currentDuel.IsReady, &currentDuel.IsOver, &currentDuel.Turn); err {
	case sql.ErrNoRows:
		return &currentDuel, errors.New("No active duel found")
	case nil:
		fmt.Println(duelPreparation)
	}
	return &currentDuel, nil
}

func GetDuelPlayers(id string) (*game.DuelPlayer, error) {
	// Select the opponents and their respective character informations
	duelPlayerRow := database.DB.QueryRow(fmt.Sprintln("SELECT challenger, challenged, challengerChar, challengedChar FROM duelPlayers WHERE preparationId=", id))
	duelPlayers := game.DuelPlayer{}

	switch err := duelPlayerRow.Scan(&duelPlayers.Challenger, &duelPlayers.Challenged, &duelPlayers.ChallengerChar, &duelPlayers.ChallengedChar); err {
	case sql.ErrNoRows:
		return &duelPlayers, errors.New("Duel not found")
	case nil:
	}
	return &duelPlayers, nil
}
