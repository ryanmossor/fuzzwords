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
	Reason	RejectionReason
}

type ExtraLifeEvent struct {
	Health	int
}

type GameOverEvent struct {
	PossibleAnswer	string
	Stats			PlayerStats
	Turns			[]Turn
}

type GameQuitEvent struct{}

type GameWonEvent struct {
	Stats	PlayerStats
	Turns	[]Turn
}

type StrikeEvent struct {
	Strikeout	bool
	StrikeCount	int
	Prompt		string
}

type PlayerDamagedEvent struct {
	Amount	int
	Health	int
}
