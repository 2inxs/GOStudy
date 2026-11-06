package inventory_test

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"gostudy/lesson4/internal/inventory"
	"gostudy/lesson4/internal/items"
)

func TestAddItemAfterLoadAndDelete(t *testing.T) {
	const record = "Potion|Мана|+20 MP|2|0.5|301\n"
	var inv inventory.Inventory
	if err := inv.Load(strings.NewReader(record)); err != nil {
		t.Fatal(err)
	}
	armor := items.NewArmor(0, "Щит", 1, 2, items.Common)
	armor.Defense = 4
	if id := inv.AddItem(armor); id != 302 || armor.ID() != id {
		t.Fatalf("AddItem() = %d, armor ID = %d, want 302", id, armor.ID())
	}
	for _, tt := range []struct {
		name string
		item items.Item
	}{
		{"nil", nil},
		{"typed nil", (*items.Weapon)(nil)},
		{"duplicate ID", items.NewPotion(301, "Дубликат", 1, "", 1)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if id := inv.AddItem(tt.item); id != 0 || len(inv.Items) != 2 {
				t.Errorf("rejected AddItem() = %d, item count = %d", id, len(inv.Items))
			}
		})
	}
	if found, ok := inv.FindItemByID(302); !ok || found != armor {
		t.Fatalf("FindItemByID(302) = %v, %t", found, ok)
	}
	if !inv.DeleteItem(302) || inv.DeleteItem(302) {
		t.Fatal("deletion should succeed only on the first call")
	}
	if found, ok := inv.FindItemByID(302); ok || found != nil {
		t.Errorf("deleted item still found: %v, %t", found, ok)
	}
	var saved bytes.Buffer
	if err := inv.Save(&saved); err != nil || saved.String() != record {
		t.Errorf("Save() = %q, %v, want original potion", saved.String(), err)
	}
}

func TestInventoryQueries(t *testing.T) {
	sword := items.NewWeapon(1, "Меч", 1, 2, items.Common)
	sword.Durability = 5
	bow := items.NewWeapon(2, "Лук", 1, 1, items.Common)
	shield := items.NewArmor(3, "Щит", 1, 2, items.Common)
	potion := items.NewPotion(4, "Мана", 0.5, "+20 MP", 0)
	inv := inventory.Inventory{Items: []items.Item{sword, potion, bow, shield, (*items.Weapon)(nil)}}

	if got := inv.Weapons(); !slices.Equal(got, []*items.Weapon{sword, bow}) {
		t.Errorf("Weapons() = %v", got)
	}
	if got := inv.BrokenItems(); !slices.Equal(got, []items.Item{potion, bow, shield}) {
		t.Errorf("BrokenItems() = %v", got)
	}
	if got := inv.ItemNames(); !slices.Equal(got, []string{"Меч", "Мана", "Лук", "Щит", "Предмет отсутствует"}) {
		t.Errorf("ItemNames() = %v", got)
	}
	if got, ok := inv.FindItemByName("Щит"); !ok || got != shield {
		t.Errorf("FindItemByName(Щит) = %v, %t", got, ok)
	}
	if got, ok := inv.FindItemByName("Свиток"); ok || got != nil {
		t.Errorf("FindItemByName(Свиток) = %v, %t", got, ok)
	}
}
