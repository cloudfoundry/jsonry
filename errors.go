package jsonry

import (
	"fmt"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
)

type UnsupportedType struct {
	TypeName string
	Context  context.Context
}

func NewUnsupportedTypeError(ctx context.Context, t reflect.Type) *UnsupportedType {
	return &UnsupportedType{
		Context:  ctx,
		TypeName: t.String(),
	}
}

func (u *UnsupportedType) Error() string {
	return fmt.Sprintf(`unsupported type "%s" %s`, u.TypeName, u.Context)
}
