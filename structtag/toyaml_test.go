package structtag_test

import (
	"reflect"

	"gopkg.in/check.v1"

	"github.com/jameinel/jsonschema-bench/structtag"
)

var _ = check.Suite(&ToYAMLSuite{})

type ToYAMLSuite struct{}

type simple struct {
	Integer int    `json:"integer"`
	String  string `json:"string"`
}

func (*ToYAMLSuite) TestSimple(c *check.C) {
	var x *simple
	t := reflect.TypeOf(x).Elem()
	yaml, err := structtag.ToYAMLSchema(t)
	c.Assert(err, check.IsNil)
	c.Check(yaml, check.Equals, `
additionalItems: false
additionalProperties: false
properties:
  integer:
    additionalItems: false
    additionalProperties: false
    type: integer
  string:
    additionalItems: false
    additionalProperties: false
    type: string
required:
- integer
- string
title: simple
type: object
`[1:])
}

type withempty struct {
	Optional int `json:"optional,omitempty"`
	Required int `json:"required"`
}

func (*ToYAMLSuite) TestOmitEmpty(c *check.C) {
	var x *withempty
	t := reflect.TypeOf(x).Elem()
	yaml, err := structtag.ToYAMLSchema(t)
	c.Assert(err, check.IsNil)
	c.Check(yaml, check.Equals, `
additionalItems: false
additionalProperties: false
properties:
  optional:
    additionalItems: false
    additionalProperties: false
    type: integer
  required:
    additionalItems: false
    additionalProperties: false
    type: integer
required:
- required
title: withempty
type: object
`[1:])
}
