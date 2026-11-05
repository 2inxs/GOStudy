package storage

import (
	"errors"
	"fmt"
)

var ErrNotStorable = errors.New("предмет не поддерживает сохранение")

type RecordError struct {
	Line int
	Err  error
}

func (e *RecordError) Error() string {
	return fmt.Sprintf("строка %d: %v", e.Line, e.Err)
}

func (e *RecordError) Unwrap() error {
	return e.Err
}
