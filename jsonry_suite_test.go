package jsonry_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type pri struct {
	private bool
	Public  bool
}

type jsm struct{ value bool }

func (j jsm) MarshalJSON() ([]byte, error) {
	if j.value {
		return nil, errors.New("ouch")
	}
	return json.Marshal("hello")
}

type jrm struct{ value bool }

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

func (j jrm) MarshalJSONry() (interface{}, error) {
	if j.value {
		return nil, errors.New("ouch")
	}
	return &pri{private: true, Public: true}, nil
}

func TestJSONry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONry Suite")
}
