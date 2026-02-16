package application

import (
	"backend/domain"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ExecuteActionRequest struct {
	GameID     string
	PlayerID   string
	ActionType domain.ActionType
	Target     *domain.Position
	Direction  domain.Direction
	Distance   int
}

type ExecuteActionResponse struct {
	GameID       string
	Turn         int
	AttackReport *domain.ReportType
	MoveReport   *domain.MoveReportType
	ErrorCode    *domain.ErrorCode
	NextPlayerID string
	WinnerID     string
	Status       domain.GameStatus
}

type GameState struct {
	GameID          string
	Turn            int
	Status          domain.GameStatus
	CurrentPlayerID string
	OpponentID      string
	Board           domain.Board
	PredictionBoard domain.PredictionBoard
	Logs            []domain.TurnLog
}

type GameService struct {
	gameRepo        domain.GameRepository
	turnLogRepo     domain.TurnLogRepository
	predictionRepo  domain.PredictionRepository
	playerGamesRepo domain.PlayerGamesIndexRepository
	cpuDecision     *CpuDecisionService
	cpuAnalysis     *CpuAnalysisService
	randSource      *rand.Rand
}

func NewGameService(
	gameRepo domain.GameRepository,
	turnLogRepo domain.TurnLogRepository,
	predictionRepo domain.PredictionRepository,
	playerGamesRepo domain.PlayerGamesIndexRepository,
	cpuDecision *CpuDecisionService,
	cpuAnalysis *CpuAnalysisService,
) *GameService {
	return &GameService{
		gameRepo:        gameRepo,
		turnLogRepo:     turnLogRepo,
		predictionRepo:  predictionRepo,
		playerGamesRepo: playerGamesRepo,
		cpuDecision:     cpuDecision,
		cpuAnalysis:     cpuAnalysis,
		randSource:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *GameService) InitializeGame(playerAID, playerBID string) (*domain.Game, error) {
	if playerAID == "" || playerBID == "" {
		return nil, errors.New("player ids are required")
	}

	game := &domain.Game{
		ID:        fmt.Sprintf("game-%d", time.Now().UnixNano()),
		Status:    domain.GameStatusWaiting,
		PlayerAID: playerAID,
		PlayerBID: playerBID,
		Board:     domain.NewBoard(),
	}

	if err := s.placeInitialSubmarines(game, playerAID); err != nil {
		return nil, err
	}
	if err := s.placeInitialSubmarines(game, playerBID); err != nil {
		return nil, err
	}
	game.Start()

	if err := s.gameRepo.Save(game); err != nil {
		return nil, err
	}
	if err := s.playerGamesRepo.AddGame(playerAID, game.ID); err != nil {
		return nil, err
	}
	if err := s.playerGamesRepo.AddGame(playerBID, game.ID); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *GameService) ExecuteAction(req ExecuteActionRequest) (ExecuteActionResponse, error) {
	game, found, err := s.gameRepo.FindByID(req.GameID)
	if err != nil {
		return ExecuteActionResponse{}, err
	}
	if !found {
		return ExecuteActionResponse{}, errors.New("game not found")
	}

	humanTurn := game.Turn
	result := game.Apply(domain.ActionCommand{
		PlayerID:  req.PlayerID,
		Type:      req.ActionType,
		Target:    req.Target,
		Direction: req.Direction,
		Distance:  req.Distance,
	})

	log := buildTurnLog(humanTurn, req, result)
	if err := s.cpuAnalysis.RecordTurn(game.ID, log); err != nil {
		return ExecuteActionResponse{}, err
	}

	if result.ErrorCode == nil {
		if err := s.gameRepo.Save(game); err != nil {
			return ExecuteActionResponse{}, err
		}
		if req.ActionType == domain.ActionTypeAttack {
			if err := s.cpuAnalysis.UpdatePrediction(game.ID, req.PlayerID, req.Target, result.AttackReport); err != nil {
				return ExecuteActionResponse{}, err
			}
		}
	}

	response := ExecuteActionResponse{
		GameID:       game.ID,
		Turn:         game.Turn,
		AttackReport: result.AttackReport,
		MoveReport:   result.MoveReport,
		ErrorCode:    result.ErrorCode,
		NextPlayerID: result.NextPlayerID,
		WinnerID:     result.WinnerID,
		Status:       result.Status,
	}

	if result.ErrorCode == nil && result.Status != domain.GameStatusFinished && isCPUPlayer(game.CurrentPlayerID) {
		cpuResult, cpuReq, cpuErr := s.executeCpuTurn(game)
		if cpuErr != nil {
			return ExecuteActionResponse{}, cpuErr
		}
		if cpuReq.ActionType == domain.ActionTypeAttack {
			if err := s.cpuAnalysis.UpdatePrediction(game.ID, cpuReq.PlayerID, cpuReq.Target, cpuResult.AttackReport); err != nil {
				return ExecuteActionResponse{}, err
			}
		}
		// Response keeps the human action report while metadata reflects latest game state.
		response.Turn = game.Turn
		response.NextPlayerID = game.CurrentPlayerID
		response.WinnerID = game.WinnerID
		response.Status = game.Status
	}

	return response, nil
}

func (s *GameService) executeCpuTurn(game *domain.Game) (domain.TurnResult, ExecuteActionRequest, error) {
	req := s.cpuDecision.DecideAction(game)
	cpuTurn := game.Turn
	result := game.Apply(domain.ActionCommand{
		PlayerID:  req.PlayerID,
		Type:      req.ActionType,
		Target:    req.Target,
		Direction: req.Direction,
		Distance:  req.Distance,
	})
	log := buildTurnLog(cpuTurn, req, result)
	if err := s.cpuAnalysis.RecordTurn(game.ID, log); err != nil {
		return domain.TurnResult{}, ExecuteActionRequest{}, err
	}
	if result.ErrorCode == nil {
		if err := s.gameRepo.Save(game); err != nil {
			return domain.TurnResult{}, ExecuteActionRequest{}, err
		}
	}
	return result, req, nil
}

func (s *GameService) GetGameState(gameID, viewerPlayerID string) (GameState, error) {
	game, found, err := s.gameRepo.FindByID(gameID)
	if err != nil {
		return GameState{}, err
	}
	if !found {
		return GameState{}, errors.New("game not found")
	}
	logs, err := s.turnLogRepo.FindByGameID(gameID)
	if err != nil {
		return GameState{}, err
	}
	predictionBoard, foundPrediction, err := s.predictionRepo.Find(gameID, viewerPlayerID)
	if err != nil {
		return GameState{}, err
	}
	if !foundPrediction {
		predictionBoard = domain.NewPredictionBoard()
	}

	return GameState{
		GameID:          game.ID,
		Turn:            game.Turn,
		Status:          game.Status,
		CurrentPlayerID: game.CurrentPlayerID,
		OpponentID:      game.OpponentOf(viewerPlayerID),
		Board:           game.Board,
		PredictionBoard: predictionBoard,
		Logs:            logs,
	}, nil
}

func (s *GameService) placeInitialSubmarines(game *domain.Game, playerID string) error {
	placed := 0
	for placed < domain.SubmarinesPerPlayer {
		position := domain.Position{
			X: s.randSource.Intn(domain.BoardSize) + 1,
			Y: s.randSource.Intn(domain.BoardSize) + 1,
		}
		submarineID := fmt.Sprintf("%s-sub-%d", playerID, placed+1)
		if game.Board.PlaceSubmarine(playerID, submarineID, position) {
			placed++
		}
	}
	return nil
}

func isCPUPlayer(playerID string) bool {
	id := strings.ToLower(playerID)
	return id == "cpu" || strings.HasPrefix(id, "cpu:") || strings.HasPrefix(id, "bot")
}

func buildTurnLog(turn int, req ExecuteActionRequest, result domain.TurnResult) domain.TurnLog {
	createdAt := time.Now().UTC().Format(time.RFC3339)
	log := domain.TurnLog{
		Turn:         turn,
		PlayerID:     req.PlayerID,
		ActionType:   req.ActionType,
		Target:       req.Target,
		AttackReport: result.AttackReport,
		MoveReport:   result.MoveReport,
		ErrorCode:    result.ErrorCode,
		CreatedAt:    createdAt,
	}
	if req.ActionType == domain.ActionTypeMove {
		dir := req.Direction
		distance := req.Distance
		log.Direction = &dir
		log.Distance = &distance
	}
	return log
}

type CpuDecisionService struct {
	randSource *rand.Rand
}

func NewCpuDecisionService() *CpuDecisionService {
	return &CpuDecisionService{randSource: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (s *CpuDecisionService) DecideAction(game *domain.Game) ExecuteActionRequest {
	playerID := game.CurrentPlayerID
	for _, submarine := range game.Board.AliveSubmarinesOf(playerID) {
		for _, target := range submarine.Position.Neighbors8() {
			targetSubmarine := game.Board.GetSubmarineAt(target)
			if targetSubmarine != nil && targetSubmarine.OwnerID == playerID {
				continue
			}
			return ExecuteActionRequest{
				GameID:     game.ID,
				PlayerID:   playerID,
				ActionType: domain.ActionTypeAttack,
				Target:     &domain.Position{X: target.X, Y: target.Y},
			}
		}
	}

	directions := []domain.Direction{domain.DirectionNorth, domain.DirectionSouth, domain.DirectionEast, domain.DirectionWest}
	return ExecuteActionRequest{
		GameID:     game.ID,
		PlayerID:   playerID,
		ActionType: domain.ActionTypeMove,
		Direction:  directions[s.randSource.Intn(len(directions))],
		Distance:   s.randSource.Intn(2) + 1,
	}
}

type CpuAnalysisService struct {
	turnLogRepo    domain.TurnLogRepository
	predictionRepo domain.PredictionRepository
}

func NewCpuAnalysisService(turnLogRepo domain.TurnLogRepository, predictionRepo domain.PredictionRepository) *CpuAnalysisService {
	return &CpuAnalysisService{turnLogRepo: turnLogRepo, predictionRepo: predictionRepo}
}

func (s *CpuAnalysisService) RecordTurn(gameID string, log domain.TurnLog) error {
	return s.turnLogRepo.Append(gameID, log)
}

func (s *CpuAnalysisService) UpdatePrediction(gameID, playerID string, target *domain.Position, report *domain.ReportType) error {
	if target == nil || report == nil {
		return nil
	}
	board, found, err := s.predictionRepo.Find(gameID, playerID)
	if err != nil {
		return err
	}
	if !found {
		board = domain.NewPredictionBoard()
	}
	x := target.X - 1
	y := target.Y - 1
	switch *report {
	case domain.ReportHit, domain.ReportHitAndSunk:
		board.ScoreGrid[y][x] += 3
	case domain.ReportWaveHigh:
		board.ScoreGrid[y][x] += 1
	default:
		board.ScoreGrid[y][x] -= 1
	}
	if board.ScoreGrid[y][x] < 0 {
		board.ScoreGrid[y][x] = 0
	}
	board.PossibleEnemyCount[y][x] = board.ScoreGrid[y][x]
	board.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.predictionRepo.Save(gameID, playerID, board)
}
