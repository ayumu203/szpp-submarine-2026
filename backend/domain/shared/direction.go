// Direction represents a movement or facing direction on the board
package shared

type Direction int

const (
	North = iota
	East
	South
	West
)

func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		return "unknown"
	}
}