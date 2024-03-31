package resp

import (
	"errors"
)

var (
	ErrProtocol = errors.New("protocol error")
)

type ErrorPrefix string

const (
	ErrorPrefixGeneric ErrorPrefix = "ERR"
)
