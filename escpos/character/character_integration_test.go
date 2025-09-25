package character_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

func TestIntegration_Character_StandardWorkflow(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("complete character formatting workflow", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetRightSideCharacterSpacing(5)...)

		fontCmd, err := cmd.SelectCharacterFont(character.FontB)
		if err != nil {
			t.Fatalf("SelectCharacterFont: %v", err)
		}
		buffer = append(buffer, fontCmd...)

		buffer = append(buffer, cmd.SetEmphasizedMode(1)...)

		underlineCmd, err := cmd.SetUnderlineMode(character.OneDot)
		if err != nil {
			t.Fatalf("SetUnderlineMode: %v", err)
		}
		buffer = append(buffer, underlineCmd...)

		sizeConfig, err := character.NewSize(2, 2) // Double width and height
		if err != nil {
			t.Fatalf("NewSize: %v", err)
		}
		buffer = append(buffer, cmd.SelectCharacterSize(sizeConfig)...)

		charsetCmd, err := cmd.SelectInternationalCharacterSet(character.USA)
		if err != nil {
			t.Fatalf("SelectInternationalCharacterSet: %v", err)
		}
		buffer = append(buffer, charsetCmd...)

		codeTableCmd, err := cmd.SelectCharacterCodeTable(character.PC437)
		if err != nil {
			t.Fatalf("SelectCharacterCodeTable: %v", err)
		}
		buffer = append(buffer, codeTableCmd...)

		// Verify
		if len(buffer) == 0 {
			t.Error("Buffer should contain commands")
		}
		expectedSpacing := []byte{common.ESC, common.SP, 5}
		if !bytes.Equal(buffer[:3], expectedSpacing) {
			t.Errorf("Buffer should start with spacing command")
		}
	})

	t.Run("all print modes combination", func(t *testing.T) {
		// Setup
		var buffer []byte
		modes := []struct {
			name string
			bits character.PrintMode
		}{
			{"normal", 0x00},
			{"emphasized", character.EmphasizedOnPm},
			{"double height", character.DoubleHeightOnPm},
			{"double width", character.DoubleWidthOnPm},
			{"underline", character.UnderlineOnPm},
			{"all effects", character.EmphasizedOnPm |
				character.DoubleHeightOnPm |
				character.DoubleWidthOnPm |
				character.UnderlineOnPm},
		}

		// Execute
		for _, m := range modes {
			cmd := cmd.SelectPrintModes(m.bits)
			buffer = append(buffer, cmd...)
		}

		// Verify
		expectedCmdCount := 6 * 3 // 6 modes × 3 bytes per command
		if len(buffer) != expectedCmdCount {
			t.Errorf("Buffer length = %d, want %d", len(buffer), expectedCmdCount)
		}
	})

	t.Run("character transformation workflow", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		rotationCmd, err := cmd.Set90DegreeClockwiseRotationMode(character.On90Dot1)
		if err != nil {
			t.Fatalf("Set90DegreeClockwiseRotationMode: %v", err)
		}
		buffer = append(buffer, rotationCmd...)

		buffer = append(buffer, cmd.SetUpsideDownMode(character.OnUdm)...)

		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(character.OnRm)...)

		buffer = append(buffer, cmd.SetSmoothingMode(character.OnSm)...)

		// Verify
		if len(buffer) < 4*3 {
			t.Error("Buffer should contain transformation commands")
		}
	})
}

