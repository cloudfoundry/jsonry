package jsonry_test

import (
	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshal", func() {
	It("marshals a basic string field", func() {
		s := struct{ Foo string }{Foo: "works"}

		json, err := jsonry.Marshal(s)
		Expect(err).NotTo(HaveOccurred())
		Expect(json).To(MatchJSON(`{"Foo": "works"}`))
	})

	Describe("paths", func() {
		It("defaults the path to the field name", func() {
			s := struct{ GUID string }{GUID: "123"}
			r, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(MatchJSON(`{"GUID":"123"}`))
		})

		It("will respect a JSON tag", func() {
			s := struct {
				GUID string `json:"guid"`
			}{GUID: "123"}
			r, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(MatchJSON(`{"guid":"123"}`))
		})

		It("will respect a JSONry tag", func() {
			s := struct {
				GUID string `jsonry:"relationships.spaces[].guid"`
			}{GUID: "123"}
			r, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(MatchJSON(`{"relationships":{"spaces":[{"guid":"123"}]}}`))
		})
	})

	Describe("types", func() {
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
			Entry("struct with private field", c{V: pri{private: true, Public: true}}, `{"V":{"Public":true}}`),
			Entry("pointer", c{V: &i}, `{"V":0}`),
			Entry("nil pointer", c{V: (*int)(nil)}, `{"V":null}`),
			Entry("slice", c{V: []interface{}{"hello", true, 42}}, `{"V":["hello",true,42]}`),
			Entry("array", c{V: [3]interface{}{"hello", true, 42}}, `{"V":["hello",true,42]}`),
			Entry("map of interfaces", c{V: map[string]interface{}{"foo": "hello", "bar": true, "baz": 42}}, `{"V":{"foo":"hello","bar":true,"baz":42}}`),
			Entry("map of strings", c{V: map[string]string{"foo": "hello", "bar": "true", "baz": "42"}}, `{"V":{"foo":"hello","bar":"true","baz":"42"}}`),
			Entry("json.Marshaler", c{V: jsm{}}, `{"V":"hello"}`),
			Entry("jsonry.Marshaler", c{V: jrm{}}, `{"V":{"Public":true}}`),
		)

		DescribeTable(
			"type failure cases",
			func(input c, message string) {
				_, err := jsonry.Marshal(input)
				Expect(err).To(MatchError(message), func() string {
					if err != nil {
						return err.Error()
					}
					return "there was no error"
				})
			},
			Entry("complex", c{V: complex(1, 2)}, `unsupported type "complex128" at field "V" (type "interface {}")`),
			Entry("channel", c{V: make(chan bool)}, `unsupported type "chan bool" at field "V" (type "interface {}")`),
			Entry("func", c{V: func() {}}, `unsupported type "func()" at field "V" (type "interface {}")`),
			Entry("map with non-string keys", c{V: map[int]interface{}{4: 3}}, `maps must only have strings keys for "map[int]interface {}" at field "V" (type "interface {}")`),
			Entry("json.Marshaler with failure", c{V: jsm{value: true}}, `error from MarshaJSON() call at field "V" (type "interface {}"): ouch`),
			Entry("jsonry.Marshaler with failure", c{V: jrm{value: true}}, `error from MarshaJSONry() call at field "V" (type "interface {}"): ouch`),
		)
	})

	Describe("omitempty", func() {
		It("omits zero values", func() {
			s := struct {
				A string `json:",omitempty"`
				B string `json:"bee,omitempty"`
				C string `jsonry:",omitempty"`
				D string `jsonry:"dee,omitempty"`
				E string
			}{}
			r, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(MatchJSON(`{"E":""}`))
		})

		It("omits nil pointers", func() {
			s := struct {
				A *string `json:",omitempty"`
				B *string `json:"bee,omitempty"`
				C *string `jsonry:",omitempty"`
				D *string `jsonry:"dee,omitempty"`
				E *string
			}{}
			r, err := jsonry.Marshal(s)
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(MatchJSON(`{"E":null}`))
		})
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
			Expect(err).To(MatchError(`the input must be a struct, not "int"`))
		})

		It("rejects a nil pointer", func() {
			type s struct{}
			var sp *s
			_, err := jsonry.Marshal(sp)
			Expect(err).To(MatchError(`the input must be a struct, not "invalid"`))
		})
	})
})
