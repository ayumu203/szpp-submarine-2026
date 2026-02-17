package domain

import (
	"testing"

	shared "backend/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestNewActionCommandSuccess(t *testing.T) {
	testList := []struct {
		name       string
		playerId   string
		actionType shared.ActionType
		targetX    int
		targetY    int
		direction  shared.Direction
		distance   int
	}{
		{
			name:       "[NewActionCommand: 移動]",
			playerId:   "p1",
			actionType: shared.Move,
			targetX:    0,
			targetY:    0,
			direction:  shared.North,
			distance:   3,
		},
		{
			name:       "[NewActionCommand: 攻撃]",
			playerId:   "p2",
			actionType: shared.Attack,
			targetX:    2,
			targetY:    3,
			direction:  shared.DirectionUnknown,
			distance:   0,
		},
	}

	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			target, err := NewPosition(tl.targetX, tl.targetY)
			if tl.actionType == shared.Attack {
				assert.NoError(t, err)
			}

			cmd, err := NewActionCommand(tl.playerId, tl.actionType, target, tl.direction, tl.distance)
			assert.NoError(t, err)

			assert.Equal(t, tl.playerId, cmd.playerId)
			assert.Equal(t, tl.actionType, cmd.actionType)
			assert.Equal(t, target, cmd.target)
			assert.Equal(t, tl.direction, cmd.direction)
			assert.Equal(t, tl.distance, cmd.distance)
		})
	}
}
func TestNewActionCommandFail(t *testing.T) {
	testList := []struct {
		name          string
		playerId      string
		actionType    shared.ActionType
		targetX       int
		targetY       int
		direction     shared.Direction
		distance      int
		expectedError error
	}{
		{
			name:          "[NewActionCommand: 移動なのにdirectionがUnknown]",
			playerId:      "p3",
			actionType:    shared.Move,
			targetX:       0,
			targetY:       0,
			direction:     shared.DirectionUnknown,
			distance:      3,
			expectedError: shared.ErrActionCommandInvalidParamCombination,
		},
		{
			name:          "[NewActionCommand: 移動なのにtargetがnilでない]",
			playerId:      "p4",
			actionType:    shared.Move,
			targetX:       1,
			targetY:       2,
			direction:     shared.DirectionUnknown,
			distance:      3,
			expectedError: shared.ErrActionCommandInvalidParamCombination,
		},
		{
			name:          "[NewActionCommand: 攻撃なのにdirectionがUnknownでない]",
			playerId:      "p5",
			actionType:    shared.Attack,
			targetX:       2,
			targetY:       3,
			direction:     shared.North,
			distance:      0,
			expectedError: shared.ErrActionCommandInvalidParamCombination,
		},
		{
			name:          "[NewActionCommand: 攻撃なのにtargetがnil]",
			playerId:      "p6",
			actionType:    shared.Attack,
			targetX:       0,
			targetY:       0,
			direction:     shared.North,
			distance:      0,
			expectedError: shared.ErrActionCommandInvalidParamCombination,
		},
	}

	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			target, _ := NewPosition(tl.targetX, tl.targetY)

			_, err := NewActionCommand(tl.playerId, tl.actionType, target, tl.direction, tl.distance)
			assert.ErrorIs(t, err, tl.expectedError)
		})
	}
}
