package items_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gostudy/lesson4/internal/items"
)

func TestDeserializeRestoresItemAndRejectsWrongRecord(t *testing.T) {
	for _, tt := range []struct {
		name      string
		item      items.Storable
		valid     string
		wrongType string
	}{
		{"weapon", &items.Weapon{}, "Weapon|Меч|12.5|5|3.5|2|1.2|101\n", "Armor|Щит|5|2|1|1|201\n"},
		{"armor", &items.Armor{}, "Armor|Щит|5|2|1|1|201\n", "Potion|Мана|+20 MP|2|0.5|301\n"},
		{"potion", &items.Potion{}, "Potion|Мана|+20 MP|2|0.5|301\n", "Weapon|Меч|12.5|5|3.5|2|1.2|101\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Deserialize(strings.NewReader(tt.valid)); err != nil {
				t.Fatal(err)
			}
			for _, invalid := range []string{tt.wrongType, "Invalid\n"} {
				if err := tt.item.Deserialize(strings.NewReader(invalid)); !errors.Is(err, items.ErrInvalidRecord) {
					t.Fatalf("Deserialize(%q) = %v, want ErrInvalidRecord", invalid, err)
				}
			}
			var saved bytes.Buffer
			if err := tt.item.Serialize(&saved); err != nil || saved.String() != tt.valid {
				t.Errorf("Serialize() = %q, %v, want %q", saved.String(), err, tt.valid)
			}
		})
	}
}

func TestDescribeItemHandlesNilAndRegularItems(t *testing.T) {
	for _, tt := range []struct {
		name string
		item items.Item
		want string
	}{
		{"nil interface", nil, "Предмет отсутствует"},
		{"typed nil", (*items.Weapon)(nil), "Предмет отсутствует"},
		{"regular item", items.NewPotion(3, "Мана", 0.5, "+20 MP", 2), "3: Мана (вес: 0.5)"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := items.DescribeItem(tt.item); got != tt.want {
				t.Errorf("DescribeItem() = %q, want %q", got, tt.want)
			}
		})
	}
}
