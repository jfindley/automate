package core

import (
	"bytes"
	"errors"
	"testing"
)

var testingTypes = map[string]interface{}{
	"string":  "test",
	"byte":    []byte("test"),
	"int":     10,
	"int16":   int16(10),
	"int32":   int32(10),
	"int64":   int64(10),
	"uint":    int16(10),
	"uint16":  int16(10),
	"uint32":  int32(10),
	"uint64":  int64(10),
	"uintptr": uintptr(10),
	"float32": float32(10),
	"float64": float64(10),
	"error":   errors.New("test"),
	"bool":    true,
}

func TestConvertByte(t *testing.T) {
	target := []byte("test")

	for k, v := range testingTypes {
		switch k {
		case "string", "byte":
			out, err := ConvertByte(v)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(out, target) != 0 {
				t.Errorf("Valid type '%s' does not match", k)
			}

		default:
			_, err := ConvertByte(v)
			if err == nil {
				t.Errorf("Invalid type '%s' did not produce an error", k)
			}
		}
	}
}

func TestConvertInt(t *testing.T) {
	var target int64 = 10

	for k, v := range testingTypes {
		switch k {
		case "int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64", "float32", "float64":
			out, err := ConvertInt(v)
			if err != nil {
				t.Error(err)
			}
			if out != target {
				t.Errorf("Valid type '%s' does not match", k)
			}

		default:
			_, err := ConvertInt(v)
			if err == nil {
				t.Errorf("Invalid type '%s' did not produce an error", k)
			}
		}
	}

}
