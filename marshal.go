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
	mapSort
	unsupportedSort
)

func Marshal(input interface{}) ([]byte, error) {
	i := underlyingValue(reflect.ValueOf(input))

	if i.Kind() != reflect.Struct {
		return nil, errors.New("the input must be a struct")
	}

	m, err := marshal(context.Context{}, i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(ctx context.Context, in reflect.Value) (r interface{}, err error) {
	switch sortOf(in) {
	case basicSort:
		r = in.Interface()
	case structSort:
		r, err = marshalStruct(ctx, underlyingValue(in))
	case listSort:
		r, err = marshalList(ctx, underlyingValue(in))
	case mapSort:
		r, err = marshalMap(ctx, underlyingValue(in))
	default:
		err = NewUnsupportedTypeError(ctx, underlyingType(in))
	}
	return
}

func marshalStruct(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	t := in.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		r, err := marshal(ctx.WithField(f.Name, f.Type), in.Field(i))
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

func marshalMap(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	iter := in.MapRange()
	for iter.Next() {
		k := iter.Key()
		if k.Kind() != reflect.String {
			return nil, NewUnsupportedKeyTypeError(ctx, underlyingType(in))
		}

		ctx := ctx.WithKey(k.String(), k.Type())

		r, err := marshal(ctx, iter.Value())
		if err != nil {
			return nil, err
		}
		out[k.String()] = r
	}

	return out, nil
}

func sortOf(v reflect.Value) sort {
	switch underlyingValue(v).Kind() {
	case reflect.Struct:
		return structSort
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return basicSort
	case reflect.Slice, reflect.Array:
		return listSort
	case reflect.Map:
		return mapSort
	default:
		return unsupportedSort
	}
}

func underlyingType(v reflect.Value) reflect.Type {
	switch v.Kind() {
	case reflect.Interface:
		return underlyingType(v.Elem())
	case reflect.Ptr:
		return reflect.PtrTo(underlyingType(v.Elem()))
	default:
		return v.Type()
	}
}

func underlyingValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		return underlyingValue(v.Elem())
	default:
		return v
	}
}
