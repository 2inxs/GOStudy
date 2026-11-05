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
			inv := &inventory.Inventory{Items: []items.Item{original}}
			if err := inv.Load(strings.NewReader(tt.input)); err != nil {
				t.Fatalf("Load() returned an unexpected error: %v", err)
			}
			if len(inv.Items) != len(tt.want) {
				t.Fatalf("inventory after Load() has %d items, want %d", len(inv.Items), len(tt.want))
			}
			for i, item := range inv.Items {
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
			inv := &inventory.Inventory{Items: []items.Item{weapon, potion}}
			wantWeapon, wantPotion := *weapon, *potion

			err := inv.Load(tt.reader)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Load() error = %v, want %v", err, tt.want)
			}
			if len(inv.Items) != 2 {
				t.Fatalf("failed Load() changed inventory length to %d, want 2", len(inv.Items))
			}
			if inv.Items[0] != weapon || inv.Items[1] != potion {
				t.Errorf("failed Load() replaced or reordered existing items: %#v", inv.Items)
			}
			if *weapon != wantWeapon || *potion != wantPotion {
				t.Error("failed Load() modified existing items")
			}
		})
	}
}
