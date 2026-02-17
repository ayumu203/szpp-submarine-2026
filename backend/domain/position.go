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
	return position.x >= 1 && position.x <= 5 && position.y >= 1 && position.y <= 5
}

func (position *Position) Neighbors8() []*Position {
	positions := []*Position{}
	d := []int{-1, 0, 1}
	x, y := position.GetPosition()
	for _, i := range d {
		for _, j := range d {
			positionNeighbor, _ := NewPosition(x+i, y+j)
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
