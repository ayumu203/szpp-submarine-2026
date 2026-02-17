package shared

type ErrorCode int

const (
	invalidTurn = iota
	invalidAction
	invalidTarget
	invalidMoveDistance
	outOfBoard
)