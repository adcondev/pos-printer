package escpos

import (
	"bytes"
	"errors"
	"testing"
)

func TestUtils_IsBufOk_ValidInput(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr error
	}{
		{"empty buffer", []byte{}, errEmptyBuffer},
		{"valid buffer", []byte{1, 2, 3}, nil},
		{"max buffer", make([]byte, MaxBuf), nil},
		{"overflow buffer", make([]byte, MaxBuf+1), errBufferOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isBufOk(tt.buf)
			if !errors.Is(tt.wantErr, err) {
				t.Errorf("isBufOk len=%d error = %v; want %v", len(tt.buf), err, tt.wantErr)
			}
		})
	}
}

func TestUtils_Format_ByteSequence(t *testing.T) {
	input := []byte("a\n\t\rB")
	want := []byte{'a', LF, HT, CR, 'B'}
	// clone input to avoid modifying original
	data := append([]byte(nil), input...)
	got := format(data)
	if !bytes.Equal(got, want) {
		t.Errorf("format(%q) = %v; want %v", input, got, want)
	}
}

func TestUtils_LengthLowHigh_LittleEndian(t *testing.T) {
	tests := []struct {
		length  int
		wantDL  byte
		wantDH  byte
		wantErr error
	}{
		{-1, 0, 0, errNegativeInt},
		{0, 0, 0, nil},
		{1, 1, 0, nil},
		{0x1234, 0x34, 0x12, nil},
		{0xFFFF, 0xFF, 0xFF, nil},
		{0x10000, 0x00, 0x00, nil},
	}
	for _, tt := range tests {
		dL, dH, err := lengthLowHigh(tt.length)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("lengthLowHigh(%d) error = %v; want %v", tt.length, err, tt.wantErr)
		}
		if err == nil {
			if dL != tt.wantDL || dH != tt.wantDH {
				t.Errorf("lengthLowHigh(%d) = (%#x,%#x); want (%#x,%#x)", tt.length, dL, dH, tt.wantDL, tt.wantDH)
			}
		}
	}
}
