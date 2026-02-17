package domain

import (
	"testing"

	share "backend/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestNewPositionSuccess(t *testing.T) {
	testList := []struct {
		name      string
		x         int
		y         int
		expectedX int
		expectedY int
	}{{
		"[NewPosition: 1, 1で初期化]",
		1,
		1,
		1,
		1,
	}, {
		"[NewPosition: 5, 5で初期化]",
		5,
		5,
		5,
		5,
	}, {
		"[NewPosition: 2, 3で初期化]",
		2,
		3,
		2,
		3,
	}}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			point, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			assert.Equal(t, point.x, tl.expectedX)
			assert.Equal(t, point.y, tl.expectedY)
		})
	}
}

func TestNewPositionFail(t *testing.T) {
	testList := []struct {
		name        string
		x           int
		y           int
		expectedErr error
	}{{
		"[NewPosition: 0, 0で初期化]",
		0,
		0,
		share.ErrOutOfBoard,
	}, {
		"[NewPosition: 6, 6で初期化]",
		6,
		6,
		share.ErrOutOfBoard,
	}, {
		"[NewPosition: -1, 3で初期化]",
		-1,
		3,
		share.ErrOutOfBoard,
	}}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			point, err := NewPosition(tl.x, tl.y)
			assert.Nil(t, point)
			assert.ErrorIs(t, err, tl.expectedErr)
		})
	}
}

func TestGetPositionSuccess(t *testing.T) {
	testList := []struct {
		name      string
		x         int
		y         int
		expectedX int
		expectedY int
	}{
		{
			"[GetPosition: 1, 1]",
			1,
			1,
			1,
			1,
		},
		{
			"[GetPosition: 5, 5]",
			5,
			5,
			5,
			5,
		},
		{
			"[GetPosition: 3, 2]",
			3,
			2,
			3,
			2,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			pos, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			x, y, err := pos.GetPosition()
			assert.NoError(t, err)
			assert.Equal(t, tl.expectedX, x)
			assert.Equal(t, tl.expectedY, y)
		})
	}
}

func TestGetPositionFail(t *testing.T) {
	t.Run("[GetPosition: nil]", func(t *testing.T) {
		position, _ := NewPosition(0, 0)
		_, _, err := position.GetPosition()
		assert.ErrorIs(t, err, share.ErrPositionIsNil)
	})
}

func TestNeighbors8Success(t *testing.T) {
	newPos := func(x int, y int) *Position {
		pos, err := NewPosition(x, y)
		assert.NoError(t, err)
		return pos
	}
	testList := []struct {
		name              string
		position          *Position
		expectedNeighbors []*Position
	}{
		{
			"[Neighbors8: 2, 3の周辺]",
			newPos(2, 3),
			[]*Position{newPos(1, 2), newPos(1, 3), newPos(1, 4), newPos(2, 2), newPos(2, 4), newPos(3, 2), newPos(3, 3), newPos(3, 4)},
		},
		{
			"[Neighbors8: 1, 1の左上コーナー]",
			newPos(1, 1),
			[]*Position{newPos(1, 2), newPos(2, 1), newPos(2, 2)},
		},
		{
			"[Neighbors8: 5, 5の右下コーナー]",
			newPos(5, 5),
			[]*Position{newPos(4, 4), newPos(4, 5), newPos(5, 4)},
		},
		{
			"[Neighbors8: 1, 5の左下コーナー]",
			newPos(1, 5),
			[]*Position{newPos(1, 4), newPos(2, 4), newPos(2, 5)},
		},
		{
			"[Neighbors8: 5, 1の右上コーナー]",
			newPos(5, 1),
			[]*Position{newPos(4, 1), newPos(4, 2), newPos(5, 2)},
		},
		{
			"[Neighbors8: 1, 3の左エッジ]",
			newPos(1, 3),
			[]*Position{newPos(1, 2), newPos(1, 4), newPos(2, 2), newPos(2, 3), newPos(2, 4)},
		},
		{
			"[Neighbors8: 5, 3の右エッジ]",
			newPos(5, 3),
			[]*Position{newPos(4, 2), newPos(4, 3), newPos(4, 4), newPos(5, 2), newPos(5, 4)},
		},
		{
			"[Neighbors8: 3, 1の上エッジ]",
			newPos(3, 1),
			[]*Position{newPos(2, 1), newPos(2, 2), newPos(3, 2), newPos(4, 1), newPos(4, 2)},
		},
		{
			"[Neighbors8: 3, 5の下エッジ]",
			newPos(3, 5),
			[]*Position{newPos(2, 4), newPos(2, 5), newPos(3, 4), newPos(4, 4), newPos(4, 5)},
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			neighbors, err := tl.position.Neighbors8()
			assert.NoError(t, err)
			assert.ElementsMatch(t, neighbors, tl.expectedNeighbors)
		})
	}
}

func TestNeighbors8Fail(t *testing.T) {
	t.Run("[Neighbors8: nil]", func(t *testing.T) {
		position, _ := NewPosition(0, 0)
		neighbors, err := position.Neighbors8()
		assert.Nil(t, neighbors)
		assert.ErrorIs(t, err, share.ErrPositionIsNil)
	})
}
