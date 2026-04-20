package game

type GameEvent any

type NewTurnEvent struct {
	Prompt	string
}

type AnswerAcceptedEvent struct {
	Answer			string
	NewLettersUsed 	[]rune
}

type AnswerRejectedEvent struct {
	Answer	string
	Reason	string
}

type ExtraLifeEvent struct {
	Health	int
}

type GameOverEvent struct {
	PossibleAnswer	string
	Stats			PlayerStats
}

type GameQuitEvent struct{}

type GameWonEvent struct {
	Stats	PlayerStats
}

type StrikeEvent struct {
	Strikeout	bool
	StrikeCount	int
	Message		string
}

type PlayerDamagedEvent struct {
	Amount	int
	Health	int
}
