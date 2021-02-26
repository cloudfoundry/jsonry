package tokenizer

func (t *Tokenizer) Next() (Token, error) {
	for t.position < len(t.Data) && blank(t.Data[t.position]) {
		t.position++
	}

	if t.position == len(t.Data) {
		return Token{Type: End, Start: t.position}, nil
	}

	switch t.Data[t.position] {
	case '{':
		return t.delimiter(ObjectOpen)
	case '}':
		return t.delimiter(ObjectClose)
	case '[':
		return t.delimiter(ArrayOpen)
	case ']':
		return t.delimiter(ArrayClose)
	case ':':
		return t.delimiter(Colon)
	case ',':
		return t.delimiter(Comma)
	case 'n':
		return t.keyword(Null, keywordNull)
	case 't':
		return t.keyword(True, keywordTrue)
	case 'f':
		return t.keyword(False, keywordFalse)
	case '"':
		return t.string()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return t.number()
	default:
		return Token{}, InvalidCharacterError{position: t.position, character: t.Data[t.position]}
	}
}

func (t *Tokenizer) delimiter(tokenType TokenType) (Token, error) {
	token := Token{
		Type:   tokenType,
		Start:  t.position,
		Length: 1,
	}
	t.position++
	return token, nil
}

func (t *Tokenizer) keyword(tokenType TokenType, keyword keyword) (Token, error) {
	var index int

	for index = 1; t.position+index < len(t.Data) && index < len(keyword); index++ {
		expect := keyword[index]
		actual := t.Data[t.position+index]
		if expect != actual {
			return Token{}, InvalidKeywordError{
				keyword:  keyword,
				actual:   actual,
				expect:   expect,
				position: t.position,
			}
		}
	}

	if index < len(keyword) {
		return Token{}, InvalidKeywordError{
			keyword:  keyword,
			actual:   ' ',
			expect:   keyword[len(t.Data)-t.position],
			position: t.position,
		}
	}

	token := Token{
		Type:   tokenType,
		Start:  t.position,
		Length: len(keyword),
	}

	t.position += len(keyword)
	return token, nil
}

func (t *Tokenizer) string() (Token, error) {
	cursor := t.position + 1
	for cursor < len(t.Data) {
		if t.Data[cursor] == '"' && t.Data[cursor-1] != '\\' {
			token := Token{
				Type:   String,
				Start:  t.position,
				Length: cursor - t.position + 1,
				Value:  t.Data[t.position+1 : cursor],
			}
			t.position = cursor + 1
			return token, nil
		}

		cursor++
	}

	return Token{}, UnexpectedEndError{position: t.position}
}

func (t *Tokenizer) number() (Token, error) {
	end := len(t.Data)
	cursor := t.position + 1
	for cursor < end && isNumber(t.Data[cursor]) {
		cursor++
	}
	token := Token{
		Type:   Number,
		Start:  t.position,
		Length: cursor - t.position,
		Value:  t.Data[t.position:cursor],
	}
	t.position = cursor
	return token, nil
}
