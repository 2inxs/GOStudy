package main

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	if isNilItem(item) {
		return
	}
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) Weapons() []*Weapon {
	weapons := Filter(inv.Items, func(item Item) bool {
		weapon, ok := item.(*Weapon)
		return ok && weapon != nil
	})
	return Map(weapons, func(item Item) *Weapon {
		return item.(*Weapon)
	})
}

func (inv *Inventory) BrokenItems() []Item {
	return Filter(inv.Items, func(item Item) bool {
		switch item := item.(type) {
		case *Weapon:
			return item != nil && item.Durability <= 0
		case *Potion:
			return item != nil && item.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) ItemNames() []string {
	return Map(inv.Items, func(item Item) string {
		if isNilItem(item) {
			return "Предмет отсутствует"
		}
		return item.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(item Item) bool {
		return !isNilItem(item) && item.Name() == name
	})
}
