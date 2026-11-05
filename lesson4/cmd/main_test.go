package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gostudy/lesson4/internal/items"
)

func TestRunLoadsUsesAndSavesInventory(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "input.txt")
	output := filepath.Join(directory, "output.txt")
	const initial = "Weapon|Сломанный лук|5|0|2|1|1|102\n" +
		"Potion|Пустое|+10 HP|0|0.5|301\n" +
		"Potion|Лечебное|+50 HP|1|0.5|302\n" +
		"Armor|Щит|5|4.5|1|1|201\n"
	if err := os.WriteFile(input, []byte(initial), 0600); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "-load", input, "-save", output)

	if err := Run(); err != nil {
		t.Fatalf("Run() returned an unexpected error: %v", err)
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("reading saved inventory: %v", err)
	}
	const want = "Weapon|Сломанный лук|5|0|2|1|1|102\n" +
		"Potion|Пустое|+10 HP|0|0.5|301\n" +
		"Potion|Лечебное|+50 HP|0|0.5|302\n" +
		"Armor|Щит|5|4.5|1|1|201\n"
	if string(got) != want {
		t.Errorf("saved inventory = %q, want %q", got, want)
	}
	unchanged, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(unchanged) != initial {
		t.Error("Run() modified the input file while saving to a different path")
	}
}

func TestRunRejectsInvalidArgumentsBeforeSaving(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"-unknown"}},
		{name: "unexpected positional argument", args: []string{"unexpected"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "inventory.txt")
			args := append([]string{"-save", output}, tt.args...)
			setArgs(t, args...)

			if err := Run(); err == nil {
				t.Fatal("Run() accepted invalid arguments")
			}
			if _, err := os.Stat(output); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("Run() created the output file for invalid arguments: %v", err)
			}
		})
	}
}

func TestRunHelpDoesNotExecuteInventory(t *testing.T) {
	output := filepath.Join(t.TempDir(), "inventory.txt")
	setArgs(t, "-save", output, "-h")

	if err := Run(); err != nil {
		t.Fatalf("Run() with -h returned %v", err)
	}
	if _, err := os.Stat(output); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Run() created an inventory file while showing help: %v", err)
	}
}

func TestRunMissingInputPreservesFileError(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "missing.txt")
	setArgs(t, "-load", filename)

	err := Run()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Run() error = %v, want os.ErrNotExist", err)
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("Run() error lost its PathError: %v", err)
	}
	if pathErr.Path != filename {
		t.Errorf("PathError.Path = %q, want %q", pathErr.Path, filename)
	}
}

func TestRunInvalidItemDoesNotOverwriteOutput(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "input.txt")
	output := filepath.Join(directory, "output.txt")
	if err := os.WriteFile(input, []byte("Weapon|Меч|-1|5|3|1|1|101\n"), 0600); err != nil {
		t.Fatal(err)
	}
	const previous = "previous inventory\n"
	if err := os.WriteFile(output, []byte(previous), 0600); err != nil {
		t.Fatal(err)
	}
	setArgs(t, "-load", input, "-save", output)

	if err := Run(); !errors.Is(err, items.ErrInvalidItem) {
		t.Fatalf("Run() error = %v, want ErrInvalidItem", err)
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != previous {
		t.Errorf("Run() overwrote the output after a fatal item error: %q", got)
	}
}

func TestMainPanicsOnFatalRunError(t *testing.T) {
	setArgs(t, "-load", filepath.Join(t.TempDir(), "missing.txt"))
	defer func() {
		value := recover()
		err, ok := value.(error)
		if !ok {
			t.Fatalf("main() panic = %v, want an error", value)
		}
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("main() panic = %v, want the original file error", err)
		}
	}()

	main()
}

func setArgs(t *testing.T, args ...string) {
	t.Helper()
	previous := os.Args
	t.Cleanup(func() { os.Args = previous })
	os.Args = append([]string{"lesson4"}, args...)
}
