package game

type Event struct {
	Id                                  int
	CategoryId							int
	EventModelId						int
	Name                                string
	Description							string
	IsAlwaysActive						bool
	WasActivated						bool			
}

type Category struct {
	Id                                  int
	Name                                string
	Events								[]Event			
}

