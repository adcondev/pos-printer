package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
)

// ============================================================================
// User Defined Commands Tests
// ============================================================================

func TestUserDefined_SelectUserDefinedCharacterSet(t *testing.T) {
	udc := &character.UserDefinedCommands{}

	tests := []struct {
		name    string
		charSet byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "user-defined off",
			charSet: 0,
			want:    []byte{0x1B, 0x25, 0},
			wantErr: false,
		},
		{
			name:    "user-defined on",
			charSet: 1,
			want:    []byte{0x1B, 0x25, 1},
			wantErr: false,
		},
		{
			name:    "any even number (LSB=0)",
			charSet: 0xFE,
			want:    []byte{0x1B, 0x25, 0xFE},
			wantErr: false,
		},
		{
			name:    "any odd number (LSB=1)",
			charSet: 0xFF,
			want:    []byte{0x1B, 0x25, 0xFF},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := udc.SelectUserDefinedCharacterSet(tt.charSet)

			// No error checking needed as this method doesn't return error
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SelectUserDefinedCharacterSet(%d) = %#v, want %#v",
					tt.charSet, got, tt.want)
			}
		})
	}
}

// TODO: Understand how to test error conditions properly and understand how type is validated.

func TestUserDefinedCommands_DefineUserDefinedCharacters(t *testing.T) {
	udc := &character.UserDefinedCommands{}

	tests := []struct {
		name        string
		height      byte
		startCode   byte
		endCode     byte
		definitions []character.UserDefinedChar
		wantPrefix  []byte
		wantErr     bool
	}{
		{
			name:      "single character definition",
			height:    3,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF}},
			},
			wantPrefix: []byte{0x1B, 0x26, 3, 65, 65},
			wantErr:    false,
		},
		{
			name:      "multiple character definitions",
			height:    2,
			startCode: 65,
			endCode:   66,
			definitions: []character.UserDefinedChar{
				{Width: 2, Data: []byte{0xFF, 0x00, 0xFF, 0x00}},
				{Width: 3, Data: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}},
			},
			wantPrefix: []byte{0x1B, 0x26, 2, 65, 66},
			wantErr:    false,
		},
		{
			name:      "zero width character",
			height:    3,
			startCode: 32,
			endCode:   32,
			definitions: []character.UserDefinedChar{
				{Width: 0, Data: nil},
			},
			wantPrefix: []byte{0x1B, 0x26, 3, 32, 32},
			wantErr:    false,
		},
		{
			name:      "invalid value",
			height:    0,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    true,
		},
		{
			name:      "invalid character code",
			height:    3,
			startCode: 31,
			endCode:   31,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    true,
		},
		{
			name:      "invalid code range",
			height:    3,
			startCode: 66,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    true,
		},
		{
			name:      "definition count mismatch",
			height:    3,
			startCode: 65,
			endCode:   67,
			definitions: []character.UserDefinedChar{
				{Width: 5, Data: []byte{0xFF}},
			},
			wantPrefix: nil,
			wantErr:    true,
		},
		{
			name:      "invalid data length",
			height:    3,
			startCode: 65,
			endCode:   65,
			definitions: []character.UserDefinedChar{
				{Width: 2, Data: []byte{0xFF}}, // Should be 6 bytes (3*2)
			},
			wantPrefix: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := udc.DefineUserDefinedCharacters(tt.height, tt.startCode, tt.endCode, tt.definitions)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("DefineUserDefinedCharacters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				// For wrapped errors with fmt.Errorf, check the base error type
				var baseErr error
				switch tt.name {
				case "invalid value":
					baseErr = character.ErrInvalidYValue
				case "invalid character code":
					baseErr = character.ErrInvalidCharacterCode
				case "invalid code range":
					baseErr = character.ErrInvalidCodeRange
				case "definition count mismatch":
					baseErr = character.ErrDefinitionMismatch
				case "invalid data length":
					baseErr = character.ErrInvalidDataLength
				default:
					baseErr = nil
				}

				if baseErr != nil && !errors.Is(err, baseErr) {
					t.Errorf("DefineUserDefinedCharacters() error = %v, want error containing %v",
						err, baseErr)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && len(got) >= len(tt.wantPrefix) {
				if !bytes.Equal(got[:len(tt.wantPrefix)], tt.wantPrefix) {
					t.Errorf("DefineUserDefinedCharacters() prefix = %#v, want %#v",
						got[:len(tt.wantPrefix)], tt.wantPrefix)
				}
			}
		})
	}
}

func TestUserDefined_CancelUserDefinedCharacter(t *testing.T) {
	udc := &character.UserDefinedCommands{}

	tests := []struct {
		name     string
		charCode byte
		want     []byte
		wantErr  bool
	}{
		{
			name:     "cancel minimum code",
			charCode: 32,
			want:     []byte{0x1B, 0x3F, 32},
			wantErr:  false,
		},
		{
			name:     "cancel typical code",
			charCode: 65,
			want:     []byte{0x1B, 0x3F, 65},
			wantErr:  false,
		},
		{
			name:     "cancel maximum code",
			charCode: 126,
			want:     []byte{0x1B, 0x3F, 126},
			wantErr:  false,
		},
		{
			name:     "invalid code too low",
			charCode: 31,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid code too high",
			charCode: 127,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := udc.CancelUserDefinedCharacter(tt.charCode)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("CancelUserDefinedCharacter(%v) error = %v, wantErr %v",
					tt.charCode, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidCharacterCode) {
					t.Errorf("CancelUserDefinedCharacter(%v) error = %v, want %v",
						tt.charCode, err, character.ErrInvalidCharacterCode)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("CancelUserDefinedCharacter(%v) = %#v, want %#v",
					tt.charCode, got, tt.want)
			}
		})
	}
}
