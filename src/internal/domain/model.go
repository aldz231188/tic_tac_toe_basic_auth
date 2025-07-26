package domain

import "github.com/google/uuid"

type GameState int
type Gametype int

const (
	PVP Gametype = iota
	PVE
)

const (
	StatusWaiting GameState = iota
	StatusTurn
	StatusDraw
	StatusWin
	// StatusBoardCorrupted
	// St
)

type Game struct {
	GameId     uuid.UUID
	Mode       Gametype
	Board      Board
	Player_X   uuid.UUID
	Player_O   uuid.UUID
	State      GameState
	CurrentPID uuid.UUID
	WinnerPID  uuid.UUID
}

type Cell int

const (
	Empty Cell = iota
	X
	O
)

type Board [3][3]Cell

type GamesList struct {
	Games uuid.UUIDs
}

type Stats struct {
	TotalGames int
	Wins       int
	Losses     int
	Draws      int
	WinRatePct float64
}
