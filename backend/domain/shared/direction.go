package shared

type Direction int

const (
	North = iota
	East
	South
	West
	DirectionUnknown
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
