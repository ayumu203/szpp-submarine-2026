package domain

import (
	shared "backend/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubmarine(t *testing.T) {
	testList := []struct {
		name     string
		id       shared.SubmarineId
		ownerId  shared.PlayerId
		position Position
		hp       int
	}{
		{"[NewSubmarine: id:1, ownerId:1, position{3, 3}, hp:3で初期化]", "submarine-1", "player-1", Position{3, 3}, 3},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine, err := NewSubmarine(tl.id, tl.ownerId, &tl.position, tl.hp)
			assert.NoError(t, err)
			assert.Equal(t, tl.id, submarine.GetId())
			assert.Equal(t, tl.ownerId, submarine.GetOwnerId())
			assert.Equal(t, tl.position, *submarine.GetPosition())
			assert.Equal(t, tl.hp, submarine.GetHp())
		})
	}
}

func TestNewSubmarineFail(t *testing.T) {
	testList := []struct {
		name        string
		id          shared.SubmarineId
		ownerId     shared.PlayerId
		position    *Position
		hp          int
		expectedErr error
	}{
		{"[NewSubmarine: idが空]", "", "player-1", &Position{1, 1}, 3, shared.ErrSubmarineIdIsEmpty},
		{"[NewSubmarine: ownerIdが空]", "submarine-1", "", &Position{1, 1}, 3, shared.ErrOwnerIdIsEmpty},
		{"[NewSubmarine: positionが不正]", "submarine-1", "player-1", &Position{0, 0}, 3, shared.ErrOutOfBoard},
		{"[NewSubmarine: hpが負]", "submarine-1", "player-1", &Position{1, 1}, -1, shared.ErrInvalidHp},
		{"[NewSubmarine: positionがnil]", "submarine-1", "player-1", nil, 3, shared.ErrPositionIsNil},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine, err := NewSubmarine(tl.id, tl.ownerId, tl.position, tl.hp)
			assert.ErrorIs(t, err, tl.expectedErr)
			assert.Nil(t, submarine)
		})
	}
}

func TestSubmarineIsSunk(t *testing.T) {
	testList := []struct {
		name     string
		hp       int
		expected bool
	}{
		{"[IsSunk: hpが1なら沈没していない]", 1, false},
		{"[IsSunk: hpが0なら沈没している]", 0, true},
	}

	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine := &Submarine{hp: tl.hp}
			assert.Equal(t, tl.expected, submarine.IsSunk(nil))
		})
	}
}

func TestTakeDamage(t *testing.T) {
	testList := []struct {
		name            string
		initialHp       int
		expectedHp      int
		expectedErr     error
		expectedIsSunk  bool
	}{
		{"[TakeDamage: hpが1ならダメージ後に0になる]", 1, 0, nil, true},
		{"[TakeDamage: hpが0ならダメージ不可]", 0, 0, shared.ErrSubmarineAlreadySunk, true},
	}

	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine := &Submarine{hp: tl.initialHp}

			err := submarine.TakeDamage(nil)
			assert.ErrorIs(t, err, tl.expectedErr)
			assert.Equal(t, tl.expectedHp, submarine.GetHp())
			assert.Equal(t, tl.expectedIsSunk, submarine.IsSunk(nil))
		})
	}
}

func TestMoveTo(t *testing.T) {
	testList := []struct {
		name           string
		initialPos     *Position
		newPosition    *Position
		expectedPos    *Position
		expectedErr    error
	}{
		{"[MoveTo: 有効な位置なら移動する]", &Position{1, 1}, &Position{2, 2}, &Position{2, 2}, nil},
		{"[MoveTo: nil位置はエラー]", &Position{1, 1}, nil, &Position{1, 1}, shared.ErrPositionIsNil},
		{"[MoveTo: 盤外位置はエラー]", &Position{1, 1}, &Position{0, 0}, &Position{1, 1}, shared.ErrOutOfBoard},
	}

	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine := &Submarine{position: tl.initialPos}

			err := submarine.MoveTo(tl.newPosition)
			assert.ErrorIs(t, err, tl.expectedErr)
			assert.Equal(t, *tl.expectedPos, *submarine.GetPosition())
		})
	}
}
