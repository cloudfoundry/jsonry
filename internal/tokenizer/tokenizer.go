package tokenizer

type Tokenizer struct {
	Data     []byte
	position int
}

func New(input []byte) *Tokenizer {
	return &Tokenizer{Data: input}
}
