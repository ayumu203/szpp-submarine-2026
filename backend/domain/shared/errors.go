package shared

import (
	"errors"
)

var (
	ErrInvalidTurn                          = errors.New("Error[TurnResult.go]: ターンの値が不正です．")
	ErrInvalidAction                        = errors.New("Error[Action.go]: アクションが不正です．")
	ErrInvalidTarget                        = errors.New("Error[Target.go]: 座標が不正です．")
	ErrInvalidMoveDistance                  = errors.New("Error[Move.go]: 移動距離が不正です．")
	ErrOutOfBoard                           = errors.New("Error[Position.go]: 場所がボードの外です．")
	ErrPositionIsNil                        = errors.New("Error[Position.go]: Positionがnilです．")
	ErrSubmarineIdIsEmpty                   = errors.New("Error[Submarine.go]: Submarineのidが空です. ")
	ErrOwnerIdIsEmpty                       = errors.New("Error[Submarine.go]: SubmarineのownerIdが空です. ")
	ErrInvalidHp                            = errors.New("Error[Submarine.go]: Submarineのhpが不正です. ")
	ErrActionCommandInvalidParamCombination = errors.New("Error[ActionCommand.go]: actionTypeと他のパラメータ間で矛盾が発生しています．")
	ErrActionCommandIsNil                   = errors.New("Error[ActionCommand.go]: ActionCommandがnilです．")
	ErrInvalidActionType                    = errors.New("Error[ActionType.go]: ActionTypeが不正です．")
	ErrBoardIsNil                           = errors.New("Error[Board.go]: Boardがnilです．")
	ErrSubmarineNotFound                    = errors.New("Error[Board.go]: Submarineが見つかりませんでした．")
	ErrInvalidDirection                     = errors.New("Error[Direction.go]: 方角が不正です．")
	ErrAllySubmarineAlreadyExists           = errors.New("Error[Board.go]: 指定した場所にはすでに味方の潜水艦がいます．")
	ErrSunkSubmarineAlreadyExists           = errors.New("Error[Board.go]: 指定した場所にはすでに沈没した潜水艦がいます．")
	ErrSubmarineIDDuplicated                = errors.New("Error[Board.go]: SubmarineIDが重複しています．")
)
