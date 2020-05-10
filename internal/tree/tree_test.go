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

	Describe("Fetch", func() {
		It("can fetch a basic value", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":"hello"}}}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal("hello"))
		})

		It("says not ok when not there", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":"hello"}}}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.not_there.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeFalse())
			Expect(v).To(BeNil())
		})

		It("can fetch a nil", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":null}}}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(BeNil())
		})

		It("can fetch an object", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":"hello"}}}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(map[string]interface{}{"d": map[string]interface{}{"e": "hello"}}))
		})

		It("can fetch a list at the leaf", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":{"d":{"e":["h","e","l","l","o"]}}}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal([]interface{}{"h", "e", "l", "l", "o"}))
		})

		It("can fetch a list at a branch", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":[{"d":{"e":"h"}},{"d":{"e":"i"}},{"d":{"e":"!"}}]}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal([]interface{}{"h", "i", "!"}))
		})

		It("inserts nils when a list has missing elements", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":[{"d":{"e":"h"}},{},{"d":{"e":"i"}},{"e":4},{"d":{"e":"!"}}]}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal([]interface{}{"h", nil, "i", nil, "!"}))
		})

		It("flattens lists of lists", func() {
			var t tree.Tree
			Expect(json.Unmarshal([]byte(`{"a":{"b":{"c":[[{"d":{"e":"h"}}],[{}],{"d":{"e":"i"}},[],{"d":{"e":"!"}}]}}}`), &t)).NotTo(HaveOccurred())
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"a.b.c.d.e"`})

			v, ok := t.Fetch(p)
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal([]interface{}{"h", nil, "i", "!"}))
		})
	})
})
