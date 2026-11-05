package inventory

import (
	"io"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
	"gostudy/lesson4/pkg"
)

type Inventory struct {
	Items []items.Item
}

func (inv *Inventory) AddItem(item items.Item) {
	if items.IsNilItem(item) {
		return
	}
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) Weapons() []*items.Weapon {
	weapons := helpers.Filter(inv.Items, func(item items.Item) bool {
		weapon, ok := item.(*items.Weapon)
		return ok && weapon != nil
	})
	return helpers.Map(weapons, func(item items.Item) *items.Weapon {
		return item.(*items.Weapon)
	})
}

func (inv *Inventory) BrokenItems() []items.Item {
	return helpers.Filter(inv.Items, func(item items.Item) bool {
		switch item := item.(type) {
		case *items.Weapon:
			return item != nil && item.Durability <= 0
		case *items.Potion:
			return item != nil && item.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) ItemNames() []string {
	return helpers.Map(inv.Items, func(item items.Item) string {
		if items.IsNilItem(item) {
			return "Предмет отсутствует"
		}
		return item.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (items.Item, bool) {
	return helpers.Find(inv.Items, func(item items.Item) bool {
		return !items.IsNilItem(item) && item.Name() == name
	})
}

func (inv *Inventory) Save(w io.Writer) error {
	return storage.Save(w, inv.Items)
}

func (inv *Inventory) Load(r io.Reader) error {
	loaded, err := storage.Load(r)
	if err != nil {
		return err
	}
	inv.Items = loaded
	return nil
}
