package jsonry_test

import (
	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshal", func() {
	Describe("supported types", func() {
		type c struct{ V interface{} }
		var i int

		DescribeTable(
			"supported types",
			func(input c, expected string) {
				out, err := jsonry.Marshal(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(MatchJSON(expected))
			},
			Entry("string", c{V: "hello"}, `{"V":"hello"}`),
			Entry("boolean", c{V: true}, `{"V":true}`),
			Entry("int", c{V: 42}, `{"V":42}`),
			Entry("float", c{V: 4.2}, `{"V":4.2}`),
		)

		DescribeTable(
			"unsupported types",
			func(input c, typeName string) {
				_, err := jsonry.Marshal(input)
				Expect(err).To(MatchError(ContainSubstring("unsupported type: %s", typeName)))
			},
			Entry("complex", c{V: complex(1, 2)}, "complex"),
			Entry("slice", c{V: []string{"hello"}}, "[]string"),
			Entry("array", c{V: [1]string{"hello"}}, "[1]string"),
			Entry("map", c{V: make(map[string]interface{})}, "map[string]interface {}"),
			Entry("channel", c{V: make(chan bool)}, "chan bool"),
			Entry("func", c{V: func() {}}, "func()"),
			Entry("pointer", c{V: &i}, "*int"),
		)
	})

	Describe("inputs", func() {
		It("accept a struct", func() {
			var s struct{}
			_, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
		})

		It("accept a struct pointer", func() {
			var s struct{}
			_, err := jsonry.Marshal(&s)
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects a non-struct value", func() {
			_, err := jsonry.Marshal(42)
			Expect(err).To(MatchError("the input must be a struct"))
		})

		It("rejects a nil pointer", func() {
			type s struct{}
			var sp *s
			_, err := jsonry.Marshal(sp)
			Expect(err).To(MatchError("the input must be a struct"))
		})
	})
})
