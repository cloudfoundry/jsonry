package jsonry_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type implementsJSONMarshaler struct{ value bool }

func (j implementsJSONMarshaler) MarshalJSON() ([]byte, error) {
	if j.value {
		return nil, errors.New("ouch")
	}
	return json.Marshal("hello")
}

func (j *implementsJSONMarshaler) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte(`"fail"`)) {
		return errors.New("ouch")
	}
	j.value = true
	return nil
}

type space struct {
	Name string `jsonry:"name,omitempty"`
	GUID string `jsonry:"guid"`
}

type nullString struct {
	value string
	null  bool
}

func (n nullString) MarshalJSON() ([]byte, error) {
	if n.null {
		return []byte("null"), nil
	}
	return json.Marshal(n.value)
}

func (n *nullString) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		n.null = true
		return nil
	}

	n.null = false
	return json.Unmarshal(input, &n.value)
}

func TestJSONry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONry Suite")
}
