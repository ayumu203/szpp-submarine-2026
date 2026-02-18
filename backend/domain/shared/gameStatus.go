// GameStatus represents the current state of a game
package shared

type GameStatus int

const (
	// waiting indicates the game is waiting to start.
	Waiting = iota
	// inProgress indicates the game is currently in progress.
	InProgress
	// finished indicates the game has been finished.
	Finished
)

func (s GameStatus) String() string {
	switch s {
	case Waiting:
		return "waiting"
	case InProgress:
		return "in_progress"
	case Finished:
		return "finished"
	default:
		return "unknown"
	}
}
