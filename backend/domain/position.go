package domain

import share "backend/domain/shared"

type Position struct {
	x int
	y int
}

func NewPosition(x int, y int) (*Position, error) {
	position := Position{
		x: x,
		y: y,
	}
	if !position.withinBoard() {
		return nil, share.ErrOutOfBoard
	}
	return &position, nil
}

func (position *Position) withinBoard() bool {
	return position.x >= share.MinPosition && position.x <= share.MaxPosition && position.y >= share.MinPosition && position.y <= share.MaxPosition
}

func (position *Position) Neighbors8() []*Position {
	positions := make([]*Position, 0, 8)
	delta := []int{-1, 0, 1}
	x, y := position.GetPosition()
	for _, dx := range delta {
		for _, dy := range delta {
			positionNeighbor, _ := NewPosition(x+dx, y+dy)
			if positionNeighbor != nil && !positionNeighbor.isEqual(position) {
				positions = append(positions, positionNeighbor)
			}
		}
	}
	return positions
}

func (position *Position) GetPosition() (int, int) {
	return position.x, position.y
}

func (position *Position) isEqual(positionExt *Position) bool {
	return position.x == positionExt.x && position.y == positionExt.y
}
