package tokenizer_test

import (
	"fmt"

	"code.cloudfoundry.org/jsonry/internal/tokenizer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokenizer", func() {
	next := func(input string) tokenizer.Token {
		tok, err := tokenizer.New([]byte(input)).Next()
		Expect(err).NotTo(HaveOccurred())
		return tok
	}

	all := func(input string) (result []tokenizer.Token) {
		tnz := tokenizer.New([]byte(input))
		for {
			tok, err := tnz.Next()
			Expect(err).NotTo(HaveOccurred())
			result = append(result, tok)
			if tok.Type == tokenizer.End {
				return
			}
		}
	}

	Describe("token identification", func() {
		DescribeTable("token types",
			func(input string, tokenType tokenizer.TokenType, value []byte) {
				Expect(next(input)).To(Equal(tokenizer.Token{
					Type:   tokenType,
					Start:  0,
					Length: len(input),
					Value:  value,
				}))
			},
			Entry("end", "", tokenizer.End, nil),
			Entry("{", "{", tokenizer.ObjectOpen, nil),
			Entry("}", "}", tokenizer.ObjectClose, nil),
			Entry("[", "[", tokenizer.ArrayOpen, nil),
			Entry("]", "]", tokenizer.ArrayClose, nil),
			Entry(":", ":", tokenizer.Colon, nil),
			Entry(",", ",", tokenizer.Comma, nil),
			Entry("null", "null", tokenizer.Null, nil),
			Entry("true", "true", tokenizer.True, nil),
			Entry("false", "false", tokenizer.False, nil),
			Entry("string", `"foo"`, tokenizer.String, []byte("foo")),
			Entry("number", "42", tokenizer.Number, []byte("42")),
		)

		It("handles multiple tokens", func() {
			Expect(all(`{}[]null,true:false"foo"42`)).To(Equal([]tokenizer.Token{
				{Type: tokenizer.ObjectOpen, Start: 0, Length: 1},
				{Type: tokenizer.ObjectClose, Start: 1, Length: 1},
				{Type: tokenizer.ArrayOpen, Start: 2, Length: 1},
				{Type: tokenizer.ArrayClose, Start: 3, Length: 1},
				{Type: tokenizer.Null, Start: 4, Length: 4},
				{Type: tokenizer.Comma, Start: 8, Length: 1},
				{Type: tokenizer.True, Start: 9, Length: 4},
				{Type: tokenizer.Colon, Start: 13, Length: 1},
				{Type: tokenizer.False, Start: 14, Length: 5},
				{Type: tokenizer.String, Start: 19, Length: 5, Value: []byte("foo")},
				{Type: tokenizer.Number, Start: 24, Length: 2, Value: []byte("42")},
				{Type: tokenizer.End, Start: 26, Length: 0},
			}))
		})

		DescribeTable("numbers",
			func(input string) {
				Expect(next(input)).To(Equal(tokenizer.Token{
					Type:   tokenizer.Number,
					Start:  0,
					Length: len(input),
					Value:  []byte(input),
				}))
			},
			Entry("integer", "123456789"),
			Entry("negative", "-42"),
			Entry("fraction", "4.2"),
			Entry("exponent lc", "1e10"),
			Entry("exponent lc+", "1e+10"),
			Entry("exponent lc-", "1e-10"),
			Entry("exponent uc", "1E10"),
			Entry("exponent uc+", "1E+10"),
			Entry("exponent uc-", "1E-10"),
		)

		DescribeTable("strings",
			func(input string, expected []string) {
				var got []string
				for _, s := range all(input) {
					if s.Type == tokenizer.String {
						got = append(got, input[s.Start:s.Start+s.Length])
					}
				}
				Expect(got).To(Equal(expected))
			},
			Entry("basic", `["hello"`, []string{`"hello"`}),
			Entry("escaped quote", `["he\"llo"`, []string{`"he\"llo"`}),
			Entry("multiple", `["he\"llo","world","!\"!"`, []string{`"he\"llo"`, `"world"`, `"!\"!"`}),
			Entry("blank", `[" ","\n","\r","\t"`, []string{`" "`, `"\n"`, `"\r"`, `"\t"`}),
		)
	})

	Describe("skipping blank space", func() {
		It("handles whitespace padding", func() {
			Expect(next("  {")).To(Equal(tokenizer.Token{
				Type:   tokenizer.ObjectOpen,
				Start:  2,
				Length: 1,
			}))
		})

		It("handles space, line feed, carriage return, and tab", func() {
			Expect(all(" \n\r\t{ \n\r\t]")).To(Equal([]tokenizer.Token{
				{Type: tokenizer.ObjectOpen, Start: 4, Length: 1},
				{Type: tokenizer.ArrayClose, Start: 9, Length: 1},
				{Type: tokenizer.End, Start: 10, Length: 0},
			}))
		})
	})

	Describe("errors", func() {
		failure := func(input string) error {
			tnz := tokenizer.New([]byte(input))
			for {
				tok, err := tnz.Next()
				switch {
				case err != nil:
					return err
				case tok.Type == tokenizer.End:
					Fail("got to end without error")
				default:
				}
			}
		}

		It("fails on unterminated strings", func() {
			err := failure(`["once","upon]`)
			Expect(err).To(MatchError("unexpected end of JSON input"), err.Error())
			Expect(err).To(BeAssignableToTypeOf(tokenizer.UnexpectedEndError{}))
			Expect(err.(tokenizer.UnexpectedEndError).Position()).To(Equal(8))
		})

		It("fails on unknown tokens", func() {
			err := failure(`[{"once":upon`)
			Expect(err).To(MatchError("invalid character 'u'"), err.Error())
			Expect(err).To(BeAssignableToTypeOf(tokenizer.InvalidCharacterError{}))
			Expect(err.(tokenizer.InvalidCharacterError).Position()).To(Equal(9))
		})

		DescribeTable("invalid keywords",
			func(input, keyword, expect, actual string, position int) {
				err := failure(input)
				msg := fmt.Sprintf("invalid character '%s' in literal %s (expecting '%s')", actual, keyword, expect)
				Expect(err).To(MatchError(msg), err.Error())
				Expect(err).To(BeAssignableToTypeOf(tokenizer.InvalidKeywordError{}))
				Expect(err.(tokenizer.InvalidKeywordError).Position()).To(Equal(position))
			},
			Entry("null", "nut", "null", "l", "t", 2),
			Entry("true", "truy", "true", "e", "y", 3),
			Entry("false", "fAlse", "false", "a", "A", 1),
			Entry("short", "fals", "false", "e", " ", 4),
		)
	})
})
