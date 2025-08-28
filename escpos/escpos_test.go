package escpos_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Commands Tests
// ============================================================================

func TestCommands_Raw_EmptyBuffer(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	_, err := cmd.Raw([]byte(""))

	if !errors.Is(err, common.ErrEmptyBuffer) {
		t.Errorf("Commands.Raw(\"\") error = %v, want %v", err, common.ErrEmptyBuffer)
	}
}

func TestCommands_Raw_ValidInput(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{
			name:  "simple text",
			input: "hello",
			want:  []byte("hello"),
		},
		{
			name:  "text with spaces",
			input: "hello world",
			want:  []byte("hello world"),
		},
		{
			name:  "special characters",
			input: "!@#$%",
			want:  []byte("!@#$%"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.Raw([]byte(tt.input))

			if err != nil {
				t.Errorf("Commands.Raw(%q) unexpected error: %v", tt.input, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Commands.Raw(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewEscposProtocol_Initialization(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	// Verify Commands struct is created
	if cmd == nil {
		t.Fatal("NewEscposCommands() returned nil")
	}

	// Verify Print capability is initialized
	if cmd.Print == nil {
		t.Fatal("NewEscposCommands() Print capability should not be nil")
	}

	// Verify LineSpace capability is initialized
	if cmd.LineSpace == nil {
		t.Fatal("NewEscposCommands() LineSpace capability should not be nil")
	}

	// Verify Print has correct type and Page capability
	pc, ok := cmd.Print.(*print.Commands)
	if !ok {
		t.Fatal("NewEscposCommands() Print should be of type *PrintCommands")
	}

	if pc.Page == nil {
		t.Fatal("NewEscposCommands() PrintCommands.Page should not be nil")
	}

	// Verify Page has correct type
	_, ok = pc.Page.(*print.PagePrint)
	if !ok {
		t.Error("NewEscposCommands() Page should be of type *PagePrint")
	}

	// Verify LineSpace has correct type
	_, ok = cmd.LineSpace.(*linespacing.Commands)
	if !ok {
		t.Error("NewEscposCommands() LineSpace should be of type *LineSpacingCommands")
	}
}

func TestCommands_Integration_PrintWithLineSpacing(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	// Set line spacing
	spacingResult := cmd.LineSpace.SetLineSpacing(40)
	expectedSpacing := []byte{common.ESC, '3', 40}

	if !bytes.Equal(spacingResult, expectedSpacing) {
		t.Errorf("Commands.LineSpace.SetLineSpacing(40) = %#v, want %#v",
			spacingResult, expectedSpacing)
	}

	// Print text
	textResult, err := cmd.Print.Text("Test")
	if err != nil {
		t.Fatalf("Commands.Print.Text(\"Test\") unexpected error: %v", err)
	}

	if !bytes.Equal(textResult, []byte("Test")) {
		t.Errorf("Commands.Print.Text(\"Test\") = %#v, want %#v",
			textResult, []byte("Test"))
	}

	// Line feed (which should use the line spacing)
	lfResult := cmd.Print.PrintAndLineFeed()
	if !bytes.Equal(lfResult, []byte{print.LF}) {
		t.Errorf("Commands.Print.PrintAndLineFeed() = %#v, want %#v",
			lfResult, []byte{print.LF})
	}
}
