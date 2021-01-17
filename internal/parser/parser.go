package parser

import (
	"fmt"

	"code.cloudfoundry.org/jsonry/internal/tokenizer"
)

func Parse(input []byte) (interface{}, error) {
	t := tokenizer.New(input)
	return parse(t)
}

func parse(t *tokenizer.Tokenizer) (interface{}, error) {
	n, err := t.Next()
	if err != nil {
		return nil, err
	}
	switch n.Type {
	case tokenizer.End:
		return nil, fmt.Errorf("unexpected end of JSON input")
	}
	return nil, nil
}
