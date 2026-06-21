package base

import (
	"reflect"
	"strings"
)

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

func hasZeroID(entity any) bool {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if !isPrimaryKeyTag(t.Field(i).Tag.Get("gorm")) {
			continue
		}
		if v.Field(i).IsZero() {
			return true
		}
	}
	return false
}

// isPrimaryKeyTag reports whether a gorm struct tag marks the field as a primary
// key (gorm accepts both "primaryKey" and the legacy "primary_key").
func isPrimaryKeyTag(tag string) bool {
	for _, part := range strings.Split(tag, ";") {
		part = strings.TrimSpace(part)
		if strings.EqualFold(part, "primaryKey") || strings.EqualFold(part, "primary_key") {
			return true
		}
	}
	return false
}
