package app

import (
	"errors"
	"math"
	"t03/internal/domain"

	"github.com/google/uuid"
)

type GameServiceImpl struct {
	repo domain.GameRepository
}

func NewGameService(repo domain.GameRepository) domain.GameService {
	return &GameServiceImpl{repo: repo}
}

func (svc *GameServiceImpl) NewGame(playerId string, gameMode string) (string, error) {
	gameID := uuid.New()
	pid, err := uuid.Parse(playerId)
	if err != nil {
		return "", err
	}
	var st domain.GameState
	var mode domain.Gametype
	switch gameMode {
	case "human":
		st = domain.StatusWaiting
		mode = domain.PVP
	case "ai":
		st = domain.StatusTurn
		mode = domain.PVE

	}

	game := &domain.Game{
		GameId:     gameID,
		Board:      domain.Board{},
		Mode:       mode,
		Player_X:   pid,
		CurrentPID: pid,
		State:      st,
	}

	err = svc.repo.SaveGame(game)
	if err != nil {
		return "", err
	}
	return gameID.String(), nil
}
func (svc *GameServiceImpl) ConnectToGame(gameId, playerId string) (*domain.Game, error) {

	game, err := svc.repo.GetGame(gameId)
	if err != nil {
		return nil, err
	}
	if game.Player_X.String() != playerId && game.State == domain.StatusWaiting {
		game.Player_O, err = uuid.Parse(playerId)
		if err != nil {
			return nil, err
		}
		game.State = domain.StatusTurn
		err = svc.repo.SaveGame(game)
		if err != nil {
			return nil, err
		}
	}
	return game, nil
}

func (svc *GameServiceImpl) GetAvailableGames(id string) (*domain.GamesList, error) {
	return svc.repo.GetAvailableGames(id)
}
func (svc *GameServiceImpl) GetPlayerStats(playerID string) (*domain.Stats, error) {
	id, err := uuid.Parse(playerID)
	if err != nil {
		return nil, err
	}
	return svc.repo.GetPlayerStats(uuid.MustParse(id.String()))
}

func (svc *GameServiceImpl) PlayerVsAi(playerMove *domain.Game, playerId string) (*domain.Game, error) {

	res, err := svc.PlayerMove(playerMove, playerId)
	if err != nil {
		return res, err
	}

	if res.State == domain.StatusTurn && res.Mode == domain.PVE {
		res, err = svc.aITurn(res, domain.O)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (svc *GameServiceImpl) PlayerMove(game *domain.Game, playerId string) (*domain.Game, error) {
	beforeMove, err := svc.repo.GetGame(game.GameId.String())
	if err != nil {
		return beforeMove, err
	}

	var turn domain.Cell
	switch beforeMove.State {
	case domain.StatusWaiting:
		return beforeMove, errors.New("waiting for another player to connect")
	case domain.StatusTurn:
		if playerId != beforeMove.Player_X.String() && playerId != beforeMove.Player_O.String() {
			return beforeMove, errors.New("not your game")
		}
		if beforeMove.CurrentPID.String() != playerId {
			return beforeMove, errors.New("wait for your turn")
		}
		if playerId == beforeMove.Player_X.String() {
			turn = domain.X
			beforeMove.CurrentPID = beforeMove.Player_O
		} else if playerId == beforeMove.Player_O.String() {
			turn = domain.O
			beforeMove.CurrentPID = beforeMove.Player_X
		}
	case domain.StatusDraw:
		return beforeMove, errors.New("played in a draw")
	case domain.StatusWin:
		return beforeMove, errors.New("player " + beforeMove.WinnerPID.String() + " win")

	}

	err = validateBoard(&beforeMove.Board, &game.Board, turn)
	if err != nil {
		return beforeMove, err
	}

	beforeMove.Board = game.Board

	over, who := checkGameOver(beforeMove.Board)
	if over {
		if who == domain.Empty {
			beforeMove.State = domain.StatusDraw
		} else {
			beforeMove.State = domain.StatusWin
			beforeMove.WinnerPID = uuid.MustParse(playerId)
		}
	}

	svc.repo.SaveGame(beforeMove)

	return beforeMove, nil
}

func (svc *GameServiceImpl) aITurn(game *domain.Game, ai domain.Cell) (*domain.Game, error) {
	beforeMove, err := svc.repo.GetGame(game.GameId.String())
	if err != nil {
		return beforeMove, err
	}

	if beforeMove.State != domain.StatusTurn {
		return beforeMove, errors.New("session finished")
	}

	bestMove := aiMove(game.Board, ai)

	if bestMove[0] == -1 {
		return game, errors.New("no moves left")
	}

	game.Board[bestMove[0]][bestMove[1]] = ai
	game.CurrentPID = beforeMove.Player_X

	over, who := checkGameOver(game.Board)
	if over {
		if who == domain.Empty {
			game.State = domain.StatusDraw
		} else {
			game.State = domain.StatusWin
			game.WinnerPID = uuid.Nil
		}
	}
	svc.repo.SaveGame(game)
	return game, nil

}

func aiMove(board domain.Board, ai domain.Cell) [2]int {
	bestScore := math.MinInt
	bestMove := [2]int{-1, -1}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				board[i][j] = ai
				score := minimax(board, 0, false, ai)
				board[i][j] = domain.Empty

				if score > bestScore {
					bestScore = score
					bestMove = [2]int{i, j}
				}
			}
		}
	}
	return bestMove
}

func validateBoard(oldBoard, newBoard *domain.Board, turn domain.Cell) error {
	moveCount := 0

	for i := range oldBoard {
		for j := range oldBoard[i] {
			oldCell := oldBoard[i][j]
			newCell := newBoard[i][j]

			switch {
			case oldCell == newCell:
				continue

			case oldCell == domain.Empty && newCell == turn:
				moveCount++

			default:
				return errors.New("board is corrupted")
			}
		}
	}

	if moveCount == 0 {
		return errors.New("your turn")
	}
	if moveCount > 1 {
		return errors.New("only one move is allowed at a time")
	}

	return nil
}

func checkGameOver(board domain.Board) (bool, domain.Cell) {
	for i := range 3 {
		if board[i][1] != domain.Empty && (board[i][1] == board[i][0] && board[i][1] == board[i][2]) {
			return true, board[i][1]
		} else if board[1][i] != domain.Empty && (board[1][i] == board[0][i] && board[1][i] == board[2][i]) {
			return true, board[1][i]
		} else {
			continue
		}
	}

	if board[1][1] != domain.Empty && ((board[0][0] == board[1][1] && board[1][1] == board[2][2]) || (board[0][2] == board[1][1] && board[1][1] == board[2][0])) {
		return true, board[1][1]
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				return false, domain.Empty
			}

		}

	}

	return true, domain.Empty
}

func minimax(board domain.Board, depth int, isMaximizing bool, ai domain.Cell) int {
	isOver, winner := checkGameOver(board)
	if isOver {
		if winner == ai {
			return 10 - depth // победа ИИ
		} else if winner != domain.Empty {
			return depth - 10 // победа игрока
		}
		return 0 // ничья
	}

	if isMaximizing {
		best := math.MinInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = ai
					score := minimax(board, depth+1, false, ai)
					board[i][j] = domain.Empty
					best = max(best, score)
				}
			}
		}
		return best
	} else {
		best := math.MaxInt
		player := domain.X
		if ai == domain.X {
			player = domain.O
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = player
					score := minimax(board, depth+1, true, ai)
					board[i][j] = domain.Empty
					best = min(best, score)
				}
			}
		}
		return best
	}
}
