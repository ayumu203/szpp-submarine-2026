package domain

import (
	shared "backend/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSubmarine(t *testing.T) {
	testList := []struct {
		name     string
		id       string
		ownerId  string
		position Position
		hp       int
	}{
		{"[NewSubmarine: id:1, ownerId:1, position{3, 3}, hp:3で初期化]", "submarine-1", "player-1", Position{3, 3}, 3},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine, err := NewSubmarine(tl.id, tl.ownerId, tl.position, tl.hp)
			assert.NoError(t, err)
			assert.Equal(t, tl.id, submarine.GetId())
			assert.Equal(t, tl.ownerId, submarine.GetOwnerId())
			assert.Equal(t, tl.position, submarine.GetPosition())
			assert.Equal(t, tl.hp, submarine.GetHp())
		})
	}
}

func TestNewSubmarineFail(t *testing.T) {
	testList := []struct {
		name        string
		id          string
		ownerId     string
		position    Position
		hp          int
		expectedErr error
	}{
		{"[NewSubmarine: idが空]", "", "player-1", Position{1, 1}, 3, shared.ErrSubmarineIdIsEmpty},
		{"[NewSubmarine: ownerIdが空]", "submarine-1", "", Position{1, 1}, 3, shared.ErrOwnerIdIsEmpty},
		{"[NewSubmarine: positionが不正]", "submarine-1", "player-1", Position{0, 0}, 3, shared.ErrOutOfBoard},
		{"[NewSubmarine: hpが負]", "submarine-1", "player-1", Position{1, 1}, -1, shared.ErrInvalidHp},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			submarine, err := NewSubmarine(tl.id, tl.ownerId, tl.position, tl.hp)
			assert.ErrorIs(t, err, tl.expectedErr)
			assert.Nil(t, submarine)
		})
	}
}
