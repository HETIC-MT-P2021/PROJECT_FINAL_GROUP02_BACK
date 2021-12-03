package utils

import "github.com/SteakBarbare/RPGBot/game"

func ResetCharacterBuffs(character game.Character) (game.Character, error) {
	character.Agility = character.AgilityMax
	character.Endurance = character.EnduranceMax
	character.Strength = character.StrengthMax
	character.Precision = character.PrecisionMax

	err := UpdateCharacterFightinhStats(character)

	if err != nil {
		return character, err
	}

	return character, nil
}