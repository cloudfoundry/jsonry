package jsonry

import (
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
	return fmt.Sprintf(`maps must only have strings keys for "%s" at %s`, u.typ, u.context)
}

type ConversionError struct {
	context context.Context
	typ     reflect.Type
	value   interface{}
}

func newConversionError(ctx context.Context, v interface{}, t reflect.Type) error {
	return &ConversionError{
		context: ctx,
		typ:     t,
		value:   v,
	}
}

func (c ConversionError) Error() string {
	msg := fmt.Sprintf(`cannot unmarshal "%+v" `, c.value)

	if c.typ != nil {
		msg = fmt.Sprintf(`%stype "%s" `, msg, c.typ)
	}

	return msg + "into " + c.context.String()
}
