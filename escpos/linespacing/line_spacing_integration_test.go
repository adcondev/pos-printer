package linespacing_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
)

func TestIntegration_LineSpacing_StandardWorkflow(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("complete line spacing workflow", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SelectDefaultLineSpacing()...)

		buffer = append(buffer, cmd.SetLineSpacing(60)...)
		// ... print header ...

		buffer = append(buffer, cmd.SetLineSpacing(30)...)
		// ... print body ...

		buffer = append(buffer, cmd.SetLineSpacing(80)...)
		// ... print footer ...

		buffer = append(buffer, cmd.SelectDefaultLineSpacing()...)

		// Verify
		if len(buffer) != 4+3*3 { // 2 defaults + 3 custom
			t.Errorf("Buffer length = %d, want 13", len(buffer))
		}

		expectedDefault := []byte{common.ESC, '2'}
		if !bytes.Equal(buffer[:2], expectedDefault) {
			t.Errorf("Should start with default command")
		}
	})

	t.Run("progressive spacing changes", func(t *testing.T) {
		// Setup
		var buffer []byte
		spacings := []linespacing.Spacing{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

		// Execute
		for _, spacing := range spacings {
			cmd := cmd.SetLineSpacing(spacing)
			buffer = append(buffer, cmd...)
		}

		// Verify
		for i, spacing := range spacings {
			offset := i * 3
			expected := []byte{common.ESC, '3', byte(spacing)}
			actual := buffer[offset : offset+3]
			if !bytes.Equal(actual, expected) {
				t.Errorf("Command %d = %#v, want %#v", i, actual, expected)
			}
		}
	})

	t.Run("alternating default and custom", func(t *testing.T) {
		// Setup
		var buffer []byte
		patterns := []struct {
			name      string
			isDefault bool
			spacing   linespacing.Spacing
		}{
			{"default", true, 0},
			{"tight", false, 20},
			{"default", true, 0},
			{"normal", false, 40},
			{"default", true, 0},
			{"wide", false, 80},
			{"default", true, 0},
		}

		// Execute
		for _, p := range patterns {
			if p.isDefault {
				buffer = append(buffer, cmd.SelectDefaultLineSpacing()...)
			} else {
				buffer = append(buffer, cmd.SetLineSpacing(p.spacing)...)
			}
		}

		// Verify
		defaultCount := bytes.Count(buffer, []byte{common.ESC, '2'})
		if defaultCount != 4 {
			t.Errorf("Default command count = %d, want 4", defaultCount)
		}
		customCount := bytes.Count(buffer, []byte{common.ESC, '3'})
		if customCount != 3 {
			t.Errorf("Custom command count = %d, want 3", customCount)
		}
	})
}

func TestIntegration_LineSpacing_EdgeCases(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("boundary values", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		minCmd := cmd.SetLineSpacing(linespacing.MinSpacing)
		buffer = append(buffer, minCmd...)

		maxCmd := cmd.SetLineSpacing(linespacing.MaxSpacing)
		buffer = append(buffer, maxCmd...)

		defaultCmd := cmd.SetLineSpacing(linespacing.NormalSpacing)
		buffer = append(buffer, defaultCmd...)

		// Verify
		expectedMin := []byte{common.ESC, '3', 0}
		if !bytes.Equal(minCmd, expectedMin) {
			t.Errorf("Min spacing = %#v, want %#v", minCmd, expectedMin)
		}
		expectedMax := []byte{common.ESC, '3', 255}
		if !bytes.Equal(maxCmd, expectedMax) {
			t.Errorf("Max spacing = %#v, want %#v", maxCmd, expectedMax)
		}
		if len(buffer) != 9 { // 3 commands × 3 bytes
			t.Errorf("Buffer length = %d, want 9", len(buffer))
		}
	})

	t.Run("rapid spacing changes", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		for i := 0; i < 50; i++ {
			spacing := linespacing.Spacing(i % 256)
			cmd := cmd.SetLineSpacing(spacing)
			buffer = append(buffer, cmd...)
		}

		// Verify
		if len(buffer) != 50*3 {
			t.Errorf("Buffer length = %d, want %d", len(buffer), 150)
		}
		testIdx := 25
		offset := testIdx * 3
		if buffer[offset] != common.ESC || buffer[offset+1] != '3' {
			t.Error("Command structure incorrect")
		}
	})

	t.Run("all possible values", func(t *testing.T) {
		// Setup
		testValues := []linespacing.Spacing{
			0, 1, 2, 5, 10, 15, 20, 25, 30, // Low values
			40, 50, 60, 70, 80, 90, 100, // Medium values
			120, 140, 160, 180, 200, // High values
			220, 240, 250, 252, 253, 254, 255, // Maximum values
		}

		// Execute
		for _, value := range testValues {
			cmd := cmd.SetLineSpacing(value)
			expected := []byte{common.ESC, '3', byte(value)}
			// Verify
			if !bytes.Equal(cmd, expected) {
				t.Errorf("SetLineSpacing(%d) = %#v, want %#v", value, cmd, expected)
			}
		}
	})
}

