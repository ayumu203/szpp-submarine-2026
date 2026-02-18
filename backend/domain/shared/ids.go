package shared

type GameId string
type PlayerId string
type SubmarineId string

func (id GameId) String() string {
	return string(id)
}

func (id PlayerId) String() string {
	return string(id)
}

func (id SubmarineId) String() string {
	return string(id)
}
