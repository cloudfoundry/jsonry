package jsonry_test

import (
	"fmt"

	"code.cloudfoundry.org/jsonry"
)

func ExampleMarshal() {
	s := struct {
		A    string
		B    string `json:"bee,omitempty"`
		GUID string `jsonry:"relationships.space.data.guid"`
		IDs  []int  `jsonry:"data.entries[].id"`
	}{
		A:    "foo",
		B:    "",
		GUID: "267758c0-985b-11ea-b9ac-48bf6bec2d78",
		IDs:  []int{1, 2, 3, 4, 5},
	}

	json, err := jsonry.Marshal(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
	// Output:
	// {"A":"foo","data":{"entries":[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5}]},"relationships":{"space":{"data":{"guid":"267758c0-985b-11ea-b9ac-48bf6bec2d78"}}}}
}

func ExampleMarshal_recursive() {
	type s struct {
		A string `jsonry:"f"`
	}
	type t struct {
		B []s `jsonry:"d[].e"`
	}
	type u struct {
		C []t `jsonry:"a.b[].c"`
	}

	data := u{C: []t{{B: []s{{A: "foo"}, {"bar"}}}, {B: []s{{A: "baz"}, {"quz"}}}}}

	json, err := jsonry.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
	// Output:
	// {"a":{"b":[{"c":{"d":[{"e":{"f":"foo"}},{"e":{"f":"bar"}}]}},{"c":{"d":[{"e":{"f":"baz"}},{"e":{"f":"quz"}}]}}]}}
}

func ExampleUnmarshal() {
	json := `
    {
      "A": "foo",
      "bee": "bar",
      "data": {
        "entries": [
          {"id": 1},
          {"id": 2},
          {"id": 3},
          {"id": 4},
          {"id": 5}
        ]
      },
      "relationships": {
        "space": {
          "data": {
            "guid":"267758c0-985b-11ea-b9ac-48bf6bec2d78"}
          }
        }
      }
    }`

	var s struct {
		A    string
		B    string `json:"bee"`
		GUID string `jsonry:"relationships.space.data.guid"`
		IDs  []int  `jsonry:"data.entries.id"`
	}

	if err := jsonry.Unmarshal([]byte(json), &s); err != nil {
		panic(err)
	}

	fmt.Printf("A: %+v\nB: %+v\nGUID: %+v\nIDs: %+v\n", s.A, s.B, s.GUID, s.IDs)
	// Output:
	// A: foo
	// B: bar
	// GUID: 267758c0-985b-11ea-b9ac-48bf6bec2d78
	// IDs: [1 2 3 4 5]
}
