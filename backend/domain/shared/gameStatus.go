// GameStatus represents the current state of a game
package shared

type GameStatus int

const (
	// waiting indicates the game is waiting to start.
	waiting = iota
	// inProgress indicates the game is currently in progress.
	inProgress
	// finished indicates the game has been finished.
)