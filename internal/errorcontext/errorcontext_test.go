package errorcontext_test

import (
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/errorcontext"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorContext", func() {
	When("it is the zero value", func() {
		It("reports being at the root", func() {
			ctx := errorcontext.ErrorContext{}
			Expect(ctx.String()).To(Equal("root path"))
		})
	})

	When("it has a field", func() {
		It("reports the field detail", func() {
			ctx := errorcontext.ErrorContext{}.WithField("Foo", reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`field "Foo" (type "string")`))
		})
	})

	When("it has an index", func() {
		It("reports the index detail", func() {
			ctx := errorcontext.ErrorContext{}.WithIndex(4, reflect.TypeOf(true))
			Expect(ctx.String()).To(Equal(`index 4 (type "bool")`))

		})
	})

	When("it has an key", func() {
		It("reports the key detail", func() {
			ctx := errorcontext.ErrorContext{}.WithKey("foo", reflect.TypeOf(true))
			Expect(ctx.String()).To(Equal(`key "foo" (type "bool")`))

		})
	})

	When("it has multiple fields", func() {
		It("reports the path details", func() {
			ctx := errorcontext.ErrorContext{}.
				WithField("Baz", reflect.TypeOf(42)).
				WithField("Bar", reflect.TypeOf(true)).
				WithField("Foo", reflect.TypeOf(""))

			Expect(ctx.String()).To(Equal(`field "Baz" (type "int") path Foo.Bar.Baz`))
		})
	})

	When("it has multiple indices", func() {
		It("reports the path details", func() {
			ctx := errorcontext.ErrorContext{}.
				WithIndex(4, reflect.TypeOf(42)).
				WithIndex(5, reflect.TypeOf(true)).
				WithIndex(3, reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`index 4 (type "int") path [3][5][4]`))
		})
	})

	When("it has multiple key", func() {
		It("reports the path details", func() {
			ctx := errorcontext.ErrorContext{}.
				WithKey("baz", reflect.TypeOf("")).
				WithKey("bar", reflect.TypeOf(42)).
				WithKey("foo", reflect.TypeOf(true))
			Expect(ctx.String()).To(Equal(`key "baz" (type "string") path ["foo"]["bar"]["baz"]`))

		})
	})

	When("it there is a mixture of fields, indices, and keys", func() {
		It("reports the details", func() {
			ctx := errorcontext.ErrorContext{}.
				WithIndex(4, reflect.TypeOf(42)).
				WithKey("bar", reflect.TypeOf(42)).
				WithField("Baz", reflect.TypeOf(42)).
				WithField("Bar", reflect.TypeOf(true)).
				WithIndex(5, reflect.TypeOf(true)).
				WithKey("foo", reflect.TypeOf(true)).
				WithField("Foo", reflect.TypeOf("")).
				WithIndex(3, reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`index 4 (type "int") path [3].Foo["foo"][5].Bar.Baz["bar"][4]`))

			ctx = errorcontext.ErrorContext{}.
				WithField("Baz", reflect.TypeOf(42)).
				WithIndex(4, reflect.TypeOf(42)).
				WithKey("bar", reflect.TypeOf(42)).
				WithField("Bar", reflect.TypeOf(true)).
				WithKey("foo", reflect.TypeOf(true)).
				WithIndex(3, reflect.TypeOf("")).
				WithIndex(5, reflect.TypeOf(true)).
				WithField("Foo", reflect.TypeOf(""))
			Expect(ctx.String()).To(Equal(`field "Baz" (type "int") path Foo[5][3]["foo"].Bar["bar"][4].Baz`))
		})
	})
})
