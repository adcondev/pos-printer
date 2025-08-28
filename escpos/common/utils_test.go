package common_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
)

func TestUtils_IsBufOk_ValidInput(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr error
	}{
		{"empty buffer", []byte{}, common.ErrEmptyBuffer},
		{"valid buffer", []byte{1, 2, 3}, nil},
		{"max buffer", make([]byte, common.MaxBuf), nil},
		{"overflow buffer", make([]byte, common.MaxBuf+1), common.ErrBufferOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.IsBufOk(tt.buf)
			if !errors.Is(tt.wantErr, err) {
				t.Errorf("IsBufOk len=%d error = %v; want %v", len(tt.buf), err, tt.wantErr)
			}
		})
	}
}

func TestUtils_Format_ByteSequence(t *testing.T) {
	input := []byte("a\n\t\rB")
	want := []byte{'a', common.LF, common.HT, common.CR, 'B'}
	// clone input to avoid modifying original
	data := append([]byte(nil), input...)
	got := common.Format(data)
	if !bytes.Equal(got, want) {
		t.Errorf("Format(%q) = %v; want %v", input, got, want)
	}
}

func TestUtils_LengthLowHigh_ValidInput(t *testing.T) {
	tests := []struct {
		length  int
		wantDL  byte
		wantDH  byte
		wantErr error
	}{
		{0, 0, 0, nil},
		{1, 1, 0, nil},
		{0x1234, 0x34, 0x12, nil},
		{0xFFFF, 0xFF, 0xFF, nil},
		{-1, 0, 0, common.ErrLengthOutOfRange},
		{0x10000, 0, 0, common.ErrLengthOutOfRange},
	}
	for _, tt := range tests {
		dL, dH, err := common.LengthLowHigh(tt.length)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("LengthLowHigh(%d) error = %v; want %v", tt.length, err, tt.wantErr)
		}
		if err == nil {
			if dL != tt.wantDL || dH != tt.wantDH {
				t.Errorf("LengthLowHigh(%d) = (%#x,%#x); want (%#x,%#x)", tt.length, dL, dH, tt.wantDL, tt.wantDH)
			}
		}
	}
}
