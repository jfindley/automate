package file

import (
	"github.com/jfindley/automate/core"
)

var schema core.InputSchema

func init() {
	var err error
	schema, err = core.NewInputSchema(schemaDef)
	if err != nil {
		panic(err)
	}
}

var schemaDef = []byte(`
path:
  required: true
  types:
  - string

action:
  types:
  - string
  values:
  - touch
  - set
  - remove

mode:
  types:
  - string
  - integer

context:
  types:
  - string

owner:
  types:
  - string
  - integer

group:
  types:
  - string
  - integer
`)
