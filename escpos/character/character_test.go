package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Utility Functions Tests
// ============================================================================

func TestBuildCharacterSize(t *testing.T) {
	tests := []struct {
		name    string
		width   byte
		height  byte
		want    byte
		wantErr bool
	}{
		{
			name:    "normal size",
			width:   1,
			height:  1,
			want:    0x00,
			wantErr: false,
		},
		{
			name:    "double width",
			width:   2,
			height:  1,
			want:    0x10,
			wantErr: false,
		},
		{
			name:    "double height",
			width:   1,
			height:  2,
			want:    0x01,
			wantErr: false,
		},
		{
			name:    "double size",
			width:   2,
			height:  2,
			want:    0x11,
			wantErr: false,
		},
		{
			name:    "maximum size",
			width:   8,
			height:  8,
			want:    0x77,
			wantErr: false,
		},
		{
			name:    "invalid width",
			width:   9,
			height:  1,
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid height",
			width:   1,
			height:  9,
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero width",
			width:   0,
			height:  1,
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero height",
			width:   1,
			height:  0,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := character.BuildCharacterSize(tt.width, tt.height)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuildCharacterSize(%d, %d) error = %v, wantErr %v",
					tt.width, tt.height, err, tt.wantErr)
				return
			}

			var baseErr error
			switch tt.name {
			case "invalid width", "zero width":
				baseErr = character.ErrInvalidCharacterWidth
			case "invalid height", "zero height":
				baseErr = character.ErrInvalidCharacterHeight
			}

			if tt.wantErr && err != nil {
				if !errors.Is(err, baseErr) {
					t.Errorf("BuildCharacterSize(%d, %d) error = %v, want %v",
						tt.width, tt.height, err, baseErr)
				}
				if !errors.Is(err, baseErr) {
					t.Errorf("BuildCharacterSize(%d, %d) error = %v, want %v",
						tt.width, tt.height, err, baseErr)
				}
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("BuildCharacterSize(%d, %d) = %v, want %v",
					tt.width, tt.height, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Character Commands Tests
// ============================================================================

func TestCommands_SetRightSideCharacterSpacing(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name    string
		spacing byte
		want    []byte
	}{
		{
			name:    "no spacing",
			spacing: 0,
			want:    []byte{common.ESC, common.SP, 0},
		},
		{
			name:    "normal spacing",
			spacing: 5,
			want:    []byte{common.ESC, common.SP, 5},
		},
		{
			name:    "maximum spacing",
			spacing: 255,
			want:    []byte{common.ESC, common.SP, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetRightSideCharacterSpacing(tt.spacing)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetRightSideCharacterSpacing(%d) = %#v, want %#v",
					tt.spacing, got, tt.want)
			}
		})
	}
}

func TestCommands_SelectPrintModes(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name string
		mode byte
		want []byte
	}{
		{
			name: "normal mode",
			mode: 0,
			want: []byte{common.ESC, '!', 0},
		},
		{
			name: "font B",
			mode: 0x01,
			want: []byte{common.ESC, '!', 0x01},
		},
		{
			name: "emphasized",
			mode: 0x08,
			want: []byte{common.ESC, '!', 0x08},
		},
		{
			name: "double height",
			mode: 0x10,
			want: []byte{common.ESC, '!', 0x10},
		},
		{
			name: "double width",
			mode: 0x20,
			want: []byte{common.ESC, '!', 0x20},
		},
		{
			name: "underline",
			mode: 0x80,
			want: []byte{common.ESC, '!', 0x80},
		},
		{
			name: "combined modes",
			mode: 0x88, // emphasized + underline
			want: []byte{common.ESC, '!', 0x88},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SelectPrintModes(tt.mode)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SelectPrintModes(%#x) = %#v, want %#v",
					tt.mode, got, tt.want)
			}
		})
	}
}

