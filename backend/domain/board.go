package domain

import (
	"backend/domain/shared"
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

func (board *Board) PlaceSubmarine(playerID shared.PlayerID, position *Position) error {
	submarine, err := NewSubmarine(playerID, position)
	if err != nil {
		return err
	}
	submarineID, err := submarine.GetID()
	if err != nil {
		return err
	}
	allySubmarine, err := board.GetAllySubmarineAt(playerID, position)
	if err != nil {
		return err
	}
	if allySubmarine != nil {
		return shared.ErrAllySubmarineAlreadyExists
	}

	board.submarines[submarineID] = submarine
	return nil
}

func (board *Board) MoveSubmarine(playerId shared.PlayerID, submarineId shared.SubmarineId, direction shared.Direction, distance int) (bool, error) {
	if board == nil {
		return false, shared.ErrBoardIsNil
	}
	submarine, ok := board.submarines[submarineId]
	if !ok {
		return false, shared.ErrSubmarineNotFound
	}
	submarinePosition, err := submarine.GetPosition()
	if err != nil {
		return false, err
	}
	moveToX, moveToY, err := submarinePosition.GetPosition()
	if err != nil {
		return false, err
	}
	if distance >= shared.MinMoveDistance && distance <= shared.MaxMoveDistance {
		// ok
	} else {
		return false, shared.ErrInvalidMoveDistance
	}
	// 一度moveToX, moveToYの値は不正になっても良いことにし，NewPosition()で範囲内かどうかを判定する
	switch direction {
	case shared.North:
		moveToY -= distance
	case shared.South:
		moveToY += distance
	case shared.East:
		moveToX += distance
	case shared.West:
		moveToX -= distance
	default:
		return false, shared.ErrInvalidDirection
	}
	moveToPosition, err := NewPosition(moveToX, moveToY)
	if err == shared.ErrOutOfBoard {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	allySubmarine, err := board.GetAllySubmarineAt(playerId, moveToPosition)
	if err != nil {
		return err
	}
	if allySubmarine != nil {
		return false, shared.ErrAllySubmarineAlreadyExists
	}
	opponentSubmarine, err := board.GetOpponentSubmarineAt(playerId, moveToPosition)
	if err != nil {
		return false, err
	}
	if opponentSubmarine != nil {
		isSunk, err := opponentSubmarine.IsSunk()
		if err != nil {
			return false, err
		}
		if isSunk {
			return false, shared.ErrSunkSubmarineAlreadyExists
		}
	}
	err = submarine.MoveTo(moveToPosition)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (board *Board) FindTargets(attackerID shared.PlayerID, center *Position) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	opponentSubmarines, err := board.GetOpponentSubmarines(attackerID)
	if err != nil {
		return nil, err
	}
	for _, submarine := range opponentSubmarines {
		submarinePosition, err := submarine.GetPosition()
		if err != nil {
			return nil, err
		}
		submarineNeighbors, err := submarinePosition.Neighbors8()
		if err != nil {
			return nil, err
		}
		for _, submarineNeighbor := range submarineNeighbors {
			if *submarineNeighbor == *center {
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
		submarinePosition, err := submarine.GetPosition()
		if err != nil {
			return false, err
		}
		if *submarinePosition == *position {
			return true, nil
		}
	}
	return false, nil
}

func (board *Board) GetAllySubmarineAt(playerId shared.PlayerID, position *Position) (*Submarine, error) {
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	allySubmarines, err := board.GetAllySubmarines(playerId)
	if err != nil {
		return nil, err
	}
	for _, submarine := range allySubmarines {
		submarinePosition, err := submarine.GetPosition()
		if err != nil {
			return false, err
		}
		if *submarinePosition != *position {
			continue
		}
	}
	return nil, nil
}

func (board *Board) GetOpponentSubmarineAt(playerId shared.PlayerID, position *Position) (*Submarine, error) {
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	opponentSubmarines, err := board.GetOpponentSubmarines(playerId)
	if err != nil {
		return nil, err
	}
	for _, submarine := range opponentSubmarines {
		submarinePosition, err := submarine.GetPosition()
		if err != nil {
			return false, err
		}
		if *submarinePosition != *position {
			continue
		}
	}
	return nil, nil
}

func (board *Board) GetAllySubmarines(playerId shared.PlayerID) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	for _, submarine := range board.submarines {
		submarineOwnerID, err := submarine.GetOnwerID()
		if err != nil {
			return nil, err
		}
		if playerId == submarineOwnerID {
			return submarine, nil
		}
	}
	return nil, nil
}
func (board *Board) GetOpponentSubmarines(playerId shared.PlayerID) ([]*Submarine, error) {
	submarines := make([]*Submarine, 0, 4)
	if board == nil {
		return nil, shared.ErrBoardIsNil
	}
	for _, submarine := range board.submarines {
		submarineOwnerID, err := submarine.GetOnwerID()
		if err != nil {
			return nil, err
		}
		if playerId != submarineOwnerID {
			return submarine, nil
		}
	}
	return nil, nil
}
