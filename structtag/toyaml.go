package structtag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func ToSchema(t reflect.Type) (*Schema, error) {
	s := &Schema{
		Type:       []Type{ObjectType},
		Title:      t.Name(),
		Properties: make(map[string]*Schema),
	}
	required := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// if field.Anonymous {
		// 	continue
		// }
		prop := &Schema{}
		switch field.Type.Kind() {
		case reflect.Bool:
			//prop.Type = BoolType
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			prop.Type = []Type{IntegerType}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			//prop.Type = []Type{IntegerType}
			//prop.Minimum = Float(0)
		case reflect.Float32, reflect.Float64:
			//prop.Type = []Type{NumberType}
			//prop.Minimum = Float(0)
		case reflect.String:
			prop.Type = []Type{StringType}
		default:
			return nil, fmt.Errorf("unhandled Type: %s", field.Type.Kind())
		}
		jsonTagString := field.Tag.Get("json")
		jsonTags := strings.Split(jsonTagString, ",")
		name := field.Name
		if len(jsonTags) >= 1 {
			name = jsonTags[0]
		}
		s.Properties[name] = prop
		isRequired := true
		for _, tag := range jsonTags[1:] {
			if tag == "omitempty" {
				isRequired = false
			}
		}
		if isRequired {
			required = append(required, name)
		}
	}
	s.Required = required
	return s, nil
}

func ToYAMLSchema(t reflect.Type) (string, error) {
	s, err := ToSchema(t)
	if err != nil {
		return "", err
	}
	// Cheat, Schema only has a JSON serialization, so round trip to JSON
	// to omit empty and then cast to YAML
	if bytes, err := json.Marshal(s); err != nil {
		return "", err
	} else {
		var v interface{}
		if err := json.Unmarshal(bytes, &v); err != nil {
			return "", err
		}
		if bytes, err := yaml.Marshal(v); err != nil {
			return "", err
		} else {
			return string(bytes), nil
		}
	}
}

//  "properties": {
//    "request-id": {
//      "type": "integer",
//      "description": "Unique identifier for this request. The response will be tagged with the\nsame value as the request id. Request identifiers should not be reused\nwithin the lifetime of a connection.\nRequest-id is mandatory and must be a valid positive integer.\n",
//      "minimum": 1
//    },
//    "type": {
//      "type": "string",
//      "description": "Type gives the name of the Facade that we will be interacting with. A\nFacade collects a set of methods, grouped together for a focused purpose.\n"
//    },
