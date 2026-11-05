package items

import (
	"errors"
	"fmt"
)

var (
	ErrNilItem         = errors.New("предмет отсутствует")
	ErrBrokenWeapon    = errors.New("оружие сломано")
	ErrNoCharges       = errors.New("зелье закончилось")
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
