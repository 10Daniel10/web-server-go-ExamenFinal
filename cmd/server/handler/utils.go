package handler

import (
	"reflect"
)

func extractJSONTag(field string, s interface{}) (value string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return field
	}

	f, ok := rt.FieldByName(field)
	if ok == false {
		return field
	}

	tag := f.Tag.Get("json")
	if tag == "" {
		return field
	}

	return tag
}
