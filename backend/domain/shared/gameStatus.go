package shared

type GameStatus int

const (
	Waiting = iota
	InProgress
	Finished
)