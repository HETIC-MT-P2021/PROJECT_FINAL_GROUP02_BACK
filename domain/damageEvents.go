package domain

import (
	"github.com/SteakBarbare/RPGBot/game"
)

func ActivateOneTimeDamageTrapEvent(character game.Character) (int) {
	newCharHitpoints := character.Hitpoints - 3

	return newCharHitpoints
}

func ActivatePermanentDamageTrapEvent(character game.Character) (int) {
	newCharHitpoints := character.Hitpoints - 1

	return newCharHitpoints
}