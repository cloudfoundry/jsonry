package jsonry_test

import (
	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Unmarshal", func() {
	It("marshals a basic string field", func() {
		var s struct{ Foo string }

		err := jsonry.Unmarshal([]byte(`{"Foo": "works"}`), &s)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).To(MatchAllFields(Fields{
			"Foo": Equal("works"),
		}))
	})

	Describe("paths", func() {
		It("defaults the path to the field name", func() {
			var s struct{ GUID string }
			err := jsonry.Unmarshal([]byte(`{"GUID":"123"}`), &s)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(MatchAllFields(Fields{
				"GUID": Equal("123"),
			}))
		})

		It("will respect a JSON tag", func() {
			var s struct {
				GUID string `json:"guid"`
			}
			err := jsonry.Unmarshal([]byte(`{"guid":"123"}`), &s)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(MatchAllFields(Fields{
				"GUID": Equal("123"),
			}))
		})

		It("will respect a JSONry tag", func() {
			var s struct {
				GUID string `jsonry:"relationships.spaces.guid"`
			}
			err := jsonry.Unmarshal([]byte(`{"relationships":{"spaces":{"guid":"123"}}}`), &s)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(MatchAllFields(Fields{
				"GUID": Equal("123"),
			}))
		})
	})

	Describe("types", func() {
		type b struct {
			S string
		}

		type c struct {
			S   string
			B   bool
			I   int
			I8  int8
			I16 int16
			I32 int32
			I64 int64
			U   uint
			U8  uint8
			U16 uint16
			U32 uint32
			U64 uint64
			F32 float32
			F64 float64
			In  interface{}
			Ss  []string
			Bs  []b
			V   b
			PI  *int
			PV  *b
		}

		i := 42

		DescribeTable(
			"supported types",
			func(input string, expected c) {
				var result c
				err := jsonry.Unmarshal([]byte(input), &result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expected))
			},
			Entry("string", `{"S":"hello"}`, c{S: "hello"}),
			Entry("boolean", `{"B":true}`, c{B: true}),
			Entry("int", `{"I":42}`, c{I: 42}),
			Entry("int8", `{"I8":42}`, c{I8: int8(42)}),
			Entry("int16", `{"I16":42}`, c{I16: int16(42)}),
			Entry("int32", `{"I32":42}`, c{I32: int32(42)}),
			Entry("int64", `{"I64":42}`, c{I64: int64(42)}),
			Entry("uint", `{"U":42}`, c{U: uint(42)}),
			Entry("uint8", `{"U8":42}`, c{U8: uint8(42)}),
			Entry("uint16", `{"U16":42}`, c{U16: uint16(42)}),
			Entry("uint32", `{"U32":42}`, c{U32: uint32(42)}),
			Entry("uint64", `{"U64":42}`, c{U64: uint64(42)}),
			Entry("float32", `{"F32":4.2}`, c{F32: float32(4.2)}),
			Entry("float64", `{"F64":4.2}`, c{F64: 4.2}),
			Entry("interface as basic", `{"In":"hello"}`, c{In: "hello"}),
			Entry("interface as nil", `{"In":null}`, c{In: nil}),
			Entry("interface as map", `{"In":{"foo":"bar"}}}`, c{In: map[string]interface{}{"foo": "bar"}}),
			Entry("pointer", `{"PI":42}`, c{PI: &i}),
			Entry("nil pointer", `{"PI":null}`, c{PI: nil}),
			Entry("struct", `{"V":{"S":"hierarchical"}}`, c{V: b{S: "hierarchical"}}),
			Entry("basic list", `{"Ss":["foo","bar","baz"]}`, c{Ss: []string{"foo", "bar", "baz"}}),
			Entry("struct list", `{"Bs":[{"S":"foo"},{"S":"bar"},{"S":"baz"}]}`, c{Bs: []b{{S: "foo"}, {S: "bar"}, {S: "baz"}}}),
		)

		DescribeTable(
			"type failure cases",
			func(input, message string) {
				var result c
				err := jsonry.Unmarshal([]byte(input), &result)
				Expect(err).To(MatchError(message), func() string {
					if err != nil {
						return err.Error()
					}
					return "there was no error"
				})
			},
			Entry("basic type mismatch", `{"I":"hello"}`, `cannot unmarshal "hello" type "string" into field "I" (type "int")`),
			Entry("nil into non-pointer", `{"I":null}`, `cannot unmarshal "<nil>" into field "I" (type "int")`),
			Entry("float into int", `{"I":4.2}`, `cannot unmarshal "4.2" into field "I" (type "int")`),
		)
	})

	Describe("outputs", func() {
		It("accept a struct pointer", func() {
			var s struct{}
			err := jsonry.Unmarshal([]byte(`{}`), &s)
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects a struct", func() {
			var s struct{}
			err := jsonry.Unmarshal([]byte(`{}`), s)
			Expect(err).To(MatchError("output must be a pointer to a struct, got a non-pointer"))
		})

		It("rejects a pointer to a non-struct", func() {
			var s int
			err := jsonry.Unmarshal([]byte(`{}`), &s)
			Expect(err).To(MatchError("output must be a pointer to a struct type, got: int"))
		})
	})
})
