package ezcli

import (
	"fmt"
	"reflect"
	"strings"
)

// TODO options
func (app *App) StructVar(s interface{}) {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Pointer {
		panic("must be a pointer to a struct")
	}
	// Get the value of the pointer
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("must be a pointer to a struct")
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		// How to handle nested structs? Recursively for sure :partyparrot:
		// if field.Type == reflect.Struct {
		// 	app.StructVar(field)
		//}

		// Get the field tag value
		fmt.Printf("%+v\n", field)
		tag := field.Tag.Get("ezcli")
		if tag == "" {
			// default
			continue
		}
		subtags := strings.Split(tag, ",")
		for _, subtag := range subtags {
			_ = subtag
		}
	}
}
