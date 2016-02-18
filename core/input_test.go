package core

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
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

	val, err := conf.Data("test")
	assert.NoError(t, err)
	assert.Equal(t, "true", val)

	err = conf.Validate(schema)
	assert.Error(t, err)

	conf.data["first"] = "test"
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["second"] = 10
	err = conf.Validate(schema)
	assert.Error(t, err)

	conf.data["second"] = []byte("test input")
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["second"] = "test input"
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["third"] = "bad"
	err = conf.Validate(schema)
	assert.Error(t, err)

	conf.data["third"] = 10
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["third"] = 10.01
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["fourth"] = "bad"
	err = conf.Validate(schema)
	assert.Error(t, err)

	conf.data["fourth"] = "create"
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["fifth"] = "bad"
	err = conf.Validate(schema)
	assert.Error(t, err)

	in, out := io.Pipe()
	conf.data["fifth"] = in
	err = conf.Validate(schema)
	assert.NoError(t, err)

	conf.data["fifth"] = out
	err = conf.Validate(schema)
	assert.NoError(t, err)

}
