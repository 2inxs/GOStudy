package items_test

import (
	"errors"
	"testing"

	"gostudy/lesson4/internal/items"
)

func TestWeaponUseConsumesLastDurability(t *testing.T) {
	weapon := items.NewWeapon(101, "Меч", 1, 3, items.Common)
	weapon.Damage, weapon.Durability = 12.5, 1

	message, err := weapon.Use()
	if err != nil {
		t.Fatalf("last use returned an unexpected error: %v", err)
	}
	if want := "Используется Меч: урон 12.50, прочность 0"; message != want {
		t.Errorf("Use() message = %q, want %q", message, want)
	}
	if weapon.Durability != 0 {
		t.Fatalf("durability after last use = %d, want 0", weapon.Durability)
	}

	message, err = weapon.Use()
	if !errors.Is(err, items.ErrBrokenWeapon) {
		t.Fatalf("Use() on a broken weapon returned %v, want ErrBrokenWeapon", err)
	}
	if message != "" {
		t.Errorf("failed Use() returned a success message: %q", message)
	}
	if weapon.Durability != 0 {
		t.Errorf("failed Use() changed durability to %d", weapon.Durability)
	}
	var itemErr *items.ItemError
	if !errors.As(err, &itemErr) {
		t.Fatalf("Use() error has no ItemError: %v", err)
	}
	if itemErr.ID != 101 || itemErr.Name != "Меч" {
		t.Errorf("ItemError identifies (%d, %q), want (101, Меч)", itemErr.ID, itemErr.Name)
	}
}

func TestPotionUseConsumesLastCharge(t *testing.T) {
	potion := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 1)

	message, err := potion.Use()
	if err != nil {
		t.Fatalf("last use returned an unexpected error: %v", err)
	}
	if want := "Используется Лечебное: +50 HP, осталось зарядов 0"; message != want {
		t.Errorf("Use() message = %q, want %q", message, want)
	}
	if potion.Charges != 0 {
		t.Fatalf("charges after last use = %d, want 0", potion.Charges)
	}

	message, err = potion.Use()
	if !errors.Is(err, items.ErrNoCharges) {
		t.Fatalf("Use() on an empty potion returned %v, want ErrNoCharges", err)
	}
	if message != "" {
		t.Errorf("failed Use() returned a success message: %q", message)
	}
	if potion.Charges != 0 {
		t.Errorf("failed Use() changed charges to %d", potion.Charges)
	}
	var itemErr *items.ItemError
	if !errors.As(err, &itemErr) {
		t.Fatalf("Use() error has no ItemError: %v", err)
	}
	if itemErr.ID != 301 || itemErr.Name != "Лечебное" {
		t.Errorf("ItemError identifies (%d, %q), want (301, Лечебное)", itemErr.ID, itemErr.Name)
	}
}

func TestInvalidWeaponUsePreservesState(t *testing.T) {
	weapon := items.NewWeapon(101, "Меч", 1, 3, items.Common)
	weapon.Damage, weapon.Durability = -1, 2
	before := *weapon

	message, err := weapon.Use()
	if !errors.Is(err, items.ErrInvalidItem) {
		t.Fatalf("Use() error = %v, want ErrInvalidItem", err)
	}
	if message != "" {
		t.Errorf("failed Use() returned a success message: %q", message)
	}
	if *weapon != before {
		t.Error("failed Use() changed the weapon")
	}
}

func TestInvalidPotionUsePreservesState(t *testing.T) {
	potion := items.NewPotion(301, "Лечебное", 0.5, "", 2)
	before := *potion

	message, err := potion.Use()
	if !errors.Is(err, items.ErrInvalidItem) {
		t.Fatalf("Use() error = %v, want ErrInvalidItem", err)
	}
	if message != "" {
		t.Errorf("failed Use() returned a success message: %q", message)
	}
	if *potion != before {
		t.Error("failed Use() changed the potion")
	}
}

func TestNilItemUse(t *testing.T) {
	var weapon *items.Weapon
	var armor *items.Armor
	var potion *items.Potion
	for _, tt := range []struct {
		name string
		item items.Usable
	}{
		{name: "weapon", item: weapon},
		{name: "armor", item: armor},
		{name: "potion", item: potion},
	} {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tt.item.Use()
			if !errors.Is(err, items.ErrNilItem) {
				t.Fatalf("Use() error = %v, want ErrNilItem", err)
			}
			if message != "" {
				t.Errorf("failed Use() returned a success message: %q", message)
			}
		})
	}
}
