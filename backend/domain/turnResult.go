package domain

import (
	share "backend/domain/shared"
)

type TurnResult struct {
	AttackReport share.AttackReportType
	MoveReport   share.MoveReportType
	errorCode    share.ErrorCode
	HitCount     int
	sunkCount    int
	nextPlayerId share.PlayerId
}

func (tr *TurnResult) GetErrorCode() share.ErrorCode {
	return tr.errorCode
}

func (tr *TurnResult) GetSunkCount() int {
	return tr.sunkCount
}

func (tr *TurnResult) GetNextPlayerId() share.PlayerId {
	return tr.nextPlayerId
}
