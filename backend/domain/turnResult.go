package domain

import (
	share "backend/domain/shared"
)

type TurnResult struct {
	AttackReport share.AttackReportType
	MoveReport   share.MoveReportType
	ErrorCode    share.ErrorCode
	HitCount     int
	SunkCount    int
	NextPlayerId string
}