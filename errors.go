package jsonry

import (
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
)

type UnsupportedType struct {
	Context  context.Context
	TypeName string
}

func NewUnsupportedTypeError(ctx context.Context, t reflect.Type) *UnsupportedType {
	return &UnsupportedType{
		Context:  ctx,
		TypeName: t.String(),
	}
}

func (u UnsupportedType) Error() string {
	return fmt.Sprintf(`unsupported type "%s" %s`, u.TypeName, u.Context)
}

type UnsupportedKeyType struct {
	Context  context.Context
	TypeName string
}

func NewUnsupportedKeyTypeError(ctx context.Context, t reflect.Type) *UnsupportedKeyType {
	return &UnsupportedKeyType{
		Context:  ctx,
		TypeName: t.String(),
	}
}

func (u UnsupportedKeyType) Error() string {
	return fmt.Sprintf(`maps must only have strings keys for "%s" %s`, u.TypeName, u.Context)
}
