package main

import (
	"fmt"
	"math/rand/v2"
)

type Weapon struct {
	name       string
	level      int
	weight     float64
	Damage     float64
	Durability int
	quality    Quality
}

func (w *Weapon) Name() string {
	return w.name
}
func (w *Weapon) Weight() float64 { return w.weight }

func (w *Weapon) Use() (string, error) {
	if w == nil {
		return "", fmt.Errorf("оружие отсутствует")
	}
	if w.Durability <= 0 {
		return "", fmt.Errorf("оружие %s сломано", w.name)
	}
	if w.Damage < 0 {
		return "", fmt.Errorf("урон оружия %s не может быть отрицательным", w.name)
	}
	w.Durability--
	return fmt.Sprintf("Используется %s: урон %.2f, прочность %d", w.name, w.Damage, w.Durability), nil
}

func (w *Weapon) CalculateDamage() (float64, error) {
	if w == nil {
		return 0, fmt.Errorf("оружие отсутствует")
	}
	if w.level <= 0 {
		return 0, fmt.Errorf("уровень оружия %s должен быть положительным", w.name)
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
		return 0, fmt.Errorf("неизвестное качество оружия %s: %g", w.name, w.quality)
	}

	baseDamage := rand.IntN(maxDamage-minDamage+1) + minDamage

	return float64(w.level) * float64(w.quality) * float64(baseDamage), nil
}
