package tokenizer

type Token struct {
	Type   TokenType
	Start  int
	Length int
	Value  []byte
}
