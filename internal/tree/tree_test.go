package tree_test

import (
	"encoding/json"
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/path"
	"code.cloudfoundry.org/jsonry/internal/tree"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	Describe("Attach", func() {
		It("attaches a branch with the right value", func() {
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})
			t := make(tree.Tree).Attach(p, "hello")
			Expect(json.Marshal(t)).To(MatchJSON(`{"a":{"b":{"c":{"d":{"e":"hello"}}}}}`))
		})

		It("can attach multiple branches", func() {
			t := make(tree.Tree).
				Attach(path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`}), "hello").
				Attach(path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.f.g"`}), "world!")
			Expect(json.Marshal(t)).To(MatchJSON(`{"a":{"b":{"c":{"d":{"e":"hello"}},"f":{"g":"world!"}}}}`))
		})

		It("creates lists according to the first list hint", func() {
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b[].c.d[].e"`})
			t := make(tree.Tree).Attach(p, []string{"hello", "world", "!"})
			Expect(json.Marshal(t)).To(MatchJSON(`{"a":{"b":[{"c":{"d":[{"e":"hello"}]}},{"c":{"d":[{"e":"world"}]}},{"c":{"d":[{"e":"!"}]}}]}}`))
		})

		When("there is no list hint", func() {
			It("creates lists at the leaf", func() {
				p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})
				t := make(tree.Tree).Attach(p, []string{"hello", "world", "!"})
				Expect(json.Marshal(t)).To(MatchJSON(`{"a":{"b":{"c":{"d":{"e":["hello","world","!"]}}}}}`))
			})
		})
	})
})
