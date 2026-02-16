package domain

import (
	"math"
	"sort"
	"strings"
	"time"
)

const (
	BoardSize           = 5
	SubmarinesPerPlayer = 4
	InitialSubmarineHP  = 3
)

type GameStatus string

const (
	GameStatusWaiting    GameStatus = "waiting"
	GameStatusInProgress GameStatus = "inProgress"
	GameStatusFinished   GameStatus = "finished"
)

type ActionType string

const (
	ActionTypeAttack ActionType = "attack"
	ActionTypeMove   ActionType = "move"
)

type Direction string

const (
	DirectionNorth Direction = "north"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
	DirectionWest  Direction = "west"
)

type ReportType string

const (
	ReportMiss       ReportType = "miss"
	ReportHit        ReportType = "hit"
	ReportHitAndSunk ReportType = "hitAndSunk"
	ReportWaveHigh   ReportType = "waveHigh"
)

type MoveReportType string

const (
	MoveReportSuccess MoveReportType = "moveSuccess"
	MoveReportBlocked MoveReportType = "moveBlocked"
)

type ErrorCode string

const (
	ErrorInvalidTurn         ErrorCode = "invalidTurn"
	ErrorInvalidAction       ErrorCode = "invalidAction"
	ErrorInvalidTarget       ErrorCode = "invalidTarget"
	ErrorInvalidMoveDistance ErrorCode = "invalidMoveDistance"
	ErrorOutOfBoard          ErrorCode = "outOfBoard"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Position) WithinBoard() bool {
	return p.X >= 1 && p.X <= BoardSize && p.Y >= 1 && p.Y <= BoardSize
}

func (p Position) Neighbors8() []Position {
	neighbors := make([]Position, 0, 8)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			next := Position{X: p.X + dx, Y: p.Y + dy}
			if next.WithinBoard() {
				neighbors = append(neighbors, next)
			}
		}
	}
	return neighbors
}

type Submarine struct {
	ID       string   `json:"id"`
	OwnerID  string   `json:"ownerId"`
	Position Position `json:"position"`
	HP       int      `json:"hp"`
}

func (s *Submarine) IsSunk() bool {
	return s.HP <= 0
}

func (s *Submarine) TakeDamage(amount int) {
	s.HP -= amount
	if s.HP < 0 {
		s.HP = 0
	}
}

func (s *Submarine) MoveTo(position Position) {
	s.Position = position
}

type Board struct {
	Submarines map[string]*Submarine `json:"submarines"`
}

func NewBoard() Board {
	return Board{Submarines: map[string]*Submarine{}}
}

func (b *Board) PlaceSubmarine(playerID, submarineID string, position Position) bool {
	if !position.WithinBoard() || b.IsOccupied(position) {
		return false
	}
	b.Submarines[submarineID] = &Submarine{
		ID:       submarineID,
		OwnerID:  playerID,
		Position: position,
		HP:       InitialSubmarineHP,
	}
	return true
}

func (b *Board) IsOccupied(position Position) bool {
	return b.GetSubmarineAt(position) != nil
}

func (b *Board) GetSubmarineAt(position Position) *Submarine {
	for _, submarine := range b.Submarines {
		if submarine.Position == position {
			return submarine
		}
	}
	return nil
}

func (b *Board) SortedSubmarines() []*Submarine {
	list := make([]*Submarine, 0, len(b.Submarines))
	for _, submarine := range b.Submarines {
		list = append(list, submarine)
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i].ID, list[j].ID) < 0
	})
	return list
}

func (b *Board) AliveSubmarinesOf(playerID string) []*Submarine {
	result := []*Submarine{}
	for _, submarine := range b.SortedSubmarines() {
		if submarine.OwnerID == playerID && !submarine.IsSunk() {
			result = append(result, submarine)
		}
	}
	return result
}

func (b *Board) RemainingHP(playerID string) int {
	total := 0
	for _, submarine := range b.Submarines {
		if submarine.OwnerID == playerID {
			total += submarine.HP
		}
	}
	return total
}

func (b *Board) CanAttackTarget(attackerID string, target Position) bool {
	for _, submarine := range b.AliveSubmarinesOf(attackerID) {
		dx := math.Abs(float64(submarine.Position.X - target.X))
		dy := math.Abs(float64(submarine.Position.Y - target.Y))
		if dx <= 1 && dy <= 1 && !(dx == 0 && dy == 0) {
			return true
		}
	}
	return false
}

func (b *Board) Attack(attackerID string, target Position) (ReportType, ErrorCode) {
	if !target.WithinBoard() {
		return "", ErrorOutOfBoard
	}
	targetSubmarine := b.GetSubmarineAt(target)
	if targetSubmarine != nil {
		if targetSubmarine.OwnerID == attackerID {
			return "", ErrorInvalidTarget
		}
		if targetSubmarine.IsSunk() {
			return ReportMiss, ""
		}
		targetSubmarine.TakeDamage(1)
		if targetSubmarine.IsSunk() {
			return ReportHitAndSunk, ""
		}
		return ReportHit, ""
	}
	if !b.CanAttackTarget(attackerID, target) {
		return "", ErrorInvalidTarget
	}

	for _, neighbor := range target.Neighbors8() {
		neighborSubmarine := b.GetSubmarineAt(neighbor)
		if neighborSubmarine != nil && neighborSubmarine.OwnerID != attackerID && !neighborSubmarine.IsSunk() {
			return ReportWaveHigh, ""
		}
	}
	return ReportMiss, ""
}

