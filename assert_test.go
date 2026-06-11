package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// assert/require helpers local to the test suite, written to keep the package
// free of a third-party test framework. require* helpers stop the test on
// failure (t.Fatalf); assert* helpers report and continue (t.Errorf).

func label(msg []any) string {
	if len(msg) == 0 {
		return ""
	}
	if s, ok := msg[0].(string); ok && len(msg) == 1 {
		return " (" + s + ")"
	}
	return " (" + fmt.Sprint(msg...) + ")"
}

func requireNoError(t *testing.T, err error, msg ...any) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error%s: %v", label(msg), err)
	}
}

func requireErrorIs(t *testing.T, err, target error, msg ...any) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error to be %v, got %v%s", target, err, label(msg))
	}
}

func requireEqual(t *testing.T, want, got any, msg ...any) {
	t.Helper()
	if !objectsEqual(want, got) {
		t.Fatalf("not equal%s:\n want: %#v\n  got: %#v", label(msg), want, got)
	}
}

// objectsEqual mirrors testify's EqualValues: a direct deep-equal, falling back
// to a type-conversion compare so values of convertible-but-distinct types
// (e.g. anonymous structs that differ only by tag order) still match.
func objectsEqual(want, got any) bool {
	if reflect.DeepEqual(want, got) {
		return true
	}
	if want == nil || got == nil {
		return false
	}
	wv, gv := reflect.ValueOf(want), reflect.ValueOf(got)
	if gv.Type().ConvertibleTo(wv.Type()) {
		return reflect.DeepEqual(want, gv.Convert(wv.Type()).Interface())
	}
	return false
}

func requireLen(t *testing.T, collection any, length int, msg ...any) {
	t.Helper()
	v := reflect.ValueOf(collection)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		if v.Len() != length {
			t.Fatalf("expected length %d, got %d%s", length, v.Len(), label(msg))
		}
	default:
		t.Fatalf("value of type %T has no length%s", collection, label(msg))
	}
}

func requireNotNil(t *testing.T, obj any, msg ...any) {
	t.Helper()
	if obj == nil {
		t.Fatalf("expected non-nil value%s", label(msg))
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			t.Fatalf("expected non-nil value%s", label(msg))
		}
	default:
		// other kinds can't be nil
	}
}

func requireErrorContains(t *testing.T, err error, contains string, msg ...any) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil%s", contains, label(msg))
	}
	if !strings.Contains(err.Error(), contains) {
		t.Fatalf("expected error %q to contain %q%s", err.Error(), contains, label(msg))
	}
}
