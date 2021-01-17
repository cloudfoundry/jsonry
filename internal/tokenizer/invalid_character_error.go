package tokenizer

import "fmt"

type InvalidCharacterError struct {
	character byte
	position  int
}

func (i InvalidCharacterError) Error() string {
	return fmt.Sprintf("invalid character '%c'", i.character)
}

func (i InvalidCharacterError) Position() int {
	return i.position
}