func TestIntegration_Character_EffectsWorkflow(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("complete color and shadow effects", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		charColorCmd, err := cmd.Effects.SelectCharacterColor(character.CharColor1)
		if err != nil {
			t.Fatalf("SelectCharacterColor: %v", err)
		}
		buffer = append(buffer, charColorCmd...)

		bgColorCmd, err := cmd.Effects.SelectBackgroundColor(character.BackgroundColor2)
		if err != nil {
			t.Fatalf("SelectBackgroundColor: %v", err)
		}
		buffer = append(buffer, bgColorCmd...)

		shadowCmd, err := cmd.Effects.SetCharacterShadowMode(
			character.ShadowModeOnByte,
			character.ShadowColor3,
		)
		if err != nil {
			t.Fatalf("SetCharacterShadowMode: %v", err)
		}
		buffer = append(buffer, shadowCmd...)

		// Verify
		expectedCharColor := []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, '1'}
		if !bytes.Contains(buffer, expectedCharColor) {
			t.Error("Buffer should contain character color command")
		}
	})

	t.Run("all color combinations", func(t *testing.T) {
		// Setup
		colorTests := []struct {
			name      string
			charColor byte
			bgColor   byte
			shadow    byte
		}{
			{"no colors", character.CharColorNone, character.BackgroundColorNone, character.ShadowColorNone},
			{"char only", character.CharColor1, character.BackgroundColorNone, character.ShadowColorNone},
			{"bg only", character.CharColorNone, character.BackgroundColor1, character.ShadowColorNone},
			{"shadow only", character.CharColorNone, character.BackgroundColorNone, character.ShadowColor1},
			{"all colors", character.CharColor1, character.BackgroundColor2, character.ShadowColor3},
		}

		for _, tc := range colorTests {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				var buffer []byte

				// Execute
				charCmd, err := cmd.Effects.SelectCharacterColor(tc.charColor)
				if err != nil {
					t.Fatalf("SelectCharacterColor: %v", err)
				}
				buffer = append(buffer, charCmd...)

				bgCmd, err := cmd.Effects.SelectBackgroundColor(tc.bgColor)
				if err != nil {
					t.Fatalf("SelectBackgroundColor: %v", err)
				}
				buffer = append(buffer, bgCmd...)

				shadowCmd, err := cmd.Effects.SetCharacterShadowMode(
					character.ShadowModeOnByte,
					tc.shadow,
				)
				if err != nil {
					t.Fatalf("SetCharacterShadowMode: %v", err)
				}
				buffer = append(buffer, shadowCmd...)

				// Verify
				if len(buffer) == 0 {
					t.Error("Buffer should contain color commands")
				}
			})
		}
	})
}

func TestIntegration_Character_CodeConversionWorkflow(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("encoding system selection", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		utf8Cmd, err := cmd.CodeConversion.SelectCharacterEncodeSystem(character.UTF8)
		if err != nil {
			t.Fatalf("SelectCharacterEncodeSystem(UTF8): %v", err)
		}
		buffer = append(buffer, utf8Cmd...)

		oneByteCmd, err := cmd.CodeConversion.SelectCharacterEncodeSystem(character.OneByte)
		if err != nil {
			t.Fatalf("SelectCharacterEncodeSystem(1Byte): %v", err)
		}
		buffer = append(buffer, oneByteCmd...)

		utf8ASCIICmd, err := cmd.CodeConversion.SelectCharacterEncodeSystem(character.UTF8Ascii)
		if err != nil {
			t.Fatalf("SelectCharacterEncodeSystem(UTF8 ASCII): %v", err)
		}
		buffer = append(buffer, utf8ASCIICmd...)

		// Verify
		expectedUTF8 := []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, 2}
		if !bytes.Contains(buffer, expectedUTF8) {
			t.Error("Buffer should contain UTF-8 encoding command")
		}
	})

	t.Run("font priority management", func(t *testing.T) {
		// Setup
		priorities := []struct {
			name     string
			priority character.FontPriority
			function character.FontFunction
		}{
			{"AnkSansSerif first", character.First, character.AnkSansSerif},
			{"Japanese second", character.Second, character.JapaneseGothic},
			{"Chinese first", character.First, character.SimplifiedChineseMincho},
			{"Korean second", character.Second, character.KoreanGothic},
		}

		for _, p := range priorities {
			t.Run(p.name, func(t *testing.T) {
				// Execute
				cmd, err := cmd.CodeConversion.SetFontPriority(p.priority, p.function)
				if err != nil {
					t.Fatalf("SetFontPriority(%s): %v", p.name, err)
				}

				// Verify
				expected := []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, byte(p.priority), byte(p.function)}
				if !bytes.Equal(cmd, expected) {
					t.Errorf("Command = %#v, want %#v", cmd, expected)
				}
			})
		}
	})
}

