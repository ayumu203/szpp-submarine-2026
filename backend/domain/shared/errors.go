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
	ErrInvalidPlayerID     = errors.New("Error[Player.go]: playerIDが不正です．")
	ErrInvalidPlayerName   = errors.New("Error[Player.go]: playerNameが不正です．")
	ErrSubmarineIdIsEmpty                   = errors.New("Error[Submarine.go]: Submarineのidが空です. ")
	ErrOwnerIdIsEmpty                       = errors.New("Error[Submarine.go]: SubmarineのownerIdが空です. ")
	ErrInvalidHp                            = errors.New("Error[Submarine.go]: Submarineのhpが不正です. ")
	ErrSubmarineAlreadySunk                 = errors.New("Error[Submarine.go]: すでに沈没している潜水艦にダメージを与えようとしています．")
	ErrInvalidDamageAmount                  = errors.New("Error[Submarine.go]: ダメージ量が不正です．")
	ErrActionCommandInvalidParamCombination = errors.New("Error[ActionCommand.go]: actionTypeと他のパラメータ間で矛盾が発生しています．")
	ErrActionCommandIsNil                   = errors.New("Error[ActionCommand.go]: ActionCommandがnilです．")
	ErrInvalidActionType                    = errors.New("Error[ActionType.go]: ActionTypeが不正です．")
)
