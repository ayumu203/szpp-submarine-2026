// ActionType represents the type of action
package shared

type ActionType int

const (
	ActionUnknown = iota
	Attack
	Move
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
