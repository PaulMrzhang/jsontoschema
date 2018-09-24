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

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {

		file := flag.String("file", "", "the path is a json file on disk")
		url := flag.String("url", "", "the path is a url to a json payload")
		flag.Parse()

		if *file == "" && *url == "" {
			fmt.Println("Reads a json input from a file, and url or via a commandline pipe")
			fmt.Println("Specify either -url followed by an url or -file followed by a filepath")
			fmt.Println("\nUsage:")
			fmt.Println("\nurl:\njsontoschema -url https://jsonplaceholder.typicode.com/todos/1")
			fmt.Println("\nfile:\njsontoschema -file dump.json")
			fmt.Println("\ncurl:\ncurl -s https://jsonplaceholder.typicode.com/todos/1 2<&1 | jsontoschema\n")
			os.Exit(1)
		}

		if *file != "" && *url != "" {
			fmt.Println("must specify either -url or -file")
			os.Exit(1)
		}

		if *file != "" {
			dat, err := ioutil.ReadFile(*file)
			if err != nil {
				panic(err)
			}

			schema, err := jsontoschema.JsonToSchema(string(dat))
			if err != nil {
				panic(err)
			}
			fmt.Println(schema)
		}

		if *url != "" {
			resp, err := http.Get(*url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				panic(err)
			}

			schema, err := jsontoschema.JsonToSchema(string(body))
			if err != nil {
				panic(err)
			}
			fmt.Println(schema)
		}

	} else {
		b := strings.Builder{}
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			fmt.Println()
			_, err := b.Write(scanner.Bytes())

			if err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				os.Exit(1)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}

		schema, err := jsontoschema.JsonToSchema(b.String())
		if err != nil {
			panic(err)
		}
		fmt.Println(schema)
	}
}
