package main

import (
	"strings"
	"testing"
)

func testingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	if (err != nil) != wantErr {
		t.Fatalf("want error: %v, got %v", wantErr, err)
	}
}

func testingAssertOutcomeToBuilderString(t *testing.T, outcome string, builder strings.Builder) {
	if outcome != builder.String() {
		t.Errorf("The outcome is different than expected.")
	}
}
