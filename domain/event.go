package domain

import "github.com/SteakBarbare/RPGBot/game"

func SetEventWasActivatedInTiles(dungeonTiles []game.DungeonTile, event game.Event, tile game.DungeonTile)[]game.DungeonTile {
	if !event.WasActivated {
		for t := 0; t < len(dungeonTiles); t++ {
			dungeonTile := &dungeonTiles[t]
	
			if dungeonTile.Id == tile.Id {
				var newEvents []game.Event
	
				for e := 0; e < len(dungeonTile.Events); e++ {
					tileEvent := dungeonTile.Events[e]
	
					if tileEvent.Id == event.Id {
						tileEvent.WasActivated = true
					}

					newEvents = append(newEvents, tileEvent)
				}
	
				(*dungeonTile).Events = newEvents
			}
		}
	}

	return dungeonTiles
}