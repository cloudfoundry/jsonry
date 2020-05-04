package jsonry

import (
	"fmt"
	"reflect"
	"strings"
)

type context []reflect.StructField

func (ctx context) String() string {
	t := ctx.tail()
	return fmt.Sprintf("field '%s' of type '%s' at path '%s'", t.Name, t.Type, ctx.path())
}

func (ctx context) withField(f reflect.StructField) context {
	return append(ctx, f)
}

func (ctx context) tail() reflect.StructField {
	return ctx[len(ctx)-1]
}

func (ctx context) path() string {
	var path []string
	for _, e := range ctx {
		path = append(path, e.Name)
	}
	return strings.Join(path, ".")
}
