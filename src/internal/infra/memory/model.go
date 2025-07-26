package memory

import (
	"github.com/google/uuid"
)

type GameEntity struct {
	GameId        uuid.UUID `db:"id"`
	Board         string    `db:"board_state"`
	Mode          int       `db:"mode"`
	Player_X      uuid.UUID `db:"player_x"`
	Player_O      uuid.UUID `db:"player_o"`
	State         int       `db:"state"`
	CurrentPID    uuid.UUID `db:"turn"`
	CurrentSimbol int       `db:"current_simbol"`
	WinnerPID     uuid.UUID `db:"winner"`
}

type UserEntity struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"user_login"`
	Password string    `db:"user_status"`
}
