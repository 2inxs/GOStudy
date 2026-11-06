package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
)

func TestCommandSequencePersistsInventory(t *testing.T) {
	runner := commandRunner{store: storage.FileStore{}}
	path := filepath.Join(t.TempDir(), "inventory.txt")
	steps := []struct {
		name string
		args []string
		out  string
		file string
	}{
		{"missing file", []string{"list"}, "", ""},
		{"add", []string{"add", "-name", "Мана", "-weight", "0.5"}, "ID 1", "Potion|Мана|без эффекта|1|0.5|1\n"},
		{"list", []string{"list"}, "1: Мана (вес: 0.5)", "Potion|Мана|без эффекта|1|0.5|1\n"},
		{"use", []string{"use", "-id", "1"}, "осталось зарядов 0", "Potion|Мана|без эффекта|0|0.5|1\n"},
		{"delete", []string{"delete", "-id", "1"}, "Удалён предмет", ""},
		{"empty file", []string{"list"}, "", ""},
	}
	for i, step := range steps {
		if !t.Run(step.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runner.run(append(step.args, "-file", path), &stdout, &stderr)
			if code != 0 || stderr.Len() != 0 {
				t.Fatalf("code = %d, stderr = %q", code, stderr.String())
			}
			if step.out == "" && stdout.Len() != 0 || !strings.Contains(stdout.String(), step.out) {
				t.Errorf("stdout = %q, want %q", stdout.String(), step.out)
			}
			got, err := os.ReadFile(path)
			if i == 0 {
				if !errors.Is(err, os.ErrNotExist) {
					t.Errorf("list on missing file: %v", err)
				}
			} else if err != nil || string(got) != step.file {
				t.Errorf("file = %q, error = %v, want %q", got, err, step.file)
			}
		}) {
			return
		}
	}
}

func TestCommandErrorCodesAndStreams(t *testing.T) {
	runner := commandRunner{store: storage.FileStore{}}
	tests := []struct {
		name string
		args []string
		file string
		code int
		want string
	}{
		{"unknown command", []string{"unknown"}, "", 2, "неизвестная команда"},
		{"unknown flag", []string{"list", "-id", "1"}, "", 2, "flag provided but not defined"},
		{"invalid weight", []string{"add", "-name", "Мана", "-weight", "NaN"}, "", 2, "положительное"},
		{"missing ID", []string{"use"}, "", 2, "положительный -id"},
		{"missing item", []string{"delete", "-id", "1"}, "", 1, "не найден"},
		{"malformed line", []string{"list"}, "Potion|Мана|+5 MP|1|0.5|1\nWeapon|Меч\n", 1, `строка 2 ("Weapon|Меч")`},
		{"duplicate ID in file", []string{"list"}, "Potion|Мана|+5 MP|1|0.5|1\nWeapon|Меч|5|1|2|1|1|1\n", 1, `строка 2 ("Weapon|Меч|5|1|2|1|1|1"): некорректная запись предмета: повторяющийся ID 1`},
		{"unassigned ID in file", []string{"list"}, "Potion|Мана|+5 MP|1|0.5|0\n", 1, "ID 0 должен быть положительным"},
		{"broken weapon", []string{"use", "-id", "1"}, "Weapon|Меч|5|0|2|1|1|1\n", 1, "Предмет сломан"},
		{"empty potion", []string{"use", "-id", "1"}, "Potion|Мана|+5 MP|0|0.5|1\n", 1, "Предмет исчерпан"},
		{"recovered panic", []string{"use", "-id", "1"}, "Potion|Проклятый предмет|без эффекта|1|0.5|1\n", 1, "сработало проклятие"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "inventory.txt")
			if err := os.WriteFile(path, []byte(tt.file), 0600); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := runner.run(append(tt.args, "-file", path), &stdout, &stderr)
			if code != tt.code || stdout.Len() != 0 || !strings.Contains(stderr.String(), tt.want) {
				t.Errorf("code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
			}
			got, err := os.ReadFile(path)
			if err != nil || string(got) != tt.file {
				t.Errorf("file changed after failed command: %q, %v", got, err)
			}
		})
	}
}

func TestCommandReportsFileErrors(t *testing.T) {
	runner := commandRunner{store: storage.FileStore{}}
	dir := t.TempDir()
	for _, tt := range []struct {
		name string
		args []string
		want string
	}{
		{"read directory", []string{"list", "-file", dir}, "инвентар"},
		{"write into missing directory", []string{"add", "-name", "Мана", "-weight", "1", "-file", filepath.Join(dir, "missing", "inventory.txt")}, "для записи"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runner.run(tt.args, &stdout, &stderr); code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), tt.want) {
				t.Errorf("code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestCommandUsesInjectedStore(t *testing.T) {
	store := &testStore{items: []items.Item{items.NewPotion(41, "Мана", 1, "", 1)}}
	runner := commandRunner{store: store}
	var stdout, stderr bytes.Buffer
	code := runner.run([]string{"add", "-name", "Здоровье", "-weight", "1", "-file", "custom"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 || !strings.Contains(stdout.String(), "ID 42") {
		t.Fatalf("code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}
	if store.loadedFrom != "custom" || store.savedTo != "custom" || len(store.items) != 2 || items.IDOf(store.items[1]) != 42 {
		t.Errorf("load = %q, save = %q, items = %v", store.loadedFrom, store.savedTo, store.items)
	}
}

func TestCommandReportsOutputErrors(t *testing.T) {
	for _, tt := range []struct {
		name       string
		args       []string
		want       string
		savedCount int
	}{
		{"list", []string{"list"}, "вывод списка предметов", 1},
		{"add", []string{"add", "-name", "Здоровье", "-weight", "1"}, "изменения сохранены", 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			store := &testStore{items: []items.Item{items.NewPotion(1, "Мана", 1, "", 1)}}
			runner := commandRunner{store: store}
			var stderr bytes.Buffer
			code := runner.run(tt.args, failingWriter{}, &stderr)
			if code != 1 || !strings.Contains(stderr.String(), tt.want) {
				t.Errorf("code = %d, stderr = %q", code, stderr.String())
			}
			if len(store.items) != tt.savedCount {
				t.Errorf("saved %d items, want %d", len(store.items), tt.savedCount)
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("output unavailable")
}

type testStore struct {
	items      []items.Item
	loadedFrom string
	savedTo    string
}

func (s *testStore) Load(filename string) ([]items.Item, error) {
	s.loadedFrom = filename
	return s.items, nil
}

func (s *testStore) Save(filename string, itemList []items.Item) error {
	s.savedTo = filename
	s.items = itemList
	return nil
}
