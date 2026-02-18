package domain

import (
	shared "backend/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	t.Run("[NewBoard]", func(t *testing.T) {
		board := NewBoard()
		assert.NotNil(t, board)
	})
}

func TestPlaceSubmarineSuccess(t *testing.T) {
	testList := []struct {
		name     string
		playerId shared.PlayerID
		x        int
		y        int
	}{}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, position)
			assert.NoError(t, err)
		})
	}
}

func TestPlaceSubmarineDuplicate(t *testing.T) {
	testList := []struct {
		name     string
		playerId shared.PlayerID
		x        int
		y        int
	}{}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, position)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, position)
			assert.ErrorIs(t, shared.ErrAllySubmarineAlreadyExists, err)
		})
	}
}

func TestMoveSubmarineSuccess(t *testing.T) {
	testList := []struct {
		name          string
		submarineX    int
		submarineY    int
		submarineID   shared.SubmarineId
		moveDirection shared.Direction
		moveDistance  int
	}{}
}

func TestMoveSubmarineInvalidDistance(t *testing.T) {
	board := NewBoard()
	pos, _ := NewPosition(5, 5)
	playerID := shared.PlayerID(1)
	board.PlaceSubmarine(playerID, pos)

	submarine := board.submarines[shared.SubmarineId(1)]
	_, err := board.MoveSubmarine(playerID, submarine.GetID(), shared.North, 100)

	if err != shared.ErrInvalidMoveDistance {
		t.Errorf("expected ErrInvalidMoveDistance, got %v", err)
	}
}

func TestMoveOutOfBoard(t *testing.T) {
	board := NewBoard()
	pos, _ := NewPosition(1, 1)
	playerID := shared.PlayerID(1)
	board.PlaceSubmarine(playerID, pos)

	submarine := board.submarines[shared.SubmarineId(1)]
	moved, err := board.MoveSubmarine(playerID, submarine.GetID(), shared.West, 5)

	if moved || err != nil {
		t.Errorf("expected (false, nil), got (%v, %v)", moved, err)
	}
}
