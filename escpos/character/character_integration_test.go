package character_test

import (
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
)

func TestIntegration_CharacterFormatting_Workflow(t *testing.T) {
	cmd := character.NewCommands()

	// Test a complete character formatting workflow
	workflow := []struct {
		name   string
		action func() ([]byte, error)
		verify func([]byte, error) error
	}{
		{
			name: "set spacing",
			action: func() ([]byte, error) {
				return cmd.SetRightSideCharacterSpacing(10), nil
			},
		},
		{
			name: "enable emphasis",
			action: func() ([]byte, error) {
				return cmd.SetEmphasizedMode(1), nil
			},
		},
		{
			name: "set underline",
			action: func() ([]byte, error) {
				return cmd.SetUnderlineMode(1)
			},
		},
		{
			name: "set font",
			action: func() ([]byte, error) {
				return cmd.SelectCharacterFont(1)
			},
		},
	}

	var allCommands []byte
	for _, w := range workflow {
		t.Run(w.name, func(t *testing.T) {
			result, err := w.action()
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", w.name, err)
			}
			allCommands = append(allCommands, result...)
		})
	}

	// Verify complete sequence
	if len(allCommands) == 0 {
		t.Error("No commands generated")
	}
}

func TestIntegration_CharacterEffects_Workflow(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("color and shadow effects", func(t *testing.T) {
		// Select character color
		colorCmd, err := cmd.Effects.SelectCharacterColor('1')
		if err != nil {
			t.Fatalf("SelectCharacterColor: %v", err)
		}

		// Set background
		bgCmd, err := cmd.Effects.SelectBackgroundColor('0')
		if err != nil {
			t.Fatalf("SelectBackgroundColor: %v", err)
		}

		// Enable shadow
		shadowCmd, err := cmd.Effects.SetCharacterShadowMode(1, '2')
		if err != nil {
			t.Fatalf("SetCharacterShadowMode: %v", err)
		}

		// Verify commands are valid
		var allCmds []byte
		allCmds = append(allCmds, colorCmd...)
		allCmds = append(allCmds, bgCmd...)
		allCmds = append(allCmds, shadowCmd...)

		if len(allCmds) == 0 {
			t.Error("No effect commands generated")
		}
	})
}
