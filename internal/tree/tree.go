package tree

import (
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/path"
)

type Tree map[string]interface{}

func (t Tree) Attach(p path.Path, v interface{}) Tree {
	switch p.Len() {
	case 0:
		panic("empty path")
	case 1:
		branch, _ := p.Pull()
		t[branch.Name] = v
	default:
		branch, stem := p.Pull()
		if branch.List {
			t[branch.Name] = spread(stem, v)
		} else {
			if _, ok := t[branch.Name].(Tree); !ok {
				t[branch.Name] = make(Tree)
			}
			t[branch.Name].(Tree).Attach(stem, v)
		}
	}

	return t
}

func spread(p path.Path, v interface{}) []interface{} {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Array && vv.Kind() != reflect.Slice {
		v = []interface{}{v}
		vv = reflect.ValueOf(v)
	}

	var s []interface{}
	for i := 0; i < vv.Len(); i++ {
		s = append(s, make(Tree).Attach(p, vv.Index(i).Interface()))
	}
	return s
}
