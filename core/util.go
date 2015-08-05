package core

import (
	"fmt"
)

func ConvertByte(in interface{}) ([]byte, error) {
	switch in.(type) {
	case string:
		return []byte(in.(string)), nil
	case []byte:
		return in.([]byte), nil
	default:
		return nil, fmt.Errorf("%s: %v", "Unable to convert value", in)
	}
}
