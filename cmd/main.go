package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fratle/jsontoschema"
)

func main() {
	filepath := "test.json"

	if len(os.Args) > 1 {
		filepath = os.Args[1]
	}

	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	schema, err := jsontoschema.JsonToSchema(string(dat))
	if err != nil {
		panic(err)
	}
	fmt.Println(schema)
}
