package game

/*
Precision			Chance to hit or defend for a melee attack
Strength			Increase melee damage every 10 points & help carrying heavy burden
Toughness			Increase damage resistance every 10 points & help resisting some effects
Agility				Increase dodge chances and help resisting some effects
Hitpoints			Character hitpoints, if it reaches 0, the character may suffer minor to lethal injuries
*/

type PlayerChar struct {
	Id                                                                                         int
	Name, Player                                                                               string
	Precision, Strength, Toughness, Agility, Hitpoints int
	IsCharAlive                                                                                bool
}

// Victories, Defeats, Temporary stat values & boolean to check if the character is fighting

type CharBattle struct {
	Id                                                                                         int
	Name, Player                                                                               string
	Precision, Strength, Toughness, Agility, Hitpoints int
	IsFighting, IsDodging, IsFleeing                                                           bool
}
