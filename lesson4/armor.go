package main

import (
	"fmt"
	"math/rand/v2"
)

type Armor struct {
	name    string
	level   int
	weight  float64
	Defense float64
	quality Quality
}

func (a *Armor) Name() string {
	return a.name
}

func (a *Armor) Weight() float64 { return a.weight }

func (a *Armor) Use() (string, error) {
	if a == nil {
		return "", fmt.Errorf("броня отсутствует")
	}
	if a.Defense < 0 {
		return "", fmt.Errorf("защита брони %s не может быть отрицательной", a.name)
	}
	return fmt.Sprintf("Используется %s: защита %.2f", a.name, a.Defense), nil
}

func (a *Armor) CalculateDefense() (float64, error) {
	if a == nil {
		return 0, fmt.Errorf("броня отсутствует")
	}
	if a.level <= 0 {
		return 0, fmt.Errorf("уровень брони %s должен быть положительным", a.name)
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
		return 0, fmt.Errorf("неизвестное качество брони %s: %g", a.name, a.quality)
	}

	baseDefense := rand.IntN(maxDefense-minDefense+1) + minDefense

	return float64(a.level) * float64(baseDefense) * float64(a.quality), nil
}
