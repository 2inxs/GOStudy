package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (w *Weapon) Serialize(writer io.Writer) error {
	if w == nil {
		return fmt.Errorf("оружие отсутствует")
	}
	if strings.ContainsAny(w.name, "|\r\n") {
		return fmt.Errorf("название оружия содержит разделитель или перенос строки")
	}
	_, err := fmt.Fprintf(writer, "Weapon|%s|%g|%d|%g|%d|%g\n",
		w.name, w.Damage, w.Durability, w.weight, w.level, w.quality)
	return err
}

func (a *Armor) Serialize(writer io.Writer) error {
	if a == nil {
		return fmt.Errorf("броня отсутствует")
	}
	if strings.ContainsAny(a.name, "|\r\n") {
		return fmt.Errorf("название брони содержит разделитель или перенос строки")
	}
	_, err := fmt.Fprintf(writer, "Armor|%s|%g|%g|%d|%g\n",
		a.name, a.Defense, a.weight, a.level, a.quality)
	return err
}

func (p *Potion) Serialize(writer io.Writer) error {
	if p == nil {
		return fmt.Errorf("зелье отсутствует")
	}
	if strings.ContainsAny(p.name, "|\r\n") || strings.ContainsAny(p.Effect, "|\r\n") {
		return fmt.Errorf("название или эффект зелья содержит разделитель или перенос строки")
	}
	_, err := fmt.Fprintf(writer, "Potion|%s|%s|%d|%g\n",
		p.name, p.Effect, p.Charges, p.weight)
	return err
}

func parseRecord(fields []string) (Item, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("пустая запись")
	}

	switch fields[0] {
	case "Weapon":
		if len(fields) != 7 {
			return nil, fmt.Errorf("для Weapon нужно 7 полей")
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
			name: fields[1], Damage: damage, Durability: durability,
			weight: weight, level: level, quality: Quality(quality),
		}, nil

	case "Armor":
		if len(fields) != 6 {
			return nil, fmt.Errorf("для Armor нужно 6 полей")
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
			name: fields[1], Defense: defense, weight: weight,
			level: level, quality: Quality(quality),
		}, nil

	case "Potion":
		if len(fields) != 5 {
			return nil, fmt.Errorf("для Potion нужно 5 полей")
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
			name: fields[1], Effect: fields[2], Charges: charges, weight: weight,
		}, nil

	default:
		return nil, fmt.Errorf("неизвестный тип предмета: %s", fields[0])
	}
}

func readItem(r io.Reader) (Item, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	line := strings.TrimSuffix(string(data), "\n")
	line = strings.TrimSuffix(line, "\r")
	if strings.ContainsAny(line, "\r\n") {
		return nil, fmt.Errorf("ожидалась одна строка с предметом")
	}
	return parseRecord(strings.Split(line, "|"))
}

func (w *Weapon) Deserialize(r io.Reader) error {
	if w == nil {
		return fmt.Errorf("оружие отсутствует")
	}
	item, err := readItem(r)
	if err != nil {
		return err
	}
	weapon, ok := item.(*Weapon)
	if !ok {
		return fmt.Errorf("ожидался Weapon")
	}
	*w = *weapon
	return nil
}

func (a *Armor) Deserialize(r io.Reader) error {
	if a == nil {
		return fmt.Errorf("броня отсутствует")
	}
	item, err := readItem(r)
	if err != nil {
		return err
	}
	armor, ok := item.(*Armor)
	if !ok {
		return fmt.Errorf("ожидался Armor")
	}
	*a = *armor
	return nil
}

func (p *Potion) Deserialize(r io.Reader) error {
	if p == nil {
		return fmt.Errorf("зелье отсутствует")
	}
	item, err := readItem(r)
	if err != nil {
		return err
	}
	potion, ok := item.(*Potion)
	if !ok {
		return fmt.Errorf("ожидался Potion")
	}
	*p = *potion
	return nil
}

func (inv *Inventory) Save(w io.Writer) error {
	for _, item := range inv.Items {
		if isNilItem(item) {
			return fmt.Errorf("нельзя сохранить отсутствующий предмет")
		}
		storable, ok := item.(Storable)
		if !ok {
			return fmt.Errorf("предмет %s не поддерживает сохранение", item.Name())
		}
		if err := storable.Serialize(w); err != nil {
			return err
		}
	}
	return nil
}

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	items := make([]Item, 0)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "|")
		item, err := parseRecord(fields)
		if err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	inv.Items = items
	return nil
}
