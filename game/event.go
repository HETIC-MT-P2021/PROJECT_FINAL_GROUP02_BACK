package game

type Event struct {
	Id                                  int
	EventType							int
	Name                                string
	Description							string
	IsAlwaysActive						bool
	WasActivated						bool			
}

