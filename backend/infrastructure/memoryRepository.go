package infrastructure

import (
	"backend/domain"
	"sync"
)

type MemoryGameRepository struct {
	mu    sync.RWMutex
	games map[string]*domain.Game
}

func NewMemoryGameRepository() *MemoryGameRepository {
	return &MemoryGameRepository{games: map[string]*domain.Game{}}
}

func (r *MemoryGameRepository) Save(game *domain.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.games[game.ID] = cloneGame(game)
	return nil
}

func (r *MemoryGameRepository) FindByID(gameID string) (*domain.Game, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	game, ok := r.games[gameID]
	if !ok {
		return nil, false, nil
	}
	return cloneGame(game), true, nil
}

type MemoryTurnLogRepository struct {
	mu   sync.RWMutex
	logs map[string][]domain.TurnLog
}

func NewMemoryTurnLogRepository() *MemoryTurnLogRepository {
	return &MemoryTurnLogRepository{logs: map[string][]domain.TurnLog{}}
}

func (r *MemoryTurnLogRepository) Append(gameID string, log domain.TurnLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs[gameID] = append(r.logs[gameID], log)
	return nil
}

func (r *MemoryTurnLogRepository) FindByGameID(gameID string) ([]domain.TurnLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	logs := r.logs[gameID]
	copied := make([]domain.TurnLog, len(logs))
	copy(copied, logs)
	return copied, nil
}

type MemoryPredictionRepository struct {
	mu          sync.RWMutex
	predictions map[string]domain.PredictionBoard
}

func NewMemoryPredictionRepository() *MemoryPredictionRepository {
	return &MemoryPredictionRepository{predictions: map[string]domain.PredictionBoard{}}
}

func (r *MemoryPredictionRepository) Save(gameID, playerID string, board domain.PredictionBoard) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.predictions[key(gameID, playerID)] = board
	return nil
}

func (r *MemoryPredictionRepository) Find(gameID, playerID string) (domain.PredictionBoard, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	board, found := r.predictions[key(gameID, playerID)]
	return board, found, nil
}

type MemoryPlayerGamesIndexRepository struct {
	mu    sync.RWMutex
	index map[string]map[string]struct{}
}

func NewMemoryPlayerGamesIndexRepository() *MemoryPlayerGamesIndexRepository {
	return &MemoryPlayerGamesIndexRepository{index: map[string]map[string]struct{}{}}
}

func (r *MemoryPlayerGamesIndexRepository) AddGame(playerID, gameID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.index[playerID]; !ok {
		r.index[playerID] = map[string]struct{}{}
	}
	r.index[playerID][gameID] = struct{}{}
	return nil
}

func (r *MemoryPlayerGamesIndexRepository) ListGames(playerID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	set := r.index[playerID]
	games := make([]string, 0, len(set))
	for gameID := range set {
		games = append(games, gameID)
	}
	return games, nil
}

func key(gameID, playerID string) string {
	return gameID + ":" + playerID
}

func cloneGame(src *domain.Game) *domain.Game {
	cloned := *src
	cloned.Board = domain.NewBoard()
	for id, submarine := range src.Board.Submarines {
		s := *submarine
		cloned.Board.Submarines[id] = &s
	}
	return &cloned
}
