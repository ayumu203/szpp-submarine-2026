package domain

import (
	shared "backend/domain/shared"
)

type Submarine struct {
	id       shared.SubmarineId
	ownerId  shared.PlayerId
	position *Position
	hp       int
}

const minDamage = 1

func NewSubmarine(id shared.SubmarineId, ownerId shared.PlayerId, position *Position, hp int) (*Submarine, error) {
	if id == "" {
		return nil, shared.ErrSubmarineIdIsEmpty
	}
	if ownerId == "" {
		return nil, shared.ErrOwnerIdIsEmpty
	}
	if position == nil {
		return nil, shared.ErrPositionIsNil
	}
	isWithin, err := position.withinBoard()
	if err != nil {
		return nil, err
	}
	if !isWithin {
		return nil, shared.ErrOutOfBoard
	}
	if hp < 0 {
		return nil, shared.ErrInvalidHp
	}
	return &Submarine{
		id:       id,
		ownerId:  ownerId,
		position: position,
		hp:       hp,
	}, nil
}

func (submarine *Submarine) IsSunk() bool {
	if submarine.hp <= 0 {
		return true
	}
	return false
}

func (submarine *Submarine) TakeDamage(amount int) error {
	if submarine.IsSunk() {
		return shared.ErrSubmarineAlreadySunk
	}
	if amount < minDamage {
		return shared.ErrInvalidDamageAmount
	}
	submarine.hp -= amount
	return nil
}

func (submarine *Submarine) MoveTo(newPosition *Position) error {
	if newPosition == nil {
		return shared.ErrPositionIsNil
	}
	isWithin, err := newPosition.withinBoard()
	if err != nil {
		return err
	}
	if !isWithin {
		return shared.ErrOutOfBoard
	}
	submarine.position = newPosition
	return nil
}

func (submarine *Submarine) GetId() shared.SubmarineId {
	return submarine.id
}

func (submarine *Submarine) GetOwnerId() shared.PlayerId {
	return submarine.ownerId
}

func (submarine *Submarine) GetPosition() *Position {
	return submarine.position
}

func (submarine *Submarine) GetHp() int {
	return submarine.hp
}
