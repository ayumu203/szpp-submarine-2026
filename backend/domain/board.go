package domain

import (
	shared "backend/domain/shared"
)

type Board struct {
	submarines map[shared.SubmarineId]*Submarine
}

func NewBoard() *Board {
	board := Board{
		submarines: make(map[shared.SubmarineId]*Submarine),
	}
	return &board
}

func (board *Board) PlaceSubmarine(playerId shared.PlayerId, submarineId shared.SubmarineId, position *Position) error {
	if board == nil {
		return shared.ErrBoardIsNil
	}
	submarine, err := NewSubmarine(submarineId, playerId, position, shared.DefaultHP)
	if err != nil {
		return err
	}
	allySubmarine, err := board.GetAllySubmarineAt(playerId, position)
	if err != nil {
		return err
	}
	if allySubmarine != nil {
		return shared.ErrAllySubmarineAlreadyExists
	}

	if board.submarines[submarineId] != nil {
		return shared.ErrSubmarineIdDuplicated
	}

	board.submarines[submarineId] = submarine
	return nil
}

func (board *Board) MoveSubmarine(playerId shared.PlayerId, submarineId shared.SubmarineId, direction shared.Direction, distance int) error {
	if board == nil {
		return shared.ErrBoardIsNil
	}
	submarine, ok := board.submarines[submarineId]
	if !ok {
		return shared.ErrSubmarineNotFound
	}
	if submarine.GetOwnerId() != playerId {
		return shared.ErrPlayerIdNotMatch
	}
	if submarine.IsSunk() {
		return shared.ErrMoveSunkSubmarine
	}
	submarinePosition := submarine.GetPosition()
	moveToX, moveToY, err := submarinePosition.GetPosition()
	if err != nil {
		return err
	}
	if distance >= shared.MinMoveDistance && distance <= shared.MaxMoveDistance {
		// ok
	} else {
		return shared.ErrInvalidMoveDistance
	}
	// 一度moveToX, moveToYの値は不正になっても良いことにし，NewPosition()で範囲内かどうかを判定する
	moveOnlyX := moveToX
	moveOnlyY := moveToY
	switch direction {
	case shared.North:
		moveToY -= distance
		if distance == 2 {
			moveOnlyY -= 1
		}
	case shared.South:
		moveToY += distance
		if distance == 2 {
			moveOnlyY += 1
		}
	case shared.East:
		moveToX += distance
		if distance == 2 {
			moveOnlyX += 1
		}
	case shared.West:
		moveToX -= distance
		if distance == 2 {
			moveOnlyX -= 1
		}
	default:
		return shared.ErrInvalidDirection
	}
	moveToPosition, err := NewPosition(moveToX, moveToY)
	if err != nil {
		return err
	}
	moveOnlyPosition, err := NewPosition(moveOnlyX, moveOnlyY)
	if err != nil {
		return err
	}
	allySubmarine, err := board.GetAllySubmarineAt(playerId, moveToPosition)
	if err != nil {
		return err
	}
	if allySubmarine != nil {
		return shared.ErrAllySubmarineAlreadyExists
	}
	allySubmarineT, err := board.GetAllySubmarineAt(playerId, moveOnlyPosition)
	if err != nil {
		return err
	}
	if allySubmarineT != nil && allySubmarineT.IsSunk() {
		return shared.ErrMovedOverSunkSubmarine
	}
	opponentSubmarineT, err := board.GetOpponentSubmarineAt(playerId, moveOnlyPosition)
	if err != nil {
		return err
	}
	if opponentSubmarineT != nil && opponentSubmarineT.IsSunk() {
		return shared.ErrMovedOverSunkSubmarine
	}
	opponentSubmarine, err := board.GetOpponentSubmarineAt(playerId, moveToPosition)
	if err != nil {
		return err
	}
	if opponentSubmarine != nil {
		isSunk := opponentSubmarine.IsSunk()
		if isSunk {
			return shared.ErrMovedOverSunkSubmarine
		}
	}
	err = submarine.MoveTo(moveToPosition)
	if err != nil {
		return err
	}
	return nil
}

func (board *Board) FindTargets(attackerId shared.PlayerId, center *Position) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	opponentSubmarines, err := board.GetOpponentSubmarines(attackerId)
	if err != nil {
		return nil, err
	}
	for _, submarine := range opponentSubmarines {
		submarinePosition := submarine.GetPosition()
		submarineNeighbors, err := submarinePosition.Neighbors8()
		if err != nil {
			return nil, err
		}
		for _, submarineNeighbor := range submarineNeighbors {
			eq, err := submarineNeighbor.isEqual(center)
			if err != nil {
				return nil, err
			}
			if eq {
				submarines = append(submarines, submarine)
				break
			}
		}
	}
	return submarines, nil
}

func (board *Board) IsOccupied(position *Position) (bool, error) {
	if board == nil {
		return false, shared.ErrBoardIsNil
	}
	for _, submarine := range board.submarines {
		submarinePosition := submarine.GetPosition()
		eq, err := submarinePosition.isEqual(position)
		if err != nil {
			return false, err
		}
		if eq {
			return true, nil
		}
	}
	return false, nil
}

func (board *Board) GetAllySubmarineAt(playerId shared.PlayerId, position *Position) (*Submarine, error) {
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	allySubmarines, err := board.GetAllySubmarines(playerId)
	if err != nil {
		return nil, err
	}
	for _, submarine := range allySubmarines {
		submarinePosition := submarine.GetPosition()
		eq, err := submarinePosition.isEqual(position)
		if err != nil {
			return nil, err
		}
		if eq {
			return submarine, nil
		}
	}
	return nil, nil
}

func (board *Board) GetOpponentSubmarineAt(playerId shared.PlayerId, position *Position) (*Submarine, error) {
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	opponentSubmarines, err := board.GetOpponentSubmarines(playerId)
	if err != nil {
		return nil, err
	}
	for _, submarine := range opponentSubmarines {
		submarinePosition := submarine.GetPosition()
		eq, err := submarinePosition.isEqual(position)
		if err != nil {
			return nil, err
		}
		if eq {
			return submarine, nil
		}
	}
	return nil, nil
}

func (board *Board) GetAllySubmarines(playerId shared.PlayerId) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	for _, submarine := range board.submarines {
		submarineOwnerId := submarine.GetOwnerId()
		if playerId == submarineOwnerId {
			submarines = append(submarines, submarine)
		}
	}
	return submarines, nil
}
func (board *Board) GetOpponentSubmarines(playerId shared.PlayerId) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	for _, submarine := range board.submarines {
		submarineOwnerId := submarine.GetOwnerId()
		if playerId != submarineOwnerId {
			submarines = append(submarines, submarine)
		}
	}
	return submarines, nil
}
