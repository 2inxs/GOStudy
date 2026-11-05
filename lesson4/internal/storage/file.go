package storage

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gostudy/lesson4/internal/items"
)

func Save(w io.Writer, itemList []items.Item) error {
	for _, item := range itemList {
		if items.IsNilItem(item) {
			return fmt.Errorf("сохранение: %w", items.ErrNilItem)
		}
		storable, ok := item.(items.Storable)
		if !ok {
			return &items.ItemError{ID: item.ID(), Name: item.Name(), Err: ErrNotStorable}
		}
		if err := storable.Serialize(w); err != nil {
			return &items.ItemError{ID: item.ID(), Name: item.Name(), Err: err}
		}
	}
	return nil
}

func Load(r io.Reader) ([]items.Item, error) {
	scanner := bufio.NewScanner(r)
	loaded := make([]items.Item, 0)
	line := 0
	for scanner.Scan() {
		line++
		item, err := items.ReadItem(strings.NewReader(scanner.Text()))
		if err != nil {
			return nil, &RecordError{Line: line, Err: err}
		}
		loaded = append(loaded, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, &RecordError{Line: line + 1, Err: err}
	}
	return loaded, nil
}
