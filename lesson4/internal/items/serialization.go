package items

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (w *Weapon) Serialize(writer io.Writer) error {
	if w == nil {
		return ErrNilItem
	}
	if strings.ContainsAny(w.name, "|\r\n") {
		return fmt.Errorf("название оружия содержит разделитель или перенос строки: %w", ErrInvalidItem)
	}
	_, err := fmt.Fprintf(writer, "Weapon|%s|%g|%d|%g|%d|%g|%d\n",
		w.name, w.Damage, w.Durability, w.weight, w.level, w.quality, w.id)
	return err
}

func (a *Armor) Serialize(writer io.Writer) error {
	if a == nil {
		return ErrNilItem
	}
	if strings.ContainsAny(a.name, "|\r\n") {
		return fmt.Errorf("название брони содержит разделитель или перенос строки: %w", ErrInvalidItem)
	}
	_, err := fmt.Fprintf(writer, "Armor|%s|%g|%g|%d|%g|%d\n",
		a.name, a.Defense, a.weight, a.level, a.quality, a.id)
	return err
}

func (p *Potion) Serialize(writer io.Writer) error {
	if p == nil {
		return ErrNilItem
	}
	if strings.ContainsAny(p.name, "|\r\n") || strings.ContainsAny(p.Effect, "|\r\n") {
		return fmt.Errorf("название или эффект зелья содержит разделитель или перенос строки: %w", ErrInvalidItem)
	}
	_, err := fmt.Fprintf(writer, "Potion|%s|%s|%d|%g|%d\n",
		p.name, p.Effect, p.Charges, p.weight, p.id)
	return err
}

func parseRecord(fields []string) (Item, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("пустая запись")
	}

	switch fields[0] {
	case "Weapon":
		if len(fields) != 8 {
			return nil, fmt.Errorf("для Weapon нужно 8 полей")
		}
		id, err := strconv.Atoi(fields[7])
		if err != nil {
			return nil, fmt.Errorf("неверный ID оружия: %w", err)
		}
		damage, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("неверный урон: %w", err)
		}
		durability, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("неверная прочность: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("неверный вес оружия: %w", err)
		}
		level, err := strconv.Atoi(fields[5])
		if err != nil {
			return nil, fmt.Errorf("неверный уровень оружия: %w", err)
		}
		quality, err := strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return nil, fmt.Errorf("неверное качество оружия: %w", err)
		}
		return &Weapon{
			id: id, name: fields[1], Damage: damage, Durability: durability,
			weight: weight, level: level, quality: Quality(quality),
		}, nil

	case "Armor":
		if len(fields) != 7 {
			return nil, fmt.Errorf("для Armor нужно 7 полей")
		}
		id, err := strconv.Atoi(fields[6])
		if err != nil {
			return nil, fmt.Errorf("неверный ID брони: %w", err)
		}
		defense, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("неверная защита: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, fmt.Errorf("неверный вес брони: %w", err)
		}
		level, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, fmt.Errorf("неверный уровень брони: %w", err)
		}
		quality, err := strconv.ParseFloat(fields[5], 64)
		if err != nil {
			return nil, fmt.Errorf("неверное качество брони: %w", err)
		}
		return &Armor{
			id: id, name: fields[1], Defense: defense, weight: weight,
			level: level, quality: Quality(quality),
		}, nil

	case "Potion":
		if len(fields) != 6 {
			return nil, fmt.Errorf("для Potion нужно 6 полей")
		}
		id, err := strconv.Atoi(fields[5])
		if err != nil {
			return nil, fmt.Errorf("неверный ID зелья: %w", err)
		}
		charges, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("неверное число зарядов: %w", err)
		}
		weight, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("неверный вес зелья: %w", err)
		}
		return &Potion{
			id: id, name: fields[1], Effect: fields[2], Charges: charges, weight: weight,
		}, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownItemType, fields[0])
	}
}

func ReadItem(r io.Reader) (Item, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	line := strings.TrimSuffix(string(data), "\n")
	line = strings.TrimSuffix(line, "\r")
	if strings.ContainsAny(line, "\r\n") {
		return nil, fmt.Errorf("ожидалась одна строка с предметом: %w", ErrInvalidRecord)
	}
	item, err := parseRecord(strings.Split(line, "|"))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRecord, err)
	}
	return item, nil
}

func (w *Weapon) Deserialize(r io.Reader) error {
	if w == nil {
		return ErrNilItem
	}
	item, err := ReadItem(r)
	if err != nil {
		return err
	}
	weapon, ok := item.(*Weapon)
	if !ok {
		return fmt.Errorf("ожидался Weapon: %w", ErrInvalidRecord)
	}
	*w = *weapon
	return nil
}

func (a *Armor) Deserialize(r io.Reader) error {
	if a == nil {
		return ErrNilItem
	}
	item, err := ReadItem(r)
	if err != nil {
		return err
	}
	armor, ok := item.(*Armor)
	if !ok {
		return fmt.Errorf("ожидался Armor: %w", ErrInvalidRecord)
	}
	*a = *armor
	return nil
}

func (p *Potion) Deserialize(r io.Reader) error {
	if p == nil {
		return ErrNilItem
	}
	item, err := ReadItem(r)
	if err != nil {
		return err
	}
	potion, ok := item.(*Potion)
	if !ok {
		return fmt.Errorf("ожидался Potion: %w", ErrInvalidRecord)
	}
	*p = *potion
	return nil
}
