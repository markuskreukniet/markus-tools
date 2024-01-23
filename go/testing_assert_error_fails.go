package main

import (
	"strings"
	"testing"
)

func testingAssertErrorToWantErrorAndOutcomeToBuilderString(t *testing.T, err error, wantErr bool, outcome string, builder strings.Builder) {
	t.Helper()
	testingAssertErrorToWantError(t, err, wantErr)
	if outcome != builder.String() {
		t.Errorf("The outcome is different than expected.")
	}
}

func testingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if (err != nil) != wantErr {
		t.Errorf("want error: %v, got %v", wantErr, err)
	}
}
