package domain

import (
	"testing"

	shared "backend/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestNewPlayerSuccess(t *testing.T) {
	testList := []struct {
		name         string
		id           string
		playerName   string
		expectedID   string
		expectedName string
	}{
		{
			name:         "[NewPlayer: p1, Aliceで初期化]",
			id:           "p1",
			playerName:   "Alice",
			expectedID:   "p1",
			expectedName: "Alice",
		},
		{
			name:         "[NewPlayer: p2, Bobで初期化]",
			id:           "p2",
			playerName:   "Bob",
			expectedID:   "p2",
			expectedName: "Bob",
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			player, err := NewPlayer(tl.id, tl.playerName)
			assert.NoError(t, err)
			assert.NotNil(t, player)
			assert.Equal(t, tl.expectedID, player.id)
			assert.Equal(t, tl.expectedName, player.name)
		})
	}
}

func TestNewPlayerFail(t *testing.T) {
	testList := []struct {
		name        string
		id          string
		playerName  string
		expectedErr error
	}{
		{
			name:        "[NewPlayer: idが空文字]",
			id:          "",
			playerName:  "Alice",
			expectedErr: share.ErrInvalidPlayerID,
		},
		{
			name:        "[NewPlayer: idが空白のみ]",
			id:          "   ",
			playerName:  "Alice",
			expectedErr: share.ErrInvalidPlayerID,
		},
		{
			name:        "[NewPlayer: nameが空文字]",
			id:          "p1",
			playerName:  "",
			expectedErr: share.ErrInvalidPlayerName,
		},
		{
			name:        "[NewPlayer: nameが空白のみ]",
			id:          "p1",
			playerName:  "   ",
			expectedErr: share.ErrInvalidPlayerName,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			player, err := NewPlayer(tl.id, tl.playerName)
			assert.Nil(t, player)
			assert.ErrorIs(t, err, tl.expectedErr)
		})
	}
}

func TestRemainingHp(t *testing.T) {
	t.Run("[RemainingHp: 初期残HPは12]", func(t *testing.T) {
		player, err := NewPlayer("p1", "Alice")
		assert.NoError(t, err)
		assert.Equal(t, 12, player.RemainingHp())
	})

	t.Run("[RemainingHp: nil receiverは0]", func(t *testing.T) {
		var player *Player
		assert.Equal(t, 0, player.RemainingHp())
	})
}
