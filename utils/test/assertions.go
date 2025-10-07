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

// AssertEmpty checks if a byte slice is empty
func AssertNotEmpty(t *testing.T, data []byte, s string) {
	if len(data) == 0 {
		t.Errorf("%s: data should not be empty", s)
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

// AssertAlphanumeric checks if all bytes are alphanumeric
func AssertAlphanumeric(t *testing.T, data []byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsAlphanumeric(data) {
		msg := formatMessage("expected all alphanumeric bytes", msgAndArgs...)
		for i, b := range data {
			if !IsAlphanumericByte(b) {
				t.Errorf("%s\nfirst invalid byte at index %d: %#x (%c)\ndata: %q",
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

// AssertInRange checks if all bytes are within specified range
func AssertInRange(t *testing.T, data []byte, min, max byte, msgAndArgs ...interface{}) {
	t.Helper()
	if !IsInRange(data, min, max) {
		msg := formatMessage(fmt.Sprintf("expected bytes in range [%d, %d]", min, max), msgAndArgs...)
		for i, b := range data {
			if b < min || b > max {
				t.Errorf("%s\nout of range byte at index %d: %d\ndata: %#v",
					msg, i, b, data)
				return
			}
		}
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
func AssertValidLength(t *testing.T, data []byte, min, max int, msgAndArgs ...interface{}) {
	t.Helper()
	if !ValidateLength(data, min, max) {
		msg := formatMessage(fmt.Sprintf("expected length between %d and %d", min, max), msgAndArgs...)
		t.Errorf("%s\nactual length: %d\ndata: %#v", msg, len(data), data)
	}
}

// AssertInvalidLength checks if data length is outside bounds
func AssertInvalidLength(t *testing.T, data []byte, min, max int, msgAndArgs ...interface{}) {
	t.Helper()
	if ValidateLength(data, min, max) {
		msg := formatMessage(fmt.Sprintf("expected length outside %d to %d", min, max), msgAndArgs...)
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

// IsAlphanumericByte Helper to check if a single byte is alphanumeric
func IsAlphanumericByte(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')
}

// AssertCommandPrefix checks if a command starts with expected prefix bytes
func AssertCommandPrefix(t *testing.T, command []byte, expectedPrefix []byte, msgAndArgs ...interface{}) {
	t.Helper()
	AssertHasPrefix(t, command, expectedPrefix, msgAndArgs...)
}

// AssertBarcodeCommand verifies a complete barcode command structure
func AssertBarcodeCommand(t *testing.T, command []byte, symbology byte, data []byte, isNullTerminated bool) {
	t.Helper()
	// Check command prefix
	AssertHasPrefix(t, command, []byte{0x1D, 0x6B, symbology}, "barcode command prefix")

	if isNullTerminated {
		// Function A format
		AssertHasNullTerminator(t, command, "Function A barcode")
		AssertContains(t, command, data, "barcode data")
	} else {
		// Function B format - has length byte
		if len(command) > 3 {
			lengthByte := command[3]
			if int(lengthByte) != len(data) {
				t.Errorf("length byte = %d, want %d", lengthByte, len(data))
			}
		}
	}
}
