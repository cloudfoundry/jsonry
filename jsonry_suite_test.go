package jsonry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJSONry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONry Suite")
}
