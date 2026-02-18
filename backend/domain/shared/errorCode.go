package shared

type ErrorCode int

const (
	InvalidTurn = iota
	InvalidAction
	InvalidTarget
	InvalidMoveDistance
	OutOfBoard
)
