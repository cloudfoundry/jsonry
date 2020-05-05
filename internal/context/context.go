package context

import (
	"fmt"
	"reflect"
)

type sort uint

const (
	field sort = iota
	index
	key
)

type segment struct {
	sort  sort
	name  string
	index int
	typ   reflect.Type
}

func (s segment) String() string {
	switch s.sort {
	case field:
		return fmt.Sprintf(`field "%s" (type "%s")`, s.name, s.typ)
	case index:
		return fmt.Sprintf(`index %d (type "%s")`, s.index, s.typ)
	default:
		return fmt.Sprintf(`key "%s" (type "%s")`, s.name, s.typ)
	}
}

type Context []segment

func (ctx Context) String() string {
	switch len(ctx) {
	case 0:
		return "at the root"
	case 1:
		return fmt.Sprintf("at %s", ctx.leaf())
	default:
		return fmt.Sprintf("at %s path %s", ctx.leaf(), ctx.path())
	}
}

func (ctx Context) leaf() segment {
	return ctx[len(ctx)-1]
}

func (ctx Context) path() string {
	var path string
	for _, s := range ctx {
		switch s.sort {
		case index:
			path = fmt.Sprintf("%s[%d]", path, s.index)
		case field:
			if len(path) > 0 {
				path = path + "."
			}
			path = path + s.name
		case key:
			path = fmt.Sprintf(`%s["%s"]`, path, s.name)
		}
	}

	return path
}

func (ctx Context) WithField(f reflect.StructField) Context {
	return append(ctx, segment{sort: field, name: f.Name, typ: f.Type})
}

func (ctx Context) WithIndex(i int, t reflect.Type) Context {
	return append(ctx, segment{sort: index, index: i, typ: t})
}

func (ctx Context) WithKey(k string, t reflect.Type) Context {
	return append(ctx, segment{sort: key, name: k, typ: t})
}
