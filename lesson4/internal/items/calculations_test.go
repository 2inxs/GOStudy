package items_test

import (
	"testing"

	"gostudy/lesson4/internal/items"
)

func TestCalculateStatsForLevelsBelow100(t *testing.T) {
	tests := []struct {
		name      string
		calculate func() (float64, error)
		min       float64
		max       float64
	}{
		{
			name:      "weapon level 1",
			calculate: items.NewWeapon(101, "Меч", 1, 3, items.Common).CalculateDamage,
			min:       0.15,
			max:       0.25,
		},
		{
			name:      "weapon level 99",
			calculate: items.NewWeapon(101, "Меч", 99, 3, items.Common).CalculateDamage,
			min:       14.85,
			max:       24.75,
		},
		{
			name:      "armor level 1",
			calculate: items.NewArmor(201, "Щит", 1, 4.5, items.Common).CalculateDefense,
			min:       0.03,
			max:       0.08,
		},
		{
			name:      "armor level 99",
			calculate: items.NewArmor(201, "Щит", 99, 4.5, items.Common).CalculateDefense,
			min:       2.97,
			max:       7.92,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.calculate()
			if err != nil {
				t.Fatalf("calculation returned an unexpected error: %v", err)
			}
			const tolerance = 1e-12
			if !(got >= tt.min-tolerance && got <= tt.max+tolerance) {
				t.Errorf("calculated stat = %g, want a value in [%g, %g]", got, tt.min, tt.max)
			}
		})
	}
}