func TestIntegration_Character_UserDefinedWorkflow(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("complete user-defined character workflow", func(t *testing.T) {
		// Setup
		var buffer []byte
		customChars := []character.UserDefinedChar{
			{
				Width: 12,
				Data:  bytes.Repeat([]byte{0xFF}, 36), // 12 width × 3 height
			},
			{
				Width: 8,
				Data:  bytes.Repeat([]byte{0xAA}, 24), // 8 width × 3 height
			},
		}

		// Execute
		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(character.UserDefinedOn)...)

		defineCmd, err := cmd.UserDefined.DefineUserDefinedCharacters(3, 64, 65, customChars)
		if err != nil {
			t.Fatalf("DefineUserDefinedCharacters: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		cancelCmd, err := cmd.UserDefined.CancelUserDefinedCharacter(64)
		if err != nil {
			t.Fatalf("CancelUserDefinedCharacter: %v", err)
		}
		buffer = append(buffer, cancelCmd...)

		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(character.UserDefinedOff)...)

		// Verify
		if len(buffer) == 0 {
			t.Error("Buffer should contain user-defined commands")
		}
	})

	t.Run("user-defined character edge cases", func(t *testing.T) {
		// Execute
		maxCmd, err := cmd.UserDefined.CancelUserDefinedCharacter(126)
		if err != nil {
			t.Fatalf("CancelUserDefinedCharacter(126): %v", err)
		}
		// Minimum character code
		minCmd, err := cmd.UserDefined.CancelUserDefinedCharacter(32)
		if err != nil {
			t.Fatalf("CancelUserDefinedCharacter(32): %v", err)
		}

		// Verify
		if len(maxCmd) != 3 {
			t.Errorf("Cancel max char length = %d, want 3", len(maxCmd))
		}
		if len(minCmd) != 3 {
			t.Errorf("Cancel min char length = %d, want 3", len(minCmd))
		}
	})
}

func TestIntegration_Character_EdgeCases(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("maximum values", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetRightSideCharacterSpacing(255)...)

		maxSize, err := character.NewSize(8, 8)
		if err != nil {
			t.Fatalf("NewSize(8,8): %v", err)
		}
		buffer = append(buffer, cmd.SelectCharacterSize(maxSize)...)

		maxTableCmd, err := cmd.SelectCharacterCodeTable(255)
		if err != nil {
			t.Fatalf("SelectCharacterCodeTable(255): %v", err)
		}
		buffer = append(buffer, maxTableCmd...)

		// Verify
		if len(buffer) == 0 {
			t.Error("Buffer should contain maximum value commands")
		}
	})

	t.Run("minimum values", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		buffer = append(buffer, cmd.SetRightSideCharacterSpacing(0)...)

		minSize, err := character.NewSize(1, 1)
		if err != nil {
			t.Fatalf("NewSize(1,1): %v", err)
		}
		buffer = append(buffer, cmd.SelectCharacterSize(minSize)...)

		minCharsetCmd, err := cmd.SelectInternationalCharacterSet(0)
		if err != nil {
			t.Fatalf("SelectInternationalCharacterSet(0): %v", err)
		}
		buffer = append(buffer, minCharsetCmd...)

		// Verify
		if len(buffer) == 0 {
			t.Error("Buffer should contain minimum value commands")
		}
	})

	t.Run("all underline modes", func(t *testing.T) {
		// Setup
		modes := []struct {
			name string
			mode character.UnderlineMode
		}{
			{"off", character.NoDot},
			{"1 dot", character.OneDot},
			{"2 dot", character.TwoDot},
			{"off ASCII", character.NoDotAscii},
			{"1 dot ASCII", character.OneDotAscii},
			{"2 dot ASCII", character.TwoDotAscii},
		}

		for _, m := range modes {
			t.Run(m.name, func(t *testing.T) {
				// Execute
				cmd, err := cmd.SetUnderlineMode(m.mode)
				if err != nil {
					t.Fatalf("SetUnderlineMode(%s): %v", m.name, err)
				}

				// Verify
				expected := []byte{common.ESC, '-', byte(m.mode)}
				if !bytes.Equal(cmd, expected) {
					t.Errorf("Command = %#v, want %#v", cmd, expected)
				}
			})
		}
	})
}

