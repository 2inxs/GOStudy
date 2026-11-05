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
	return fmt.Sprintf("%d: %s (вес: %g)", i.ID(), i.Name(), i.Weight())
}
