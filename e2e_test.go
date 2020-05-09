package jsonry_test

import (
	"code.cloudfoundry.org/jsonry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("end to end", func() {
	It("marshals", func() {
		o := struct {
			Num         float32               `json:"number"`
			Spaces      []space               `jsonry:"relationships.space.data"`
			Orgs        []string              `jsonry:"relationships.orgs.data.guids"`
			OtherSpace  []space               `jsonry:"other_space[].data"`
			Type        string                `jsonry:"authentication.type"`
			Username    string                `jsonry:"authentication.credentials.username"`
			Password    string                `jsonry:"authentication.credentials.password"`
			Labels      map[string]nullString `jsonry:"metadata.labels"`
			Annotations map[string]space      `jsonry:"metadata.annotations"`
		}{
			Num:         12,
			Spaces:      []space{{GUID: "foo"}, {Name: "Bar", GUID: "bar"}},
			Orgs:        []string{"baz", "quz"},
			OtherSpace:  []space{{GUID: "alpha"}, {Name: "Beta", GUID: "beta"}},
			Type:        "basic",
			Username:    "fake-user",
			Password:    "fake secret",
			Labels:      map[string]nullString{"first": {value: "one"}, "second": {null: true}},
			Annotations: map[string]space{"alpha": {GUID: "foo"}, "beta": {Name: "Bar", GUID: "bar"}},
		}
		r, err := jsonry.Marshal(o)
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(MatchJSON(`
 		{
			"number": 12,
			"relationships": {
			  "orgs": {
				"data": {
				  "guids": ["baz","quz"]
				}
			  },
			  "space": {
				"data": [
				  {"guid": "foo"},
				  {"guid": "bar","name": "Bar"}
				]
			  }
			},
			"authentication": {
				"type": "basic",
				"credentials": {
					"username": "fake-user",
					"password": "fake secret"
				}
			},
			"metadata": {
			  "annotations": {
				"alpha": {
				  "guid": "foo"
				},
				"beta": {
				  "guid": "bar",
				  "name": "Bar"
				}
			  },
			  "labels": {
				"first": "one",
				"second": null
			  }
			},
			"other_space": [
			  {
				"data": {"guid": "alpha"}
			  },
			  {
				"data": {"guid": "beta","name": "Beta"}
			  }
			]
		}`))
	})
})
