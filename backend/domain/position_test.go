package domain

import (
	"testing"

	shared "backend/domain/shared"

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
			assert.Equal(t, tl.expectedX, point.x)
			assert.Equal(t, tl.expectedY, point.y)
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
		shared.ErrOutOfBoard,
	}, {
		"[NewPosition: 6, 6で初期化]",
		6,
		6,
		shared.ErrOutOfBoard,
	}, {
		"[NewPosition: -1, 3で初期化]",
		-1,
		3,
		shared.ErrOutOfBoard,
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
		assert.ErrorIs(t, err, shared.ErrPositionIsNil)
	})
}

func TestNeighbors8Success(t *testing.T) {
	testList := []struct {
		name              string
		x                 int
		y                 int
		expectedNeighbors []struct {
			x int
			y int
		}
	}{
		{
			"[Neighbors8: 2, 3の周辺]",
			2,
			3,
			[]struct {
				x int
				y int
			}{{1, 2}, {1, 3}, {1, 4}, {2, 2}, {2, 4}, {3, 2}, {3, 3}, {3, 4}},
		},
		{
			"[Neighbors8: 1, 1の左上コーナー]",
			1,
			1,
			[]struct {
				x int
				y int
			}{{1, 2}, {2, 1}, {2, 2}},
		},
		{
			"[Neighbors8: 5, 5の右下コーナー]",
			5,
			5,
			[]struct {
				x int
				y int
			}{{4, 4}, {4, 5}, {5, 4}},
		},
		{
			"[Neighbors8: 1, 5の左下コーナー]",
			1,
			5,
			[]struct {
				x int
				y int
			}{{1, 4}, {2, 4}, {2, 5}},
		},
		{
			"[Neighbors8: 5, 1の右上コーナー]",
			5,
			1,
			[]struct {
				x int
				y int
			}{{4, 1}, {4, 2}, {5, 2}},
		},
		{
			"[Neighbors8: 1, 3の左エッジ]",
			1,
			3,
			[]struct {
				x int
				y int
			}{{1, 2}, {1, 4}, {2, 2}, {2, 3}, {2, 4}},
		},
		{
			"[Neighbors8: 5, 3の右エッジ]",
			5,
			3,
			[]struct {
				x int
				y int
			}{{4, 2}, {4, 3}, {4, 4}, {5, 2}, {5, 4}},
		},
		{
			"[Neighbors8: 3, 1の上エッジ]",
			3,
			1,
			[]struct {
				x int
				y int
			}{{2, 1}, {2, 2}, {3, 2}, {4, 1}, {4, 2}},
		},
		{
			"[Neighbors8: 3, 5の下エッジ]",
			3,
			5,
			[]struct {
				x int
				y int
			}{{2, 4}, {2, 5}, {3, 4}, {4, 4}, {4, 5}},
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			position, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)

			neighbors, err := position.Neighbors8()
			assert.NoError(t, err)

			// 期待される座標のPositionリストを作成
			expectedPositions := make([]*Position, 0, len(tl.expectedNeighbors))
			for _, expected := range tl.expectedNeighbors {
				pos, err := NewPosition(expected.x, expected.y)
				assert.NoError(t, err)
				expectedPositions = append(expectedPositions, pos)
			}

			assert.ElementsMatch(t, expectedPositions, neighbors)
		})
	}
}

func TestNeighbors8Fail(t *testing.T) {
	t.Run("[Neighbors8: nil]", func(t *testing.T) {
		position, _ := NewPosition(0, 0)
		neighbors, err := position.Neighbors8()
		assert.Nil(t, neighbors)
		assert.ErrorIs(t, err, shared.ErrPositionIsNil)
	})
}