func TestIntegration_Character_ErrorConditions(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("invalid parameters", func(t *testing.T) {
		// Invalid underline mode
		_, err := cmd.SetUnderlineMode(99)
		if err == nil {
			t.Error("SetUnderlineMode(99) should return error")
		}

		// Invalid character font
		_, err = cmd.SelectCharacterFont(99)
		if err == nil {
			t.Error("SelectCharacterFont(99) should return error")
		}

		// Invalid international character set
		_, err = cmd.SelectInternationalCharacterSet(200)
		if err == nil {
			t.Error("SelectInternationalCharacterSet(200) should return error")
		}

		// Invalid rotation mode
		_, err = cmd.Set90DegreeClockwiseRotationMode(99)
		if err == nil {
			t.Error("Set90DegreeClockwiseRotationMode(99) should return error")
		}

		// Invalid print color
		_, err = cmd.SelectPrintColor(99)
		if err == nil {
			t.Error("SelectPrintColor(99) should return error")
		}

		// Invalid code table page
		_, err = cmd.SelectCharacterCodeTable(100)
		if err == nil {
			t.Error("SelectCharacterCodeTable(100) should return error")
		}

		// Invalid character size
		_, err = character.NewSize(0, 1)
		if err == nil {
			t.Error("NewSize(0,1) should return error")
		}

		_, err = character.NewSize(1, 9)
		if err == nil {
			t.Error("NewSize(1,9) should return error")
		}
	})

	t.Run("effects invalid parameters", func(t *testing.T) {
		// Invalid character color
		_, err := cmd.Effects.SelectCharacterColor('4')
		if err == nil {
			t.Error("SelectCharacterColor('4') should return error")
		}

		// Invalid background color
		_, err = cmd.Effects.SelectBackgroundColor('5')
		if err == nil {
			t.Error("SelectBackgroundColor('5') should return error")
		}

		// Invalid shadow mode
		_, err = cmd.Effects.SetCharacterShadowMode(99, character.ShadowColor1)
		if err == nil {
			t.Error("SetCharacterShadowMode(99) should return error")
		}

		// Invalid shadow color
		_, err = cmd.Effects.SetCharacterShadowMode(character.ShadowModeOnByte, '9')
		if err == nil {
			t.Error("SetCharacterShadowMode with invalid color should return error")
		}
	})

	t.Run("code conversion invalid parameters", func(t *testing.T) {
		// Invalid encoding
		_, err := cmd.CodeConversion.SelectCharacterEncodeSystem(99)
		if err == nil {
			t.Error("SelectCharacterEncodeSystem(99) should return error")
		}

		// Invalid font priority
		_, err = cmd.CodeConversion.SetFontPriority(2, character.AnkSansSerif)
		if err == nil {
			t.Error("SetFontPriority(2) should return error")
		}

		// Invalid font type
		_, err = cmd.CodeConversion.SetFontPriority(0, 99)
		if err == nil {
			t.Error("SetFontPriority with invalid font should return error")
		}
	})

	t.Run("user-defined invalid parameters", func(t *testing.T) {
		// Invalid character code (too low)
		_, err := cmd.UserDefined.CancelUserDefinedCharacter(31)
		if err == nil {
			t.Error("CancelUserDefinedCharacter(31) should return error")
		}

		// Invalid character code (too high)
		_, err = cmd.UserDefined.CancelUserDefinedCharacter(127)
		if err == nil {
			t.Error("CancelUserDefinedCharacter(127) should return error")
		}

		// Invalid y value
		_, err = cmd.UserDefined.DefineUserDefinedCharacters(0, 32, 33, nil)
		if err == nil {
			t.Error("DefineUserDefinedCharacters with y=0 should return error")
		}

		// Mismatched definition count
		def := []character.UserDefinedChar{{Width: 8, Data: []byte{0xFF}}}
		_, err = cmd.UserDefined.DefineUserDefinedCharacters(1, 32, 34, def)
		if err == nil {
			t.Error("DefineUserDefinedCharacters with wrong count should return error")
		}
	})
}

