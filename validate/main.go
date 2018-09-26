package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/fratle/jsontoschema/jsonreader"
	"github.com/xeipuuv/gojsonschema"
)

var usage = ``

func main() {

	argFiles := flag.String("files", "", "comma seperated list of json files on disk")
	argUrls := flag.String("urls", "", "comma seperated list of urls of json payloads")
	schemaFile := flag.String("schema", "", "filepath of the json schema to validate the input against")
	schemaURL := flag.String("schema_url", "", "url of the json schema to validate the input against")
	flag.Parse()

	files := strings.Split(*argFiles, ",")
	urls := strings.Split(*argUrls, ",")

	var schemaStr string

	fmt.Println("argFiles:", *argFiles)
	fmt.Println("argUrls:", *argUrls)
	fmt.Println("schemaFile:", *schemaFile)
	fmt.Println("schemaURL:", *schemaURL)

	if *schemaFile != "" {
		var err error
		schemaStr, err = jsonreader.FromFile(*schemaFile)

		if err != nil {
			panic(err)
		}
	}

	if *schemaURL != "" {
		var err error
		schemaStr, err = jsonreader.FromUrl(*schemaURL)
		if err != nil {
			panic(err)
		}
	}

	jsonStrings := []string{}
	for _, s := range files {
		if s == "" {
			continue
		}
		json, err := jsonreader.FromFile(s)
		if err != nil {
			panic(err)
		}

		jsonStrings = append(jsonStrings, json)
	}

	for _, s := range urls {
		if s == "" {
			continue
		}
		json, err := jsonreader.FromUrl(s)
		if err != nil {
			panic(err)
		}
		jsonStrings = append(jsonStrings, json)
	}

	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	schema, err := gojsonschema.NewSchema(schemaLoader)

	if err != nil {
		panic(err)
	}

	for _, doc := range jsonStrings {
		documentLoader := gojsonschema.NewStringLoader(doc)
		result, err := schema.Validate(documentLoader)

		if err != nil {
			fmt.Println(doc)
			panic(err)
		}

		if result.Valid() {
			fmt.Printf("The document is valid\n")
		} else {
			fmt.Printf("The document is not valid. see errors :\n")
			for _, err := range result.Errors() {
				// Err implements the ResultError interface
				fmt.Printf("- %s\n", err)
			}
		}
	}
}
