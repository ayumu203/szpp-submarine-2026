package shared

import (
	"errors"
)

var (
	ErrOutOfBoard = errors.New("Error[Position.go]: 場所がボードの外です．")
)
