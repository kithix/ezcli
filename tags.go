package ezcli

import (
	"reflect"
)

const (
	tagEnv = "env"
)

func parseTags(field reflect.StructField) []CommandOptions {
	val, exists := field.Tag.Lookup(tagEnv)
	if exists {
		_ = val
	}
	_ = val
	return make([]CommandOptions, 0)
}
