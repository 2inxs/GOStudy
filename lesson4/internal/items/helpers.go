package items

import (
	"fmt"
	"reflect"
)

func IsNilItem(i Item) bool {
	if i == nil {
		return true
	}
	value := reflect.ValueOf(i)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func DescribeItem(i Item) string {
	if IsNilItem(i) {
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%d: %s (вес: %g)", IDOf(i), i.Name(), i.Weight())
}

// IDOf возвращает ID предмета или 0, если предмет не предоставляет ID.
func IDOf(item Item) int {
	identified, ok := item.(interface{ ID() int })
	if !ok || IsNilItem(item) {
		return 0
	}
	return identified.ID()
}

// RegisterID проверяет ID и отмечает его занятым в переданном наборе.
func RegisterID(used map[int]struct{}, id int) error {
	if id <= 0 {
		return fmt.Errorf("ID %d должен быть положительным: %w", id, ErrInvalidItem)
	}
	if _, exists := used[id]; exists {
		return fmt.Errorf("повторяющийся ID %d: %w", id, ErrInvalidItem)
	}
	used[id] = struct{}{}
	return nil
}

// UseDangerousItem показывает, как перехватить панику на границе операции.
// Поломка предмета ожидаема, поэтому вызывающий код получает обычный error.
// При panic строка return item.Use() не выполняется. Начинается раскрутка стека.
// recover срабатывает при прямом вызове из отложенной функции во время раскрутки.
func UseDangerousItem(item Item) (message string, err error) {
	if IsNilItem(item) {
		return "", ErrNilItem
	}
	var name string
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("использование предмета %q: %v", name, recovered)
		}
	}()

	name = item.Name()
	if name == "Проклятый предмет" {
		panic("сработало проклятие")
	}
	return item.Use()
}
