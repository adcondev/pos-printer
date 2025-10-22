package character_test

import (
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
)

func TestIntegration_Effects_ColorCombinations(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("promotional text with all effects", func(t *testing.T) {
		var buffer []byte

		// Apply character color
		charColor, err := cmd.Effects.SelectCharacterColor(character.CharColor2)
		if err != nil {
			t.Fatalf("SelectCharacterColor: %v", err)
		}
		buffer = append(buffer, charColor...)

		// Apply background color
		bgColor, err := cmd.Effects.SelectBackgroundColor(character.BackgroundColor1)
		if err != nil {
			t.Fatalf("SelectBackgroundColor: %v", err)
		}
		buffer = append(buffer, bgColor...)

		// Enable shadow
		shadow, err := cmd.Effects.SetCharacterShadowMode(
			character.ShadowModeOnByte,
			character.ShadowColor3,
		)
		if err != nil {
			t.Fatalf("SetCharacterShadowMode: %v", err)
		}
		buffer = append(buffer, shadow...)

		// Combine with reverse mode
		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(character.OnRm)...)

		if len(buffer) != 25 { // 7 + 7 + 8 + 3 bytes
			t.Errorf("Buffer length = %d, want 25", len(buffer))
		}
	})

	t.Run("effect reset workflow", func(t *testing.T) {
		// Turn off all effects
		charCmd, _ := cmd.Effects.SelectCharacterColor(character.CharColorNone)
		bgCmd, _ := cmd.Effects.SelectBackgroundColor(character.BackgroundColorNone)
		shadowCmd, _ := cmd.Effects.SetCharacterShadowMode(
			character.ShadowModeOffByte,
			character.ShadowColorNone,
		)

		totalLen := len(charCmd) + len(bgCmd) + len(shadowCmd)
		if totalLen != 22 { // 7 + 7 + 8 bytes
			t.Errorf("Reset commands length = %d, want 22", totalLen)
		}
	})
}
