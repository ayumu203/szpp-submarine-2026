package domain

import (
	share "backend/domain/shared"
	"strings"
)

type Player struct {
	id   string
	name string
}

func NewPlayer(id string, name string) (*Player, error) {
	if strings.TrimSpace(id) == "" {
		return nil, share.ErrInvalidPlayerID
	}
	if strings.TrimSpace(name) == "" {
		return nil, share.ErrInvalidPlayerName
	}
	return &Player{
		id:   id,
		name: strings.TrimSpace(name),
	}, nil
}

func (player *Player) RemainingHp() int {
	if player == nil {
		return 0
	}
	// submarineマージまでの対応.
	// interface を呼び出す(sugaくんの実装).
	return 12
}

func (player *Player) GetId() string {
	if player == nil {
		return ""
	}
	return player.id
}

func (player *Player) GetName() string {
	if player == nil {
		return ""
	}
	return player.name
}