func TestIntegration_Character_RealWorldScenarios(t *testing.T) {
	cmd := character.NewCommands()

	t.Run("receipt header formatting", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		titleSize, _ := character.NewSize(2, 3)
		buffer = append(buffer, cmd.SelectCharacterSize(titleSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(1)...)
		// ... print title ...

		normalSize, _ := character.NewSize(1, 1)
		buffer = append(buffer, cmd.SelectCharacterSize(normalSize)...)
		buffer = append(buffer, cmd.SetEmphasizedMode(0)...)
		underlineCmd, _ := cmd.SetUnderlineMode(character.OneDot)
		buffer = append(buffer, underlineCmd...)
		// ... print subtitle ...

		resetUnderline, _ := cmd.SetUnderlineMode(character.NoDot)
		buffer = append(buffer, resetUnderline...)

		// Verify
		if len(buffer) < 6*3 {
			t.Error("Buffer should contain multiple formatting commands")
		}
	})

	t.Run("multi-language support", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		utf8Cmd, _ := cmd.CodeConversion.SelectCharacterEncodeSystem(character.UTF8)
		buffer = append(buffer, utf8Cmd...)

		chinesePriority, _ := cmd.CodeConversion.SetFontPriority(
			character.First,
			character.SimplifiedChineseMincho,
		)
		buffer = append(buffer, chinesePriority...)

		chinaCharset, _ := cmd.SelectInternationalCharacterSet(character.China)
		buffer = append(buffer, chinaCharset...)

		codeTable, _ := cmd.SelectCharacterCodeTable(character.PC437)
		buffer = append(buffer, codeTable...)

		// Verify
		if !bytes.Contains(buffer, []byte{common.FS, '(', 'C'}) {
			t.Error("Buffer should contain encoding commands")
		}
	})

	t.Run("promotional text with effects", func(t *testing.T) {
		// Setup
		var buffer []byte

		// Execute
		charColor, _ := cmd.Effects.SelectCharacterColor(character.CharColor2)
		buffer = append(buffer, charColor...)

		bgColor, _ := cmd.Effects.SelectBackgroundColor(character.BackgroundColor1)
		buffer = append(buffer, bgColor...)

		shadow, _ := cmd.Effects.SetCharacterShadowMode(
			character.ShadowModeOnByte,
			character.ShadowColor3,
		)
		buffer = append(buffer, shadow...)

		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(character.OnRm)...)

		promoSize, _ := character.NewSize(3, 2)
		buffer = append(buffer, cmd.SelectCharacterSize(promoSize)...)

		buffer = append(buffer, cmd.SetWhiteBlackReverseMode(character.OffRm)...)
		shadowOff, _ := cmd.Effects.SetCharacterShadowMode(
			character.ShadowModeOffByte,
			character.ShadowColorNone,
		)
		buffer = append(buffer, shadowOff...)

		// Verify
		if len(buffer) < 7 {
			t.Error("Buffer should contain multiple effect commands")
		}
	})

	t.Run("custom logo characters", func(t *testing.T) {
		// Setup
		var buffer []byte
		logoChars := make([]character.UserDefinedChar, 4)
		for i := range logoChars {
			logoChars[i] = character.UserDefinedChar{
				Width: 12,
				Data:  bytes.Repeat([]byte{byte(0x01 << i)}, 36), // Pattern for each part
			}
		}

		// Execute
		defineCmd, err := cmd.UserDefined.DefineUserDefinedCharacters(3, 64, 67, logoChars)
		if err != nil {
			t.Fatalf("DefineUserDefinedCharacters for logo: %v", err)
		}
		buffer = append(buffer, defineCmd...)

		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(character.UserDefinedOn)...)

		// ... print logo using characters 64-67 ...

		buffer = append(buffer, cmd.UserDefined.SelectUserDefinedCharacterSet(character.UserDefinedOff)...)

		// Verify
		if !bytes.Contains(buffer, []byte{0x1B, 0x26}) {
			t.Error("Buffer should contain define characters command")
		}
	})
}
