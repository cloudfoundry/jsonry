package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/jsonry/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: cmd <json file>")
		os.Exit(2)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
		os.Exit(2)
	}

	_, err = parser.Parse(data)
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
