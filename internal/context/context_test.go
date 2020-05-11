package context_test

import (
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	When("it is the zero value", func() {
		It("reports being at the root", func() {
			ctx := context.Context{}
			Expect(ctx.String()).To(Equal("root path"))
		})
	})

	When("it has a field", func() {
		It("reports the field detail", func() {
			ctx := context.Context{}.WithField("Foo", reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`field "Foo" (type "string")`))
		})
	})

	When("it has an index", func() {
		It("reports the index detail", func() {
			ctx := context.Context{}.WithIndex(4, reflect.TypeOf(true))
			Expect(ctx.String()).To(Equal(`index 4 (type "bool")`))

		})
	})

	When("it has an key", func() {
		It("reports the key detail", func() {
			ctx := context.Context{}.WithKey("foo", reflect.TypeOf(true))
			Expect(ctx.String()).To(Equal(`key "foo" (type "bool")`))

		})
	})

	When("it has multiple fields", func() {
		It("reports the path details", func() {
			ctx := context.Context{}.
				WithField("Foo", reflect.TypeOf("")).
				WithField("Bar", reflect.TypeOf(true)).
				WithField("Baz", reflect.TypeOf(42))

			Expect(ctx.String()).To(Equal(`field "Baz" (type "int") path Foo.Bar.Baz`))
		})
	})

	When("it has multiple indices", func() {
		It("reports the path details", func() {
			ctx := context.Context{}.
				WithIndex(3, reflect.TypeOf("")).
				WithIndex(5, reflect.TypeOf(true)).
				WithIndex(4, reflect.TypeOf(42))
			Expect(ctx.String()).To(Equal(`index 4 (type "int") path [3][5][4]`))
		})
	})

	When("it has multiple key", func() {
		It("reports the path details", func() {
			ctx := context.Context{}.
				WithKey("foo", reflect.TypeOf(true)).
				WithKey("bar", reflect.TypeOf(42)).
				WithKey("baz", reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`key "baz" (type "string") path ["foo"]["bar"]["baz"]`))

		})
	})

	When("it there is a mixture of fields, indices, and keys", func() {
		It("reports the details", func() {
			ctx := context.Context{}.
				WithIndex(3, reflect.TypeOf("")).
				WithField("Foo", reflect.TypeOf("")).
				WithKey("foo", reflect.TypeOf(true)).
				WithIndex(5, reflect.TypeOf(true)).
				WithField("Bar", reflect.TypeOf(true)).
				WithField("Baz", reflect.TypeOf(42)).
				WithKey("bar", reflect.TypeOf(42)).
				WithIndex(4, reflect.TypeOf(42))
			Expect(ctx.String()).To(Equal(`index 4 (type "int") path [3].Foo["foo"][5].Bar.Baz["bar"][4]`))

			ctx = context.Context{}.
				WithField("Foo", reflect.TypeOf("")).
				WithIndex(5, reflect.TypeOf(true)).
				WithIndex(3, reflect.TypeOf("")).
				WithKey("foo", reflect.TypeOf(true)).
				WithField("Bar", reflect.TypeOf(true)).
				WithKey("bar", reflect.TypeOf(42)).
				WithIndex(4, reflect.TypeOf(42)).
				WithField("Baz", reflect.TypeOf(42))
			Expect(ctx.String()).To(Equal(`field "Baz" (type "int") path Foo[5][3]["foo"].Bar["bar"][4].Baz`))
		})
	})

})
