package jsonry

import (
	"encoding/json"
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/tree"

	"code.cloudfoundry.org/jsonry/internal/context"
	"code.cloudfoundry.org/jsonry/internal/path"
)

func Marshal(input interface{}) ([]byte, error) {
	iv := inspectValue(reflect.ValueOf(input))

	if iv.kind != reflect.Struct {
		return nil, fmt.Errorf(`the input must be a struct, not "%s"`, iv.kind)
	}

	m, err := marshalStruct(context.Context{}, iv.value, iv.typ)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
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

func marshal(ctx context.Context, in reflect.Value) (r interface{}, err error) {
	input := inspectValue(in)

	switch {
	case implements(input.typ, (*json.Marshaler)(nil)):
		r, err = marshalJSONMarshaler(ctx, input.value)
	case basicType(input.kind):
		r = in.Interface()
	case input.kind == reflect.Invalid:
		r = nil
	case input.kind == reflect.Struct:
		r, err = marshalStruct(ctx, input.value, input.typ)
	case input.kind == reflect.Slice || input.kind == reflect.Array:
		r, err = marshalList(ctx, input.value)
	case input.kind == reflect.Map:
		r, err = marshalMap(ctx, input.value)
	default:
		err = newUnsupportedTypeError(ctx, input.typ)
	}

	return
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
		return nil, fmt.Errorf("error from MarshaJSON() call at %s: %w", ctx, toError(t[1]))
	}

	var r interface{}
	err := json.Unmarshal(t[0].Bytes(), &r)
	if err != nil {
		return nil, fmt.Errorf(`error parsing MarshaJSON() output "%s" at %s: %w`, t[0].Bytes(), ctx, err)
	}

	return r, nil
}
