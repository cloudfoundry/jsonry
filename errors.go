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

func newUnsupportedTypeError(ctx context.Context, t reflect.Type) *UnsupportedType {
	return &UnsupportedType{
		context: ctx,
		typ:     t,
	}
}

func (u UnsupportedType) Error() string {
	return fmt.Sprintf(`unsupported type "%s" %s`, u.typ, u.context)
}

type UnsupportedKeyType struct {
	context context.Context
	typ     reflect.Type
}

func newUnsupportedKeyTypeError(ctx context.Context, t reflect.Type) *UnsupportedKeyType {
	return &UnsupportedKeyType{
		context: ctx,
		typ:     t,
	}
}

func (u UnsupportedKeyType) Error() string {
	return fmt.Sprintf(`maps must only have strings keys for "%s" %s`, u.typ, u.context)
}
