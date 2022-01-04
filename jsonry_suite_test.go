package jsonry_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type implementsJSONMarshaler struct {
	bytes []byte
	err   error
}

func (i implementsJSONMarshaler) MarshalJSON() ([]byte, error) {
	return i.bytes, i.err
}

type implementsJSONUnmarshaler struct {
	hasBeenSet bool
}

func (i *implementsJSONUnmarshaler) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte(`"fail"`)) {
		return errors.New("ouch")
	}
	i.hasBeenSet = true
	return nil
}

type implementsOmissible string

func (i implementsOmissible) OmitJSONry() bool {
	return i == "omit"
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
