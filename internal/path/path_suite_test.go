package path_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONry Internal Path Suite")
}
