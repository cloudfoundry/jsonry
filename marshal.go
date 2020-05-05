package jsonry

import (
	"encoding/json"
	"errors"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
)

type sort uint

const (
	basicSort sort = iota
	structSort
	listSort
	unsupportedSort
)

func Marshal(input interface{}) ([]byte, error) {
	i := actualValue(reflect.ValueOf(input))

	if i.Kind() != reflect.Struct {
		return nil, errors.New("the input must be a struct")
	}

	m, err := marshal(context.Context{}, i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(ctx context.Context, in reflect.Value) (interface{}, error) {
	switch sortOf(in) {
	case basicSort:
		return in.Interface(), nil
	case structSort:
		r, err := marshalStruct(ctx, actualValue(in))
		if err != nil {
			return nil, err
		}
		return r, nil
	case listSort:
		r, err := marshalList(ctx, actualValue(in))
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		return nil, NewUnsupportedTypeError(ctx, actualType(in))
	}
	return nil, nil
}

func marshalStruct(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	t := actualType(in)
	for i := 0; i < t.NumField(); i++ {
		v := in.Field(i)
		f := t.Field(i)
		ctx := ctx.WithField(f)

		r, err := marshal(ctx, v)
		if err != nil {
			return nil, err
		}

		out[f.Name] = r
	}

	return out, nil
}

func marshalList(ctx context.Context, in reflect.Value) ([]interface{}, error) {
	var out []interface{}

	for i := 0; i < in.Len(); i++ {
		ctx := ctx.WithIndex(i, in.Type())
		r, err := marshal(ctx, in.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

func sortOf(v reflect.Value) sort {
	switch actualValue(v).Kind() {
	case reflect.Struct:
		return structSort
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return basicSort
	case reflect.Slice, reflect.Array:
		return listSort
	default:
		return unsupportedSort
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
	case reflect.Interface, reflect.Ptr:
		return actualValue(v.Elem())
	default:
		return v
	}
}
