package game

import (
	"time"

	"github.com/SteakBarbare/RPGBot/database"
)

type Dungeon struct {
	Id                                                      int
	CreatedAt												time.Time
	CreatedBy												int64
	SelectedCharacterId                          			database.NullInt32
	HasStarted, HasEnded, IsPaused                          bool
}

type DungeonTile struct {
	Id                                                      int
	DungeonId												int
	X, Y                                                    int
	IsDiscovered, IsExit, IsImpassable                      bool
	Entities                                                []EntityInstance
	Characters                                              []Character
	Events                                                  []Event
}

type DungeonInstance struct {
	Dungeon Dungeon
	DungeonTiles []DungeonTile
}