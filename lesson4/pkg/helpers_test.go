package helpers

import (
	"slices"
	"strconv"
	"testing"
)

func TestFilter(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  []int
	}{
		{"keeps order and duplicates", []int{4, 1, 2, 4}, []int{4, 2, 4}},
		{"no matches", []int{1, 3}, nil},
		{"empty input", nil, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, func(n int) bool { return n%2 == 0 })
			if !slices.Equal(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  []string
	}{
		{"converts each value", []int{4, 0, -2}, []string{"4", "0", "-2"}},
		{"empty input", nil, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.input, strconv.Itoa); !slices.Equal(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  int
		found bool
	}{
		{"first of two matches", []int{3, 4, 6}, 4, true},
		{"zero is a valid match", []int{3, 0}, 0, true},
		{"no match", []int{1, 3}, 0, false},
		{"empty input", nil, 0, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, found := Find(tt.input, func(n int) bool { return n%2 == 0 })
			if got != tt.want || found != tt.found {
				t.Errorf("Find() = %d, %t, want %d, %t", got, found, tt.want, tt.found)
			}
		})
	}
}
