package domain

import "github.com/SteakBarbare/RPGBot/game"

func IsCharacterDead(newCharHitpoints int, character game.Character)(game.Character){
	if newCharHitpoints <= 0 {
		character.Hitpoints = 0
		character.IsAlive = false
	} else {
		character.Hitpoints = newCharHitpoints
	}

	return character
}