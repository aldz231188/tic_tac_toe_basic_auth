package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"t03/internal/api"
	"t03/internal/api/dto"
	"t03/internal/domain"
)

type GameHandler struct {
	GameService domain.GameService
	UserService domain.UserService
}

func NewGameHandler(gameService domain.GameService, userService domain.UserService) *GameHandler {
	return &GameHandler{
		GameService: gameService,
		UserService: userService,
	}
}

func (h *GameHandler) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	playerId, ok := UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Mode == "" {
		req.Mode = "human"
	}

	id, err := h.GameService.NewGame(playerId, req.Mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GameHandler) HandleGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleConnectToGame(w, r)
	case http.MethodPost:
		h.HandleGameMove(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func (h *GameHandler) HandleConnectToGame(w http.ResponseWriter, r *http.Request) {
	playerId, ok := UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/game/")
	game, err := h.GameService.ConnectToGame(id, playerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := api.ToGameResponse(game)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GameHandler) HandleGamesList(w http.ResponseWriter, r *http.Request) {
	playerId, ok := UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	gamesList, err := h.GameService.GetAvailableGames(playerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := api.ToGamesListResponse(gamesList)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func (h *GameHandler) HandlePlayerStats(w http.ResponseWriter, r *http.Request) {
	_, ok := UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/stats/")
	stats, err := h.GameService.GetPlayerStats(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resStats := api.ToStats(stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resStats)

}

func (h *GameHandler) HandleSignUpRequest(w http.ResponseWriter, r *http.Request) {
	var data dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	_, err := h.UserService.Register(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *GameHandler) HandleSignInRequest(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	playerId, err := h.UserService.AuthenticateBasic(auth)
	if err != nil {
		http.Error(w, "authorization error "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"player_id": playerId,
	})
}

func (h *GameHandler) HandleGameMove(w http.ResponseWriter, r *http.Request) {

	playerId, ok := UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/game/")

	var gameReq dto.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&gameReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	playerMove, err := api.ToDomainGame(id, gameReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerMove, err = h.GameService.PlayerVsAi(playerMove, playerId)

	response := api.ToGameResponse(playerMove)
	if err != nil {
		response.Status = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
