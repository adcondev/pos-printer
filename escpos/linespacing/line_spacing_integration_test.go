package linespacing_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/linespacing"
)

func TestIntegration_LineSpacing_Workflow(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("spacing changes workflow", func(t *testing.T) {
		// Set custom spacing
		custom := cmd.SetLineSpacing(50)

		// Reset to default
		defaultSpacing := cmd.SelectDefaultLineSpacing()

		// Set another custom
		custom2 := cmd.SetLineSpacing(100)

		// Verify sequence
		var fullSequence []byte
		fullSequence = append(fullSequence, custom...)
		fullSequence = append(fullSequence, defaultSpacing...)
		fullSequence = append(fullSequence, custom2...)

		if len(fullSequence) != 8 { // 3 + 2 + 3 bytes
			t.Errorf("Expected 8 bytes, got %d", len(fullSequence))
		}
	})
}

func TestIntegration_LineSpacing_EdgeCases(t *testing.T) {
	cmd := linespacing.NewCommands()

	t.Run("minimum to maximum spacing", func(t *testing.T) {
		// Test extreme values
		minSpacing := cmd.SetLineSpacing(0)
		if len(minSpacing) != 3 {
			t.Errorf("SetLineSpacing(0) length = %d, want 3", len(minSpacing))
		}

		maxSpacing := cmd.SetLineSpacing(255)
		if len(maxSpacing) != 3 {
			t.Errorf("SetLineSpacing(255) length = %d, want 3", len(maxSpacing))
		}
	})

	t.Run("rapid spacing changes", func(t *testing.T) {
		var commands []byte

		// Simulate rapid changes
		for i := byte(10); i <= 50; i += 10 {
			cmd := cmd.SetLineSpacing(i)
			commands = append(commands, cmd...)
		}

		// Should have 5 commands * 3 bytes each
		if len(commands) != 15 {
			t.Errorf("Rapid changes resulted in %d bytes, want 15", len(commands))
		}
	})

	t.Run("alternating default and custom", func(t *testing.T) {
		custom1 := cmd.SetLineSpacing(20)
		default1 := cmd.SelectDefaultLineSpacing()
		custom2 := cmd.SetLineSpacing(40)
		default2 := cmd.SelectDefaultLineSpacing()

		// Verify each command is distinct
		if bytes.Equal(custom1, custom2) {
			t.Error("Different custom spacings should produce different commands")
		}
		if !bytes.Equal(default1, default2) {
			t.Error("Default commands should be identical")
		}
	})
}

func TestIntegration_LineSpacing_WithMockAndFake(t *testing.T) {
	t.Run("mock vs real implementation", func(t *testing.T) {
		original := linespacing.NewCommands()
		mock := NewMockCapability()

		realResult := original.SetLineSpacing(100)
		mockResult := mock.SetLineSpacing(100)

		if !bytes.Equal(realResult, mockResult) {
			t.Errorf("Mock and real differ: mock=%#v, real=%#v", mockResult, realResult)
		}
	})

	t.Run("fake state consistency", func(t *testing.T) {
		fake := NewFakeCapability()

		// Multiple operations
		fake.SetLineSpacing(25)
		fake.SetLineSpacing(50)
		fake.SelectDefaultLineSpacing()
		fake.SetLineSpacing(75)

		state := fake.GetState()
		if state["currentSpacing"].(byte) != 75 {
			t.Errorf("Final spacing = %d, want 75", state["currentSpacing"])
		}
		if state["timesChanged"].(int) != 4 {
			t.Errorf("Times changed = %d, want 4", state["timesChanged"])
		}

		history := fake.GetCommandHistory()
		if len(history) != 4 {
			t.Errorf("History length = %d, want 4", len(history))
		}
	})
}
