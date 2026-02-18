package shared

type GameStatus int

const (
	Waiting = iota
	InProgress
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
