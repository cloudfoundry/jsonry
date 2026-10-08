package jsonry_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"code.cloudfoundry.org/jsonry"
)

type container struct {
	A []int  `jsonry:"a"`
	C any    `jsonry:"b.c"`
	D string `jsonry:"b.d"`
}

const marshaled = `{"a":[1,2,3],"b":{"c":null,"d":"hello"}}`

func unmarshaled() container {
	return container{
		A: []int{1, 2, 3},
		C: nil,
		D: "hello",
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("JSONry", func(b *testing.B) {
		var receiver container
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			jsonry.Unmarshal([]byte(marshaled), &receiver)
		}

		b.StopTimer()
		if !reflect.DeepEqual(receiver, unmarshaled()) {
			b.Fatalf("received struct not equal")
		}
	})

	b.Run("JSON", func(b *testing.B) {
		var result container
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var receiver struct {
				A []int `json:"a"`
				B struct {
					C any    `json:"c"`
					D string `json:"d"`
				} `json:"b"`
			}
			json.Unmarshal([]byte(marshaled), &receiver)
			result = container{
				A: receiver.A,
				C: receiver.B.C,
				D: receiver.B.D,
			}
		}

		b.StopTimer()
		if !reflect.DeepEqual(result, unmarshaled()) {
			b.Fatalf("received struct not equal")
		}
	})
}

func BenchmarkMarshal(b *testing.B) {
	b.Run("JSONry", func(b *testing.B) {
		u := unmarshaled()
		var result []byte
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			result, _ = jsonry.Marshal(u)
		}

		b.StopTimer()
		if !bytes.Equal(result, []byte(marshaled)) {
			b.Fatalf("received bytes not equal")
		}
	})

	b.Run("JSON", func(b *testing.B) {
		u := unmarshaled()
		var result []byte
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var transmitter struct {
				A []int `json:"a"`
				B struct {
					C any    `json:"c"`
					D string `json:"d"`
				} `json:"b"`
			}
			transmitter.A = u.A
			transmitter.B.C = u.C
			transmitter.B.D = u.D

			result, _ = json.Marshal(transmitter)
		}

		b.StopTimer()
		if !bytes.Equal(result, []byte(marshaled)) {
			b.Fatalf("received bytes not equal")
		}
	})
}
