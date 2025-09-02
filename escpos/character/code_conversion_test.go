package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Code Conversion Commands Tests
// ============================================================================

func TestCodeConversionCommands_SelectCharacterEncodeSystem(t *testing.T) {
	cc := character.NewCodeConversionCommands()

	tests := []struct {
		name     string
		encoding byte
		want     []byte
		wantErr  bool
	}{
		{
			name:     "1-byte encoding",
			encoding: 1,
			want:     []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, 1},
			wantErr:  false,
		},
		{
			name:     "UTF-8 encoding",
			encoding: 2,
			want:     []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, 2},
			wantErr:  false,
		},
		{
			name:     "1-byte encoding ASCII",
			encoding: '1',
			want:     []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, '1'},
			wantErr:  false,
		},
		{
			name:     "UTF-8 encoding ASCII",
			encoding: '2',
			want:     []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, '2'},
			wantErr:  false,
		},
		{
			name:     "invalid encoding",
			encoding: 3,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.SelectCharacterEncodeSystem(tt.encoding)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCharacterEncodeSystem(%v) error = %v, wantErr %v",
					tt.encoding, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidEncoding) {
					t.Errorf("SelectCharacterEncodeSystem(%v) error = %v, want %v",
						tt.encoding, err, character.ErrInvalidEncoding)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectCharacterEncodeSystem(%v) = %#v, want %#v",
					tt.encoding, got, tt.want)
			}
		})
	}
}

func TestCodeConversionCommands_SetFontPriority(t *testing.T) {
	cc := character.NewCodeConversionCommands()

	tests := []struct {
		name     string
		priority byte
		fontType byte
		want     []byte
		wantErr  bool
	}{
		{
			name:     "first priority ANK font",
			priority: 0,
			fontType: 0,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 0, 0},
			wantErr:  false,
		},
		{
			name:     "second priority Japanese",
			priority: 1,
			fontType: 11,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 1, 11},
			wantErr:  false,
		},
		{
			name:     "first priority Simplified Chinese",
			priority: 0,
			fontType: 20,
			want:     []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, 0, 20},
			wantErr:  false,
		},
		{
			name:     "invalid priority",
			priority: 2,
			fontType: 0,
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid font type",
			priority: 0,
			fontType: 99,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.SetFontPriority(tt.priority, tt.fontType)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetFontPriority(%v, %v) error = %v, wantErr %v",
					tt.priority, tt.fontType, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidFontPriority) {
					t.Errorf("SetFontPriority(%v, %v) error = %v, want %v",
						tt.priority, tt.fontType, err, character.ErrInvalidFontPriority)
				}
				if !errors.Is(err, character.ErrInvalidFontPriority) {
					t.Errorf("SetFontPriority(%v, %v) error = %v, want %v",
						tt.priority, tt.fontType, err, character.ErrInvalidFontPriority)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetFontPriority(%v, %v) = %#v, want %#v",
					tt.priority, tt.fontType, got, tt.want)
			}
		})
	}
}
