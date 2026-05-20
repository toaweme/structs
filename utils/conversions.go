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
			// spew.Dump("sintax.ToAnySlice", i, val)
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
