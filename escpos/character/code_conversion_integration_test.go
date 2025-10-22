package character_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/shared"
)

func TestIntegration_CodeConversion_MultiLanguageSupport(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("UTF-8 with font priorities", func(t *testing.T) {
		var buffer []byte

		// Enable UTF-8
		utf8Cmd, err := cmd.CodeConversion.SelectCharacterEncodeSystem(character.UTF8)
		if err != nil {
			t.Fatalf("SelectCharacterEncodeSystem(UTF8): %v", err)
		}
		buffer = append(buffer, utf8Cmd...)

		// Set Chinese as primary font
		chinesePriority, err := cmd.CodeConversion.SetFontPriority(
			character.First,
			character.SimplifiedChineseMincho,
		)
		if err != nil {
			t.Fatalf("SetFontPriority(Chinese): %v", err)
		}
		buffer = append(buffer, chinesePriority...)

		// Set Japanese as secondary font
		japanesePriority, err := cmd.CodeConversion.SetFontPriority(
			character.Second,
			character.JapaneseGothic,
		)
		if err != nil {
			t.Fatalf("SetFontPriority(Japanese): %v", err)
		}
		buffer = append(buffer, japanesePriority...)

		// Verify commands were generated
		if !bytes.Contains(buffer, []byte{shared.FS, '(', 'C'}) {
			t.Error("Buffer should contain encoding commands")
		}

		if len(buffer) != 23 { // 7 + 8 + 8 bytes
			t.Errorf("Buffer length = %d, want 23", len(buffer))
		}
	})

	t.Run("encoding switch workflow", func(t *testing.T) {
		// Switch from 1-byte to UTF-8 and back
		oneByteCmd, _ := cmd.CodeConversion.SelectCharacterEncodeSystem(character.OneByte)
		utf8Cmd, _ := cmd.CodeConversion.SelectCharacterEncodeSystem(character.UTF8)

		if len(oneByteCmd) != 7 || len(utf8Cmd) != 7 {
			t.Error("Encoding commands should be 7 bytes each")
		}
	})
}
