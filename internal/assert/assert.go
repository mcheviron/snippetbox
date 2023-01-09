package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	// NOTE: this signals to the Go test runner that this is a helper function
	// and when Errorf is invoked, the file that called this function will
	// be reproted instead of this file and function
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func StringContains(t *testing.T, got, subString string) {
	t.Helper()
	if !strings.Contains(got, subString) {
		t.Errorf("got: %q, expected to contain: %q", got, subString)
	}
}
