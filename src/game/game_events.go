package game

import "time"

type GameEvent any

type NewTurnEvent struct {
	Prompt	string
}

type AnswerAcceptedEvent struct {
	Answer	string
}

type AnswerRejectedEvent struct {
	Answer	string
	Reason	string
}

type ExtraLifeEvent struct {
	Health		uint
}

type GameOverEvent struct {
	PossibleAnswer	string
}

type GameWonEvent struct{}

type StrikeEvent struct {
	Strikeout	bool
	StrikeCount	int
	Message		string
}

type TimerTickEvent struct {
	TimerId		uint
	Duration	time.Duration
}

type PlayerDamagedEvent struct {
	Amount		uint
	Health		uint
}
