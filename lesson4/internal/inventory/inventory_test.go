package inventory_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"gostudy/lesson4/internal/inventory"
	"gostudy/lesson4/internal/items"
)

func TestLoadReplacesInventory(t *testing.T) {
	original := items.NewWeapon(101, "Старый меч", 1, 3, items.Common)
	tests := []struct {
		name  string
		input string
		want  []items.Item
	}{
		{
			name:  "replace existing items",
			input: "Potion|Лечебное|+50 HP|3|0.5|301\n",
			want:  []items.Item{items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)},
		},
		{name: "empty input clears inventory", want: []items.Item{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inv inventory.Inventory
			inv.AddItem(original)
			if err := inv.Load(strings.NewReader(tt.input)); err != nil {
				t.Fatalf("Load() returned an unexpected error: %v", err)
			}
			itemList := inv.Items()
			if len(itemList) != len(tt.want) {
				t.Fatalf("inventory after Load() has %d items, want %d", len(itemList), len(tt.want))
			}
			for i, item := range itemList {
				if !reflect.DeepEqual(item, tt.want[i]) {
					t.Errorf("item %d = %#v, want %#v", i, item, tt.want[i])
				}
			}
		})
	}
}

func TestFailedLoadPreservesInventory(t *testing.T) {
	readErr := errors.New("connection lost")
	const validRecord = "Potion|Лечебное|+50 HP|3|0.5|301\n"
	tests := []struct {
		name   string
		reader io.Reader
		want   error
	}{
		{
			name:   "corruption after a valid record",
			reader: strings.NewReader(validRecord + "Weapon|Меч\n"),
			want:   items.ErrInvalidRecord,
		},
		{
			name:   "reader failure after a valid record",
			reader: io.MultiReader(strings.NewReader(validRecord), iotest.ErrReader(readErr)),
			want:   readErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weapon := items.NewWeapon(101, "Меч", 2, 3, items.Rare)
			weapon.Damage, weapon.Durability = 12.5, 5
			potion := items.NewPotion(302, "Мана", 0.5, "+20 MP", 2)
			var inv inventory.Inventory
			inv.AddItem(weapon)
			inv.AddItem(potion)
			wantWeapon, wantPotion := *weapon, *potion

			err := inv.Load(tt.reader)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Load() error = %v, want %v", err, tt.want)
			}
			itemList := inv.Items()
			if len(itemList) != 2 {
				t.Fatalf("failed Load() changed inventory length to %d, want 2", len(itemList))
			}
			if itemList[0] != weapon || itemList[1] != potion {
				t.Errorf("failed Load() replaced or reordered existing items: %#v", itemList)
			}
			if *weapon != wantWeapon || *potion != wantPotion {
				t.Error("failed Load() modified existing items")
			}
		})
	}
}
