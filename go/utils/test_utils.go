package utils

import (
	"testing"
)

func TMust[T any](t *testing.T, v T, err error) T {
	t.Helper()

	TMustErr(t, err)

	return v
}

func TMustErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
