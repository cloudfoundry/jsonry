package jsonry

import (
	"encoding/json"
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
)

type sort uint

const (
	basicSort sort = iota
	structSort
	listSort
	mapSort
	jsonMarshalerSort
	jsonryMarshalerSort
	unsupportedSort
)

var (
	jsonMarshalerInterface   = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	jsonryMarshalerInterface = reflect.TypeOf((*Marshaler)(nil)).Elem()
)

// Marshaler is the interface implemented by types that
// can marshal themselves into a Go type that JSONry can handle.
type Marshaler interface {
	MarshalJSONry() (interface{}, error)
}

func Marshal(input interface{}) ([]byte, error) {
	i := underlyingValue(reflect.ValueOf(input))

	if i.Kind() != reflect.Struct {
		return nil, fmt.Errorf(`the input must be a struct, not "%s"`, i.Kind())
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
		r, err = marshalStruct(ctx, in)
	case listSort:
		r, err = marshalList(ctx, in)
	case mapSort:
		r, err = marshalMap(ctx, in)
	case jsonMarshalerSort:
		r, err = marshalJSONMarshaler(ctx, in)
	default:
		err = NewUnsupportedTypeError(ctx, underlyingType(in))
	}
	return
}

func marshalStruct(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	s := underlyingValue(in)
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		private := f.PkgPath != ""

		if !private {
			r, err := marshal(ctx.WithField(f.Name, f.Type), s.Field(i))
			if err != nil {
				return nil, err
			}

			out[f.Name] = r
		}
	}

	return out, nil
}

func marshalList(ctx context.Context, in reflect.Value) ([]interface{}, error) {
	var out []interface{}

	list := underlyingValue(in)
	for i := 0; i < list.Len(); i++ {
		ctx := ctx.WithIndex(i, list.Type())
		r, err := marshal(ctx, list.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

func marshalMap(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	iter := underlyingValue(in).MapRange()
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

func marshalJSONMarshaler(ctx context.Context, in reflect.Value) (interface{}, error) {
	t := underlyingValue(in).MethodByName("MarshalJSON").Call(nil)

	if !t[1].IsNil() {
		return nil, fmt.Errorf("error from MarshaJSON() call %s: %w", ctx, toError(t[1]))
	}

	var r interface{}
	err := json.Unmarshal(t[0].Bytes(), &r)
	if err != nil {
		return nil, fmt.Errorf(`error parsing MarshaJSON() output "%s" %s: %w`, t[0].Bytes(), ctx, err)
	}

	return r, nil
}

func sortOf(v reflect.Value) sort {
	switch {
	case underlyingType(v).Implements(jsonMarshalerInterface):
		return jsonMarshalerSort
	case v.Type().Implements(jsonryMarshalerInterface):
		return jsonryMarshalerSort
	}

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

func toError(v reflect.Value) error {
	if v.CanInterface() {
		if err, ok := v.Interface().(error); ok {
			return err
		}
		return fmt.Errorf("could not cast to error: %+v", v)
	}
	r := v.MethodByName("Error").Call(nil)
	return fmt.Errorf("%s", r[0])
}
