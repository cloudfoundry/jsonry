package jsonry_test

import (
	"fmt"

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
			Entry("int8", c{V: int8(42)}, `{"V":42}`),
			Entry("int16", c{V: int16(42)}, `{"V":42}`),
			Entry("int32", c{V: int32(42)}, `{"V":42}`),
			Entry("int64", c{V: int64(42)}, `{"V":42}`),
			Entry("uint", c{V: uint(42)}, `{"V":42}`),
			Entry("uint8", c{V: uint8(42)}, `{"V":42}`),
			Entry("uint16", c{V: uint16(42)}, `{"V":42}`),
			Entry("uint32", c{V: uint32(42)}, `{"V":42}`),
			Entry("uint64", c{V: uint64(42)}, `{"V":42}`),
			Entry("float32", c{V: float32(4.2)}, `{"V":4.2}`),
			Entry("float64", c{V: 4.2}, `{"V":4.2}`),
			Entry("struct", c{V: c{V: "hierarchical"}}, `{"V":{"V":"hierarchical"}}`),
			Entry("pointer", c{V: &i}, `{"V":0}`),
			Entry("slice", c{V: []interface{}{"hello", true, 42}}, `{"V":["hello",true,42]}`),
			Entry("array", c{V: [3]interface{}{"hello", true, 42}}, `{"V":["hello",true,42]}`),
		)

		DescribeTable(
			"unsupported types",
			func(input c, typeName string) {
				_, err := jsonry.Marshal(input)
				Expect(err).To(MatchError(fmt.Sprintf(`unsupported type "%s" at field "V" (type "interface {}")`, typeName)), err.Error())
			},
			Entry("complex", c{V: complex(1, 2)}, "complex128"),
			Entry("map", c{V: make(map[string]interface{})}, "map[string]interface {}"),
			Entry("channel", c{V: make(chan bool)}, "chan bool"),
			Entry("func", c{V: func() {}}, "func()"),
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
