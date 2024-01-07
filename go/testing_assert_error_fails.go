package main

import (
	"strings"
	"testing"
)

func testingAssertErrorToWantErrorAndOutcomeToBuilderString(t *testing.T, err error, wantErr bool, outcome string, builder strings.Builder) {
	testingAssertErrorToWantError(t, err, wantErr)
	if outcome != builder.String() {
		t.Errorf("The outcome is different than expected.")
	}
}

func testingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	if (err != nil) != wantErr {
		t.Fatalf("want error: %v, got %v", wantErr, err)
	}
}
