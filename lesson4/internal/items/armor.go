package items

import (
	"fmt"
	"math/rand/v2"
)

type Armor struct {
	id      int
	name    string
	level   int
	weight  float64
	Defense float64
	quality Quality
}

func (a *Armor) ID() int { return a.id }

func (a *Armor) Name() string {
	return a.name
}

func (a *Armor) Weight() float64 { return a.weight }

func (a *Armor) Use() (string, error) {
	if a == nil {
		return "", ErrNilItem
	}
	if a.Defense < 0 {
		return "", &ItemError{ID: a.id, Name: a.name, Err: fmt.Errorf("защита не может быть отрицательной: %w", ErrInvalidItem)}
	}
	return fmt.Sprintf("Используется %s: защита %.2f", a.name, a.Defense), nil
}

func (a *Armor) CalculateDefense() (float64, error) {
	if a == nil {
		return 0, ErrNilItem
	}
	if a.level <= 0 {
		return 0, &ItemError{ID: a.id, Name: a.name, Err: fmt.Errorf("уровень должен быть положительным: %w", ErrInvalidItem)}
	}
	var minDefense, maxDefense int

	switch a.quality {
	case Common:
		minDefense, maxDefense = 3, 8
	case Rare:
		minDefense, maxDefense = 6, 12
	case Epic:
		minDefense, maxDefense = 10, 17
	case Legendary:
		minDefense, maxDefense = 16, 25
	default:
		return 0, &ItemError{ID: a.id, Name: a.name, Err: fmt.Errorf("неизвестное качество %g: %w", a.quality, ErrInvalidItem)}
	}

	baseDefense := rand.IntN(maxDefense-minDefense+1) + minDefense

	return float64(a.level) / 100 * float64(baseDefense) * float64(a.quality), nil
}
