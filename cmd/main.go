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
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	schema, err := jsontoschema.JsonToSchema(string(dat))
	if err != nil {
		return "", err
	}

	return schema, nil
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

	schema, err := jsontoschema.JsonToSchema(string(body))
	if err != nil {
		return "", err
	}
	return schema, nil
}

func getSchemaFromStdin() (string, error) {
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

	schema, err := jsontoschema.JsonToSchema(b.String())
	if err != nil {

		return "", err
	}
	return schema, nil
}

func handlePipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	isPipe := info.Mode()&os.ModeNamedPipe != 0

	if isPipe {
		schema, err := getSchemaFromStdin()
		if err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		return true
	}
	return false
}

func main() {
	handlePipe()

	file := flag.String("file", "", "the path is a json file on disk")
	url := flag.String("url", "", "the path is a url to a json payload")
	flag.Parse()

	if *file == "" && *url == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	if *file != "" {
		schema, err := jsonFromFile(*file)
		if err != nil {
			fmt.Printf("An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	if *url != "" {
		schema, err := jsonFromUrl(*url)
		if err != nil {
			fmt.Printf("An error occured while getting json from url: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(schema)
		os.Exit(0)
	}

	return
}
