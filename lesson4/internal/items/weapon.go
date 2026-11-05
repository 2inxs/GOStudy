package items

import (
	"fmt"
	"math/rand/v2"
)

type Weapon struct {
	id         int
	name       string
	level      int
	weight     float64
	Damage     float64
	Durability int
	quality    Quality
}

func (w *Weapon) ID() int { return w.id }

func (w *Weapon) Name() string {
	return w.name
}
func (w *Weapon) Weight() float64 { return w.weight }

func (w *Weapon) Use() (string, error) {
	if w == nil {
		return "", ErrNilItem
	}
	if w.Durability <= 0 {
		return "", &ItemError{ID: w.id, Name: w.name, Err: ErrBrokenWeapon}
	}
	if w.Damage < 0 {
		return "", &ItemError{ID: w.id, Name: w.name, Err: fmt.Errorf("урон не может быть отрицательным: %w", ErrInvalidItem)}
	}
	w.Durability--
	return fmt.Sprintf("Используется %s: урон %.2f, прочность %d", w.name, w.Damage, w.Durability), nil
}

func (w *Weapon) CalculateDamage() (float64, error) {
	if w == nil {
		return 0, ErrNilItem
	}
	if w.level <= 0 {
		return 0, &ItemError{ID: w.id, Name: w.name, Err: fmt.Errorf("уровень должен быть положительным: %w", ErrInvalidItem)}
	}
	var minDamage, maxDamage int

	switch w.quality {
	case Common:
		minDamage, maxDamage = 15, 25
	case Rare:
		minDamage, maxDamage = 20, 30
	case Epic:
		minDamage, maxDamage = 25, 40
	case Legendary:
		minDamage, maxDamage = 35, 50
	default:
		return 0, &ItemError{ID: w.id, Name: w.name, Err: fmt.Errorf("неизвестное качество %g: %w", w.quality, ErrInvalidItem)}
	}

	baseDamage := rand.IntN(maxDamage-minDamage+1) + minDamage

	return float64(w.level) / 100 * float64(w.quality) * float64(baseDamage), nil
}
