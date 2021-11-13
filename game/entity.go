package game

type EntityModel struct {
	Id                                                                             int
	Name                                                                           string
	Precision, Strength, Endurance, Agility, Hitpoints 							   int
}

// Victories, Defeats, Temporary stat values & boolean to check if the character is fighting

type EntityInstance struct {
	Id                                                                                         int
	ModelId                                                                                    int
	Name                                                                           			   string
	Precision, Strength, Endurance, Agility, Hitpoints 										   int
	ChosenActionId                                                           				   int
	IsAlive                                                                        			   bool
}
