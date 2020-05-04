package jsonry

import (
	"fmt"
	"reflect"
)

type UnsupportedType struct {
	TypeName string
	Context  context
}

func NewUnsupportedTypeError(ctx context, t reflect.Type) *UnsupportedType {
	return &UnsupportedType{
		Context:  ctx,
		TypeName: t.String(),
	}
}

func (u *UnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type '%s' at %s", u.TypeName, u.Context)
}
