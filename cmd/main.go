package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fratle/jsontoschema"
)

var usage = `Reads a json input from a file, and url or via a commandline pipe
Specify either -url followed by an url or -file followed by a filepath

Usage:
url:  jsontoschema -url https://jsonplaceholder.typicode.com/todos/1
file: jsontoschema -file dump.json
curl: curl -s https://jsonplaceholder.typicode.com/todos/1 2<&1 | jsontoschema`

func jsonFromFile(file string) (string, error) {
	fb, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(fb), nil
}

func jsonFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func jsonFromStdin() (string, error) {
	b := strings.Builder{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		_, err := b.Write(scanner.Bytes())

		if err != nil {
			return "", err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func isPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeNamedPipe != 0
}

func jsonFromPipe(_ string) (string, error) {
	json, err := jsonFromStdin()
	if err != nil {
		return "", err
	}
	return json, nil
}

type jsonF func(s string) (string, error)

func createSchemaFromJSON(f jsonF, s string) (string, error) {
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

func main() {

	if isPipe() {
		schema, err := createSchemaFromJSON(jsonFromPipe, "")
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
		schema, err := createSchemaFromJSON(jsonFromFile, *file)

		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	if *url != "" {
		schema, err := createSchemaFromJSON(jsonFromUrl, *url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	return
}