func TestIntegration_LineSpacing_ModeTransitions(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("standard to page mode considerations", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		standardSpacing := cmd.SetLineSpacing(40)
		buffer = append(buffer, standardSpacing...)

		pageModeSpacing := cmd.SetLineSpacing(50)
		buffer = append(buffer, pageModeSpacing...)

		returnSpacing := cmd.SetLineSpacing(40)
		buffer = append(buffer, returnSpacing...)

		// Verify
		if !bytes.Equal(standardSpacing, returnSpacing) {
			t.Error("Same spacing spacing should produce same command")
		}
		if bytes.Equal(standardSpacing, pageModeSpacing) {
			t.Error("Different spacing values should produce different commands")
		}
		if len(buffer) != 9 { // 3 commands × 3 bytes
			t.Errorf("Buffer length = %d, want 9", len(buffer))
		}
	})

	t.Run("motion unit independence", func(t *testing.T) {
		// Setup
		spacing := linespacing.Spacing(100)

		// Execute
		cmd1 := cmd.SetLineSpacing(spacing)
		cmd2 := cmd.SetLineSpacing(spacing)

		// Verify
		if !bytes.Equal(cmd1, cmd2) {
			t.Error("Same spacing spacing should be consistent")
		}
		if cmd1[2] != byte(spacing) {
			t.Errorf("Spacing spacing = %d, want %d", cmd1[2], spacing)
		}
	})
}

func TestIntegration_LineSpacing_RealWorldScenarios(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("receipt with variable spacing", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetLineSpacing(25)...)
		// ... print store info ...

		buffer = append(buffer, cmd.SetLineSpacing(40)...)
		// ... print items ...

		buffer = append(buffer, cmd.SetLineSpacing(60)...)
		// ... print separator ...

		buffer = append(buffer, cmd.SetLineSpacing(40)...)
		// ... print totals ...

		buffer = append(buffer, cmd.SetLineSpacing(80)...)
		// ... print footer ...

		buffer = append(buffer, cmd.SelectDefaultLineSpacing()...)

		// Verify
		if bytes.Count(buffer, []byte{common.ESC, '3'}) != 5 {
			t.Error("Should have 5 custom spacing commands")
		}
		if bytes.Count(buffer, []byte{common.ESC, '2'}) != 1 {
			t.Error("Should have 1 default spacing command")
		}
	})

	t.Run("table with consistent row spacing", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		tableSpacing := linespacing.Spacing(35)
		buffer = append(buffer, cmd.SetLineSpacing(tableSpacing)...)

		// Simulate printing 10 rows
		for i := 0; i < 10; i++ {
			// Row content would go here
			// Spacing remains consistent
		}

		buffer = append(buffer, cmd.SelectDefaultLineSpacing()...)

		// Verify
		expectedSpacing := []byte{common.ESC, '3', byte(tableSpacing)}
		if !bytes.Equal(buffer[:3], expectedSpacing) {
			t.Error("Table should use consistent spacing")
		}
	})

	t.Run("document with sections", func(t *testing.T) {
		// Setup
		var buffer []byte
		sections := []struct {
			name    string
			spacing linespacing.Spacing
		}{
			{"title", 70},
			{"subtitle", 50},
			{"body", 35},
			{"list", 30},
			{"footer", 60},
		}

		// Execute
		for _, section := range sections {
			cmd := cmd.SetLineSpacing(section.spacing)
			buffer = append(buffer, cmd...)
		}

		// Verify
		for i, section := range sections {
			offset := i * 3
			if buffer[offset+2] != byte(section.spacing) {
				t.Errorf("Section %s spacing = %d, want %d",
					section.name, buffer[offset+2], section.spacing)
			}
		}
	})

	t.Run("dense listing mode", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetLineSpacing(linespacing.MinSpacing)...)

		// Simulate printing many lines
		// ...

		buffer = append(buffer, cmd.SetLineSpacing(30)...)

		// Verify
		if buffer[2] != 0 {
			t.Error("Dense mode should use minimum spacing")
		}
	})
}

func TestIntegration_LineSpacing_CompatibilityPatterns(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("model-specific default values", func(t *testing.T) {
		// Setup
		modelDefaults := []linespacing.Spacing{30, 40, 50, 60, 70, 80}

		for _, defaultVal := range modelDefaults {
			t.Run(fmt.Sprintf("default_%d", defaultVal), func(t *testing.T) {
				// Execute
				setCmd := cmd.SetLineSpacing(defaultVal)

				defaultCmd := cmd.SelectDefaultLineSpacing()

				// Verify
				if bytes.Equal(setCmd, defaultCmd) {
					t.Error("Set and default commands should differ")
				}
				if len(setCmd) != 3 || len(defaultCmd) != 2 {
					t.Error("Command lengths incorrect")
				}
			})
		}
	})

	t.Run("maximum spacing limit handling", func(t *testing.T) {
		// Test that maximum physical limit (1016mm) is handled
		// Even though we can send 255, printer may limit internally

		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetLineSpacing(255)...)

		buffer = append(buffer, cmd.SetLineSpacing(200)...)

		// Verify
		if len(buffer) != 6 {
			t.Error("Commands should be generated for high values")
		}
	})
}
