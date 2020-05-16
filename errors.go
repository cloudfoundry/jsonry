package jsonry

import (
	"encoding/json"
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
)

type UnsupportedType struct {
	context context.Context
	typ     reflect.Type
}

func newUnsupportedTypeError(ctx context.Context, t reflect.Type) error {
	return &UnsupportedType{
		context: ctx,
		typ:     t,
	}
}

func (u UnsupportedType) Error() string {
	return fmt.Sprintf(`unsupported type "%s" at %s`, u.typ, u.context)
}

type UnsupportedKeyType struct {
	context context.Context
	typ     reflect.Type
}

func newUnsupportedKeyTypeError(ctx context.Context, t reflect.Type) error {
	return &UnsupportedKeyType{
		context: ctx,
		typ:     t,
	}
}

func (u UnsupportedKeyType) Error() string {
	return fmt.Sprintf(`maps must only have string keys for "%s" at %s`, u.typ, u.context)
}

type ConversionError struct {
	context context.Context
	value   interface{}
}

func newConversionError(ctx context.Context, value interface{}) error {
	return &ConversionError{
		context: ctx,
		value:   value,
	}
}

func (c ConversionError) Error() string {
	var t string
	switch c.value.(type) {
	case nil:
	case json.Number:
		t = "number"
	default:
		t = reflect.TypeOf(c.value).String()
	}

	msg := fmt.Sprintf(`cannot unmarshal "%+v" `, c.value)

	if t != "" {
		msg = fmt.Sprintf(`%stype "%s" `, msg, t)
	}

	return msg + "into " + c.context.String()
}
