package main

import (
	"strings"
	"testing"
)

func testingAssertErrorToWantErrorAndOutcomeToBuilderString(t *testing.T, err error, wantErr bool, builder strings.Builder, got string) {
	t.Helper()
	testingAssertErrorToWantError(t, err, wantErr)
	testingAssertEqualStrings(t, builder.String(), got)
}

func testingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if (err != nil) != wantErr {
		t.Errorf("want error: %v, got error: %v", wantErr, err)
	}
}

func testingAssertEqualStrings(t *testing.T, want string, got string) {
	t.Helper()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
