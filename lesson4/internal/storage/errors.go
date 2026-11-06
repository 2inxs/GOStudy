package storage

import (
	"errors"
	"fmt"
)

var ErrNotStorable = errors.New("предмет не поддерживает сохранение")

type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("строка %d: %v", e.Line, e.Err)
}

func (e *ParseItemError) Unwrap() error {
	return e.Err
}
