package game

import "time"

type GameEvent any

type AnswerAcceptedEvent struct {
	Answer		string
}

type AnswerRejectedEvent struct {
	Answer	string
	Reason	string
}

type ExtraLifeEvent struct{}

type GameQuitEvent struct{}

type GameOverEvent struct{}

type GameWonEvent struct{}

type StrikeEvent struct {
	Strikeout	bool
	Message		string
}

type TimerTickEvent struct {
	TimerId		uint
	Duration	time.Duration
}

type PlayerDamagedEvent struct {
	EventType	string
	Amount		uint
}
