package domain

import "github.com/SteakBarbare/RPGBot/game"

func ActivateSmallHealEvent(character game.Character) (int) {
	newCharHitpoints := character.Hitpoints + 1

	if newCharHitpoints > character.HitpointsMax {
		return character.Hitpoints
	}

	return newCharHitpoints
}

func ActivateMediumHealEvent(character game.Character) (int) {
	newCharHitpoints := character.Hitpoints + 3

	if newCharHitpoints > character.HitpointsMax {
		return character.Hitpoints
	}

	return newCharHitpoints
}