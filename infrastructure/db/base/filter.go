package base

import "reflect"

func IsEmptyFilter[T any](filter T) bool {
	v := reflect.ValueOf(filter)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}

		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				return false
			}
		} else if !field.IsZero() {
			return false
		}
	}

	return true
}
