package utils

import (
	"testing"
)

func TMust[T any](t *testing.T, v T, err error) T {
	t.Helper()

	TMustErr(t, err)

	return v
}

func TMust2[TI, TJ any](t *testing.T, vI TI, vJ TJ, err error) (TI, TJ) {
	t.Helper()

	TMustErr(t, err)

	return vI, vJ
}

func TMustErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
