package consts

import "github.com/SteakBarbare/RPGBot/game"

var SmallMeal = game.Event{
	EventType: 0,
	Name: "Small meal",
	Description: "Lucky find ! You recover 2 hp",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var PermanentTrap = game.Event{
	EventType: 1,
	Name: "Trap",
	Description: "It's a trap ! You lose 2hp",
	IsAlwaysActive:	true,
	WasActivated: false,
}


var EventModels = []game.Event{
	SmallMeal,
	PermanentTrap,
}

