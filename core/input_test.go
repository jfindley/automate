package core

import (
	"io"
	"testing"
)

var testSchema = []byte(`
first:
  required: true

second:
  required: false
  types:
  - string
  - data

third:
  required: false
  types:
  - integer
  - float

fourth:
  required: false
  values:
  - create
  - delete

fifth:
  types:
  - pipe

`)

func TestInput(t *testing.T) {

	schema, err := NewInputSchema(testSchema)
	if err != nil {
		t.Fatal(err)
	}

	testInput := map[string]interface{}{"test": "true"}

	conf := NewConfigInput(testInput)

	if conf.Data()["test"] != "true" {
		t.Error("Data method did not return data")
	}

	if err := conf.Validate(schema); err == nil {
		t.Error("Config validated without required parameter")
	}

	conf.data["first"] = "test"
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with all required parameters")
	}

	conf.data["second"] = 10
	if err := conf.Validate(schema); err == nil {
		t.Error("Config validated with improper types")
	}

	conf.data["second"] = []byte("test input")
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	conf.data["second"] = "test input"
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	conf.data["third"] = "bad"
	if err := conf.Validate(schema); err == nil {
		t.Error("Config validated with improper types")
	}

	conf.data["third"] = 10
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	conf.data["third"] = 10.01
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	conf.data["fourth"] = "bad"
	if err := conf.Validate(schema); err == nil {
		t.Error("Config validated with improper types")
	}

	conf.data["fourth"] = "create"
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	in, out := io.Pipe()
	conf.data["fifth"] = "bad"
	if err := conf.Validate(schema); err == nil {
		t.Error("Config validated with improper types")
	}

	conf.data["fifth"] = in
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

	conf.data["fifth"] = out
	if err := conf.Validate(schema); err != nil {
		t.Error("Config failed to validate with proper types")
	}

}
