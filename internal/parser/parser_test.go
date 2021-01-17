package parser_test

import (
	"code.cloudfoundry.org/jsonry/internal/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	XDescribe("scalar types", func() {
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
			Entry("number", "42", parser.Number("42")),
			Entry("string", `"hello"`, parser.String(`"hello"`)),
		)
	})

	XDescribe("array", func() {
		It("parses an array", func() {
			output, err := parser.Parse([]byte(`[true,false,null]`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal([]interface{}{true, false, nil}))
		})
	})

	XDescribe("object", func() {
		It("parses an object", func() {
			output, err := parser.Parse([]byte(`{"foo":true,"bar":null,"baz":false}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(map[string]interface{}{
				"foo": true,
				"bar": nil,
				"baz": false,
			}))
		})

		// non-string keys
		// trailing ,
	})

	Describe("errors", func() {
		Context("blank input", func() {
			It("fails", func() {
				_, err := parser.Parse([]byte(" "))
				Expect(err).To(MatchError(`unexpected end of JSON input`))
			})
		})
		// stuff at the end
		// wrapper tokenizer errors
	})
})
