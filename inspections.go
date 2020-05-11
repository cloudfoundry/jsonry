package jsonry

import "reflect"

func public(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func basicType(k reflect.Kind) bool {
	switch k {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

type valueDetails struct {
	value   reflect.Value
	typ     reflect.Type
	kind    reflect.Kind
	pointer bool
}

func inspectValue(v reflect.Value) valueDetails {
	k := v.Kind()
	switch k {
	case reflect.Ptr:
		r := inspectValue(v.Elem())
		r.pointer = true
		return r
	case reflect.Interface:
		return inspectValue(v.Elem())
	case reflect.Invalid:
		return valueDetails{
			value:   v,
			typ:     nil,
			kind:    k,
			pointer: false,
		}
	default:
		return valueDetails{
			value:   v,
			typ:     v.Type(),
			kind:    k,
			pointer: false,
		}
	}
}
