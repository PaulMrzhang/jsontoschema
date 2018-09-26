package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fratle/jsontoschema"
	"github.com/fratle/jsontoschema/jsonreader"
)

var usage = `Reads a json input from a file, and url or via a commandline pipe
Specify either -url followed by an url or -file followed by a filepath

Usage:
url:  jsontoschema -url https://jsonplaceholder.typicode.com/todos/1
file: jsontoschema -file dump.json
curl: curl -s https://jsonplaceholder.typicode.com/todos/1 2<&1 | jsontoschema`

func createSchemaFromJSON(f jsonreader.JsonFunc, s string) (string, error) {
	json, err := f(s)

	if err != nil {
		return "", err
	}

	schema, err := jsontoschema.JsonToSchema(json)

	if err != nil {
		return "", err
	}

	return schema, nil
}

func isPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeNamedPipe != 0
}

func main() {
	if isPipe() {
		schema, err := createSchemaFromJSON(jsonreader.FromPipe, "")
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	file := flag.String("file", "", "the path is a json file on disk")
	url := flag.String("url", "", "the path is a url to a json payload")
	flag.Parse()

	if *file == "" && *url == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	if *file != "" {
		schema, err := createSchemaFromJSON(jsonreader.FromFile, *file)

		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	if *url != "" {
		schema, err := createSchemaFromJSON(jsonreader.FromUrl, *url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	return
}
