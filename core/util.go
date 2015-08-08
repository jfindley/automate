package core

import (
	"fmt"
	"strconv"
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

func ConvertInt(in interface{}) (int64, error) {
	switch in.(type) {
	case string:
		return strconv.ParseInt(in.(string), 10, 64)
	case int:
		return int64(in.(int)), nil
	case int16:
		return int64(in.(int16)), nil
	case int32:
		return int64(in.(int32)), nil
	case int64:
		return in.(int64), nil
	case uint:
		return int64(in.(uint)), nil
	case uint16:
		return int64(in.(uint16)), nil
	case uint32:
		return int64(in.(uint32)), nil
	case uint64:
		return int64(in.(uint64)), nil
	case float32:
		return int64(in.(float32)), nil
	case float64:
		return int64(in.(float64)), nil
	default:
		return 0, fmt.Errorf("%s: %v", "Unable to convert value", in)
	}
}
