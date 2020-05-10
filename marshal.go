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
	iv := inspectValue(reflect.ValueOf(input))

	if iv.realKind != reflect.Struct {
		return nil, fmt.Errorf(`the input must be a struct, not "%s"`, iv.realKind)
	}

	m, err := marshalStruct(context.Context{}, iv.realValue, iv.realType)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(ctx context.Context, in reflect.Value) (r interface{}, err error) {
	input := inspectValue(in)

	switch {
	case implements(input.realType, (*json.Marshaler)(nil)):
		r, err = marshalJSONMarshaler(ctx, input.realValue)
	case implements(input.realType, (*Marshaler)(nil)):
		r, err = marshalJSONryMarshaler(ctx, input.realValue)
	case basicType(input.realKind):
		r = in.Interface()
	case input.realKind == reflect.Invalid:
		r = nil
	case input.realKind == reflect.Struct:
		r, err = marshalStruct(ctx, input.realValue, input.realType)
	case input.realKind == reflect.Slice || input.realKind == reflect.Array:
		r, err = marshalList(ctx, input.realValue)
	case input.realKind == reflect.Map:
		r, err = marshalMap(ctx, input.realValue)
	default:
		err = newUnsupportedTypeError(ctx, input.realType)
	}

	return
}

func marshalStruct(ctx context.Context, in reflect.Value, t reflect.Type) (map[string]interface{}, error) {
	out := make(tree.Tree)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if public(f) {
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
