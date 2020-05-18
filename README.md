[![CI](https://github.com/cloudfoundry/jsonry/workflows/Go/badge.svg)](https://github.com/cloudfoundry/jsonry/actions?query=workflow%3AGo)
[![GoDoc](https://godoc.org/code.cloudfoundry.org/jsonry?status.png)](https://godoc.org/code.cloudfoundry.org/jsonry)
# JSONry

A Go library and notation for converting between a Go `struct` and JSON.

```go
s := struct {
  GUID string `jsonry:"relationships.space.data.guid"`
}{
  GUID: "267758c0-985b-11ea-b9ac-48bf6bec2d78",
}

json, _ := jsonry.Marshal(s)
fmt.Println(string(json))
```
Will generate the following JSON:
```json
{
  "relationships": {
    "space": {
      "data": {
        "guid": "267758c0-985b-11ea-b9ac-48bf6bec2d78"
      }
    }
  }
}
```
The operation is reversible using `Unmarshal()`. The key advantage is that nested JSON can be generated and parsed without
the need to create intermediate Go structures. Check out [the documentation](https://godoc.org/code.cloudfoundry.org/jsonry) for details.

JSONry started life in the [Cloud Foundry CLI](https://github.com/cloudfoundry/cli) project. It has been extracted so
that it can be used in other projects too.
