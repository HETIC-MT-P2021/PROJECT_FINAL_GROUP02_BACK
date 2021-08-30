package game

type DuelPreparation struct {
	Id              string
	SelectingPlayer string
	IsReady         int
	IsOver          bool
	Turn            int
}

type DuelPlayer struct {
	PreparationId  string
	Challenger     string
	Challenged     string
	ChallengerChar string
	ChallengedChar string
}

type DuelBattle struct {
	Id            string
	Challengers   []string
	Characters    []string
	IsOver        bool
	Turn          int
	ActiveFighter string
}
