package jsonry

import (
	"encoding/json"
	"errors"
	"reflect"
)

func Marshal(input interface{}) ([]byte, error) {
	i := reflect.ValueOf(input)

	if i.Kind()==reflect.Ptr {
		i = i.Elem()
	}

	if i.Kind() != reflect.Struct {
		return nil, errors.New("the input must be a struct")
	}

	m, err := marshal(i)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func marshal(in reflect.Value) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	t := in.Type()
	for i := 0; i<t.NumField(); i++ {
		f := t.Field(i)
		out[f.Name] = in.Field(i).Interface()
	}

	return out, nil
}
