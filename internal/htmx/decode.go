package htmx

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func StringOrSlice(str interface{}) ([]string, error) {
	var err error
	result := make([]string, 0)
	v := reflect.ValueOf(str)

	switch v.Kind() {
	case reflect.String:
		result = append(result, str.(string))
		break

	case reflect.Slice:
		result, err = stringSlice(str)
		if err != nil {
			return nil, errors.Wrap(err, "htmx.StringOrSlice")
		}
		break

	default:
		return nil, nil
	}

	return result, nil
}

func stringSlice(input interface{}) ([]string, error) {
	interfaces, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("stringSlice: failed to cast input to slice")
	}

	result := make([]string, len(interfaces))
	for i, val := range interfaces {
		strVal, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("stringSlice: failed to cast value to string")
		}

		result[i] = strVal
	}

	return result, nil
}
