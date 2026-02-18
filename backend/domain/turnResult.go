package domain

import (
	shared "backend/domain/shared"
)

type TurnResult struct {
	AttackReport shared.AttackReportType
	MoveReport   shared.MoveReportType
	errorCode    shared.ErrorCode
	HitCount     int
	sunkCount    int
	nextPlayerId shared.PlayerId
}

func (tr *TurnResult) GetErrorCode() shared.ErrorCode {
	return tr.errorCode
}

func (tr *TurnResult) GetSunkCount() int {
	return tr.sunkCount
}

func (tr *TurnResult) GetNextPlayerId() shared.PlayerId {
	return tr.nextPlayerId
}
