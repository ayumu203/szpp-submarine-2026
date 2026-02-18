package domain

import (
	shared "backend/domain/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	t.Run("[NewBoard]", func(t *testing.T) {
		board := NewBoard()
		assert.NotNil(t, board)
	})
}

func TestPlaceSubmarineSuccess(t *testing.T) {
	testList := []struct {
		name        string
		submarineId shared.SubmarineId
		playerId    shared.PlayerId
		x           int
		y           int
	}{
		{
			name:        "[PlaceSubmarine: 1, 1に置く]",
			submarineId: shared.SubmarineId("s1"),
			playerId:    shared.PlayerId("p1"),
			x:           1,
			y:           1,
		},
		{
			name:        "[PlaceSubmarine: 5, 5に置く]",
			submarineId: shared.SubmarineId("s2"),
			playerId:    shared.PlayerId("p2"),
			x:           5,
			y:           5,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, position)
			assert.NoError(t, err)
		})
	}
}

func TestPlaceSubmarineDuplicate(t *testing.T) {
	testList := []struct {
		name         string
		playerId     shared.PlayerId
		submarineId1 shared.SubmarineId
		submarineId2 shared.SubmarineId
		x            int
		y            int
	}{
		{
			name:         "[PlaceSubmarine: 1, 1に置く]",
			submarineId1: shared.SubmarineId("s1a"),
			submarineId2: shared.SubmarineId("s1b"),
			playerId:     shared.PlayerId("p1"),
			x:            1,
			y:            1,
		},
		{
			name:         "[PlaceSubmarine: 5, 5に置く]",
			submarineId1: shared.SubmarineId("s2a"),
			submarineId2: shared.SubmarineId("s2b"),
			playerId:     shared.PlayerId("p2"),
			x:            5,
			y:            5,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.x, tl.y)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId1, position)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId2, position)
			assert.ErrorIs(t, err, shared.ErrAllySubmarineAlreadyExists)
		})
	}
}
func TestMoveSubmarineSuccess(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		submarineId shared.SubmarineId
		startX      int
		startY      int
		direction   shared.Direction
		distance    int
	}{
		{
			name:        "[MoveSubmarine: 北に1移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      3,
			startY:      3,
			direction:   shared.North,
			distance:    1,
		},
		{
			name:        "[MoveSubmarine: 南に1移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      3,
			startY:      3,
			direction:   shared.South,
			distance:    1,
		},
		{
			name:        "[MoveSubmarine: 東に1移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      3,
			startY:      3,
			direction:   shared.East,
			distance:    1,
		},
		{
			name:        "[MoveSubmarine: 西に1移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      3,
			startY:      3,
			direction:   shared.West,
			distance:    1,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.startX, tl.startY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, position)
			assert.NoError(t, err)
			err = board.MoveSubmarine(tl.playerId, tl.submarineId, tl.direction, tl.distance)
			assert.NoError(t, err)
		})
	}
}

func TestMoveSubmarineOutOfBoard(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		submarineId shared.SubmarineId
		startX      int
		startY      int
		direction   shared.Direction
		distance    int
	}{
		{
			name:        "[MoveSubmarine: 北端から北に移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      1,
			startY:      1,
			direction:   shared.North,
			distance:    1,
		},
		{
			name:        "[MoveSubmarine: 西端から西に移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      1,
			startY:      1,
			direction:   shared.West,
			distance:    1,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.startX, tl.startY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, position)
			assert.NoError(t, err)
			err = board.MoveSubmarine(tl.playerId, tl.submarineId, tl.direction, tl.distance)
			assert.ErrorIs(t, err, shared.ErrOutOfBoard)
		})
	}
}

func TestMoveSubmarineInvalidDistance(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		submarineId shared.SubmarineId
		startX      int
		startY      int
		direction   shared.Direction
		distance    int
	}{
		{
			name:        "[MoveSubmarine: 移動距離0]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      5,
			startY:      5,
			direction:   shared.North,
			distance:    0,
		},
		{
			name:        "[MoveSubmarine移動距離が最大値超過]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      5,
			startY:      5,
			direction:   shared.North,
			distance:    shared.MaxMoveDistance + 1,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.startX, tl.startY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, position)
			assert.NoError(t, err)
			err = board.MoveSubmarine(tl.playerId, tl.submarineId, tl.direction, tl.distance)
			assert.ErrorIs(t, err, shared.ErrInvalidMoveDistance)
		})
	}
}

