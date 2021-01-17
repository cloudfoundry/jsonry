package tokenizer

type TokenType int

const (
	End TokenType = iota
	ObjectOpen
	ObjectClose
	ArrayOpen
	ArrayClose
	Colon
	Comma
	Null
	True
	False
	Number
	String
)

type keyword string

const (
	keywordTrue  keyword = "true"
	keywordFalse keyword = "false"
	keywordNull  keyword = "null"
)

func isNumber(symbol byte) bool {
	switch symbol {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E', '+', '-':
		return true
	default:
		return false
	}
}

func blank(symbol byte) bool {
	switch symbol {
	case '\u0020', '\u000A', '\u000D', '\u0009':
		return true
	default:
		return false
	}
}
