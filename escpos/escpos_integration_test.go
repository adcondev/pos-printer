package escpos_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Cross-Capability Integration Tests
// ============================================================================

func TestIntegration_PrintWithLineSpacing_RealImplementations(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	// Set custom line spacing
	spacingCmd := cmd.LineSpace.SetLineSpacing(50)
	expectedSpacing := []byte{common.ESC, '3', 50}

	if !bytes.Equal(spacingCmd, expectedSpacing) {
		t.Errorf("SetLineSpacing(50) = %#v, want %#v", spacingCmd, expectedSpacing)
	}

	// Print text with the new spacing
	textCmd, err := cmd.Print.Text("Hello World")
	if err != nil {
		t.Fatalf("Text() unexpected error: %v", err)
	}

	// Line feed (uses the spacing)
	lfCmd := cmd.Print.PrintAndLineFeed()

	// Verify commands
	if !bytes.Equal(textCmd, []byte("Hello World")) {
		t.Errorf("Text() = %#v, want %#v", textCmd, []byte("Hello World"))
	}
	if !bytes.Equal(lfCmd, []byte{print.LF}) {
		t.Errorf("PrintAndLineFeed() = %#v, want %#v", lfCmd, []byte{print.LF})
	}
}

func TestIntegration_CompleteReceiptFlow(t *testing.T) {
	cmd := escpos.NewEscposCommands()

	// Build a complete receipt flow
	commands := []struct {
		name   string
		action func() ([]byte, error)
		verify func([]byte) error
	}{
		{
			name: "set line spacing",
			action: func() ([]byte, error) {
				return cmd.LineSpace.SetLineSpacing(40), nil
			},
			verify: func(result []byte) error {
				expected := []byte{common.ESC, '3', 40}
				if !bytes.Equal(result, expected) {
					t.Errorf("SetLineSpacing = %#v, want %#v", result, expected)
				}
				return nil
			},
		},
		{
			name: "print header",
			action: func() ([]byte, error) {
				return cmd.Print.Text("RECEIPT")
			},
			verify: func(result []byte) error {
				if !bytes.Equal(result, []byte("RECEIPT")) {
					t.Errorf("Text header = %#v, want RECEIPT", result)
				}
				return nil
			},
		},
		{
			name: "line feed",
			action: func() ([]byte, error) {
				return cmd.Print.PrintAndLineFeed(), nil
			},
			verify: func(result []byte) error {
				if !bytes.Equal(result, []byte{print.LF}) {
					t.Errorf("LineFeed = %#v, want LF", result)
				}
				return nil
			},
		},
		{
			name: "print item",
			action: func() ([]byte, error) {
				return cmd.Print.Text("Item 1: $10.00")
			},
			verify: func(result []byte) error {
				if !bytes.Contains(result, []byte("Item 1")) {
					t.Error("Should contain item text")
				}
				return nil
			},
		},
		{
			name: "form feed",
			action: func() ([]byte, error) {
				return cmd.Print.FormFeed(), nil
			},
			verify: func(result []byte) error {
				if !bytes.Equal(result, []byte{print.FF}) {
					t.Errorf("FormFeed = %#v, want FF", result)
				}
				return nil
			},
		},
	}

	// Execute and verify each command
	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.action()
			if err != nil {
				t.Fatalf("%s unexpected error: %v", tc.name, err)
			}
			if tc.verify != nil {
				_ = tc.verify(result)
			}
		})
	}
}

// TestIntegration_CustomCapabilities tests using custom capability implementations
func TestIntegration_CustomCapabilities(t *testing.T) {
	// Create a custom Commands with specific implementations
	customCmd := &escpos.Protocol{
		Print:     &print.Commands{},
		LineSpace: &linespacing.Commands{},
	}

	// Test that custom configuration works
	t.Run("custom print capability", func(t *testing.T) {
		text, err := customCmd.Print.Text("Custom")
		if err != nil {
			t.Fatalf("Text() error: %v", err)
		}
		if !bytes.Equal(text, []byte("Custom")) {
			t.Errorf("Text() = %#v, want 'Custom'", text)
		}
	})

	t.Run("custom line spacing capability", func(t *testing.T) {
		spacing := customCmd.LineSpace.SetLineSpacing(25)
		expected := []byte{common.ESC, '3', 25}
		if !bytes.Equal(spacing, expected) {
			t.Errorf("SetLineSpacing() = %#v, want %#v", spacing, expected)
		}
	})
}