func TestMoveSubmarineInvalidPlayer(t *testing.T) {
	testList := []struct {
		name        string
		playerId1   shared.PlayerId
		playerId2   shared.PlayerId
		submarineId shared.SubmarineId
		startX      int
		startY      int
		direction   shared.Direction
		distance    int
	}{
		{
			name:        "[MoveSubmarine: 自分のではない潜水艦を移動]",
			playerId1:   shared.PlayerId("p1"),
			playerId2:   shared.PlayerId("p2"),
			submarineId: shared.SubmarineId("s1"),
			startX:      5,
			startY:      5,
			direction:   shared.North,
			distance:    shared.MaxMoveDistance,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.startX, tl.startY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId1, tl.submarineId, position)
			assert.NoError(t, err)
			err = board.MoveSubmarine(tl.playerId2, tl.submarineId, tl.direction, tl.distance)
			assert.ErrorIs(t, err, shared.ErrPlayerIdNotMatch)
		})
	}
}

func TestMoveSunkSubmarine(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		submarineId shared.SubmarineId
		startX      int
		startY      int
		direction   shared.Direction
		distance    int
	}{
		{
			name:        "[MoveSubmarine: 沈没した潜水艦を移動]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			startX:      5,
			startY:      5,
			direction:   shared.North,
			distance:    shared.MaxMoveDistance,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			position, err := NewPosition(tl.startX, tl.startY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, position)
			assert.NoError(t, err)
			err = board.submarines[tl.submarineId].TakeDamage(3)
			assert.NoError(t, err)
			err = board.MoveSubmarine(tl.playerId, tl.submarineId, tl.direction, tl.distance)
			assert.ErrorIs(t, err, shared.ErrMoveSunkSubmarine)
		})
	}
}

func TestMoveSubmarineInvalidDirection(t *testing.T) {
	t.Run("[MoveSubmarine: 不正な方向]", func(t *testing.T) {
		board := NewBoard()
		playerId := shared.PlayerId("p1")
		submarineId := shared.SubmarineId("s1")
		position, err := NewPosition(5, 5)
		assert.NoError(t, err)
		err = board.PlaceSubmarine(playerId, submarineId, position)
		assert.NoError(t, err)
		err = board.MoveSubmarine(playerId, submarineId, shared.Direction(99), 1)
		assert.ErrorIs(t, err, shared.ErrInvalidDirection)
	})
}

func TestMoveSubmarineNotFound(t *testing.T) {
	t.Run("[MoveSubmarine: 存在しない潜水艦を移動]", func(t *testing.T) {
		board := NewBoard()
		err := board.MoveSubmarine(shared.PlayerId("s1"), shared.SubmarineId("s99"), shared.North, 1)
		assert.ErrorIs(t, err, shared.ErrSubmarineNotFound)
	})
}

