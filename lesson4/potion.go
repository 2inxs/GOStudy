package main

import "fmt"

type Potion struct {
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string {
	return p.name
}

func (p *Potion) Weight() float64 { return p.weight }

func (p *Potion) Use() (string, error) {
	if p == nil {
		return "", fmt.Errorf("зелье отсутствует")
	}
	if p.Charges <= 0 {
		return "", fmt.Errorf("зелье %s закончилось", p.name)
	}
	if p.Effect == "" {
		return "", fmt.Errorf("у зелья %s не задан эффект", p.name)
	}
	p.Charges--
	return fmt.Sprintf("Используется %s: %s, осталось зарядов %d", p.name, p.Effect, p.Charges), nil
}
