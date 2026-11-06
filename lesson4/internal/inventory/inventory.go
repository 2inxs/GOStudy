package inventory

import (
	"io"
	"math"
	"slices"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
	"gostudy/lesson4/pkg"
)

type Inventory struct {
	Items []items.Item
}

func (inv *Inventory) AddItem(item items.Item) int {
	if items.IsNilItem(item) {
		return 0
	}

	id := items.IDOf(item)
	if id > 0 && inv.indexByID(id) >= 0 {
		return 0
	}
	if id <= 0 {
		assignable, ok := item.(interface {
			ID() int
			SetID(int)
		})
		if !ok {
			return 0
		}
		maxID := 0
		for _, existing := range inv.Items {
			maxID = max(maxID, items.IDOf(existing))
		}
		if maxID == math.MaxInt {
			return 0
		}
		id = maxID + 1
		assignable.SetID(id)
	}

	inv.Items = append(inv.Items, item)
	return id
}

func (inv *Inventory) FindItemByID(id int) (items.Item, bool) {
	index := inv.indexByID(id)
	if index < 0 {
		return nil, false
	}
	return inv.Items[index], true
}

func (inv *Inventory) DeleteItem(id int) bool {
	index := inv.indexByID(id)
	if index < 0 {
		return false
	}
	inv.Items = slices.Delete(inv.Items, index, index+1)
	return true
}

func (inv *Inventory) indexByID(id int) int {
	if id <= 0 {
		return -1
	}
	return slices.IndexFunc(inv.Items, func(item items.Item) bool {
		return items.IDOf(item) == id
	})
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
		case *items.Armor:
			return item != nil && item.Defense == 0
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
