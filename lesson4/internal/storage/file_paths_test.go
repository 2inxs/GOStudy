package storage_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
)

func TestSaveFileReportsPermissionError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.txt")
	const original = "Potion|Мана|+10 MP|1|0.5|1\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0400); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(path, 0600); err != nil {
			t.Error(err)
		}
	})
	err := storage.SaveFile(path, nil)
	if err == nil {
		t.Skip("the current account can write to read-only files")
	}
	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("SaveFile error = %v, want permission error", err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != original {
		t.Errorf("original file = %q, read error = %v", got, err)
	}
}

func TestSaveFilePreservesDataOnSerializationError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.txt")
	const original = "Potion|Мана|+10 MP|1|0.5|1\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	valid := items.NewPotion(2, "Здоровье", 1, "+10 HP", 1)
	invalid := items.NewPotion(3, "Яд|неверная запись", 1, "-10 HP", 1)
	if err := storage.SaveFile(path, []items.Item{valid, invalid}); !errors.Is(err, items.ErrInvalidItem) {
		t.Fatalf("SaveFile error = %v, want invalid item", err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != original {
		t.Errorf("original file = %q, read error = %v", got, err)
	}
}