func TestCommands_SetUnderlineMode(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name    string
		mode    byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "underline off",
			mode:    0,
			want:    []byte{common.ESC, '-', 0},
			wantErr: false,
		},
		{
			name:    "underline 1 dot",
			mode:    1,
			want:    []byte{common.ESC, '-', 1},
			wantErr: false,
		},
		{
			name:    "underline 2 dots",
			mode:    2,
			want:    []byte{common.ESC, '-', 2},
			wantErr: false,
		},
		{
			name:    "underline off ASCII",
			mode:    '0',
			want:    []byte{common.ESC, '-', '0'},
			wantErr: false,
		},
		{
			name:    "underline 1 dot ASCII",
			mode:    '1',
			want:    []byte{common.ESC, '-', '1'},
			wantErr: false,
		},
		{
			name:    "underline 2 dots ASCII",
			mode:    '2',
			want:    []byte{common.ESC, '-', '2'},
			wantErr: false,
		},
		{
			name:    "invalid mode",
			mode:    3,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SetUnderlineMode(tt.mode)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUnderlineMode(%v) error = %v, wantErr %v",
					tt.mode, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidUnderlineMode) {
					t.Errorf("SetUnderlineMode(%v) error = %v, want %v",
						tt.mode, err, character.ErrInvalidUnderlineMode)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SetUnderlineMode(%v) = %#v, want %#v",
					tt.mode, got, tt.want)
			}
		})
	}
}

func TestCommands_SetEmphasizedMode(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name string
		mode byte
		want []byte
	}{
		{
			name: "emphasized off",
			mode: 0,
			want: []byte{common.ESC, 'E', 0},
		},
		{
			name: "emphasized on",
			mode: 1,
			want: []byte{common.ESC, 'E', 1},
		},
		{
			name: "any even number (LSB=0)",
			mode: 0xFE,
			want: []byte{common.ESC, 'E', 0xFE},
		},
		{
			name: "any odd number (LSB=1)",
			mode: 0xFF,
			want: []byte{common.ESC, 'E', 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SetEmphasizedMode(tt.mode)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetEmphasizedMode(%d) = %#v, want %#v",
					tt.mode, got, tt.want)
			}
		})
	}
}

func TestCommands_SelectCharacterFont(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name    string
		font    byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "font A",
			font:    0,
			want:    []byte{common.ESC, 'M', 0},
			wantErr: false,
		},
		{
			name:    "font B",
			font:    1,
			want:    []byte{common.ESC, 'M', 1},
			wantErr: false,
		},
		{
			name:    "font C",
			font:    2,
			want:    []byte{common.ESC, 'M', 2},
			wantErr: false,
		},
		{
			name:    "font A ASCII",
			font:    '0',
			want:    []byte{common.ESC, 'M', '0'},
			wantErr: false,
		},
		{
			name:    "special font A",
			font:    97,
			want:    []byte{common.ESC, 'M', 97},
			wantErr: false,
		},
		{
			name:    "special font B",
			font:    98,
			want:    []byte{common.ESC, 'M', 98},
			wantErr: false,
		},
		{
			name:    "invalid font",
			font:    99,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SelectCharacterFont(tt.font)

			// Standardized error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCharacterFont(%v) error = %v, wantErr %v",
					tt.font, err, tt.wantErr)
				return
			}

			// Check specific error type if expecting error
			if tt.wantErr && err != nil {
				if !errors.Is(err, character.ErrInvalidCharacterFont) {
					t.Errorf("SelectCharacterFont(%v) error = %v, want %v",
						tt.font, err, character.ErrInvalidCharacterFont)
				}
				return
			}

			// Check result if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("SelectCharacterFont(%v) = %#v, want %#v",
					tt.font, got, tt.want)
			}
		})
	}
}

func TestCommands_SelectCharacterSize(t *testing.T) {
	cmd := character.NewCommands()

	tests := []struct {
		name string
		size byte
		want []byte
	}{
		{
			name: "normal size (1x1)",
			size: 0x00,
			want: []byte{common.GS, '!', 0x00},
		},
		{
			name: "double width (2x1)",
			size: 0x10,
			want: []byte{common.GS, '!', 0x10},
		},
		{
			name: "double height (1x2)",
			size: 0x01,
			want: []byte{common.GS, '!', 0x01},
		},
		{
			name: "double size (2x2)",
			size: 0x11,
			want: []byte{common.GS, '!', 0x11},
		},
		{
			name: "maximum size (8x8)",
			size: 0x77,
			want: []byte{common.GS, '!', 0x77},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.SelectCharacterSize(tt.size)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SelectCharacterSize(%#x) = %#v, want %#v",
					tt.size, got, tt.want)
			}
		})
	}
}
