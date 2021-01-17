package tokenizer

import "fmt"

type InvalidKeywordError struct {
	keyword  keyword
	actual   byte
	expect   byte
	position int
}

func (i InvalidKeywordError) Error() string {
	return fmt.Sprintf("invalid character '%c' in literal %s (expecting '%c')", i.actual, i.keyword, i.expect)
}

func (i InvalidKeywordError) Position() int {
	return i.position
}
