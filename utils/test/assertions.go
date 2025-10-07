package test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

// AssertBytes compares byte slices and provides detailed error messages
func AssertBytes(t *testing.T, got, want []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !bytes.Equal(got, want) {
		msg := "byte comparison failed"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
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

// AssertContains checks if a byte slice contains a subsequence
func AssertContains(t *testing.T, haystack, needle []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !bytes.Contains(haystack, needle) {
		msg := "byte slice does not contain expected subsequence"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\nhaystack: %#v\nneedle:   %#v", msg, haystack, needle)
	}
}

// AssertNotContains checks if a byte slice does not contain a subsequence
func AssertNotContains(t *testing.T, haystack, needle []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if bytes.Contains(haystack, needle) {
		msg := "byte slice contains unexpected subsequence"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\nhaystack: %#v\nneedle:   %#v", msg, haystack, needle)
	}
}

// AssertByteCount verifies the number of occurrences of a byte sequence
func AssertByteCount(t *testing.T, data []byte, pattern []byte, expectedCount int, msgAndArgs ...interface{}) {
	t.Helper()
	count := bytes.Count(data, pattern)
	if count != expectedCount {
		msg := "byte pattern count mismatch"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\npattern: %#v\ngot count: %d, want: %d", msg, pattern, count, expectedCount)
	}
}

// AssertLength checks if a byte slice has the expected length
func AssertLength(t *testing.T, data []byte, expectedLength int, msgAndArgs ...interface{}) {
	t.Helper()
	if len(data) != expectedLength {
		msg := "byte slice length mismatch"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\ngot length: %d, want: %d\ndata: %#v", msg, len(data), expectedLength, data)
	}
}

// AssertHasPrefix checks if a byte slice starts with a specific prefix
func AssertHasPrefix(t *testing.T, data, prefix []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !bytes.HasPrefix(data, prefix) {
		msg := "byte slice does not have expected prefix"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\ndata:   %#v\nprefix: %#v", msg, data, prefix)
	}
}

// AssertHasSuffix checks if a byte slice ends with a specific suffix
func AssertHasSuffix(t *testing.T, data, suffix []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !bytes.HasSuffix(data, suffix) {
		msg := "byte slice does not have expected suffix"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			} else {
				msg = fmt.Sprint(msgAndArgs[0])
			}
		}
		t.Errorf("%s\ndata:   %#v\nsuffix: %#v", msg, data, suffix)
	}
}
