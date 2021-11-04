package utils

import (
	"errors"
	"math/rand"

	"github.com/SteakBarbare/RPGBot/consts"
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
	numberOfDifferentEvents := len(consts.EntityModels)

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

			tile.Events = append(tile.Events, event)
		}
	}
	
	return tile, nil
}

func HandleNewTileEvents(newDungeonTiles []game.DungeonTile, s *discordgo.Session, m *discordgo.MessageCreate, authorId int64) error {
	playerTile, _, isFound := findPlayerPosInDungeon(authorId, newDungeonTiles)

	if isFound {
		if playerTile.IsExit {
			s.ChannelMessageSend(m.ChannelID, "This is the exit, you can leave now or later !")

			return nil
		}
		if len(playerTile.Entities) == 0 && len(playerTile.Events) == 0 {
			s.ChannelMessageSend(m.ChannelID, "The room is empty !")

			return nil
		}
		
		if len(playerTile.Events) > 0 {
			stringToDisplay := "Something is happening ! \n\n"

			for _, event := range playerTile.Events {
				stringToDisplay += " -" + event.Description + "\n\n"
			}

			s.ChannelMessageSend(m.ChannelID, stringToDisplay)
		} 

		if len(playerTile.Entities) > 0 {
			stringToDisplay := "Someone is lurking is the shadows"

			for i, entity := range playerTile.Entities {
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

		return nil
	}
	
	return errors.New("Character not found after position update")
}