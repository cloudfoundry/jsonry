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
	if output.realKind != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct type, got: %s", output.realType)
	}

	var tree tree.Tree

	d := json.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	if err := d.Decode(&tree); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	return unmarshalIntoStruct(context.Context{}, output.realValue, output.realType, tree)
}

func unmarshal(ctx context.Context, out reflect.Value, r interface{}) error {
	v, t, k, _ := inspectTarget(out)

	var err error
	switch {
	case basicType(k):
		err = assign(ctx, v, r)
	case k == reflect.Struct:
		if m, ok := r.(map[string]interface{}); ok {
			err = unmarshalIntoStruct(ctx, v, t, m)
		}
	}
	return err
}

func unmarshalIntoStruct(ctx context.Context, out reflect.Value, t reflect.Type, tree tree.Tree) error {
	for i := 0; i < out.NumField(); i++ {
		f := t.Field(i)

		if public(f) {
			p := path.ComputePath(f)
			if r, ok := tree.Fetch(p); ok {
				if err := unmarshal(ctx.WithField(f.Name, f.Type), out.Field(i), r); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func assign(ctx context.Context, out reflect.Value, in interface{}) error {
	inValue := reflect.ValueOf(in)
	if !inValue.IsValid() || inValue.IsZero() {
		return nil
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
		return fmt.Errorf(`cannot convert value of type "%s" to type "%s" %s`, inType, outType, ctx)
	}
}

func assignNumber(ctx context.Context, out reflect.Value, outType reflect.Type, in json.Number) error {
	switch outType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := in.Int64()
		if err != nil {
			return fmt.Errorf(`cannot convert integer "%s" to type "%s" %s`, in, outType, ctx)
		}
		out.Set(reflect.ValueOf(i).Convert(outType))

	case reflect.Float32, reflect.Float64:
		i, err := in.Float64()
		if err != nil {
			return fmt.Errorf(`cannot convert float "%s" to type "%s" %s`, in, outType, ctx)
		}
		out.Set(reflect.ValueOf(i).Convert(outType))

	default:
		return fmt.Errorf(`cannot convert number "%s" to type "%s" %s`, in, outType, ctx)
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

func inspectTarget(v reflect.Value) (reflect.Value, reflect.Type, reflect.Kind, bool) {
	k := v.Kind()
	switch k {
	case reflect.Ptr:
		return v, v.Type().Elem(), v.Type().Elem().Kind(), true
	case reflect.Invalid:
		return v, nil, k, false
	default:
		return v, v.Type(), k, false
	}
}