func TestMoveSubmarineOverSunkSubmarineOccupied(t *testing.T) {
	testList := []struct {
		name            string
		playerId1       shared.PlayerId
		playerId2       shared.PlayerId
		moveSubmarineId shared.SubmarineId
		moveSubmarineX  int
		moveSubmarineY  int
		sunkSubmarineId shared.SubmarineId
		sunkSubmarineX  int
		sunkSubmarineY  int
		moveDirection   shared.Direction
		moveDistance    int
	}{
		{
			name:            "[MoveSubmarine: 沈没した潜水艦がある場所に移動]",
			playerId1:       shared.PlayerId("p1"),
			playerId2:       shared.PlayerId("p2"),
			moveSubmarineId: shared.SubmarineId("s1"),
			moveSubmarineX:  3,
			moveSubmarineY:  4,
			sunkSubmarineId: shared.SubmarineId("s2"),
			sunkSubmarineX:  3,
			sunkSubmarineY:  3,
			moveDirection:   shared.North,
			moveDistance:    1,
		},
		{
			name:            "[MoveSubmarine: 沈没した潜水艦がある場所を通過]",
			playerId1:       shared.PlayerId("p1"),
			playerId2:       shared.PlayerId("p2"),
			moveSubmarineId: shared.SubmarineId("s1"),
			moveSubmarineX:  3,
			moveSubmarineY:  4,
			sunkSubmarineId: shared.SubmarineId("s2"),
			sunkSubmarineX:  3,
			sunkSubmarineY:  3,
			moveDirection:   shared.North,
			moveDistance:    2,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			moveSubmarine, err := NewPosition(tl.moveSubmarineX, tl.moveSubmarineY)
			assert.NoError(t, err)
			sunkSubmarine, err := NewPosition(tl.sunkSubmarineX, tl.sunkSubmarineY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId1, tl.moveSubmarineId, moveSubmarine)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId2, tl.sunkSubmarineId, sunkSubmarine)
			assert.NoError(t, err)
			err = board.submarines[tl.sunkSubmarineId].TakeDamage(3)
			assert.NoError(t, err)
			assert.True(t, board.submarines[tl.sunkSubmarineId].IsSunk())
			err = board.MoveSubmarine(tl.playerId1, tl.moveSubmarineId, tl.moveDirection, tl.moveDistance)
			assert.ErrorIs(t, err, shared.ErrMovedOverSunkSubmarine)
		})
	}
}

func TestFindTargets(t *testing.T) {
	testList := []struct {
		name              string
		attackerId        shared.PlayerId
		defenderId        shared.PlayerId
		attackerSubId     shared.SubmarineId
		defenderSubId     shared.SubmarineId
		attackerX         int
		attackerY         int
		defenderX         int
		defenderY         int
		centerX           int
		centerY           int
		expectedTargetNum int
	}{
		{
			name:              "[FindTarget: 攻撃範囲内に敵潜水艦が1隻]",
			attackerId:        shared.PlayerId("p1"),
			defenderId:        shared.PlayerId("p2"),
			attackerSubId:     shared.SubmarineId("s1"),
			defenderSubId:     shared.SubmarineId("s2"),
			attackerX:         4,
			attackerY:         4,
			defenderX:         5,
			defenderY:         5,
			centerX:           4,
			centerY:           5,
			expectedTargetNum: 1,
		},
		{
			name:              "攻撃範囲外に敵潜水艦",
			attackerId:        shared.PlayerId("p1"),
			defenderId:        shared.PlayerId("p2"),
			attackerSubId:     shared.SubmarineId("s1"),
			defenderSubId:     shared.SubmarineId("s2"),
			attackerX:         1,
			attackerY:         1,
			defenderX:         5,
			defenderY:         5,
			centerX:           1,
			centerY:           2,
			expectedTargetNum: 0,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			attackerPos, err := NewPosition(tl.attackerX, tl.attackerY)
			assert.NoError(t, err)
			defenderPos, err := NewPosition(tl.defenderX, tl.defenderY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.attackerId, tl.attackerSubId, attackerPos)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.defenderId, tl.defenderSubId, defenderPos)
			assert.NoError(t, err)
			center, err := NewPosition(tl.centerX, tl.centerY)
			assert.NoError(t, err)
			targets, err := board.FindTargets(tl.attackerId, center)
			assert.NoError(t, err)
			assert.Len(t, targets, tl.expectedTargetNum)
		})
	}
}

