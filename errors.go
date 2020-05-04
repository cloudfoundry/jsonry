package jsonry

import (
	"fmt"
	"reflect"
)

type UnsupportedType struct {
	TypeName string
}

func NewUnsupportedTypeError(t reflect.Type) *UnsupportedType {
	return &UnsupportedType{TypeName: t.String()}
}

func (u *UnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type: %s", u.TypeName)
}
