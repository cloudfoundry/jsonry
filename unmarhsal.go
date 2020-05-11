package jsonry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/path"

	"code.cloudfoundry.org/jsonry/internal/context"
	"code.cloudfoundry.org/jsonry/internal/tree"
)

func Unmarshal(data []byte, out interface{}) error {
	output := inspectValue(reflect.ValueOf(out))

	if !output.pointer {
		return errors.New("output must be a pointer to a struct, got a non-pointer")
	}
	if output.kind != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct type, got: %s", output.typ)
	}

	var target tree.Tree

	d := json.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	if err := d.Decode(&target); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	return unmarshalIntoStruct(context.Context{}, output.value, output.typ, target)
}

func unmarshal(ctx context.Context, target reflect.Value, source interface{}) error {
	t, k, _ := inspectTarget(target.Type())

	var err error
	switch {
	case basicType(k), k == reflect.Interface:
		err = assign(ctx, target, source)
	case k == reflect.Struct:
		if m, ok := source.(map[string]interface{}); ok {
			err = unmarshalIntoStruct(ctx, target, t, m)
		} else {
			panic("no")
		}
	case k == reflect.Slice, k == reflect.Array:
		err = unmarshalIntoSlice(ctx, target, t, source)
	default:
		err = newUnsupportedTypeError(ctx, t)
	}
	return err
}

func unmarshalIntoStruct(ctx context.Context, target reflect.Value, t reflect.Type, source tree.Tree) error {
	for i := 0; i < target.NumField(); i++ {
		f := t.Field(i)

		if public(f) {
			p := path.ComputePath(f)
			if r, ok := source.Fetch(p); ok {
				if err := unmarshal(ctx.WithField(f.Name, f.Type), target.Field(i), r); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func unmarshalIntoSlice(ctx context.Context, target reflect.Value, targetType reflect.Type, source interface{}) error {
	sv := reflect.ValueOf(source)

	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		sv = reflect.ValueOf([]interface{}{source})
	}

	result := reflect.MakeSlice(targetType, sv.Len(), sv.Len())
	for i := 0; i < sv.Len(); i++ {
		elem := reflect.New(targetType.Elem()).Elem()
		if err := unmarshal(ctx.WithIndex(i, targetType.Elem()), elem, sv.Index(i).Interface()); err != nil {
			return err
		}
		result.Index(i).Set(elem)
	}
	target.Set(result)
	return nil
}

func assign(ctx context.Context, out reflect.Value, in interface{}) error {
	inValue := reflect.ValueOf(in)
	if !inValue.IsValid() || inValue.IsZero() {
		switch out.Kind() {
		case reflect.Ptr, reflect.Interface:
			return nil
		default:
			return newConversionError(ctx, in, nil)
		}
	}

	inType := inValue.Type()
	out, outType := depointerify(out)

	switch {
	case inType == reflect.TypeOf((*json.Number)(nil)).Elem():
		return assignNumber(ctx, out, outType, in.(json.Number))
	case inType.AssignableTo(outType):
		out.Set(inValue)
		return nil
	case inType.ConvertibleTo(outType):
		out.Set(inValue.Convert(outType))
		return nil
	default:
		return newConversionError(ctx, in, inType)
	}
}

func assignNumber(ctx context.Context, out reflect.Value, outType reflect.Type, in json.Number) error {
	switch outType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := in.Int64()
		if err != nil {
			return newConversionError(ctx, in, nil)
		}
		out.Set(reflect.ValueOf(i).Convert(outType))

	case reflect.Float32, reflect.Float64:
		f, err := in.Float64()
		if err != nil {
			return newConversionError(ctx, in, nil)
		}
		out.Set(reflect.ValueOf(f).Convert(outType))

	default:
		return newConversionError(ctx, in, nil)
	}

	return nil
}

func depointerify(v reflect.Value) (reflect.Value, reflect.Type) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		n := reflect.New(t)
		v.Set(n)
		v = n.Elem()
	}

	return v, t
}

func inspectTarget(t reflect.Type) (reflect.Type, reflect.Kind, bool) {
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		pt, pk, _ := inspectTarget(t.Elem())
		return pt, pk, true
	default:
		return t, k, false
	}
}
