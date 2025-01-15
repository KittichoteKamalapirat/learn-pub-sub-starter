package routing

import "time"

// ExchangePerilDirect is the name of the direct exchange.
const ExchangePerilDirect = "peril_direct"

// PauseKey is the routing key for pause messages.
// const PauseKey = "pause"

// PlayingState represents the state of the game.
type PlayingState struct {
	IsPaused bool `json:"isPaused"`
}

type GameLog struct {
	CurrentTime time.Time
	Message     string
	Username    string
}
