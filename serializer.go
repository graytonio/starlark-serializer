package starlarkserializer

import (
	"errors"
	"reflect"
	"strings"

	"go.starlark.net/starlark"
)

// Marshals data to a starlark value
func Marshal(i any) (starlark.Value, error) {
	v := reflect.ValueOf(i)

	switch v.Type().Kind() {
	case reflect.String:
		return starlark.String(i.(string)), nil
	case reflect.Int:
		return starlark.MakeInt(i.(int)), nil
	case reflect.Bool:
		return starlark.Bool(i.(bool)), nil
	case reflect.Struct:
		return marshalStruct(i)
	default:
		return nil, errors.New("invalid reflect kind")
	}
}

func marshalStruct(v any) (*starlark.Dict, error) {
	val := reflect.ValueOf(v)
	results := starlark.NewDict(val.NumField())

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(val.Type().Field(i).Name)
		value, err := Marshal(val.Field(i).Interface())
		if err != nil {
			return nil, err
		}

		results.SetKey(starlark.String(key), value)
	}

	return results, nil
}
