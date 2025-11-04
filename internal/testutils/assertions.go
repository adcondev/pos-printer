// Package testutils provides utility functions for testing byte slices and error handling.
package testutils

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

// AssertNotEmpty checks if a byte slice is empty
func AssertNotEmpty(t *testing.T, data []byte, s string) {
	if len(data) == 0 {
		t.Errorf("%s: data should not be empty", s)
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

// AssertNumeric checks if all bytes are numeric and provides detailed error
func AssertNumeric(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsNumeric(data) {
		msg := formatMessage("expected all numeric bytes", msgAndArgs...)
		// Find first non-numeric byte for better error reporting
		for i, b := range data {
			if b < '0' || b > '9' {
				t.Errorf("%s\nfirst non-numeric byte at index %d: %#x (%c)\ndata: %q",
					msg, i, b, b, data)
				return
			}
		}
	}
}

// AssertUppercase verifies all alphabetic bytes are uppercase
func AssertUppercase(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsUppercase(data) {
		msg := formatMessage("expected all uppercase letters", msgAndArgs...)
		for i, b := range data {
			if b >= 'a' && b <= 'z' {
				t.Errorf("%s\nfirst lowercase at index %d: %c\ndata: %q",
					msg, i, b, data)
				return
			}
		}
	}
}

// AssertPrintableASCII checks if all bytes are printable ASCII
func AssertPrintableASCII(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsPrintableASCII(data) {
		msg := formatMessage("expected all printable ASCII", msgAndArgs...)
		for i, b := range data {
			if b < 32 || b > 126 {
				t.Errorf("%s\nnon-printable byte at index %d: %#x\ndata: %#v",
					msg, i, b, data)
				return
			}
		}
	}
}

// AssertHasNullTerminator checks if data ends with null byte
func AssertHasNullTerminator(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !HasNullTerminator(data) {
		msg := formatMessage("expected null terminator", msgAndArgs...)
		if len(data) == 0 {
			t.Errorf("%s\ndata is empty", msg)
		} else {
			t.Errorf("%s\nlast byte: %#x, expected 0x00\ndata: %#v",
				msg, data[len(data)-1], data)
		}
	}
}

// AssertEvenLength checks if byte slice has even length
func AssertEvenLength(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsEvenLength(data) {
		msg := formatMessage("expected even length", msgAndArgs...)
		t.Errorf("%s\nlength: %d (odd)\ndata: %#v", msg, len(data), data)
	}
}

// AssertContainsOnly checks if data contains only bytes from allowed set
func AssertContainsOnly(t *testing.T, data []byte, allowed []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !ContainsOnly(data, allowed) {
		msg := formatMessage("expected only allowed bytes", msgAndArgs...)
		allowedMap := make(map[byte]bool)
		for _, b := range allowed {
			allowedMap[b] = true
		}
		for i, b := range data {
			if !allowedMap[b] {
				t.Errorf("%s\nunallowed byte at index %d: %#x (%c)\nallowed: %#v\ndata: %q",
					msg, i, b, b, allowed, data)
				return
			}
		}
	}
}

// AssertValidLength checks if data length is within bounds
func AssertValidLength(t *testing.T, data []byte, minim, maxim int, msgAndArgs ...interface{}) {
	t.Helper()
	if !ValidateLength(data, minim, maxim) {
		msg := formatMessage(fmt.Sprintf("expected length between %d and %d", minim, maxim), msgAndArgs...)
		t.Errorf("%s\nactual length: %d\ndata: %#v", msg, len(data), data)
	}
}

// Helper to format messages
func formatMessage(defaultMsg string, msgAndArgs ...interface{}) string {
	if len(msgAndArgs) > 0 {
		if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 {
			return fmt.Sprintf(format, msgAndArgs[1:]...)
		}
		return fmt.Sprint(msgAndArgs[0])
	}
	return defaultMsg
}
