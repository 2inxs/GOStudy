package items_test

import (
	"errors"
	"math"
	"strings"
	"testing"

	"gostudy/lesson4/internal/items"
)

func TestArmorUseChangesDefense(t *testing.T) {
	tests := []struct {
		name    string
		defense float64
		want    float64
		wantErr error
	}{
		{"normal wear", 3, 2, nil},
		{"last fraction", 0.5, 0, nil},
		{"broken", 0, 0, items.ErrBroken},
		{"negative defense", -1, -1, items.ErrInvalidItem},
		{"infinite defense", math.Inf(1), math.Inf(1), items.ErrInvalidItem},
		{"NaN defense", math.NaN(), math.NaN(), items.ErrInvalidItem},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			armor := items.NewArmor(1, "Щит", 1, 2, items.Common)
			armor.Defense = tt.defense
			message, err := armor.Use()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Use() error = %v, want %v", err, tt.wantErr)
			}
			if (message == "") != (tt.wantErr != nil) {
				t.Errorf("message = %q, error = %v", message, err)
			}
			if armor.Defense != tt.want && !(math.IsNaN(armor.Defense) && math.IsNaN(tt.want)) {
				t.Errorf("defense = %v, want %v", armor.Defense, tt.want)
			}
		})
	}
}

func TestUseDangerousItem(t *testing.T) {
	tests := []struct {
		name        string
		itemName    string
		charges     int
		wantCharges int
		wantErr     error
		panicText   string
	}{
		{"normal use", "Мана", 1, 0, nil, ""},
		{"ordinary error", "Мана", 0, 0, items.ErrExhausted, ""},
		{"recovered panic", "Проклятый предмет", 1, 1, nil, "сработало проклятие"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			potion := items.NewPotion(1, tt.itemName, 0.5, "+10 MP", tt.charges)
			message, err := items.UseDangerousItem(potion)
			switch {
			case tt.panicText != "":
				if err == nil || !strings.Contains(err.Error(), tt.panicText) || message != "" {
					t.Errorf("message = %q, error = %v, want recovered panic", message, err)
				}
			default:
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("error = %v, want %v", err, tt.wantErr)
				}
			}
			if potion.Charges != tt.wantCharges {
				t.Errorf("charges = %d, want %d", potion.Charges, tt.wantCharges)
			}
		})
	}
}

func TestUseDangerousItemRecoversPanicInName(t *testing.T) {
	if _, err := items.UseDangerousItem(panickingName{}); err == nil {
		t.Fatal("UseDangerousItem() returned no error after Name panicked")
	}
}

type panickingName struct{}

func (panickingName) Name() string         { panic("name unavailable") }
func (panickingName) Weight() float64      { return 0 }
func (panickingName) Use() (string, error) { return "", nil }
