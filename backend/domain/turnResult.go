package domain

import (
	"backend/domain/shared"
)

type TurnResult struct {
	AttackReport shared.AttackReportType
	MoveReport   shared.MoveReportType
	errorCode    shared.ErrorCode
	HitCount     int
	sunkCount    int
	nextPlayerId string
}

func (tr *TurnResult) GetErrorCode() shared.ErrorCode {
	return tr.errorCode
}

func (tr *TurnResult) GetSunkCount() int {
	return tr.sunkCount
}

func (tr *TurnResult) GetNextPlayerId() string {
	return tr.nextPlayerId
}
