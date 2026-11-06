package items

import "fmt"

type Potion struct {
	id      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) ID() int { return p.id }

func (p *Potion) Name() string {
	return p.name
}

func (p *Potion) Weight() float64 { return p.weight }

func (p *Potion) Use() (string, error) {
	if p == nil {
		return "", ErrNilItem
	}
	if p.Charges <= 0 {
		return "", &ItemError{ID: p.id, Name: p.name, Err: ErrExhausted}
	}
	if p.Effect == "" {
		return "", &ItemError{ID: p.id, Name: p.name, Err: fmt.Errorf("не задан эффект: %w", ErrInvalidItem)}
	}
	p.Charges--
	return fmt.Sprintf("Используется %s: %s, осталось зарядов %d", p.name, p.Effect, p.Charges), nil
}
