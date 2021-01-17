package tokenizer

type UnexpectedEndError struct {
	position int
}

func (UnexpectedEndError) Error() string {
	return "unexpected end of JSON input"
}

func (u UnexpectedEndError) Position() int {
	return u.position
}
