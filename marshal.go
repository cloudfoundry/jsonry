package jsonry

import (
	"encoding/json"
	"errors"
	"reflect"
)

type metakind uint

const (
	metakindBasic metakind = iota
	metakindStruct
	metakindUnsupported
)

func Marshal(input interface{}) ([]byte, error) {
	i := reflect.ValueOf(input)

	// TODO: use actualMetakind() or similar
	if i.Kind() == reflect.Ptr {
		i = i.Elem()
	}

	if i.Kind() != reflect.Struct {
		return nil, errors.New("the input must be a struct")
	}

	m, err := marshal(context{}, i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(ctx context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	t := actualType(in)
	for i := 0; i < t.NumField(); i++ {
		v := in.Field(i)
		f := t.Field(i)
		ctx = ctx.withField(f)

		switch actualMetakind(v) {
		case metakindBasic:
			out[f.Name] = v.Interface()
		case metakindStruct:
			r, err := marshal(ctx, actualValue(v))
			if err != nil {
				return nil, err
			}
			out[f.Name] = r
		default:
			return nil, NewUnsupportedTypeError(ctx, actualType(v))
		}
	}

	return out, nil
}

func actualMetakind(v reflect.Value) metakind {
	switch v.Kind() {
	case reflect.Interface:
		return actualMetakind(v.Elem())
	case reflect.Struct:
		return metakindStruct
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

func actualValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface:
		return v.Elem()
	default:
		return v
	}
}
