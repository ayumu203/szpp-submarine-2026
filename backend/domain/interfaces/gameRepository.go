package interfaces

import "context"

type GameRepository interface {
	// Save persists a game identified by gameID.
	Save(ctx context.Context, gameID string) error
	// FindByID retrieves a game by its ID.
	FindByID(ctx context.Context, gameID string) (interface{}, error)
}