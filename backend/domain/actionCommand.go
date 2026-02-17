package domain

import shared "backend/domain/shared"

type ActionCommand struct {
	playerId   string
	actionType shared.ActionType
	target     Position
	direction  shared.Direction
	distance   int
}
