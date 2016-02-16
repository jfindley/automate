package core

import (
	"fmt"
	"reflect"
	"regexp"

	"errors"
	"gopkg.in/yaml.v2"
)

type InputSchema map[string]struct {
	Types    []string
	Required bool
	Values   []interface{}
}

func NewInputSchema(in []byte) (InputSchema, error) {
	var schema InputSchema
	err := yaml.Unmarshal(in, &schema)
	return schema, err
}

type ConfigInput struct {
	data  map[string]interface{}
	types map[string]string
}

func (c *ConfigInput) Data(val string) (interface{}, error) {
	if _, ok := c.data[val]; ok {
		return c.data[val], nil
	}
	return nil, errors.New("Value undefined")
}

func (c *ConfigInput) Type(val string) (string, error) {
	if _, ok := c.types[val]; ok {
		return c.types[val], nil
	}
	return "", errors.New("Type undefined")
}

func (c *ConfigInput) Validate(schema InputSchema) error {
	for k, thisSchema := range schema {
		thisData, ok := c.data[k]

		if !ok {
			if thisSchema.Required {
				return fmt.Errorf("Missing required parameter: %s", k)
			} else {
				continue
			}
		}

		if len(thisSchema.Values) > 0 {
			if !containsValues(thisData, thisSchema.Values) {
				return fmt.Errorf("Invalid value of paramter %s: %v", k, thisData)
			}
		}

		if len(thisSchema.Types) > 0 {
			thisType := reflect.TypeOf(thisData).String()

			if t, ok := containsTypes(thisType, thisSchema.Types); !ok {
				return fmt.Errorf("Invalid type of paramter %s: %s", k, thisType)
			} else {
				c.types[k] = t
			}
		}

	}
	return nil
}

func NewConfigInput(in map[string]interface{}) *ConfigInput {
	c := new(ConfigInput)
	c.data = in
	c.types = make(map[string]string)
	return c
}

func containsValues(needle interface{}, haystack []interface{}) bool {

	for i := range haystack {
		if reflect.DeepEqual(needle, haystack[i]) {
			return true
		}
	}

	return false
}

func containsTypes(needle string, haystack []string) (string, bool) {

	integer := regexp.MustCompile(`^(u)?int(16)?(32)?(64)?$`)
	float := regexp.MustCompile(`^float(32|64)$`)

	for i := range haystack {

		switch haystack[i] {
		case needle:
			return needle, true

		case "data":
			if needle == "[]uint8" {
				return "data", true
			}

		case "integer":
			if integer.MatchString(needle) {
				return "integer", true
			}

		case "float", "number":
			if integer.MatchString(needle) || float.MatchString(needle) {
				return "number", true
			}

		case "pipe":
			if needle == "*io.PipeWriter" || needle == "*io.PipeReader" {
				return "pipe", true
			}
		}

	}

	return "", false
}
