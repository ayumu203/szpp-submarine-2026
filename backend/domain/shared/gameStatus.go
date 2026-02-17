package shared

type gameStatus int

const (
	waiting = iota
	inProgress
	finished
)