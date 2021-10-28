package game

/*
Precision			Chance to hit or defend for a melee attack
Strength			Increase melee damage every 10 points & help carrying heavy burden
Endurance			Increase damage resistance every 10 points & help resisting some effects
Agility				Increase dodge chances and help resisting some effects
Hitpoints			Character hitpoints, if it reaches 0, the character may suffer minor to lethal injuries
*/

type Character struct {
	Id                                                                                          int
	Name 																						string
	PlayerId                                                                               	    int64   
	Precision, Strength, Endurance, Agility, Hitpoints 										    int
	PrecisionMax, StrengthMax, EnduranceMax, AgilityMax, HitpointsMax							int
	IsAlive, IsOccupied            bool
	ChosenActionId                                         int                    
}

