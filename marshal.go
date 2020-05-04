package jsonry

import (
	"encoding/json"
	"errors"
	"reflect"
)

type metakind uint

const (
	metakindBasic metakind = iota
	metakindUnsupported
)

func Marshal(input interface{}) ([]byte, error) {
	i := reflect.ValueOf(input)

	if i.Kind() == reflect.Ptr {
		i = i.Elem()
	}

	if i.Kind() != reflect.Struct {
		return nil, errors.New("the input must be a struct")
	}

	m, err := marshal(i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	t := in.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch actual(in.Field(i)) {
		case metakindBasic:
			out[f.Name] = in.Field(i).Interface()
		default:
			return nil, NewUnsupportedTypeError(actualType(in.Field(i)))
		}
	}

	return out, nil
}

func actual(v reflect.Value) metakind {
	switch v.Kind() {
	case reflect.Interface:
		return actual(v.Elem())
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return metakindBasic
	default:
		return metakindUnsupported
	}
}

func actualType(v reflect.Value) reflect.Type {
	switch v.Kind() {
	case reflect.Interface:
		return v.Elem().Type()
	default:
		return v.Type()
	}
}
