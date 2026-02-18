package domain

import "backend/domain/shared"

type ActionCommand struct {
	playerId   string
	actionType shared.ActionType
	target     *Position
	direction  shared.Direction
	distance   int
}

func NewActionCommand(
	playerId string,
	actionType shared.ActionType,
	target *Position,
	direction shared.Direction,
	distance int,
) (*ActionCommand, error) {
	if actionType == shared.Move && (distance < shared.MinMoveDistance || distance > shared.MaxMoveDistance) {
		return nil, shared.ErrInvalidMoveDistance
	}
	actionCommand := ActionCommand{
		playerId:   playerId,
		actionType: actionType,
		target:     target,
		direction:  direction,
		distance:   distance,
	}
	switch actionCommand.actionType {
	case shared.Attack:
		if target == nil || direction != shared.DirectionUnknown {
			return nil, shared.ErrActionCommandInvalidParamCombination
		}
	case shared.Move:
		if target != nil || direction == shared.DirectionUnknown {
			return nil, shared.ErrActionCommandInvalidParamCombination
		}
	default:
		return nil, shared.ErrInvalidActionType
	}
	return &actionCommand, nil
}

func (actionCommand *ActionCommand) GetPlayerId() (string, error) {
	if actionCommand == nil {
		return "", shared.ErrActionCommandIsNil
	}
	return actionCommand.playerId, nil
}

func (actionCommand *ActionCommand) GetActionType() (shared.ActionType, error) {
	if actionCommand == nil {
		return shared.ActionUnknown, shared.ErrActionCommandIsNil
	}
	return actionCommand.actionType, nil
}

func (actionCommand *ActionCommand) GetTarget() (*Position, error) {
	if actionCommand == nil {
		return nil, shared.ErrActionCommandIsNil
	}
	return actionCommand.target, nil
}

func (actionCommand *ActionCommand) GetDirection() (shared.Direction, error) {
	if actionCommand == nil {
		return shared.DirectionUnknown, shared.ErrActionCommandIsNil
	}
	return actionCommand.direction, nil
}

func (actionCommand *ActionCommand) GetDistance() (int, error) {
	if actionCommand == nil {
		return 0, shared.ErrActionCommandIsNil
	}
	return actionCommand.distance, nil
}
