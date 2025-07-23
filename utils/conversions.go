package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	// "github.com/davecgh/go-spew/spew"
)

func ToFloat(value any) (float64, error) {
	switch v := value.(type) {
	case float32, float64:
		return value.(float64), nil
	case string:
		float, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse float64 default value: %s: %w", v, err)
		}
		return float, nil
	default:
		return 0, fmt.Errorf("unsupported float type: %T", value)
	}
}

func ParseBool(val string) bool {
	switch strings.ToLower(val) {
	case "true", "yes", "1":
		return true
	case "false", "no", "0":
		return false
	default:
		return false
	}
}

func ToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return value.(int), nil
	case string:
		integer, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("failed to parse int value: %s: %w", v, err)
		}
		return integer, nil
	default:
		return 0, fmt.Errorf("failed to parse int from type: %T, value: %v", value, value)
	}
}

func ToString(value any) (string, error) {
	// spew.Dump("sintax.ToString", value)
	switch value.(type) {
	case string:
		return value.(string), nil
	case int:
		return strconv.Itoa(value.(int)), nil
	case float32, float64:
		return fmt.Sprintf("%f", value), nil
	default:
		return "", fmt.Errorf("unsupported string type: %T", value)
	}
}

func ToAnySlice(value any) ([]any, error) {
	switch v := value.(type) {
	case []any:
		return value.([]any), nil
	case string:
		if v == "" {
			return []any{}, nil
		}
		return []any{v}, nil
	case []string:
		anySlice := make([]any, len(v))
		for i, val := range v {
			anySlice[i] = val
		}
		return anySlice, nil
	default:
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Slice {
			return nil, fmt.Errorf("unsupported slice type: %T: %v", value, value)
		}

		anySlice := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			anySlice[i] = val.Index(i).Interface()
		}

		return anySlice, nil
	}
}