func (b *Board) MoveSubmarine(playerID string, direction Direction, distance int) (MoveReportType, ErrorCode) {
	if distance != 1 && distance != 2 {
		return "", ErrorInvalidMoveDistance
	}
	vector := Position{}
	switch direction {
	case DirectionNorth:
		vector = Position{X: 0, Y: -1}
	case DirectionSouth:
		vector = Position{X: 0, Y: 1}
	case DirectionEast:
		vector = Position{X: 1, Y: 0}
	case DirectionWest:
		vector = Position{X: -1, Y: 0}
	default:
		return "", ErrorInvalidAction
	}

	alive := b.AliveSubmarinesOf(playerID)
	for _, submarine := range alive {
		blockedBySunk := false
		for step := 1; step <= distance; step++ {
			cell := Position{X: submarine.Position.X + vector.X*step, Y: submarine.Position.Y + vector.Y*step}
			if !cell.WithinBoard() {
				blockedBySunk = true
				break
			}
			atCell := b.GetSubmarineAt(cell)
			if atCell != nil && atCell.IsSunk() {
				blockedBySunk = true
				break
			}
		}
		if blockedBySunk {
			continue
		}

		target := Position{X: submarine.Position.X + vector.X*distance, Y: submarine.Position.Y + vector.Y*distance}
		if !target.WithinBoard() {
			continue
		}
		occupied := b.GetSubmarineAt(target)
		if occupied != nil && occupied.ID != submarine.ID {
			continue
		}
		submarine.MoveTo(target)
		return MoveReportSuccess, ""
	}
	return MoveReportBlocked, ""
}

type ActionCommand struct {
	PlayerID  string
	Type      ActionType
	Target    *Position
	Direction Direction
	Distance  int
}

type TurnResult struct {
	AttackReport *ReportType
	MoveReport   *MoveReportType
	ErrorCode    *ErrorCode
	NextPlayerID string
	WinnerID     string
	Status       GameStatus
}

type TurnLog struct {
	Turn         int
	PlayerID     string
	ActionType   ActionType
	Target       *Position
	Direction    *Direction
	Distance     *int
	AttackReport *ReportType
	MoveReport   *MoveReportType
	ErrorCode    *ErrorCode
	CreatedAt    string
}

type PredictionBoard struct {
	ScoreGrid          [BoardSize][BoardSize]int `json:"scoreGrid"`
	PossibleEnemyCount [BoardSize][BoardSize]int `json:"possibleEnemyCount"`
	UpdatedAt          string                    `json:"updatedAt"`
}

func NewPredictionBoard() PredictionBoard {
	return PredictionBoard{UpdatedAt: time.Now().UTC().Format(time.RFC3339)}
}

type Game struct {
	ID              string     `json:"id"`
	Status          GameStatus `json:"status"`
	Turn            int        `json:"turn"`
	PlayerAID       string     `json:"playerAId"`
	PlayerBID       string     `json:"playerBId"`
	CurrentPlayerID string     `json:"currentPlayerId"`
	WinnerID        string     `json:"winnerId"`
	Board           Board      `json:"board"`
	CreatedAt       string     `json:"createdAt"`
	UpdatedAt       string     `json:"updatedAt"`
}

func (g *Game) Start() {
	g.Status = GameStatusInProgress
	g.Turn = 1
	g.CurrentPlayerID = g.PlayerAID
	now := time.Now().UTC().Format(time.RFC3339)
	if g.CreatedAt == "" {
		g.CreatedAt = now
	}
	g.UpdatedAt = now
}

func (g *Game) IsFinished() bool {
	return g.Status == GameStatusFinished
}

func (g *Game) OpponentOf(playerID string) string {
	if playerID == g.PlayerAID {
		return g.PlayerBID
	}
	return g.PlayerAID
}

func (g *Game) Apply(command ActionCommand) TurnResult {
	result := TurnResult{
		NextPlayerID: g.CurrentPlayerID,
		Status:       g.Status,
	}
	if g.IsFinished() {
		errorCode := ErrorInvalidAction
		result.ErrorCode = &errorCode
		return result
	}
	if command.PlayerID != g.CurrentPlayerID {
		errorCode := ErrorInvalidTurn
		result.ErrorCode = &errorCode
		return result
	}

	switch command.Type {
	case ActionTypeAttack:
		if command.Target == nil {
			errorCode := ErrorInvalidTarget
			result.ErrorCode = &errorCode
			return result
		}
		report, errCode := g.Board.Attack(command.PlayerID, *command.Target)
		if errCode != "" {
			result.ErrorCode = &errCode
			result.Status = g.Status
			return result
		}
		result.AttackReport = &report
	case ActionTypeMove:
		report, errCode := g.Board.MoveSubmarine(command.PlayerID, command.Direction, command.Distance)
		if errCode != "" {
			result.ErrorCode = &errCode
			result.Status = g.Status
			return result
		}
		result.MoveReport = &report
	default:
		errorCode := ErrorInvalidAction
		result.ErrorCode = &errorCode
		return result
	}

	opponent := g.OpponentOf(command.PlayerID)
	if g.Board.RemainingHP(opponent) == 0 {
		g.Status = GameStatusFinished
		g.WinnerID = command.PlayerID
		result.WinnerID = g.WinnerID
	}
	if g.Status != GameStatusFinished {
		g.CurrentPlayerID = opponent
		g.Turn++
	}
	g.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	result.NextPlayerID = g.CurrentPlayerID
	result.Status = g.Status
	return result
}
