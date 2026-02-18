package interfaces

import (
	"backend/domain/shared"
	"context"
)

type GameRepository interface {
	// Save persists a game identified by gameID.
	Save(ctx context.Context, gameID shared.GameId) error
	// FindByID retrieves a game by its ID.
	FindByID(ctx context.Context, gameID shared.GameId) (interface{}, error)
}
