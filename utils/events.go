package utils

import (
	"errors"
	"log"
	"math/rand"
	"strconv"

	"github.com/SteakBarbare/RPGBot/consts"
	"github.com/SteakBarbare/RPGBot/domain"
	"github.com/SteakBarbare/RPGBot/game"
	"github.com/bwmarrin/discordgo"
)

func findNumberOfEvents() int {
	randomEventCoef := rand.Intn(100)

	if randomEventCoef <= 30 {
		return 0
	} else if randomEventCoef <= 70 {
		return 1
	} else if randomEventCoef <= 90 {
		return 2
	} 
	
	return 3
}

func getRandomEvent() game.Event {
	numberOfDifferentEvents := len(consts.EventModels)

	eventIndex := rand.Intn(numberOfDifferentEvents)

	return consts.EventModels[eventIndex]
}

func GenerateEventsAndEntity(numberOfEvents int)(game.EntityInstance, []game.Event, bool) {
	var entity game.EntityInstance
	var events []game.Event

	wasEntityGenerated := false

	for i := 0; i < numberOfEvents; i++ {
		isEventOrEntityCoef := rand.Intn(100)

		if isEventOrEntityCoef <= 70 || wasEntityGenerated {
			event := getRandomEvent()

			events = append(events, event)
		} else {
			entity = getRandomEntity()

			wasEntityGenerated = true
		}
	}

	return entity, events, wasEntityGenerated
}


func HandleEventsOnDiscover(tile game.DungeonTile) (game.DungeonTile, error ){
	var entity game.EntityInstance
	var events []game.Event

	wasEntityGenerated := false

	numberOfEvents := findNumberOfEvents()

	if numberOfEvents > 0 {
		entity, events, wasEntityGenerated = GenerateEventsAndEntity(numberOfEvents)
	}

	if wasEntityGenerated {
		entityId, err := InsertEntityInstance(entity)

		if err != nil {
			return tile, err
		}

		err = LinkEntityTile(entityId, tile.Id)

		if err != nil  {
			return tile, err
		}

		entity.Id = entityId

		tile.Entities = append(tile.Entities, entity)
	}

	if len(events) > 0 {
		for _, event := range events {
			eventId, err := InsertEvent(event)

			if err != nil {
				return tile, err
			}

			err = LinkEventTile(eventId, tile.Id)

			if err != nil  {
				return tile, err
			}

			event.Id = eventId

			tile.Events = append(tile.Events, event)
		}
	}
	
	return tile, nil
}

func UpdateCharacterInDungeonTiles(character game.Character, characterTileId int, dungeonTiles []game.DungeonTile)([]game.DungeonTile){
	for t := 0; t < len(dungeonTiles); t+=1 {
		tile := &dungeonTiles[t]

		if tile.Id == characterTileId {
			for c := 0; c < len(tile.Characters); c+=1 {
				if tile.Characters[c].Id == character.Id {
					(*tile).Characters[c] = character

					break
				}
			}

			break
		}
	}

	return dungeonTiles
}

func ActivateBuffEvent(event game.Event, character game.Character, characterTile game.DungeonTile, tiles []game.DungeonTile)(game.Character, error){
	var updatedCharacter game.Character
	var err error

	updatedCharacter = character

	switch event.EventModelId {
	case consts.BuffPrecisionTriggerEvent.EventModelId:
		closeByTiles := GetCloseTiles(characterTile, tiles)

		err = InsertEventAndLinkToTiles(consts.BuffPrecisionEffectEvent, closeByTiles)
	
		if err !=nil {
			return updatedCharacter, err
		}

		updatedCharacter.Precision += 15

		err = UpdateCharPrecision(updatedCharacter)

		if err != nil {
			return updatedCharacter, err
		}

		break
	case consts.BuffPrecisionEffectEvent.EventModelId:
		updatedCharacter.Precision += 15

		err = UpdateCharPrecision(updatedCharacter)

		if err !=nil {
			return updatedCharacter, err
		}
		
		break
	default:
		errMsg := "Event does not exist:" + strconv.Itoa(event.EventModelId)
		log.Println(errMsg)

		return character, errors.New(errMsg)
	}


	return updatedCharacter, err
}

