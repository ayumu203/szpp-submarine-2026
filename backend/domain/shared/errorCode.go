package shared

type ErrorCode int

const (
	InvalidTurn = iota
	InvalidAction
	InvalidTarget
	InvalidMoveDistance
	OutOfBoard
	BoardIsNil
	SubmarineNotFound
	InvalidDirection
	AllySubmarineAlreadyExists
	SunkSubmarineAlreadyExists
)
