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

		switch fType.Type.Kind() {
		case reflect.Bool:
			v := fVal.Bool()
			a.genericVar(&v, optFns...)
			a.postLoadFuncs = append(a.postLoadFuncs, func() {
				fVal.SetBool(v)
			})

		// Handle ints
		case reflect.Int:
			setInt[int](a, fVal, optFns)
		case reflect.Int8:
			setInt[int8](a, fVal, optFns)
		case reflect.Int16:
			setInt[int16](a, fVal, optFns)
		case reflect.Int32:
			setInt[int32](a, fVal, optFns)
		case reflect.Int64:
			setInt[int64](a, fVal, optFns)

		// Handle uints
		case reflect.Uint:
			setUint[uint](a, fVal, optFns)
		case reflect.Uint8:
			setUint[uint8](a, fVal, optFns)
		case reflect.Uint16:
			setUint[uint16](a, fVal, optFns)
		case reflect.Uint32:
			setUint[uint32](a, fVal, optFns)
		case reflect.Uint64:
			setUint[uint64](a, fVal, optFns)

		case reflect.String:
			v := fVal.String()
			a.genericVar(&v, optFns...)
			a.postLoadFuncs = append(a.postLoadFuncs, func() {
				fVal.SetString(v)
			})
		default:
			panic("unable to use struct value")
			// Do we skip struct values we can't use?
			// Maybe make it an option?
		}

	}
}

func setUint[T uint | uint8 | uint16 | uint32 | uint64](a *App, val reflect.Value, optFns []varOptFn) {
	v := T(val.Uint())
	a.genericVar(&v, optFns...)
	a.postLoadFuncs = append(a.postLoadFuncs, func() {
		val.SetUint(uint64(v))
	})
}

func setInt[T int | int8 | int16 | int32 | int64](a *App, val reflect.Value, optFns []varOptFn) {
	v := T(val.Int())
	a.genericVar(&v, optFns...)
	a.postLoadFuncs = append(a.postLoadFuncs, func() {
		val.SetInt(int64(v))
	})
}
