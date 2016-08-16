package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/juju/gojsonschema"
	"launchpad.net/gnuflag"
)

func main() {
	os.Exit(_main())
}

func usage() {
	fmt.Printf("jsonschema [schema file] [target file]\n")
}

func _main() int {
	gnuflag.Usage = usage
	gnuflag.Parse(true)
	if len(os.Args) < 2 {
		usage()
		return 1
	}

	schemaf, err := os.Open(os.Args[1])
	if err != nil {
		log.Printf("failed to open schema: %s", err)
		return 1
	}
	defer schemaf.Close()

	schemacontent, err := ioutil.ReadAll(schemaf)
	if err != nil {
		log.Printf("failed to read schema: %s", err)
		return 1
	}

	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(schemacontent)))
	if err != nil {
		log.Printf("failed to read schema: %s", err)
		return 1
	}

	if len(os.Args) < 3 {
		return 0
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		log.Printf("failed to open data: %s", err)
		return 1
	}
	defer f.Close()

	in, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("failed to read data: %s", err)
		return 1
	}

	loader := gojsonschema.NewStringLoader(string(in))
	result, err := schema.Validate(loader)
	if err != nil {
		log.Printf("failed to validate: %s", err)
		return 1
	}
	if result == nil {
		log.Printf("did not get a Result object")
		return 1
	}
	if !result.Valid() {
		for _, resulterr := range result.Errors() {
			log.Printf(resulterr.String())
		}
		return 1
	}
	return 0
}
