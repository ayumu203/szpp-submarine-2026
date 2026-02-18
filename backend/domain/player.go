package domain

import (
	shared "backend/domain/shared"
	"strings"
)

type Player struct {
	id         string
	name       string
	submarines []*Submarine
}

func NewPlayer(id string, name string, submarines []*Submarine) (*Player, error) {
	if strings.TrimSpace(id) == "" {
		return nil, shared.ErrInvalidPlayerID
	}
	if strings.TrimSpace(name) == "" {
		return nil, shared.ErrInvalidPlayerName
	}
	return &Player{
		id:         id,
		name:       strings.TrimSpace(name),
		submarines: submarines,
	}, nil
}

func (player *Player) RemainingHp() int {
	if player == nil {
		return 0
	}
	hp := 0
	for _, s := range player.submarines {
		hp += s.GetHp()
	}
	return hp
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

func (player *Player) GetSubmarines() []*Submarine {
	if player == nil {
		return nil
	}
	cp := make([]*Submarine, len(player.submarines))
	copy(cp, player.submarines)
	return cp
}
