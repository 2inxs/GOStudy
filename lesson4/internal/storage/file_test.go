package storage_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
)

func TestSave(t *testing.T) {
	weapon := items.NewWeapon(101, "Меч", 2, 3.5, items.Rare)
	weapon.Damage, weapon.Durability = 12.5, 5
	armor := items.NewArmor(201, "Кольчуга", 3, 12.6, items.Epic)
	armor.Defense = 8.5
	potion := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)

	tests := []struct {
		name  string
		items []items.Item
		want  string
	}{
		{name: "empty inventory"},
		{
			name:  "all item types in inventory order",
			items: []items.Item{potion, weapon, armor},
			want: "Potion|Лечебное|+50 HP|3|0.5|301\n" +
				"Weapon|Меч|12.5|5|3.5|2|1.2|101\n" +
				"Armor|Кольчуга|8.5|12.6|3|1.5|201\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := storage.Save(&output, tt.items); err != nil {
				t.Fatalf("Save() returned an unexpected error: %v", err)
			}
			if got := output.String(); got != tt.want {
				t.Errorf("Save() wrote %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSaveStopsAtInvalidItem(t *testing.T) {
	first := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)
	last := items.NewPotion(302, "Мана", 0.5, "+20 MP", 1)
	var nilWeapon *items.Weapon
	var nilArmor *items.Armor
	var nilPotion *items.Potion

	tests := []struct {
		name string
		item items.Item
		want error
	}{
		{name: "nil interface", want: items.ErrNilItem},
		{name: "nil weapon", item: nilWeapon, want: items.ErrNilItem},
		{name: "nil armor", item: nilArmor, want: items.ErrNilItem},
		{name: "nil potion", item: nilPotion, want: items.ErrNilItem},
		{name: "item without serialization", item: unstorableItem{}, want: storage.ErrNotStorable},
		{name: "delimiter in name", item: items.NewWeapon(101, "Меч|лук", 1, 3, items.Common), want: items.ErrInvalidItem},
		{name: "newline in effect", item: items.NewPotion(303, "Яд", 1, "урон\n10", 1), want: items.ErrInvalidItem},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := storage.Save(&output, []items.Item{first, tt.item, last})
			if !errors.Is(err, tt.want) {
				t.Fatalf("Save() error = %v, want %v", err, tt.want)
			}
			if got, want := output.String(), "Potion|Лечебное|+50 HP|3|0.5|301\n"; got != want {
				t.Errorf("Save() wrote %q, want only the record before the failure: %q", got, want)
			}
			if errors.Is(err, items.ErrNilItem) {
				return
			}
			var itemErr *items.ItemError
			if !errors.As(err, &itemErr) {
				t.Fatalf("Save() error has no ItemError: %v", err)
			}
			if itemErr.ID != tt.item.ID() || itemErr.Name != tt.item.Name() {
				t.Errorf("ItemError identifies (%d, %q), want (%d, %q)", itemErr.ID, itemErr.Name, tt.item.ID(), tt.item.Name())
			}
		})
	}
}

func TestSaveWriterFailure(t *testing.T) {
	writeErr := errors.New("disk full")
	writer := &limitedWriter{remaining: len("Potion|Лечебное|+50 HP|3|0.5|301\n") + 4, err: writeErr}
	first := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)
	second := items.NewWeapon(101, "Меч", 1, 3, items.Common)
	third := items.NewArmor(201, "Щит", 1, 4.5, items.Common)

	err := storage.Save(writer, []items.Item{first, second, third})
	if !errors.Is(err, writeErr) {
		t.Fatalf("Save() error = %v, want the original writer error", err)
	}
	var itemErr *items.ItemError
	if !errors.As(err, &itemErr) {
		t.Fatalf("Save() error has no ItemError: %v", err)
	}
	if itemErr.ID != 101 || itemErr.Name != "Меч" {
		t.Errorf("ItemError identifies (%d, %q), want (101, Меч)", itemErr.ID, itemErr.Name)
	}
	if got, want := writer.output.String(), "Potion|Лечебное|+50 HP|3|0.5|301\nWeap"; got != want {
		t.Errorf("written data = %q, want %q", got, want)
	}
	if writer.calls != 2 {
		t.Errorf("writer was called %d times; writing must stop at the failed item", writer.calls)
	}
}