func ActivateDebuffEvent(event game.Event, character game.Character, characterTile game.DungeonTile, tiles []game.DungeonTile)(game.Character, error){
	var updatedCharacter game.Character
	var err error

	updatedCharacter = character

	switch event.EventModelId {
	case consts.DebuffStrengthTriggerEvent.EventModelId:
		closeByTiles := GetCloseTiles(characterTile, tiles)

		err = InsertEventAndLinkToTiles(consts.DebuffStrengthEffectEvent, closeByTiles)
	
		if err !=nil {
			return updatedCharacter, err
		}

		updatedCharacter.Strength -= 15

		err = UpdateCharPrecision(updatedCharacter)

		if err != nil {
			return updatedCharacter, err
		}

		break
	case consts.DebuffStrengthEffectEvent.EventModelId:
		updatedCharacter.Strength -= 15

		err = UpdateCharPrecision(updatedCharacter)

		if err !=nil {
			return updatedCharacter, err
		}
		
		break
	default:
		errMsg := "Event does not exist:" + strconv.Itoa(event.EventModelId)
		log.Println(errMsg)

		return character, errors.New(errMsg)
	}


	return updatedCharacter, err
}

func ActivateHealEvent(event game.Event, character game.Character)(game.Character, error){
	var err error
	var newCharHitpoints int

	switch event.EventModelId {
	case consts.SmallHealEvent.EventModelId:
		newCharHitpoints = domain.ActivateSmallHealEvent(character)

		break
	case consts.MediumHealEvent.EventModelId:
		newCharHitpoints = domain.ActivateMediumHealEvent(character)

		break
	default:
		errMsg := "Event does not exist:" + strconv.Itoa(event.EventModelId)
		log.Println(errMsg)

		return character, errors.New(errMsg)
	}

	character.Hitpoints = newCharHitpoints

	err = UpdateCharHitpoints(character)

	return character, err
}

func ActivateDamageEvent(event game.Event, character game.Character)(game.Character, error){
	var updatedCharacter game.Character
	var err error
	var newCharHitpoints int

	switch event.EventModelId {
	case consts.OneTimeDamageTrapEvent.EventModelId:
		newCharHitpoints = domain.ActivateOneTimeDamageTrapEvent(character)

		break
	case consts.PermanentDamageTrapEvent.EventModelId:
		newCharHitpoints = domain.ActivatePermanentDamageTrapEvent(character)

		break
	default:
		errMsg := "Event does not exist:" + strconv.Itoa(event.EventModelId)
		log.Println(errMsg)

		return character, errors.New(errMsg)
	}

	updatedCharacter = domain.IsCharacterDead(newCharHitpoints, character)
	
	err = UpdateCharHitpoints(updatedCharacter)

	if !updatedCharacter.IsAlive {
		err = UpdateCharIsDead(character.Id)
	}

	return updatedCharacter, err
}

