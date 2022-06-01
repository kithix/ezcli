package ezcli

import (
	"reflect"
)

const (
	tagEnv  = "env"
	tagFlag = "flag"
)

func (a *App) parseTags(field reflect.StructField) []varOptFn {
	varOptFns := make([]varOptFn, 0)

	envVal, exists := field.Tag.Lookup(tagEnv)
	if exists {
		// Use the field name if there was no custom name provided
		if envVal == "" {
			envVal = field.Name
		}
		varOptFns = append(varOptFns, VarEnv(envVal))
	}
	flagVal, exists := field.Tag.Lookup(tagFlag)
	if exists && flagVal != "" {
		varOptFns = append(varOptFns, VarName(flagVal))
	} else {
		// default to the field name if we have no value
		varOptFns = append(varOptFns, VarName(field.Name))
	}

	return varOptFns
}

func (a *App) StructVar(s any) {
	sType := reflect.TypeOf(s)
	sVal := reflect.ValueOf(s)
	// Get the value of any pointers
	if sType.Kind() == reflect.Pointer {
		sType = sType.Elem()
		sVal = sVal.Elem()
	}
	if sType.Kind() != reflect.Struct {
		panic("Must iterate over fields of a struct")
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < sType.NumField(); i++ {
		fType := sType.Field(i)
		fVal := sVal.Field(i)

		// Skip any unexported or anonymous fields
		if !fType.IsExported() {
			continue
		}

		// Recurse over structs
		if fType.Type.Kind() == reflect.Struct {
			a.StructVar(fVal)
			continue
		}

		// Parse the tags as options
		optFns := a.parseTags(fType)

		// Use our structs set value as the default
		optFns = append(optFns, VarDefaultValue(fVal.Interface()))

		// Get the value out of the field
		switch fType.Type.String() {
		case "bool":
			v := fVal.Bool()
			a.genericVar(&v, optFns...)
			a.postLoadFuncs = append(a.postLoadFuncs, func() {
				fVal.SetBool(v)
			})
		case "int":
			v := int(fVal.Int())
			a.genericVar(&v, optFns...)
			a.postLoadFuncs = append(a.postLoadFuncs, func() {
				fVal.SetInt(int64(v))
			})
		case "string":
			v := fVal.String()
			a.genericVar(&v, optFns...)
			a.postLoadFuncs = append(a.postLoadFuncs, func() {
				fVal.SetString(v)
			})
		default:
			// Do we skip struct values we can't use?
			// Maybe make it an option?
		}

	}
}
