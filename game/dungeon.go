package game

import "time"

type Dungeon struct {
	Id                                                      int
	CreatedAt												time.Time
	CreatedBy												string
	SelectedCharacterId                          			int
	HasStarted, HasEnded                                    bool
}

type DungeonTile struct {
	Id                                                      int
	DungeonId												int
	X, Y                                                    int
	IsDiscovered, IsExit, IsImpassable                      bool
	Entities                                                []EntityInstance
	Characters                                              []CharacterInstance
	Events                                                  []EventModel
}

type DungeonInstance struct {
	Dungeon Dungeon
	DungeonTiles []DungeonTile
}