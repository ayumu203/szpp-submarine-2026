package domain

import (
	shared "backend/domain/shared"
	"errors"
)

type Position struct {
	x int
	y int
}

func NewPosition(x int, y int) (*Position, error) {
	position := Position{
		x: x,
		y: y,
	}
	withinResult, err := position.withinBoard()
	if err != nil {
		return nil, err
	}
	if !withinResult {
		return nil, shared.ErrOutOfBoard
	}
	return &position, nil
}

func (position *Position) withinBoard() (bool, error) {
	if position == nil {
		return false, shared.ErrPositionIsNil
	}
	return position.x >= shared.MinPosition && position.x <= shared.MaxPosition && position.y >= shared.MinPosition && position.y <= shared.MaxPosition, nil
}

func (position *Position) Neighbors8() ([]*Position, error) {
	if position == nil {
		return nil, shared.ErrPositionIsNil
	}
	positions := make([]*Position, 0, 8)
	delta := []int{-1, 0, 1}
	x, y, err := position.GetPosition()
	if err != nil {
		return nil, err
	}
	for _, dx := range delta {
		for _, dy := range delta {
			positionNeighbor, err := NewPosition(x+dx, y+dy)
			if errors.Is(err, shared.ErrOutOfBoard) {
				continue
			}
			if err != nil {
				return nil, err
			}
			isEqual, err := positionNeighbor.isEqual(position)
			if err != nil {
				return nil, err
			}
			if !isEqual {
				positions = append(positions, positionNeighbor)
			}
		}
	}
	return positions, nil
}

func (position *Position) GetPosition() (int, int, error) {
	if position == nil {
		return 0, 0, shared.ErrPositionIsNil
	}
	return position.x, position.y, nil
}

func (position *Position) isEqual(positionExt *Position) (bool, error) {
	if position == nil || positionExt == nil {
		return false, shared.ErrPositionIsNil
	}
	return position.x == positionExt.x && position.y == positionExt.y, nil
}
