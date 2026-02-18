package shared

import (
	"errors"
)

var (
	ErrInvalidTurn         = errors.New("Error[TurnResult.go]: ターンの値が不正です．")
	ErrInvalidAction       = errors.New("Error[Action.go]: アクションが不正です．")
	ErrInvalidTarget       = errors.New("Error[Target.go]: 座標が不正です．")
	ErrInvalidMoveDistance = errors.New("Error[Move.go]: 移動距離が不正です．")
	ErrOutOfBoard          = errors.New("Error[Position.go]: 場所がボードの外です．")
	ErrPositionIsNil       = errors.New("Error[Position.go]: Positionがnilです．")
	ErrBoardIsNil          = errors.New("Error[Board.go]: Boardがnilです．")
	ErrSubmarineNotFound   = errors.New("Error[Board.go]: Submarineが見つかりませんでした．")
	ErrInvalidDirection    = errors.New("Error[Direction.go]: 方角が不正です．")
)
