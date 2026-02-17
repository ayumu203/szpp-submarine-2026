package shared

type AttackReportType int

const (
	InvalidAttack = iota
	Miss
	Hit
	HitAndSunk
	WaveHigh
)