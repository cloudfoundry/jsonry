package jsonry

import (
	"encoding/json"
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/tree"

	"code.cloudfoundry.org/jsonry/internal/context"
	"code.cloudfoundry.org/jsonry/internal/path"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into a Go type that JSONry can handle.
type Marshaler interface {
	MarshalJSONry() (interface{}, error)
}

func Marshal(input interface{}) ([]byte, error) {
	i, _, k := underlying(reflect.ValueOf(input))

	if k != reflect.Struct {
		return nil, fmt.Errorf(`the input must be a struct, not "%s"`, i.Kind())
	}

	m, err := marshal(context.Context{}, i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(ctx context.Context, in reflect.Value) (r interface{}, err error) {
	uv, ut, k := underlying(in)

	switch {
	case implements(ut, (*json.Marshaler)(nil)):
		r, err = marshalJSONMarshaler(ctx, uv)
	case implements(ut, (*Marshaler)(nil)):
		r, err = marshalJSONryMarshaler(ctx, uv)
	case isBasicType(k):
		r = in.Interface()
	case k == reflect.Invalid:
		r = nil
	case k == reflect.Struct:
		r, err = marshalStruct(ctx, uv)
	case k == reflect.Slice || k == reflect.Array:
		r, err = marshalList(ctx, uv)
	case k == reflect.Map:
		r, err = marshalMap(ctx, uv)
	default:
		err = newUnsupportedTypeError(ctx, ut)
	}

	return
}

func marshalStruct(ctx context.Context, in reflect.Value) (map[string]interface{}, error) {
	out := make(tree.Tree)

	t := in.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		private := f.PkgPath != ""

		if !private {
			p := path.ComputePath(f)
			if !p.OmitEmpty || !in.Field(i).IsZero() {
				r, err := marshal(ctx.WithField(f.Name, f.Type), in.Field(i))
				if err != nil {
					return nil, err
				}

				out.Attach(p, r)
			}
		}
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
			return nil, newUnsupportedKeyTypeError(ctx, in.Type())
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
	t := in.MethodByName("MarshalJSON").Call(nil)

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

func marshalJSONryMarshaler(ctx context.Context, in reflect.Value) (interface{}, error) {
	t := in.MethodByName("MarshalJSONry").Call(nil)

	if !t[1].IsNil() {
		return nil, fmt.Errorf("error from MarshaJSONry() call %s: %w", ctx, toError(t[1]))
	}

	return marshal(ctx, t[0])
}

func underlying(v reflect.Value) (reflect.Value, reflect.Type, reflect.Kind) {
	k := v.Kind()
	switch k {
	case reflect.Interface, reflect.Ptr:
		return underlying(v.Elem())
	case reflect.Invalid:
		return v, nil, k
	default:
		return v, v.Type(), k
	}
}

func isBasicType(k reflect.Kind) bool {
	switch k {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
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

func implements(t reflect.Type, iptr interface{}) bool {
	if t == nil {
		return false
	}
	return t.Implements(reflect.TypeOf(iptr).Elem())
}
