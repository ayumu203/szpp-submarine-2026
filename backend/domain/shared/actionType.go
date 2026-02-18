package shared

type ActionType int

const (
	Attack = iota
	Move
	ActionUnknown
)

func (a ActionType) String() string {
	switch a {
	case Attack:
		return "attack"
	case Move:
		return "move"
	default:
		return "unknown"
	}
}
