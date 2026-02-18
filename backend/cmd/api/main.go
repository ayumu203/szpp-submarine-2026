package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== DTO Definitions ====================

// Position represents a coordinate on the board
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// InitializeGameRequest ゲーム初期化リクエスト
type InitializeGameRequest struct {
	PlayerAId          string     `json:"playerAId"`
	PlayerBId          string     `json:"playerBId"`
	SubmarinePositions []Position `json:"submarinePositions"`
}

// InitializeGameResponse ゲーム初期化レスポンス
type InitializeGameResponse struct {
	GameId          string  `json:"gameId"`
	Status          string  `json:"status"` // waiting, inProgress, finished
	Turn            int     `json:"turn"`
	CurrentPlayerId string  `json:"currentPlayerId"`
	Error           *string `json:"error,omitempty"` // invalidPosition
}

// ExecuteActionRequest アクション実行リクエスト
type ExecuteActionRequest struct {
	GameId     string    `json:"gameId"`
	PlayerId   string    `json:"playerId"`
	ActionType string    `json:"actionType"` // attack, move
	Target     *Position `json:"target,omitempty"`
	Direction  *string   `json:"direction,omitempty"` // north, south, east, west
	Distance   *int      `json:"distance,omitempty"`  // 1 or 2
}

// ExecuteActionResponse アクション実行レスポンス
type ExecuteActionResponse struct {
	GameId       string  `json:"gameId"`
	Turn         int     `json:"turn"`
	AttackReport *string `json:"attackReport,omitempty"` // invalidAttack, miss, hit, hitAndSunk, waveHigh
	MoveReport   *string `json:"moveReport,omitempty"`   // moveSuccess, moveFailed
	ErrorCode    *string `json:"errorCode,omitempty"`    // invalidTurn, invalidAction, invalidTarget, invalidMoveDistance, outOfBoard
	NextPlayerId string  `json:"nextPlayerId"`
	WinnerId     *string `json:"winnerId,omitempty"`
	Status       string  `json:"status"` // inProgress, finished
}

// GetGameStateRequest ゲーム状態取得リクエスト
type GetGameStateRequest struct {
	GameId         string `json:"gameId"`
	ViewerPlayerId string `json:"viewerPlayerId"`
}

// BoardViewDto ボード表示DTO
type BoardViewDto struct {
	Cells      [][]string              `json:"cells"` // 5x5 grid
	Submarines map[string]SubmarineDto `json:"submarines"`
}

// SubmarineDto 潜水艦情報DTO
type SubmarineDto struct {
	OwnerId string `json:"ownerId"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	HP      int    `json:"hp"`
	Sunk    bool   `json:"sunk"`
}

// PredictionBoardDto 予測ボードDTO
type PredictionBoardDto struct {
	ScoreGrid          [][]int `json:"scoreGrid"`          // 5x5
	PossibleEnemyCount [][]int `json:"possibleEnemyCount"` // 5x5
	UpdatedAt          string  `json:"updatedAt"`
}

// TurnLogDto ターンログDTO
type TurnLogDto struct {
	Turn         int       `json:"turn"`
	PlayerId     string    `json:"playerId"`
	ActionType   string    `json:"actionType"` // attack, move
	Target       *Position `json:"target,omitempty"`
	Direction    *string   `json:"direction,omitempty"`
	Distance     *int      `json:"distance,omitempty"`
	AttackReport *string   `json:"attackReport,omitempty"`
	MoveReport   *string   `json:"moveReport,omitempty"`
	ErrorCode    *string   `json:"errorCode,omitempty"`
	CreatedAt    string    `json:"createdAt"`
}

// GetGameStateResponse ゲーム状態取得レスポンス
type GetGameStateResponse struct {
	GameId          string             `json:"gameId"`
	Turn            int                `json:"turn"`
	Status          string             `json:"status"` // inProgress, finished
	CurrentPlayerId string             `json:"currentPlayerId"`
	OpponentId      string             `json:"opponentId"`
	AllyBoard       BoardViewDto       `json:"allyBoard"`
	EnemyBoard      BoardViewDto       `json:"enemyBoard"`
	PredictionBoard PredictionBoardDto `json:"predictionBoard"`
	Logs            []TurnLogDto       `json:"logs"`
}

// ==================== Handlers ====================

// POST /initialize
func handleInitializeGame(c *gin.Context) {
	var req InitializeGameRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	// TODO: バリデーション、ゲーム初期化ロジック実装
	// 仮の値で応答
	gameId := "game_" + time.Now().Format("20060102150405")
	resp := InitializeGameResponse{
		GameId:          gameId,
		Status:          "inProgress",
		Turn:            1,
		CurrentPlayerId: req.PlayerAId,
	}

	c.JSON(http.StatusOK, resp)
}

// POST /action
func handleExecuteAction(c *gin.Context) {
	var req ExecuteActionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	// TODO: アクション検証、実行ロジック実装、CPU手の処理
	// 仮の値で応答
	moveSuccess := "moveSuccess"
	resp := ExecuteActionResponse{
		GameId:       req.GameId,
		Turn:         1,
		MoveReport:   &moveSuccess,
		NextPlayerId: "playerB",
		Status:       "inProgress",
	}

	c.JSON(http.StatusOK, resp)
}

// POST /state
func handleGetGameState(c *gin.Context) {
	var req GetGameStateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	// TODO: ゲーム状態取得ロジック実装
	// 仮の5x5ボード作成
	emptyBoard := make([][]string, 5)
	for i := range emptyBoard {
		emptyBoard[i] = make([]string, 5)
		for j := range emptyBoard[i] {
			emptyBoard[i][j] = "."
		}
	}

	scoreGrid := make([][]int, 5)
	for i := range scoreGrid {
		scoreGrid[i] = make([]int, 5)
	}

	possibleEnemyCount := make([][]int, 5)
	for i := range possibleEnemyCount {
		possibleEnemyCount[i] = make([]int, 5)
	}

	resp := GetGameStateResponse{
		GameId:          req.GameId,
		Turn:            1,
		Status:          "inProgress",
		CurrentPlayerId: req.ViewerPlayerId,
		OpponentId:      "opponentId",
		AllyBoard: BoardViewDto{
			Cells:      emptyBoard,
			Submarines: make(map[string]SubmarineDto),
		},
		EnemyBoard: BoardViewDto{
			Cells:      emptyBoard,
			Submarines: make(map[string]SubmarineDto),
		},
		PredictionBoard: PredictionBoardDto{
			ScoreGrid:          scoreGrid,
			PossibleEnemyCount: possibleEnemyCount,
			UpdatedAt:          time.Now().Format(time.RFC3339),
		},
		Logs: []TurnLogDto{},
	}

	c.JSON(http.StatusOK, resp)
}

func main() {
	e := gin.Default()

	// Routes
	e.POST("/initialize", handleInitializeGame)
	e.POST("/action", handleExecuteAction)
	e.GET("/state", handleGetGameState)

	e.Run(":8080")
}