func TestLoad(t *testing.T) {
	weapon := items.NewWeapon(101, "Меч", 2, 3.5, items.Rare)
	weapon.Damage, weapon.Durability = 12.5, 5
	armor := items.NewArmor(201, "Кольчуга", 3, 12.6, items.Epic)
	armor.Defense = 8.5
	potion := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)

	tests := []struct {
		name  string
		input string
		want  []items.Item
	}{
		{name: "empty input", want: []items.Item{}},
		{
			name: "all fields and item order",
			input: "Armor|Кольчуга|8.5|12.6|3|1.5|201\n" +
				"Weapon|Меч|12.5|5|3.5|2|1.2|101\n" +
				"Potion|Лечебное|+50 HP|3|0.5|301\n",
			want: []items.Item{armor, weapon, potion},
		},
		{
			name:  "final record without newline",
			input: "Potion|Лечебное|+50 HP|3|0.5|301",
			want:  []items.Item{potion},
		},
		{
			name:  "windows line endings",
			input: "Weapon|Меч|12.5|5|3.5|2|1.2|101\r\nPotion|Лечебное|+50 HP|3|0.5|301\r\n",
			want:  []items.Item{weapon, potion},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.Load(iotest.OneByteReader(strings.NewReader(tt.input)))
			if err != nil {
				t.Fatalf("Load() returned an unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Load() returned %d items, want %d", len(got), len(tt.want))
			}
			for i, item := range got {
				if !reflect.DeepEqual(item, tt.want[i]) {
					t.Errorf("item %d = %#v, want %#v", i, item, tt.want[i])
				}
			}
		})
	}
}

func TestLoadRejectsCorruptRecords(t *testing.T) {
	const validRecord = "Potion|Лечебное|+50 HP|3|0.5|301\n"
	tests := []struct {
		name   string
		record string
		cause  error
	}{
		{name: "blank line", record: "\n"},
		{name: "missing fields", record: "Weapon|Меч\n"},
		{name: "extra field", record: "Potion|Лечебное|+50 HP|3|0.5|301|extra\n"},
		{name: "unknown item type", record: "Scroll|Свиток|401\n", cause: items.ErrUnknownItemType},
		{name: "invalid ID", record: "Weapon|Меч|10|5|3|1|1|bad-id\n", cause: strconv.ErrSyntax},
		{name: "invalid number", record: "Armor|Щит|bad-defense|4.5|1|1|201\n", cause: strconv.ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.Load(strings.NewReader(validRecord + tt.record + validRecord))
			if !errors.Is(err, items.ErrInvalidRecord) {
				t.Fatalf("Load() error = %v, want ErrInvalidRecord", err)
			}
			if got != nil {
				t.Errorf("Load() returned partial data after a corrupt record: %#v", got)
			}
			var recordErr *storage.RecordError
			if !errors.As(err, &recordErr) {
				t.Fatalf("Load() error has no RecordError: %v", err)
			}
			if recordErr.Line != 2 {
				t.Errorf("RecordError.Line = %d, want 2", recordErr.Line)
			}
			if tt.cause != nil && !errors.Is(err, tt.cause) {
				t.Errorf("Load() error lost its cause %v: %v", tt.cause, err)
			}
		})
	}
}

func TestLoadReaderFailure(t *testing.T) {
	readErr := errors.New("read failed")
	reader := io.MultiReader(
		strings.NewReader("Potion|Лечебное|+50 HP|3|0.5|301\n"),
		iotest.ErrReader(readErr),
	)
	got, err := storage.Load(reader)
	if !errors.Is(err, readErr) {
		t.Fatalf("Load() error = %v, want the original reader error", err)
	}
	if got != nil {
		t.Errorf("Load() returned partial data after a reader failure: %#v", got)
	}
	var recordErr *storage.RecordError
	if !errors.As(err, &recordErr) {
		t.Fatalf("Load() error has no RecordError: %v", err)
	}
	if recordErr.Line != 2 {
		t.Errorf("RecordError.Line = %d, want 2", recordErr.Line)
	}
}

func TestLoadRecordExceedsScannerLimit(t *testing.T) {
	record := "Potion|" + strings.Repeat("x", 128*1024) + "|+50 HP|3|0.5|301\n"
	got, err := storage.Load(strings.NewReader(record))
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("Load() error = %v, want bufio.ErrTooLong", err)
	}
	if got != nil {
		t.Errorf("Load() returned data for an unreadable record: %#v", got)
	}
	var recordErr *storage.RecordError
	if !errors.As(err, &recordErr) {
		t.Fatalf("Load() error has no RecordError: %v", err)
	}
	if recordErr.Line != 1 {
		t.Errorf("RecordError.Line = %d, want 1", recordErr.Line)
	}
}

type unstorableItem struct{}

func (unstorableItem) ID() int         { return 401 }
func (unstorableItem) Name() string    { return "Свиток" }
func (unstorableItem) Weight() float64 { return 0.1 }

type limitedWriter struct {
	output    bytes.Buffer
	remaining int
	err       error
	calls     int
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	w.calls++
	if len(p) <= w.remaining {
		w.remaining -= len(p)
		return w.output.Write(p)
	}
	n, _ := w.output.Write(p[:w.remaining])
	w.remaining = 0
	return n, w.err
}
