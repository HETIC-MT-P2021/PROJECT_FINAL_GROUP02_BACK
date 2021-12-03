package utils

import (
	"math/rand"

	"github.com/SteakBarbare/RPGBot/consts"
	"github.com/SteakBarbare/RPGBot/game"
)

func createEntityFromModels(entityModel game.EntityModel) game.EntityInstance {
	return game.EntityInstance{
		ModelId: entityModel.Id,
		Name: entityModel.Name,
		Precision: entityModel.Precision,
		Strength: entityModel.Strength,
		Endurance: entityModel.Endurance,
		Agility: entityModel.Agility,
		Hitpoints: entityModel.Hitpoints,
		IsAlive: true,
		ChosenActionId: 0,
	} 
}

func getRandomEntity() game.EntityInstance {
	numberOfDifferentEntity := len(consts.EntityModels)

	entityIndex := rand.Intn(numberOfDifferentEntity)

	return createEntityFromModels(consts.EntityModels[entityIndex])
}
