package parser_test

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/jsonry/internal/tokenizer"

	"code.cloudfoundry.org/jsonry/internal/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("scalar types", func() {
		It("parses null", func() {
			output, err := parser.Parse([]byte("null"))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(BeNil())
		})

		DescribeTable("types",
			func(input string, expected interface{}) {
				output, err := parser.Parse([]byte(input))
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(Equal(expected))
			},
			Entry("true", "true", true),
			Entry("false", "false", false),
			Entry("number", "42", json.Number(`42`)),
			Entry("string", `"hello"`, `hello`),
		)
	})

	Describe("array", func() {
		It("parses an empty array", func() {
			output, err := parser.Parse([]byte(`[]`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(SatisfyAll(
				HaveLen(0),
				BeAssignableToTypeOf([]interface{}{}),
			))
		})

		It("parses an array a value", func() {
			output, err := parser.Parse([]byte(`[true]`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal([]interface{}{true}))
		})

		It("parses an array with many values", func() {
			output, err := parser.Parse([]byte(`[true,false,null]`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal([]interface{}{true, false, nil}))
		})

		It("fails when a value is missing", func() {
			testError(`[,]`, `invalid character ',' looking for beginning of value`)
			testError(`[true,]`, `invalid character ']' looking for beginning of value`)
			testError(`[true,`, `unexpected end of JSON input`)
		})

		It("fails when a comma is missing", func() {
			testError(`[true`, `unexpected end of JSON input`)
			testError(`[true false]`, `invalid character 'f' after array element`)
		})

		When("the array is is not closed", func() {
			It("fails when a delimiter is missing", func() {
				testError(`[true, false`, `unexpected end of JSON input`)
			})

			It("fails when a value is missing", func() {
				testError(`[true, false,`, `unexpected end of JSON input`)
			})
		})
	})

	Describe("object", func() {
		It("parses an empty object", func() {
			output, err := parser.Parse([]byte(`{}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(map[string]interface{}{}))
		})

		It("parses an object with a key:value pair", func() {
			output, err := parser.Parse([]byte(`{"key":true}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(map[string]interface{}{"key": true}))
		})

		It("parses an object with many key:value pairs", func() {
			output, err := parser.Parse([]byte(`{"foo":true,"bar":null,"baz":false}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(map[string]interface{}{
				"foo": true,
				"bar": nil,
				"baz": false,
			}))
		})

		Context("key errors", func() {
			It("fails when a key is missing", func() {
				testError(`{`, `unexpected end of JSON input`)
			})

			It("fails when a key is non-string value", func() {
				testError(`{12`, `invalid character '1' looking for beginning of object key string`)
			})

			It("fails when a key is not a token", func() {
				testError(`{*`, `invalid character '*' looking for beginning of object key string`)
				testError(`{foo`, `invalid character 'f' looking for beginning of object key string`)
			})
		})

		Context("colon separator errors", func() {
			It("fails when a colon is missing", func() {
				testError(`{"foo"`, `unexpected end of JSON input`)
			})

			It("fails when a colon is replaced by a value", func() {
				testError(`{"foo" 42`, `invalid character '4' after object key`)
			})

			It("fails when a colon is replaced by a non-token", func() {
				testError(`{"foo" *`, `invalid character '*' after object key`)
				testError(`{"foo" foo`, `invalid character 'f' after object key`)
			})
		})

		Context("value errors", func() {
			It("fails when value is missing", func() {
				testError(`{"foo":`, `unexpected end of JSON input`)
			})

			It("fails when value is not a token", func() {
				testError(`{"foo":*`, `invalid character '*' looking for beginning of value`)
				testError(`{"foo":foo`, `invalid character 'o' in literal false (expecting 'a')`)
			})

			It("fails when the value is not a value", func() {
				testError(`{"foo":,}`, `invalid character ',' looking for beginning of value`)
			})
		})

		Context("comma separator errors", func() {
			It("fails when a value is not followed by anything", func() {
				testError(`{"foo":true`, `unexpected end of JSON input`)
			})

			It("fails when a value is followed by another value", func() {
				testError(`{"foo":true false`, `invalid character 'f' after object key:value pair`)
			})

			It("fails when a value is followed by a non-token", func() {
				testError(`{"foo":true foo`, `invalid character 'f' after object key:value pair`)
				testError(`{"foo":true*`, `invalid character '*' after object key:value pair`)
			})
		})
	})

	Describe("errors", func() {
		It("fails on blank input", func() {
			_, err := parser.Parse([]byte(" "))
			Expect(err).To(MatchError(`unexpected end of JSON input`))
		})

		DescribeTable(
			"unexpected characters",
			func(input string) {
				_, err := parser.Parse([]byte(input))
				Expect(err).To(MatchError(fmt.Sprintf(`invalid character '%s' looking for beginning of value`, input)))
			},
			Entry("non-token", "*"),
			Entry("comma", ","),
			Entry("colon", ":"),
			Entry("array close", "]"),
			Entry("object close", "}"),
		)

		DescribeTable(
			"invalid keywords",
			func(input, keyword, expect, actual string, position int) {
				_, err := parser.Parse([]byte(input))
				msg := fmt.Sprintf("invalid character '%s' in literal %s (expecting '%s')", actual, keyword, expect)
				Expect(err).To(MatchError(msg), err.Error())
				Expect(err).To(BeAssignableToTypeOf(tokenizer.InvalidKeywordError{}))
				Expect(err.(tokenizer.InvalidKeywordError).Position()).To(Equal(position))
			},
			Entry("null", " nut", "null", "l", "t", 1),
			Entry("true", "  truy", "true", "e", "y", 2),
			Entry("false", "fAlse", "false", "a", "A", 0),
			Entry("short", "  fals", "false", "e", " ", 2),
		)

		It("fails where the input unexpectedly continues", func() {
			testError("1 2", `invalid character '2' after top-level value`)
			testError("1*", `invalid character '*' after top-level value`)
			testError("1 fao", `invalid character 'f' after top-level value`)
		})
	})
})

func testError(input, message string) {
	_, err := parser.Parse([]byte(input))
	ExpectWithOffset(1, err).To(MatchError(message), err.Error())

	err = json.Unmarshal([]byte(input), nil)
	ExpectWithOffset(1, err).To(MatchError(message), err.Error())
}
