package test

import (
	"bytes"
	"errors"
	"testing"
)

// AssertBytes compares byte slices and provides detailed error messages
func AssertBytes(t *testing.T, got, want []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !bytes.Equal(got, want) {
		msg := "byte comparison failed"
		if len(msgAndArgs) > 0 {
			msg = msgAndArgs[0].(string)
		}
		t.Errorf("%s\ngot:  %#v\nwant: %#v", msg, got, want)
	}
}

// AssertError checks if an error matches the expected error
func AssertError(t *testing.T, got, want error) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Errorf("unexpected error: %v", got)
		}
		return
	}
	if !errors.Is(got, want) {
		t.Errorf("error = %v, want %v", got, want)
	}
}

// AssertErrorOccurred checks that an error occurred when expected
func AssertErrorOccurred(t *testing.T, err error, wantErr bool, methodName string) bool {
	t.Helper()
	if (err != nil) != wantErr {
		t.Errorf("%s error = %v, wantErr %v", methodName, err, wantErr)
		return false
	}
	return true
}

// AssertUint16Bytes verifies uint16 little-endian byte conversion
func AssertUint16Bytes(t *testing.T, value uint16, wantLow, wantHigh byte) {
	t.Helper()
	gotLow := byte(value & 0xFF)
	gotHigh := byte((value >> 8) & 0xFF)
	if gotLow != wantLow || gotHigh != wantHigh {
		t.Errorf("uint16(%d) bytes = (%#x, %#x), want (%#x, %#x)",
			value, gotLow, gotHigh, wantLow, wantHigh)
	}
}
