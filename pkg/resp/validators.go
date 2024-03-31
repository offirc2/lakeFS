package resp

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")
)

func NArgs(n int) ValidateFn {
	return func(args [][]byte) error {
		if len(args) != n {
			return fmt.Errorf("%w: expected '%d' args", ErrValidation, n)
		}
		return nil
	}
}

func NoArgs() ValidateFn {
	return NArgs(0)
}

func MinArgs(n int) ValidateFn {
	return func(args [][]byte) error {
		if len(args) < n {
			return fmt.Errorf("%w: expected at least '%d' args", ErrValidation, n)
		}
		return nil
	}
}

func ArgsBetween(min, max int) ValidateFn {
	return func(args [][]byte) error {
		if len(args) > max || len(args) < min {
			return fmt.Errorf("%w: expected %d-%d args", ErrValidation, min, max)
		}
		return nil
	}
}
