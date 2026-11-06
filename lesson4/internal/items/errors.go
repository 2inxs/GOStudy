package items

import (
	"errors"
	"fmt"
)

var (
	ErrNilItem         = errors.New("предмет отсутствует")
	ErrBroken          = errors.New("item is broken")
	ErrExhausted       = errors.New("item is exhausted")
	ErrInvalidItem     = errors.New("некорректные параметры предмета")
	ErrInvalidRecord   = errors.New("некорректная запись предмета")
	ErrUnknownItemType = errors.New("неизвестный тип предмета")
)

type ItemError struct {
	ID   int
	Name string
	Err  error
}

func (e *ItemError) Error() string {
	return fmt.Sprintf("предмет %d (%s): %v", e.ID, e.Name, e.Err)
}

func (e *ItemError) Unwrap() error {
	return e.Err
}
