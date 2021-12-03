package consts

import "github.com/SteakBarbare/RPGBot/game"

var Warrior = game.EntityModel{
	Id: 0,
	Name: "Warrior",
	Precision: 20,
	Strength: 40,
	Endurance: 40,
	Agility: 10,
	Hitpoints: 7,
}

var Archer = game.EntityModel{
	Id: 1,
	Name: "Archer",
	Precision: 50,
	Strength: 30,
	Endurance: 20,
	Agility: 50,
	Hitpoints: 4,
}

var EntityModels = []game.EntityModel{
	Warrior,
	Archer,
}