func TestIsOccupied(t *testing.T) {
	testList := []struct {
		name       string
		playerId   shared.PlayerId
		subId      shared.SubmarineId
		placeX     int
		placeY     int
		checkX     int
		checkY     int
		isOccupied bool
	}{
		{
			name:       "潜水艦が存在するマス",
			playerId:   shared.PlayerId("p1"),
			subId:      shared.SubmarineId("s1"),
			placeX:     3,
			placeY:     3,
			checkX:     3,
			checkY:     3,
			isOccupied: true,
		},
		{
			name:       "潜水艦が存在しないマス",
			playerId:   shared.PlayerId("p1"),
			subId:      shared.SubmarineId("s1"),
			placeX:     3,
			placeY:     3,
			checkX:     4,
			checkY:     4,
			isOccupied: false,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			placePos, err := NewPosition(tl.placeX, tl.placeY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.subId, placePos)
			assert.NoError(t, err)
			checkPos, err := NewPosition(tl.checkX, tl.checkY)
			assert.NoError(t, err)
			occupied, err := board.IsOccupied(checkPos)
			assert.NoError(t, err)
			assert.Equal(t, tl.isOccupied, occupied)
		})
	}
}
func TestGetAllySubmarines(t *testing.T) {
	testList := []struct {
		name              string
		playerId          shared.PlayerId
		opponentId        shared.PlayerId
		allySubIds        []shared.SubmarineId
		opponentSubIds    []shared.SubmarineId
		allyPositions     [][2]int
		opponentPositions [][2]int
		expectedCount     int
	}{
		{
			name:              "[GetAllySubmarines: 味方潜水艦2隻を取得]",
			playerId:          shared.PlayerId("p1"),
			opponentId:        shared.PlayerId("p2"),
			allySubIds:        []shared.SubmarineId{"s1", "s2"},
			opponentSubIds:    []shared.SubmarineId{"s3"},
			allyPositions:     [][2]int{{1, 1}, {2, 2}},
			opponentPositions: [][2]int{{5, 5}},
			expectedCount:     2,
		},
		{
			name:              "[GetAllySubmarines: 味方潜水艦0隻]",
			playerId:          shared.PlayerId("p1"),
			opponentId:        shared.PlayerId("p2"),
			allySubIds:        []shared.SubmarineId{},
			opponentSubIds:    []shared.SubmarineId{"s3"},
			allyPositions:     [][2]int{},
			opponentPositions: [][2]int{{5, 5}},
			expectedCount:     0,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			for i, subId := range tl.allySubIds {
				pos, err := NewPosition(tl.allyPositions[i][0], tl.allyPositions[i][1])
				assert.NoError(t, err)
				err = board.PlaceSubmarine(tl.playerId, subId, pos)
				assert.NoError(t, err)
			}
			for i, subId := range tl.opponentSubIds {
				pos, err := NewPosition(tl.opponentPositions[i][0], tl.opponentPositions[i][1])
				assert.NoError(t, err)
				err = board.PlaceSubmarine(tl.opponentId, subId, pos)
				assert.NoError(t, err)
			}
			allies, err := board.GetAllySubmarines(tl.playerId)
			assert.NoError(t, err)
			assert.Len(t, allies, tl.expectedCount)
		})
	}
}

func TestGetOpponentSubmarines(t *testing.T) {
	testList := []struct {
		name              string
		playerId          shared.PlayerId
		opponentId        shared.PlayerId
		allySubIds        []shared.SubmarineId
		opponentSubIds    []shared.SubmarineId
		allyPositions     [][2]int
		opponentPositions [][2]int
		expectedCount     int
	}{
		{
			name:              "[GetOpponentSubmarines: 敵潜水艦2隻を取得]",
			playerId:          shared.PlayerId("p1"),
			opponentId:        shared.PlayerId("p2"),
			allySubIds:        []shared.SubmarineId{"s1"},
			opponentSubIds:    []shared.SubmarineId{"s2", "s3"},
			allyPositions:     [][2]int{{1, 1}},
			opponentPositions: [][2]int{{5, 5}, {4, 4}},
			expectedCount:     2,
		},
		{
			name:              "[GetOpponentSubmarines: 敵潜水艦0隻]",
			playerId:          shared.PlayerId("p1"),
			opponentId:        shared.PlayerId("p2"),
			allySubIds:        []shared.SubmarineId{"s1"},
			opponentSubIds:    []shared.SubmarineId{},
			allyPositions:     [][2]int{{1, 1}},
			opponentPositions: [][2]int{},
			expectedCount:     0,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			for i, subId := range tl.allySubIds {
				pos, err := NewPosition(tl.allyPositions[i][0], tl.allyPositions[i][1])
				assert.NoError(t, err)
				err = board.PlaceSubmarine(tl.playerId, subId, pos)
				assert.NoError(t, err)
			}
			for i, subId := range tl.opponentSubIds {
				pos, err := NewPosition(tl.opponentPositions[i][0], tl.opponentPositions[i][1])
				assert.NoError(t, err)
				err = board.PlaceSubmarine(tl.opponentId, subId, pos)
				assert.NoError(t, err)
			}
			opponents, err := board.GetOpponentSubmarines(tl.playerId)
			assert.NoError(t, err)
			assert.Len(t, opponents, tl.expectedCount)
		})
	}
}

