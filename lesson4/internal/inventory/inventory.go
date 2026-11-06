package inventory

import (
	"fmt"
	"io"
	"math"
	"slices"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
	"gostudy/lesson4/pkg"
)

type Inventory struct {
	itemList []items.Item
}

// Items возвращает копию списка. Сами предметы доступны для использования.
func (inv *Inventory) Items() []items.Item {
	return slices.Clone(inv.itemList)
}

// LoadItems заменяет список только после проверки всех предметов.
func (inv *Inventory) LoadItems(itemList []items.Item) error {
	usedIDs := make(map[int]struct{}, len(itemList))
	for index, item := range itemList {
		if items.IsNilItem(item) {
			return fmt.Errorf("предмет %d: %w", index+1, items.ErrNilItem)
		}
		if err := items.RegisterID(usedIDs, items.IDOf(item)); err != nil {
			return fmt.Errorf("предмет %d: %w", index+1, err)
		}
	}
	inv.itemList = slices.Clone(itemList)
	return nil
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
		for _, existing := range inv.itemList {
			maxID = max(maxID, items.IDOf(existing))
		}
		if maxID == math.MaxInt {
			return 0
		}
		id = maxID + 1
		assignable.SetID(id)
	}

	inv.itemList = append(inv.itemList, item)
	return id
}

func (inv *Inventory) FindItemByID(id int) (items.Item, bool) {
	index := inv.indexByID(id)
	if index < 0 {
		return nil, false
	}
	return inv.itemList[index], true
}

func (inv *Inventory) DeleteItem(id int) bool {
	index := inv.indexByID(id)
	if index < 0 {
		return false
	}
	inv.itemList = slices.Delete(inv.itemList, index, index+1)
	return true
}

func (inv *Inventory) indexByID(id int) int {
	if id <= 0 {
		return -1
	}
	return slices.IndexFunc(inv.itemList, func(item items.Item) bool {
		return items.IDOf(item) == id
	})
}

func (inv *Inventory) Weapons() []*items.Weapon {
	weapons := helpers.Filter(inv.itemList, func(item items.Item) bool {
		weapon, ok := item.(*items.Weapon)
		return ok && weapon != nil
	})
	return helpers.Map(weapons, func(item items.Item) *items.Weapon {
		return item.(*items.Weapon)
	})
}

func (inv *Inventory) BrokenItems() []items.Item {
	return helpers.Filter(inv.itemList, func(item items.Item) bool {
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
	return helpers.Map(inv.itemList, func(item items.Item) string {
		if items.IsNilItem(item) {
			return "Предмет отсутствует"
		}
		return item.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (items.Item, bool) {
	return helpers.Find(inv.itemList, func(item items.Item) bool {
		return !items.IsNilItem(item) && item.Name() == name
	})
}

func (inv *Inventory) Save(w io.Writer) error {
	return storage.Save(w, inv.itemList)
}

func (inv *Inventory) Load(r io.Reader) error {
	loaded, err := storage.Load(r)
	if err != nil {
		return err
	}
	return inv.LoadItems(loaded)
}
