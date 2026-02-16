package presentation

import (
	"backend/application"
	"backend/domain"
	"encoding/json"
	"net/http"
)

type GameHandler struct {
	service *application.GameService
}

func NewGameHandler(service *application.GameService) *GameHandler {
	return &GameHandler{service: service}
}

type InitializeGameRequest struct {
	PlayerAID string `json:"playerAId"`
	PlayerBID string `json:"playerBId"`
}

type InitializeGameResponse struct {
	GameID          string            `json:"gameId"`
	Status          domain.GameStatus `json:"status"`
	Turn            int               `json:"turn"`
	CurrentPlayerID string            `json:"currentPlayerId"`
}

type ExecuteActionRequest struct {
	GameID     string            `json:"gameId"`
	PlayerID   string            `json:"playerId"`
	ActionType domain.ActionType `json:"actionType"`
	Target     *domain.Position  `json:"target,omitempty"`
	Direction  domain.Direction  `json:"direction,omitempty"`
	Distance   int               `json:"distance,omitempty"`
}

type ExecuteActionResponse struct {
	GameID       string                 `json:"gameId"`
	Turn         int                    `json:"turn"`
	AttackReport *domain.ReportType     `json:"attackReport,omitempty"`
	MoveReport   *domain.MoveReportType `json:"moveReport,omitempty"`
	ErrorCode    *domain.ErrorCode      `json:"errorCode,omitempty"`
	NextPlayerID string                 `json:"nextPlayerId"`
	WinnerID     string                 `json:"winnerId,omitempty"`
	Status       domain.GameStatus      `json:"status"`
}

type GetGameStateResponse struct {
	GameID          string             `json:"gameId"`
	Turn            int                `json:"turn"`
	Status          domain.GameStatus  `json:"status"`
	CurrentPlayerID string             `json:"currentPlayerId"`
	OpponentID      string             `json:"opponentId"`
	Board           BoardViewDTO       `json:"board"`
	PredictionBoard PredictionBoardDTO `json:"predictionBoard"`
	Logs            []domain.TurnLog   `json:"logs"`
}

type BoardViewDTO struct {
	Cells      [domain.BoardSize][domain.BoardSize]*string `json:"cells"`
	Submarines map[string]SubmarineDTO                     `json:"submarines"`
}

type SubmarineDTO struct {
	OwnerID string `json:"ownerId"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	HP      int    `json:"hp"`
	Sunk    bool   `json:"sunk"`
}

type PredictionBoardDTO struct {
	ScoreGrid          [domain.BoardSize][domain.BoardSize]int `json:"scoreGrid"`
	PossibleEnemyCount [domain.BoardSize][domain.BoardSize]int `json:"possibleEnemyCount"`
	UpdatedAt          string                                  `json:"updatedAt"`
}

func (h *GameHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/initialize", h.handleInitialize)
	mux.HandleFunc("/action", h.handleAction)
	mux.HandleFunc("/state", h.handleState)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}

func (h *GameHandler) handleInitialize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req InitializeGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	game, err := h.service.InitializeGame(req.PlayerAID, req.PlayerBID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, InitializeGameResponse{
		GameID:          game.ID,
		Status:          game.Status,
		Turn:            game.Turn,
		CurrentPlayerID: game.CurrentPlayerID,
	})
}

func (h *GameHandler) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ExecuteActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	res, err := h.service.ExecuteAction(application.ExecuteActionRequest{
		GameID:     req.GameID,
		PlayerID:   req.PlayerID,
		ActionType: req.ActionType,
		Target:     req.Target,
		Direction:  req.Direction,
		Distance:   req.Distance,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, ExecuteActionResponse{
		GameID:       res.GameID,
		Turn:         res.Turn,
		AttackReport: res.AttackReport,
		MoveReport:   res.MoveReport,
		ErrorCode:    res.ErrorCode,
		NextPlayerID: res.NextPlayerID,
		WinnerID:     res.WinnerID,
		Status:       res.Status,
	})
}

func (h *GameHandler) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	gameID := r.URL.Query().Get("gameId")
	viewerPlayerID := r.URL.Query().Get("viewerPlayerId")
	if gameID == "" || viewerPlayerID == "" {
		http.Error(w, "gameId and viewerPlayerId are required", http.StatusBadRequest)
		return
	}
	state, err := h.service.GetGameState(gameID, viewerPlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	submarines := map[string]SubmarineDTO{}
	cells := [domain.BoardSize][domain.BoardSize]*string{}
	for id, submarine := range state.Board.Submarines {
		submarines[id] = SubmarineDTO{
			OwnerID: submarine.OwnerID,
			X:       submarine.Position.X,
			Y:       submarine.Position.Y,
			HP:      submarine.HP,
			Sunk:    submarine.IsSunk(),
		}
		label := id
		cells[submarine.Position.Y-1][submarine.Position.X-1] = &label
	}

	writeJSON(w, http.StatusOK, GetGameStateResponse{
		GameID:          state.GameID,
		Turn:            state.Turn,
		Status:          state.Status,
		CurrentPlayerID: state.CurrentPlayerID,
		OpponentID:      state.OpponentID,
		Board: BoardViewDTO{
			Cells:      cells,
			Submarines: submarines,
		},
		PredictionBoard: PredictionBoardDTO{
			ScoreGrid:          state.PredictionBoard.ScoreGrid,
			PossibleEnemyCount: state.PredictionBoard.PossibleEnemyCount,
			UpdatedAt:          state.PredictionBoard.UpdatedAt,
		},
		Logs: state.Logs,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
