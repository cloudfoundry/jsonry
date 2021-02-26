package parser

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/jsonry/internal/tokenizer"
)

type afterArrayElementContextError struct{ cause error }

func (e afterArrayElementContextError) Error() string {
	return fmt.Sprintf("%s after array element", e.cause)
}

type afterKeyValuePairContextError struct{ cause error }

func (e afterKeyValuePairContextError) Error() string {
	return fmt.Sprintf("%s after object key:value pair", e.cause)
}

type afterKeyContextError struct{ cause error }

func (e afterKeyContextError) Error() string { return fmt.Sprintf("%s after object key", e.cause) }

type beginningOfKeyContextError struct{ cause error }

func (e beginningOfKeyContextError) Error() string {
	return fmt.Sprintf("%s looking for beginning of object key string", e.cause)
}

type nonValueTokenError struct{ Token tokenizer.Token }

func (e nonValueTokenError) Error() string { return "" }

type unexpectedEndOfInput struct{}

func (e unexpectedEndOfInput) Error() string { return "unexpected end of JSON input" }

func Parse(input []byte) (interface{}, error) {
	t := tokenizer.New(input)
	r, err := parseValue(t)
	switch e := err.(type) {
	case nil:
	case nonValueTokenError:
		return nil, fmt.Errorf(`invalid character '%c' looking for beginning of value`, input[e.Token.Start])
	default:
		return nil, err
	}

	a, err := t.Next()
	switch e := err.(type) {
	case nil:
	case tokenizer.InvalidCharacterError:
		return nil, fmt.Errorf(`%s after top-level value`, err)
	case tokenizer.InvalidKeywordError:
		return nil, fmt.Errorf(`invalid character '%c' after top-level value`, input[e.Position()])
	default:
		return nil, err
	}
	if a.Type != tokenizer.End {
		return nil, fmt.Errorf(`invalid character '%c' after top-level value`, input[a.Start])
	}

	return r, nil
}

func parseValue(t *tokenizer.Tokenizer) (interface{}, error) {
	n, err := t.Next()
	switch err.(type) {
	case nil:
	case tokenizer.InvalidCharacterError:
		return nil, fmt.Errorf("%s looking for beginning of value", err)
	default:
		return nil, err
	}

	switch n.Type {
	case tokenizer.End:
		return nil, unexpectedEndOfInput{}
	case tokenizer.Null:
		return nil, nil
	case tokenizer.True:
		return true, nil
	case tokenizer.False:
		return false, nil
	case tokenizer.Number:
		return json.Number(n.Value), nil
	case tokenizer.String:
		return string(n.Value), nil
	case tokenizer.ArrayOpen:
		return parseArray(t)
	case tokenizer.ObjectOpen:
		return parseObject(t)
	default:
		return nil, nonValueTokenError{Token: n}
	}
}

func parseArray(t *tokenizer.Tokenizer) ([]interface{}, error) {
	var result []interface{}

	for {
		v, err := parseValue(t)
		if err != nil {
			if e, ok := err.(nonValueTokenError); ok && e.Token.Type == tokenizer.ArrayClose && len(result) == 0 {
				return result, nil
			}
			return nil, err
		}
		result = append(result, v)

		sep, err := advance(t, "after array element", tokenizer.Comma, tokenizer.ArrayClose)
		if err != nil {
			return nil, err
		}
		if sep.Type == tokenizer.ArrayClose {
			return result, nil
		}
	}
}

func parseObject(t *tokenizer.Tokenizer) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for {
		key, err := advance(t, "looking for beginning of object key string", tokenizer.ObjectClose, tokenizer.String)
		if err != nil {
			return nil, err
		}

		if key.Type == tokenizer.ObjectClose && len(result) == 0 {
			return result, nil
		}

		_, err = advance(t, "after object key", tokenizer.Colon)
		if err != nil {
			return nil, err
		}

		value, err := parseValue(t)
		if err != nil {
			return nil, err
		}
		result[string(key.Value)] = value

		sep, err := advance(t, "after object key:value pair", tokenizer.Comma, tokenizer.ObjectClose)
		if err != nil {
			return nil, err
		}

		if sep.Type == tokenizer.ObjectClose {
			return result, nil
		}
	}
}

func advance(t *tokenizer.Tokenizer, context string, allowed ...tokenizer.TokenType) (tokenizer.Token, error) {
	tok, err := t.Next()
	switch e := err.(type) {
	case nil:
	case tokenizer.InvalidKeywordError:
		return tokenizer.Token{}, fmt.Errorf(`invalid character '%c' %s`, t.Data[e.Position()], context)
	case tokenizer.InvalidCharacterError:
		return tokenizer.Token{}, fmt.Errorf("%s %s", err, context)
	default:
		return tokenizer.Token{}, err
	}

	if tok.Type == tokenizer.End {
		return tokenizer.Token{}, unexpectedEndOfInput{}
	}

	for _, a := range allowed {
		if tok.Type == a {
			return tok, nil
		}
	}

	return tokenizer.Token{}, fmt.Errorf(`invalid character '%c' %s`, t.Data[tok.Start], context)
}
