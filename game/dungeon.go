package game

import (
	"time"

	"github.com/SteakBarbare/RPGBot/database"
)

type Dungeon struct {
	Id                                                      int
	CreatedAt												time.Time
	CreatedBy												string
	SelectedCharacterId                          			database.NullInt32
	HasStarted, HasEnded                                    bool
}

type DungeonTile struct {
	Id                                                      int
	DungeonId												int
	X, Y                                                    int
	IsDiscovered, IsExit, IsImpassable                      bool
	Entities                                                []EntityInstance
	Characters                                              []Character
	Events                                                  []EventModel
}

type DungeonInstance struct {
	Dungeon Dungeon
	DungeonTiles []DungeonTile
}