// Package utils provides type-coercion helpers used to convert loosely-typed
// input values (from tags, env vars, decoded config) into concrete field types.
package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ToFloat converts a numeric or numeric-string value to a float64.
// It returns an error for unsupported types or unparseable strings.
func ToFloat(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		float, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse float64 value: %s: %w", v, err)
		}
		return float, nil
	default:
		return 0, fmt.Errorf("unsupported float type: %T", value)
	}
}

// ParseBool reports whether val is a truthy string
// ("true", "yes", "1", case-insensitive)
// everything else, including unrecognized input, is false.
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

// ToInt converts a numeric or numeric-string value to an int.
// Floats with a fractional part and unparseable strings return an error.
func ToInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		if float32(int(v)) != v {
			return 0, fmt.Errorf("failed to parse int from float32 with fractional part: %v", v)
		}
		return int(v), nil
	case float64:
		if float64(int(v)) != v {
			return 0, fmt.Errorf("failed to parse int from float64 with fractional part: %v", v)
		}
		return int(v), nil
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

// ToString converts a string, int, or float value to its string form
func ToString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float32, float64:
		return fmt.Sprintf("%f", value), nil
	default:
		return "", fmt.Errorf("unsupported string type: %T", value)
	}
}

// ToAnySlice converts a slice value to []any.
// []any and []string are handled directly
// a non-empty string becomes a single-element slice
// any other slice is converted element-wise via reflection.
// Non-slice types return an error.
func ToAnySlice(value any) ([]any, error) {
	switch v := value.(type) {
	case []any:
		return v, nil
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
		for i := range val.Len() {
			anySlice[i] = val.Index(i).Interface()
		}

		return anySlice, nil
	}
}
