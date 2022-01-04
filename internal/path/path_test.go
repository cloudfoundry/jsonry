package path_test

import (
	"reflect"

	"code.cloudfoundry.org/jsonry/internal/path"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	It("computes a path from the field name", func() {
		p := path.ComputePath(reflect.StructField{Name: "foo"})
		Expect(p.String()).To(Equal("foo"))
		Expect(p.Len()).To(Equal(1))
	})

	It("computes a path from a JSON tag", func() {
		p := path.ComputePath(reflect.StructField{Tag: `json:"foo"`})
		Expect(p.String()).To(Equal("foo"))
		Expect(p.Len()).To(Equal(1))
	})

	It("computes a path from a JSONry tag", func() {
		p := path.ComputePath(reflect.StructField{Tag: `jsonry:"foo.bar[].baz.quz"`})
		Expect(p.String()).To(Equal("foo.bar[].baz.quz"))
		Expect(p.Len()).To(Equal(4))
	})

	It("implements Pull()", func() {
		p := path.ComputePath(reflect.StructField{Tag: `jsonry:"foo.bar[].baz.quz"`})
		Expect(p.Len()).To(Equal(4))

		s, p := p.Pull()
		Expect(s).To(Equal(path.Segment{Name: "foo", List: false}))
		Expect(p.Len()).To(Equal(3))

		s, p = p.Pull()
		Expect(s).To(Equal(path.Segment{Name: "bar", List: true}))
		Expect(p.Len()).To(Equal(2))
	})

	Context("omitempty", func() {
		It("picks it up from a JSON tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `json:",omitempty"`})
			Expect(p.OmitEmpty).To(BeTrue())
		})

		It("picks it up from a JSONry tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:",omitempty"`})
			Expect(p.OmitEmpty).To(BeTrue())
		})
	})

	Context("always omit", func() {
		It("picks it up from a JSON tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `json:"-"`})
			Expect(p.OmitAlways).To(BeTrue())
		})

		It("picks it up from a JSONry tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"-"`})
			Expect(p.OmitAlways).To(BeTrue())
		})

		It("allows a literal name `-` from a JSON tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `json:"-,"`})
			Expect(p.OmitAlways).To(BeFalse())
		})

		It("allows a literal name `-` from a JSONry tag", func() {
			p := path.ComputePath(reflect.StructField{Tag: `jsonry:"-,"`})
			Expect(p.OmitAlways).To(BeFalse())
		})
	})
})
