package shared

type AttackReportType int

const (
	invalidAttack = iota
	miss
	hit
	hitAndSunk
	waveHigh
)