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
		submarines   []*Submarine
		expectedID   string
		expectedName string
	}{
		{
			name:         "[NewPlayer: p1, Aliceで初期化]",
			id:           "p1",
			playerName:   "Alice",
			submarines:   []*Submarine{},
			expectedID:   "p1",
			expectedName: "Alice",
		},
		{
			name:         "[NewPlayer: p2, Bobで初期化]",
			id:           "p2",
			playerName:   "Bob",
			submarines:   []*Submarine{{id: "s1", hp: 3}},
			expectedID:   "p2",
			expectedName: "Bob",
		},
		{
			name:         "[NewPlayer: nameの前後空白はtrimされる]",
			id:           "p3",
			playerName:   "  Carol  ",
			submarines:   nil,
			expectedID:   "p3",
			expectedName: "Carol",
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			player, err := NewPlayer(tl.id, tl.playerName, tl.submarines)
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
			expectedErr: shared.ErrInvalidPlayerID,
		},
		{
			name:        "[NewPlayer: idが空白のみ]",
			id:          "   ",
			playerName:  "Alice",
			expectedErr: shared.ErrInvalidPlayerID,
		},
		{
			name:        "[NewPlayer: nameが空文字]",
			id:          "p1",
			playerName:  "",
			expectedErr: shared.ErrInvalidPlayerName,
		},
		{
			name:        "[NewPlayer: nameが空白のみ]",
			id:          "p1",
			playerName:  "   ",
			expectedErr: shared.ErrInvalidPlayerName,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			player, err := NewPlayer(tl.id, tl.playerName, nil)
			assert.Nil(t, player)
			assert.ErrorIs(t, err, tl.expectedErr)
		})
	}
}

func TestRemainingHp(t *testing.T) {
	t.Run("[RemainingHp: 初期残HPは12]", func(t *testing.T) {
		submarines := []*Submarine{
			{id: "s1", hp: 3},
			{id: "s2", hp: 3},
			{id: "s3", hp: 3},
			{id: "s4", hp: 3},
		}
		player, err := NewPlayer("p1", "Alice", submarines)
		assert.NoError(t, err)
		assert.Equal(t, 12, player.RemainingHp())
	})

	t.Run("[RemainingHp: nil receiverは0]", func(t *testing.T) {
		var player *Player
		assert.Equal(t, 0, player.RemainingHp())
	})
}

func TestGetSubmarines(t *testing.T) {
	t.Run("[GetSubmarines: nil receiverならnil]", func(t *testing.T) {
		var player *Player
		assert.Nil(t, player.GetSubmarines())
	})

	t.Run("[GetSubmarines: 返却スライスはコピー]", func(t *testing.T) {
		s1 := &Submarine{id: "s1", hp: 3}
		s2 := &Submarine{id: "s2", hp: 4}
		player := &Player{
			id:         "p1",
			name:       "Alice",
			submarines: []*Submarine{s1, s2},
		}

		got := player.GetSubmarines()
		assert.Equal(t, 2, len(got))
		assert.Equal(t, s1, got[0])
		assert.Equal(t, s2, got[1])

		got[0] = &Submarine{id: "other", hp: 9}
		assert.Equal(t, s1, player.submarines[0])
	})

	t.Run("[GetSubmarines: 内部がnilスライスなら空スライス]", func(t *testing.T) {
		player := &Player{id: "p1", name: "Alice", submarines: nil}
		got := player.GetSubmarines()
		assert.NotNil(t, got)
		assert.Len(t, got, 0)
	})
}