func TestGetAllySubmarineAt(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		submarineId shared.SubmarineId
		placeX      int
		placeY      int
		checkX      int
		checkY      int
		expectFound bool
	}{
		{
			name:        "[GetAllySubmarineAt: 該当位置に味方潜水艦あり]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			placeX:      3,
			placeY:      3,
			checkX:      3,
			checkY:      3,
			expectFound: true,
		},
		{
			name:        "[GetAllySubmarineAt: 該当位置に味方潜水艦なし]",
			playerId:    shared.PlayerId("p1"),
			submarineId: shared.SubmarineId("s1"),
			placeX:      3,
			placeY:      3,
			checkX:      4,
			checkY:      4,
			expectFound: false,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			placePos, err := NewPosition(tl.placeX, tl.placeY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.submarineId, placePos)
			assert.NoError(t, err)
			checkPos, err := NewPosition(tl.checkX, tl.checkY)
			assert.NoError(t, err)
			submarine, err := board.GetAllySubmarineAt(tl.playerId, checkPos)
			assert.NoError(t, err)
			if tl.expectFound {
				assert.NotNil(t, submarine)
			} else {
				assert.Nil(t, submarine)
			}
		})
	}
}

func TestGetOpponentSubmarineAt(t *testing.T) {
	testList := []struct {
		name        string
		playerId    shared.PlayerId
		opponentId  shared.PlayerId
		subId       shared.SubmarineId
		opSubId     shared.SubmarineId
		placeX      int
		placeY      int
		checkX      int
		checkY      int
		expectFound bool
	}{
		{
			name:        "[GetOpponentSubmarineAt: 該当位置に敵潜水艦あり]",
			playerId:    shared.PlayerId("p1"),
			opponentId:  shared.PlayerId("p2"),
			subId:       shared.SubmarineId("s1"),
			opSubId:     shared.SubmarineId("s2"),
			placeX:      5,
			placeY:      5,
			checkX:      5,
			checkY:      5,
			expectFound: true,
		},
		{
			name:        "[GetOpponentSubmarineAt: 該当位置に敵潜水艦なし]",
			playerId:    shared.PlayerId("p1"),
			opponentId:  shared.PlayerId("p2"),
			subId:       shared.SubmarineId("s1"),
			opSubId:     shared.SubmarineId("s2"),
			placeX:      5,
			placeY:      5,
			checkX:      3,
			checkY:      3,
			expectFound: false,
		},
	}
	for _, tl := range testList {
		t.Run(tl.name, func(t *testing.T) {
			board := NewBoard()
			allyPos, err := NewPosition(1, 1)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.playerId, tl.subId, allyPos)
			assert.NoError(t, err)
			placePos, err := NewPosition(tl.placeX, tl.placeY)
			assert.NoError(t, err)
			err = board.PlaceSubmarine(tl.opponentId, tl.opSubId, placePos)
			assert.NoError(t, err)
			checkPos, err := NewPosition(tl.checkX, tl.checkY)
			assert.NoError(t, err)
			submarine, err := board.GetOpponentSubmarineAt(tl.playerId, checkPos)
			assert.NoError(t, err)
			if tl.expectFound {
				assert.NotNil(t, submarine)
			} else {
				assert.Nil(t, submarine)
			}
		})
	}
}

func TestPlaceSubmarineDuplicateId(t *testing.T) {
	t.Run("[PlaceSubmarine: 同じSubmarineIdで2回配置]", func(t *testing.T) {
		board := NewBoard()
		playerId := shared.PlayerId("p1")
		submarineId := shared.SubmarineId("s1")
		pos1, err := NewPosition(1, 1)
		assert.NoError(t, err)
		pos2, err := NewPosition(2, 2)
		assert.NoError(t, err)
		err = board.PlaceSubmarine(playerId, submarineId, pos1)
		assert.NoError(t, err)
		err = board.PlaceSubmarine(playerId, submarineId, pos2)
		assert.ErrorIs(t, err, shared.ErrSubmarineIdDuplicated)
	})
}
