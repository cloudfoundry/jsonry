package jsonry_test

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/gomega/gstruct"

	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unmarshal", func() {
	unmarshal := func(receiver interface{}, json string) {
		err := jsonry.Unmarshal([]byte(json), receiver)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	expectToFail := func(receiver interface{}, json, message string) {
		err := jsonry.Unmarshal([]byte(json), receiver)
		ExpectWithOffset(1, err).To(MatchError(message), func() string {
			if err != nil {
				return err.Error()
			}
			return fmt.Sprintf("there was no error, unmarshaled: %+v", receiver)
		})
	}

	It("unmarshals a basic string field", func() {
		var s struct{ Foo string }
		unmarshal(&s, `{"Foo": "works"}`)
		Expect(s.Foo).To(Equal("works"))
	})

	Describe("paths", func() {
		It("defaults the path to the field name", func() {
			var s struct{ GUID string }
			unmarshal(&s, `{"GUID":"123"}`)
			Expect(s).To(MatchAllFields(Fields{"GUID": Equal("123")}))
		})

		It("will respect a JSON tag", func() {
			var s struct {
				GUID string `json:"guid"`
			}
			unmarshal(&s, `{"guid":"123"}`)
			Expect(s).To(MatchAllFields(Fields{"GUID": Equal("123")}))
		})

		It("will respect a JSONry tag", func() {
			var s struct {
				GUID string `jsonry:"relationships.spaces.guid"`
			}
			unmarshal(&s, `{"relationships":{"spaces":{"guid":"123"}}}`)
			Expect(s).To(MatchAllFields(Fields{"GUID": Equal("123")}))
		})
	})

	Describe("types", func() {
		It("unmarshals into a string field", func() {
			var s struct{ S string }
			unmarshal(&s, `{"S": "works"}`)
			Expect(s.S).To(Equal("works"))

			expectToFail(&s, `{"S": 12}`, `cannot unmarshal "12" type "number" into field "S" (type "string")`)
		})

		It("unmarshals into a bool field", func() {
			var s struct{ T, F bool }
			unmarshal(&s, `{"T":true,"F":false}`)
			Expect(s.T).To(BeTrue())
			Expect(s.F).To(BeFalse())

			expectToFail(&s, `{"T": 12}`, `cannot unmarshal "12" type "number" into field "T" (type "bool")`)
		})

		It("unmarshals into an int field", func() {
			var s struct{ I int }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(42))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "int")`)
		})

		It("unmarshals into an int8 field", func() {
			var s struct{ I int8 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(int8(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "int8")`)
		})

		It("unmarshals into an int16 field", func() {
			var s struct{ I int16 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(int16(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "int16")`)
		})

		It("unmarshals into an int32 field", func() {
			var s struct{ I int32 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(int32(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "int32")`)
		})

		It("unmarshals into an int64 field", func() {
			var s struct{ I int64 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(int64(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "int64")`)
		})

		It("unmarshals into a uint field", func() {
			var s struct{ I uint }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(uint(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "uint")`)
		})

		It("unmarshals into a uint8 field", func() {
			var s struct{ I uint8 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(uint8(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "uint8")`)
		})

		It("unmarshals into a uint16 field", func() {
			var s struct{ I uint16 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(uint16(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "uint16")`)
		})

		It("unmarshals into a uint32 field", func() {
			var s struct{ I uint32 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(uint32(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "uint32")`)
		})

		It("unmarshals into a uint64 field", func() {
			var s struct{ I uint64 }
			unmarshal(&s, `{"I":42}`)
			Expect(s.I).To(Equal(uint64(42)))

			expectToFail(&s, `{"I":"foo"}`, `cannot unmarshal "foo" type "string" into field "I" (type "uint64")`)
		})

		It("unmarshals into a float32 field", func() {
			var s struct{ A, B float32 }
			unmarshal(&s, `{"A":42,"B":4.2}`)
			Expect(s.A).To(Equal(float32(42)))
			Expect(s.B).To(Equal(float32(4.2)))

			expectToFail(&s, `{"A":"foo"}`, `cannot unmarshal "foo" type "string" into field "A" (type "float32")`)
		})

		It("unmarshals into a float64 field", func() {
			var s struct{ A, B float64 }
			unmarshal(&s, `{"A":42,"B":4.2}`)
			Expect(s.A).To(Equal(float64(42)))
			Expect(s.B).To(Equal(4.2))

			expectToFail(&s, `{"A":"foo"}`, `cannot unmarshal "foo" type "string" into field "A" (type "float64")`)
		})

		It("rejects a complex64 field", func() {
			var s struct{ C complex64 }
			expectToFail(&s, `{}`, `unsupported type "complex64" at field "C" (type "complex64")`)
		})

		It("rejects a complex128 field", func() {
			var s struct{ C complex128 }
			expectToFail(&s, `{}`, `unsupported type "complex128" at field "C" (type "complex128")`)
		})

		It("unmarshals into an interface{} field", func() {
			var s struct{ N, B, S, I, U, F, L, M interface{} }
			unmarshal(&s, `{"N":null,"B":true,"S":"foo","I":-42,"U":12,"F":4.2,"L":[1,2],"M":{"f":"b"}}`)
			Expect(s).To(MatchAllFields(Fields{
				"N": BeNil(),
				"B": BeTrue(),
				"S": Equal("foo"),
				"I": Equal(-42),
				"U": Equal(12),
				"F": Equal(4.2),
				"L": Equal([]interface{}{json.Number("1"), json.Number("2")}),
				"M": Equal(map[string]interface{}{"f": "b"}),
			}))
		})

		It("unmarshals into pointers of basic types", func() {
			var s struct {
				S, T *string
				I, J *int
			}
			unmarshal(&s, `{"S":"foo","T":null,"I":12,"J":null}`)
			Expect(s).To(MatchAllFields(Fields{
				"S": PointTo(Equal("foo")),
				"T": BeNil(),
				"I": PointTo(Equal(12)),
				"J": BeNil(),
			}))

			expectToFail(&s, `{"J":"foo"}`, `cannot unmarshal "foo" type "string" into field "J" (type "*int")`)
		})

		It("unmarshals into a slice field", func() {
			var s struct {
				S []string
				N []int
				I []interface{}
				E []string
			}
			unmarshal(&s, `{"S":["a","b","c"],"N":[1,2,3],"I":["a",2,true]}`)

			Expect(s).To(MatchAllFields(Fields{
				"S": Equal([]string{"a", "b", "c"}),
				"N": Equal([]int{1, 2, 3}),
				"I": Equal([]interface{}{"a", 2, true}),
				"E": BeEmpty(),
			}))

			expectToFail(&s, `{"S":"foo"}`, `cannot unmarshal "foo" type "string" into field "S" (type "[]string")`)
		})

		It("unmarshals into a slice pointer field", func() {
			var s struct{ S *[]string }
			unmarshal(&s, `{"S":["a","b","c"]}`)
			Expect(s.S).To(PointTo(Equal([]string{"a", "b", "c"})))

			expectToFail(&s, `{"S":"foo"}`, `cannot unmarshal "foo" type "string" into field "S" (type "*[]string")`)
		})

		It("rejects an array field", func() {
			var s struct{ S [3]string }
			expectToFail(&s, `{}`, `unsupported type "[3]string" at field "S" (type "[3]string")`)
		})

		Context("maps", func() {
			It("unmarshals maps with interface values", func() {
				By("map", func() {
					var s struct{ I map[string]interface{} }
					unmarshal(&s, `{"I":{"a":"b","c":5,"d":true}}`)
					Expect(s.I).To(Equal(map[string]interface{}{"a": "b", "c": 5, "d": true}))
				})

				By("pointer", func() {
					var s struct{ I *map[string]interface{} }
					unmarshal(&s, `{"I":{"a":"b","c":5,"d":true}}`)
					Expect(s.I).To(PointTo(Equal(map[string]interface{}{"a": "b", "c": 5, "d": true})))
				})
			})

			It("unmarshals maps with string values", func() {
				By("map", func() {
					var s struct{ S map[string]string }
					unmarshal(&s, `{"S":{"a":"b","c":"d"}}`)
					Expect(s.S).To(Equal(map[string]string{"a": "b", "c": "d"}))
				})

				By("pointer", func() {
					var s struct{ S *map[string]string }
					unmarshal(&s, `{"S":{"a":"b","c":"d"}}`)
					Expect(s.S).To(PointTo(Equal(map[string]string{"a": "b", "c": "d"})))
				})
			})

			It("unmarshals maps with number values", func() {
				By("map", func() {
					var s struct{ N map[string]int }
					unmarshal(&s, `{"N":{"f":5}}`)
					Expect(s.N).To(Equal(map[string]int{"f": 5}))
				})

				By("pointer", func() {
					var s struct{ N *map[string]int }
					unmarshal(&s, `{"N":{"f":5}}`)
					Expect(s.N).To(PointTo(Equal(map[string]int{"f": 5})))
				})
			})

			It("unmarshals omitted maps", func() {
				By("map", func() {
					var s struct{ I map[string]interface{} }
					unmarshal(&s, `{}`)
					Expect(s.I).To(BeNil())
					Expect(s.I).To(BeEmpty())
				})

				By("pointer", func() {
					var s struct{ I *map[string]interface{} }
					unmarshal(&s, `{}`)
					Expect(s.I).To(BeNil())
				})
			})

			It("unmarshals null maps", func() {
				By("map", func() {
					var s struct{ I map[string]interface{} }
					unmarshal(&s, `{"I": null}`)
					Expect(s.I).To(BeNil())
					Expect(s.I).To(BeEmpty())
				})

				By("pointer", func() {
					var s struct{ I *map[string]interface{} }
					unmarshal(&s, `{"I": null}`)
					Expect(s.I).To(BeNil())
				})
			})

			It("unmarshals empty maps", func() {
				By("map", func() {
					var s struct{ I map[string]interface{} }
					unmarshal(&s, `{"I": {}}`)
					Expect(s.I).NotTo(BeNil())
					Expect(s.I).To(BeEmpty())
				})

				By("pointer", func() {
					var s struct{ I *map[string]interface{} }
					unmarshal(&s, `{"I": {}}`)
					Expect(s.I).NotTo(BeNil())
					Expect(s.I).To(PointTo(Not(BeNil())))
					Expect(s.I).To(PointTo(BeEmpty()))
				})
			})
		})

		It("rejects an map field that does not have string keys", func() {
			var s struct{ S map[int]string }
			expectToFail(&s, `{}`, `maps must only have string keys for "int" at field "S" (type "map[int]string")`)
		})

		It("unmarshals into json.Unmarshaler field", func() {
			var s struct{ S implementsJSONUnmarshaler }
			unmarshal(&s, `{"S":"ok"}`)
			Expect(s.S).To(Equal(implementsJSONUnmarshaler{hasBeenSet: true}))

			expectToFail(&s, `{"S":"fail"}`, `error from UnmarshalJSON() call at field "S" (type "jsonry_test.implementsJSONUnmarshaler"): ouch`)
		})

		It("unmarshals into named types and type aliases", func() {
			type alias = string
			type named string
			var s struct {
				A alias
				N named
			}
			unmarshal(&s, `{"A":"foo","N":"bar"}`)
			Expect(s.A).To(Equal("foo"))
			Expect(s.N).To(Equal(named("bar")))

			expectToFail(&s, `{"A":12}`, `cannot unmarshal "12" type "number" into field "A" (type "string")`)
			expectToFail(&s, `{"N":13}`, `cannot unmarshal "13" type "number" into field "N" (type "jsonry_test.named")`)
		})

		When("unmarshalling null", func() {
			It("leaves basic types untouched", func() {
				s := struct {
					S string
					T int
					U uint
					V float64
					W bool
				}{
					S: "foo",
					T: -65,
					U: 12,
					V: 3.14,
					W: true,
				}
				unmarshal(&s, `{"S": null, "T": null, "U": null, "V": null, "W": null}`)
				Expect(s.S).To(Equal("foo"))
				Expect(s.T).To(Equal(-65))
				Expect(s.U).To(Equal(uint(12)))
				Expect(s.V).To(Equal(3.14))
				Expect(s.W).To(BeTrue())
			})

			It("overwrites a pointer as nil", func() {
				v := "hello"
				s := struct{ S *string }{S: &v}
				unmarshal(&s, `{"S": null}`)
				Expect(s.S).To(BeNil())
			})

			It("overwrites an interface{} as nil", func() {
				s := struct{ S interface{} }{S: "hello"}
				unmarshal(&s, `{"S": null}`)
				Expect(s.S).To(BeNil())
			})
		})
	})

	Describe("recursive composition", func() {
		It("unmarshals into a struct field", func() {
			type t struct{ S string }
			var s struct{ T t }
			unmarshal(&s, `{"T":{"S":"foo"}}`)
			Expect(s.T.S).To(Equal("foo"))

			expectToFail(&s, `{"T":"foo"}`, `cannot unmarshal "foo" type "string" into field "T" (type "jsonry_test.t")`)
		})

		It("unmarshals into a struct pointer field", func() {
			type t struct{ S string }
			var s struct{ T *t }
			unmarshal(&s, `{"T":{"S":"foo"}}`)
			Expect(s.T.S).To(Equal("foo"))

			expectToFail(&s, `{"T":"foo"}`, `cannot unmarshal "foo" type "string" into field "T" (type "*jsonry_test.t")`)
		})

		It("unmarshals a slice of structs", func() {
			type t struct{ S string }
			var s struct{ T []t }
			unmarshal(&s, `{"T":[{"S":"foo"},{"S":"bar"},{},{"S":"baz"}]}`)
			Expect(s.T).To(Equal([]t{{S: "foo"}, {S: "bar"}, {}, {S: "baz"}}))

			expectToFail(&s, `{"T":[null]}`, `cannot unmarshal "<nil>" into index 0 (type "jsonry_test.t") path T[0]`)
		})

		It("unmarshals a map of structs", func() {
			type t struct{ S string }
			var s struct{ T map[string]t }
			unmarshal(&s, `{"T":{"foo":{"S":"alpha"},"bar":{"S":"beta"}}}`)
			Expect(s.T).To(Equal(map[string]t{"foo": {S: "alpha"}, "bar": {S: "beta"}}))

			expectToFail(&s, `{"T":5}`, `cannot unmarshal "5" type "number" into field "T" (type "map[string]jsonry_test.t")`)
		})
	})

	Describe("receiver", func() {
		It("accept a struct pointer", func() {
			var s struct{}
			err := jsonry.Unmarshal([]byte(`{}`), &s)
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects a struct", func() {
			var s struct{}
			err := jsonry.Unmarshal([]byte(`{}`), s)
			Expect(err).To(MatchError("receiver must be a pointer to a struct, got a non-pointer"))
		})

		It("rejects a pointer to a non-struct", func() {
			var s int
			err := jsonry.Unmarshal([]byte(`{}`), &s)
			Expect(err).To(MatchError("receiver must be a pointer to a struct type, got: int"))
		})
	})
})
