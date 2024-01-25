package test

import (
	"strings"
	"testing"
)

func TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t *testing.T, err error, wantErr bool, builder strings.Builder, got string) {
	t.Helper()
	TestingAssertErrorToWantError(t, err, wantErr)
	TestingAssertEqualStrings(t, builder.String(), got)
}

func TestingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if (err != nil) != wantErr {
		t.Errorf("want error: %v, got error: %v", wantErr, err)
	}
}

func TestingAssertEqualStrings(t *testing.T, want string, got string) {
	t.Helper()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
