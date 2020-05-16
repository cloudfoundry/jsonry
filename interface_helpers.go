package jsonry

import (
	"fmt"
	"reflect"
)

func toError(v reflect.Value) error {
	if v.CanInterface() {
		if err, ok := v.Interface().(error); ok {
			return err
		}
		return fmt.Errorf("could not cast to error: %+v", v)
	}
	r := v.MethodByName("Error").Call(nil)
	return fmt.Errorf("%s", r[0])
}

func implements(t reflect.Type, iptr interface{}) bool {
	if t == nil {
		return false
	}
	return t.Implements(reflect.TypeOf(iptr).Elem())
}
