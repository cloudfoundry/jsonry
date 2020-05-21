package jsonry_test

import (
	"errors"

	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marshal", func() {
	expectToMarshal := func(input interface{}, expected string) {
		out, err := jsonry.Marshal(input)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, out).To(MatchJSON(expected))
	}

	expectToFail := func(input interface{}, message string) {
		_, err := jsonry.Marshal(input)
		ExpectWithOffset(1, err).To(MatchError(message), func() string {
			if err != nil {
				return err.Error()
			}
			return "there was no error"
		})
	}

	It("marshals a basic string field", func() {
		expectToMarshal(struct{ Foo string }{Foo: "works"}, `{"Foo": "works"}`)
	})

	Describe("paths", func() {
		It("defaults the path to the field name", func() {
			s := struct{ GUID string }{GUID: "123"}
			expectToMarshal(s, `{"GUID":"123"}`)
		})

		It("will respect a JSON tag", func() {
			s := struct {
				GUID string `json:"guid"`
			}{GUID: "123"}
			expectToMarshal(s, `{"guid":"123"}`)
		})

		It("will respect a JSONry tag", func() {
			s := struct {
				GUID string `jsonry:"relationships.spaces[].guid"`
			}{GUID: "123"}
			expectToMarshal(s, `{"relationships":{"spaces":[{"guid":"123"}]}}`)
		})
	})

	Describe("types", func() {
		It("marshals a string", func() {
			expectToMarshal(struct{ S string }{S: "hello"}, `{"S":"hello"}`)
		})

		It("marshals a boolean", func() {
			expectToMarshal(struct{ T, F bool }{T: true, F: false}, `{"T":true, "F":false}`)
		})

		It("marshals an int", func() {
			expectToMarshal(struct{ I int }{I: 42}, `{"I":42}`)
		})

		It("marshals an int8", func() {
			expectToMarshal(struct{ I int8 }{I: 42}, `{"I":42}`)
		})

		It("marshals an int16", func() {
			expectToMarshal(struct{ I int16 }{I: 42}, `{"I":42}`)
		})

		It("marshals an int32", func() {
			expectToMarshal(struct{ I int32 }{I: 42}, `{"I":42}`)
		})

		It("marshals an int64", func() {
			expectToMarshal(struct{ I int64 }{I: 42}, `{"I":42}`)
		})

		It("marshals a uint", func() {
			expectToMarshal(struct{ I uint }{I: 42}, `{"I":42}`)
		})

		It("marshals a uint8", func() {
			expectToMarshal(struct{ I uint8 }{I: 42}, `{"I":42}`)
		})

		It("marshals a uint16", func() {
			expectToMarshal(struct{ I uint16 }{I: 42}, `{"I":42}`)
		})

		It("marshals a uint32", func() {
			expectToMarshal(struct{ I uint32 }{I: 42}, `{"I":42}`)
		})

		It("marshals a uint64", func() {
			expectToMarshal(struct{ I uint64 }{I: 42}, `{"I":42}`)
		})

		It("marshals a float32", func() {
			expectToMarshal(struct{ F float32 }{F: 4.2}, `{"F":4.2}`)
		})

		It("marshals a float64", func() {
			expectToMarshal(struct{ F float64 }{F: 4.2}, `{"F":4.2}`)
		})

		It("does not marshal a complex64", func() {
			expectToFail(struct{ C complex64 }{C: complex(1, 2)}, `unsupported type "complex64" at field "C" (type "complex64")`)
		})

		It("does not marshal a complex64", func() {
			expectToFail(struct{ C complex128 }{C: complex(1, 2)}, `unsupported type "complex128" at field "C" (type "complex128")`)
		})

		It("does not marshal a channel", func() {
			expectToFail(struct{ C chan bool }{C: make(chan bool)}, `unsupported type "chan bool" at field "C" (type "chan bool")`)
		})

		It("does not marshal a function", func() {
			expectToFail(struct{ F func() }{F: func() {}}, `unsupported type "func()" at field "F" (type "func()")`)
		})

		It("marshals via a pointer", func() {
			i := 42
			expectToMarshal(struct{ P *int }{P: &i}, `{"P":42}`)
		})

		It("marshals a nil interface", func() {
			expectToMarshal(struct{ N interface{} }{N: nil}, `{"N":null}`)
		})

		It("marshals a nil pointer", func() {
			expectToMarshal(struct{ N *string }{N: nil}, `{"N":null}`)
		})

		It("marshals a struct with a private field", func() {
			expectToMarshal(struct{ P, p string }{P: "foo", p: "bar"}, `{"P":"foo"}`)
		})

		It("marshals a slice", func() {
			s := []interface{}{"hello", true, 42}
			expectToMarshal(struct{ S []interface{} }{S: s}, `{"S":["hello",true,42]}`)
			expectToMarshal(struct {
				S *[]interface{}
			}{S: &s}, `{"S":["hello",true,42]}`)
		})

		It("marshals an array", func() {
			s := [3]interface{}{"hello", true, 42}
			expectToMarshal(struct{ S [3]interface{} }{S: s}, `{"S":["hello",true,42]}`)
			expectToMarshal(struct {
				S *[3]interface{}
			}{S: &s}, `{"S":["hello",true,42]}`)
		})

		It("marshals a map", func() {
			mi := map[string]interface{}{"foo": "hello", "bar": true, "baz": 42}
			ms := map[string]string{"foo": "hello", "bar": "true", "baz": "42"}
			mn := map[int]interface{}{4: 3}

			expectToMarshal(struct{ M map[string]interface{} }{M: mi}, `{"M":{"foo":"hello","bar":true,"baz":42}}`)
			expectToMarshal(struct{ M *map[string]interface{} }{M: &mi}, `{"M":{"foo":"hello","bar":true,"baz":42}}`)
			expectToMarshal(struct{ M map[string]string }{M: ms}, `{"M":{"foo":"hello","bar":"true","baz":"42"}}`)
			expectToMarshal(struct{ M *map[string]string }{M: &ms}, `{"M":{"foo":"hello","bar":"true","baz":"42"}}`)

			expectToFail(struct{ M map[int]interface{} }{M: mn}, `maps must only have string keys for "map[int]interface {}" at field "M" (type "map[int]interface {}")`)
		})

		It("marshals a json.Marshaler", func() {
			expectToMarshal(struct{ I implementsJSONMarshaler }{I: implementsJSONMarshaler{bytes: []byte(`"hello"`)}}, `{"I":"hello"}`)
			expectToMarshal(struct{ I *implementsJSONMarshaler }{I: &implementsJSONMarshaler{bytes: []byte(`"hello"`)}}, `{"I":"hello"}`)
			expectToMarshal(struct{ I *implementsJSONMarshaler }{I: (*implementsJSONMarshaler)(nil)}, `{"I":null}`)

			expectToFail(struct{ I implementsJSONMarshaler }{I: implementsJSONMarshaler{err: errors.New("ouch")}}, `error from MarshaJSON() call at field "I" (type "jsonry_test.implementsJSONMarshaler"): ouch`)
			expectToFail(struct{ I implementsJSONMarshaler }{I: implementsJSONMarshaler{}}, `error parsing MarshaJSON() output "" at field "I" (type "jsonry_test.implementsJSONMarshaler"): unexpected end of JSON input`)
		})

		It("marshals from named types and type aliases", func() {
			type alias = string
			type named string
			s := struct {
				A alias
				N named
			}{
				A: "foo",
				N: named("bar"),
			}
			expectToMarshal(s, `{"A":"foo","N":"bar"}`)
		})
	})

	Describe("recursive composition", func() {
		type i struct{ S string }

		It("marshals a struct within a struct", func() {
			expectToMarshal(struct{ T i }{T: i{S: "foo"}}, `{"T":{"S":"foo"}}`)
		})

		It("marshals a struct within a slice", func() {
			expectToMarshal(struct{ T []i }{T: []i{{S: "foo"}, {S: "bar"}}}, `{"T":[{"S":"foo"},{"S":"bar"}]}`)
		})

		It("marshals a struct within a map", func() {
			expectToMarshal(struct{ T map[string]i }{T: map[string]i{"A": {S: "foo"}, "B": {S: "bar"}}}, `{"T":{"A":{"S":"foo"},"B":{"S":"bar"}}}`)
		})
	})

	Describe("omitempty", func() {
		It("omits zero values of basic types", func() {
			s := struct {
				A string `json:",omitempty"`
				B string `json:"bee,omitempty"`
				C string `jsonry:",omitempty"`
				D string `jsonry:"dee,omitempty"`
				E string
			}{}
			expectToMarshal(s, `{"E":""}`)
		})

		It("omits zero value structs", func() {
			type t struct{ A string }
			s := struct {
				B t `jsonry:",omitempty"`
			}{}
			expectToMarshal(s, `{}`)
		})

		It("omits empty lists", func() {
			s := struct {
				A []string  `jsonry:",omitempty"`
				D [0]string `jsonry:",omitempty"`
			}{}
			expectToMarshal(s, `{}`)
		})

		It("omits empty maps", func() {
			s := struct {
				A map[interface{}]interface{} `jsonry:",omitempty"`
				D map[int]int                 `jsonry:",omitempty"`
			}{}
			expectToMarshal(s, `{}`)
		})

		It("omits nil pointers", func() {
			s := struct {
				A *string `json:",omitempty"`
				B *string `json:"bee,omitempty"`
				C *string `jsonry:",omitempty"`
				D *string `jsonry:"dee,omitempty"`
				E *string
			}{}
			expectToMarshal(s, `{"E":null}`)
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
