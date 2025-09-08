// Package tree is a Go representation of a JSON object
package tree

import (
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/path"
)

type Tree map[string]any

func (t Tree) Attach(p path.Path, v any) Tree {
	switch p.Len() {
	case 0:
		panic("empty path")
	case 1:
		leaf, _ := p.Pull()
		t[leaf.Name] = v
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

func (t Tree) Fetch(p path.Path) (any, bool) {
	switch p.Len() {
	case 0:
		panic("empty path")
	case 1:
		leaf, _ := p.Pull()
		v, ok := t[leaf.Name]
		return v, ok
	default:
		branch, stem := p.Pull()
		v, ok := t[branch.Name]
		if !ok {
			return nil, false
		}

		switch vt := v.(type) {
		case map[string]any:
			return Tree(vt).Fetch(stem)
		case []any:
			return unspread(vt, stem), true
		default:
			return nil, false
		}
	}
}

func spread(p path.Path, v any) []any {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Array && vv.Kind() != reflect.Slice {
		v = []any{v}
		vv = reflect.ValueOf(v)
	}

	var s []any
	for i := 0; i < vv.Len(); i++ {
		s = append(s, make(Tree).Attach(p, vv.Index(i).Interface()))
	}
	return s
}

func unspread(v []any, stem path.Path) []any {
	l := make([]any, 0, len(v))
	for i := range v {
		switch vt := v[i].(type) {
		case map[string]any:
			if r, ok := Tree(vt).Fetch(stem); ok {
				l = append(l, r)
			} else {
				l = append(l, nil)
			}
		case []any:
			l = append(l, unspread(vt, stem)...)
		default:
			l = append(l, v[i])
		}
	}

	return l
}
