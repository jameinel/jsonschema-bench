package main_test

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/juju/gojsonschema"
	schema "github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsschema/validator"
)

var _cached []string
func allObjects() []string {
	if len(_cached) != 0 {
		return _cached
	}
	f, err := os.Open("all-requests.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var all []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) < 1 || txt == "[" || txt == "{}" || txt == "]" {
			continue
		}
		if txt[len(txt)-1] == ',' {
			txt = txt[:len(txt)-1]
		}
		all = append(all, txt)
	}
	_cached = all
	return all
}

func BenchmarkJSONSchema(b *testing.B) {
	all := allObjects()
	f, err := os.Open("message-schema.json")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}
	checker, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(content)))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, content := range all {
			loader := gojsonschema.NewStringLoader(content)
			res, err := checker.Validate(loader)
			if err != nil {
				b.Fatalf("failed to validate %d\n%s\nbecause\n%s", j, content, err)
			}
			if !res.Valid() {
				b.Fatal(res.Errors())
			}
		}
	}
}

func BenchmarkJSSchema(b *testing.B) {
	all := allObjects()
	s, err := schema.ReadFile("message-schema.json")
	if err != nil {
		b.Fatal(err)
	}
	checker := validator.New(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, content := range all {
			var v interface{}
			err = json.Unmarshal([]byte(content), &v)
			if err != nil {
				b.Fatalf("failed to unmarshal %d\n%s\nbecause\n%s", j, content, err)
			}
			err = checker.Validate(v)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
