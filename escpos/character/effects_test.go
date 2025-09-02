package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Effects Commands Tests
// ============================================================================

func TestEffectsCommands_SelectCharacterColor(t *testing.T) {
	ec := character.NewEffectsCommands()

	tests := []struct {
		name    string
		color   byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "no color",
			color:   '0',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, '0'},
			wantErr: false,
		},
		{
			name:    "color 1",
			color:   '1',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, '1'},
			wantErr: false,
		},
		{
			name:    "color 2",
			color:   '2',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, '2'},
			wantErr: false,
		},
		{
			name:    "color 3",
			color:   '3',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, '3'},
			wantErr: false,
		},
		{
			name:    "invalid color",
			color:   '4',
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ec.SelectCharacterColor(tt.color)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCharacterColor(%v) error = %v, wantErr %v",
					tt.color, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidCharacterColor) {
					t.Errorf("SelectCharacterColor(%v) error = %v, want %v",
						tt.color, err, character.ErrInvalidCharacterColor)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectCharacterColor(%v) = %#v, want %#v",
					tt.color, got, tt.want)
			}
		})
	}
}

func TestEffectsCommands_SelectBackgroundColor(t *testing.T) {
	ec := character.NewEffectsCommands()

	tests := []struct {
		name    string
		color   byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "no background",
			color:   '0',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x31, '0'},
			wantErr: false,
		},
		{
			name:    "background color 1",
			color:   '1',
			want:    []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x31, '1'},
			wantErr: false,
		},
		{
			name:    "invalid background",
			color:   '5',
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ec.SelectBackgroundColor(tt.color)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectBackgroundColor(%v) error = %v, wantErr %v",
					tt.color, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidBackgroundColor) {
					t.Errorf("SelectBackgroundColor(%v) error = %v, want %v",
						tt.color, err, character.ErrInvalidBackgroundColor)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectBackgroundColor(%v) = %#v, want %#v",
					tt.color, got, tt.want)
			}
		})
	}
}

func TestEffectsCommands_SetCharacterShadowMode(t *testing.T) {
	ec := &character.EffectsCommands{}

	tests := []struct {
		name        string
		shadowMode  byte
		shadowColor byte
		want        []byte
		wantErr     bool
	}{
		{
			name:        "shadow off no color",
			shadowMode:  0,
			shadowColor: '0',
			want:        []byte{common.GS, '(', 'N', 0x03, 0x00, 0x32, 0, '0'},
			wantErr:     false,
		},
		{
			name:        "shadow on color 1",
			shadowMode:  1,
			shadowColor: '1',
			want:        []byte{common.GS, '(', 'N', 0x03, 0x00, 0x32, 1, '1'},
			wantErr:     false,
		},
		{
			name:        "shadow off ASCII",
			shadowMode:  '0',
			shadowColor: '2',
			want:        []byte{common.GS, '(', 'N', 0x03, 0x00, 0x32, '0', '2'},
			wantErr:     false,
		},
		{
			name:        "shadow on ASCII",
			shadowMode:  '1',
			shadowColor: '3',
			want:        []byte{common.GS, '(', 'N', 0x03, 0x00, 0x32, '1', '3'},
			wantErr:     false,
		},
		{
			name:        "invalid shadow mode",
			shadowMode:  2,
			shadowColor: '1',
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "invalid shadow color",
			shadowMode:  0,
			shadowColor: '4',
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ec.SetCharacterShadowMode(tt.shadowMode, tt.shadowColor)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCharacterShadowMode(%v, %v) error = %v, wantErr %v",
					tt.shadowMode, tt.shadowColor, err, tt.wantErr)
				return
			}

			var baseErr error
			switch tt.name {
			case "invalid shadow mode":
				baseErr = character.ErrInvalidShadowMode
			case "invalid shadow color":
				baseErr = character.ErrInvalidShadowColor
			default:
				baseErr = nil
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, baseErr) {
					t.Errorf("SetCharacterShadowMode(%v, %v) error = %v, want %v",
						tt.shadowMode, tt.shadowColor, err, baseErr)
				}
				if !errors.Is(err, baseErr) {
					t.Errorf("SetCharacterShadowMode(%v, %v) error = %v, want %v",
						tt.shadowMode, tt.shadowColor, err, baseErr)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetCharacterShadowMode(%v, %v) = %#v, want %#v",
					tt.shadowMode, tt.shadowColor, got, tt.want)
			}
		})
	}
}