func ActivateTileEvent(characterTile game.DungeonTile, tiles []game.DungeonTile, event game.Event, character game.Character) (game.Character, []game.DungeonTile, string, error){
	var updatedCharacter game.Character
	var err error
	var eventActivationDescription string

	log.Println("Activate tile event")
	log.Println(event)
	log.Println(character)

	updatedDungeonTiles := tiles

	switch event.CategoryId {
	case consts.DamageCategoryId:
		updatedCharacter, err = ActivateDamageEvent(event, character)

		if !updatedCharacter.IsAlive {
			eventActivationDescription = "You have no HP left, you died ! \n"
		} else {
			eventActivationDescription = "You have " + strconv.Itoa(updatedCharacter.Hitpoints) + " Hitpoints left\n"
		}

		break
	case consts.HealCategoryId:
		updatedCharacter, err = ActivateHealEvent(event, character)

		if updatedCharacter.Hitpoints == character.Hitpoints {
			eventActivationDescription = "You are already full health"
		} else {
			eventActivationDescription = "You have " + strconv.Itoa(updatedCharacter.Hitpoints) + " Hitpoints left\n"
		}

		break
	case consts.BuffCategoryId:
		updatedCharacter, err = ActivateBuffEvent(event, character, characterTile, updatedDungeonTiles)

		eventActivationDescription = "You have " + strconv.Itoa(updatedCharacter.Precision) + " Precision in this room\n"

		break
	case consts.DebuffCategoryId:
		updatedCharacter, err = ActivateDebuffEvent(event, character, characterTile, updatedDungeonTiles)
		
		eventActivationDescription = "You have " + strconv.Itoa(updatedCharacter.Strength) + " Strength in this room\n"

		break
	default:
		errMsg := "Category does not exist:" + strconv.Itoa(event.CategoryId)
		log.Println(errMsg)

		return updatedCharacter,updatedDungeonTiles, "",  errors.New(errMsg)
	}

	if err != nil {
		return updatedCharacter, updatedDungeonTiles, eventActivationDescription, err
	}

	updatedDungeonTiles = UpdateCharacterInDungeonTiles(character, characterTile.Id, updatedDungeonTiles)

	if !event.WasActivated {
		log.Println("Activate")
		log.Println(event.Id)
		updatedDungeonTiles = domain.SetEventWasActivatedInTiles(updatedDungeonTiles, event, characterTile)

		err := UpdateEventWasActivated(event.Id)

		if err != nil {
			return updatedCharacter, updatedDungeonTiles, eventActivationDescription, err
		}
	}

	return updatedCharacter, updatedDungeonTiles, eventActivationDescription, nil
}

func ActivateTileEvents(characterTile game.DungeonTile, tiles []game.DungeonTile, character game.Character)(string, game.Character, []game.DungeonTile, error){
	updatedCharacter := character
	var err error
	var updatedDungeonTiles []game.DungeonTile
	
	eventDisplayString :=  "Something is happening ! \n\n"

	for _, event := range characterTile.Events {
		if !event.WasActivated || event.IsAlwaysActive{
			var eventActivationDescription string

			eventDisplayString += " -" + event.Description + "\n\n"

			updatedCharacter, updatedDungeonTiles, eventActivationDescription, err = ActivateTileEvent(characterTile, tiles, event, character)

			if err != nil {
				eventDisplayString += "An error happend: " + err.Error() + "\n\n"

				break
			}

			eventDisplayString += eventActivationDescription + "\n"

			if !updatedCharacter.IsAlive{
				break
			}
		}
	}
	
	return eventDisplayString, updatedCharacter, updatedDungeonTiles, err
}

func HandleNewTileEvents(dungeonTiles []game.DungeonTile, s *discordgo.Session, m *discordgo.MessageCreate, authorId int64) (error, game.Character) {
	var updatedCharacter game.Character

	characterTile, character, isFound := findPlayerPosInDungeon(authorId, dungeonTiles)

	if isFound {
		updatedCharacter, err := ResetCharacterBuffs(character)
		
		if err != nil {
			return err, updatedCharacter
		}

		if characterTile.IsExit {
			s.ChannelMessageSend(m.ChannelID, "This is the exit, you can leave now or later !")

			return nil, updatedCharacter
		}

		if len(characterTile.Entities) == 0 && len(characterTile.Events) == 0 {
			s.ChannelMessageSend(m.ChannelID, "The room is empty !")

			return nil, updatedCharacter
		}
		
		if len(characterTile.Events) > 0 {
			eventDisplayString, updatedCharacter, _, err := ActivateTileEvents(characterTile, dungeonTiles, updatedCharacter)

			if err != nil {
				log.Println(err.Error())
				return err, updatedCharacter
			}
			
			s.ChannelMessageSend(m.ChannelID, eventDisplayString)

			if !updatedCharacter.IsAlive {
				return nil, updatedCharacter
			}
		} 
		
		if len(characterTile.Entities) > 0 {
			stringToDisplay := "Someone is lurking is the shadows"

			for i, entity := range characterTile.Entities {
				log.Println(entity)

				if i == 0 {
					stringToDisplay += ", it's a "
				} else {
					stringToDisplay += " and a "
				}

				stringToDisplay += entity.Name
			}

			stringToDisplay += "\n\n Prepare for Battle !"

			s.ChannelMessageSend(m.ChannelID, stringToDisplay)
		}

		return nil, updatedCharacter
	}
	
	return errors.New("Character not found after position update"), updatedCharacter
}