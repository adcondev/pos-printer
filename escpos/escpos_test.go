package escpos

import (
	"bytes"
	"errors"
	"testing"
)

// ============================================================================
// Commands Tests
// ============================================================================

func TestCommands_Raw_EmptyBuffer(t *testing.T) {
	cmd := NewEscposProtocol()

	_, err := cmd.Raw("")

	if !errors.Is(err, errEmptyBuffer) {
		t.Errorf("Commands.Raw(\"\") error = %v, want %v", err, errEmptyBuffer)
	}
}

func TestCommands_Raw_ValidInput(t *testing.T) {
	cmd := NewEscposProtocol()

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
			got, err := cmd.Raw(tt.input)

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
	cmd := NewEscposProtocol()

	// Verify Commands struct is created
	if cmd == nil {
		t.Fatal("NewEscposProtocol() returned nil")
	}

	// Verify Print capability is initialized
	if cmd.Print == nil {
		t.Fatal("NewEscposProtocol() Print capability should not be nil")
	}

	// Verify LineSpace capability is initialized
	if cmd.LineSpace == nil {
		t.Fatal("NewEscposProtocol() LineSpace capability should not be nil")
	}

	// Verify Print has correct type and Page capability
	pc, ok := cmd.Print.(*PrintCommands)
	if !ok {
		t.Fatal("NewEscposProtocol() Print should be of type *PrintCommands")
	}

	if pc.Page == nil {
		t.Fatal("NewEscposProtocol() PrintCommands.Page should not be nil")
	}

	// Verify Page has correct type
	_, ok = pc.Page.(*PagePrint)
	if !ok {
		t.Error("NewEscposProtocol() Page should be of type *PagePrint")
	}

	// Verify LineSpace has correct type
	_, ok = cmd.LineSpace.(*LineSpacingCommands)
	if !ok {
		t.Error("NewEscposProtocol() LineSpace should be of type *LineSpacingCommands")
	}
}

func TestCommands_Integration_PrintWithLineSpacing(t *testing.T) {
	cmd := NewEscposProtocol()

	// Set line spacing
	spacingResult := cmd.LineSpace.SetLineSpacing(40)
	expectedSpacing := []byte{ESC, '3', 40}

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
	if !bytes.Equal(lfResult, []byte{LF}) {
		t.Errorf("Commands.Print.PrintAndLineFeed() = %#v, want %#v",
			lfResult, []byte{LF})
	}
}
