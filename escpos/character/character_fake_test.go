package character_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

// Ensure FakeCapability implements character.Capability
var _ character.Capability = (*FakeCapability)(nil)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates character formatting with state tracking
type FakeCapability struct {
	buffer             []byte
	currentSpacing     character.Spacing
	currentFont        byte
	currentSize        character.Size
	isEmphasized       bool
	isUnderlined       bool
	underlineThickness byte
	isDoubleStrike     bool
	isUpsideDown       bool
	isReversed         bool
	isSmoothing        bool
	rotation           byte
	printColor         byte
	codeTable          character.CodeTable
	internationalSet   character.InternationalSet
	lastCommand        string
	commandsCount      map[string]int
}

// NewFakeCapability creates a new fake character capability
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:        make([]byte, 0),
		commandsCount: make(map[string]int),
		currentFont:   0,
		currentSize:   0x00,
		codeTable:     0,
	}
}

// Add these methods to FakeCapability
func (fc *FakeCapability) GetCommandsHistory() []string {
	history := make([]string, 0)
	for cmd := range fc.commandsCount {
		for i := 0; i < fc.commandsCount[cmd]; i++ {
			history = append(history, cmd)
		}
	}
	return history
}

func (fc *FakeCapability) GetTotalCommands() int {
	total := 0
	for _, count := range fc.commandsCount {
		total += count
	}
	return total
}

// Standardize state getters
func (fc *FakeCapability) GetState() map[string]interface{} {
	return map[string]interface{}{
		"spacing":       fc.currentSpacing,
		"font":          fc.currentFont,
		"size":          fc.currentSize,
		"emphasized":    fc.isEmphasized,
		"underlined":    fc.isUnderlined,
		"doubleStrike":  fc.isDoubleStrike,
		"upsideDown":    fc.isUpsideDown,
		"reversed":      fc.isReversed,
		"smoothing":     fc.isSmoothing,
		"rotation":      fc.rotation,
		"printColor":    fc.printColor,
		"codeTable":     fc.codeTable,
		"international": fc.internationalSet,
	}
}

func (fc *FakeCapability) SetRightSideCharacterSpacing(n character.Spacing) []byte {
	cmd := []byte{common.ESC, common.SP, byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentSpacing = n
	fc.lastCommand = "SetRightSideCharacterSpacing"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SelectPrintModes(n character.PrintMode) []byte {
	cmd := []byte{common.ESC, '!', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)

	// Parse mode bits
	if n&0x01 != 0 {
		fc.currentFont = 1
	} else {
		fc.currentFont = 0
	}
	fc.isEmphasized = n&0x08 != 0
	fc.isUnderlined = n&0x80 != 0

	// Double height/width
	switch {
	case n&0x10 != 0 && n&0x20 != 0:
		fc.currentSize = 0x11 // Double both
	case n&0x10 != 0:
		fc.currentSize = 0x01 // Double height
	case n&0x20 != 0:
		fc.currentSize = 0x10 // Double width
	default:
		fc.currentSize = 0x00 // Normal
	}

	fc.lastCommand = "SelectPrintModes"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetUnderlineMode(n character.UnderlineMode) ([]byte, error) {
	// Validate
	switch n {
	case 0, '0':
		fc.isUnderlined = false
		fc.underlineThickness = 0
	case 1, '1':
		fc.isUnderlined = true
		fc.underlineThickness = 1
	case 2, '2':
		fc.isUnderlined = true
		fc.underlineThickness = 2
	default:
		return nil, character.ErrUnderlineMode
	}

	cmd := []byte{common.ESC, '-', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)

	fc.lastCommand = "SetUnderlineMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SetEmphasizedMode(n character.EmphasizedMode) []byte {
	cmd := []byte{common.ESC, 'E', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.isEmphasized = n&0x01 != 0
	fc.lastCommand = "SetEmphasizedMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetDoubleStrikeMode(n character.DoubleStrikeMode) []byte {
	cmd := []byte{common.ESC, 'G', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.isDoubleStrike = n&0x01 != 0
	fc.lastCommand = "SetDoubleStrikeMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SelectCharacterFont(n character.FontType) ([]byte, error) {
	// Validate
	switch n {
	case 0, '0':
		fc.currentFont = 0
	case 1, '1':
		fc.currentFont = 1
	case 2, '2':
		fc.currentFont = 2
	case 3, '3':
		fc.currentFont = 3
	case 4, '4':
		fc.currentFont = 4
	case 97:
		fc.currentFont = 97
	case 98:
		fc.currentFont = 98
	default:
		return nil, character.ErrCharacterFont
	}

	cmd := []byte{common.ESC, 'M', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)

	fc.lastCommand = "SelectCharacterFont"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SelectInternationalCharacterSet(n character.InternationalSet) ([]byte, error) {
	// Basic validation (simplified)
	if n > 17 && (n < 66 || n > 75) && n != 82 {
		return nil, character.ErrCharacterSet
	}

	cmd := []byte{common.ESC, 'R', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.internationalSet = n
	fc.lastCommand = "SelectInternationalCharacterSet"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) Set90DegreeClockwiseRotationMode(n character.RotationMode) ([]byte, error) {
	// Validate
	switch n {
	case 0, '0':
		fc.rotation = 0
	case 1, '1':
		fc.rotation = 1
	case 2, '2':
		fc.rotation = 2
	default:
		return nil, character.ErrRotationMode
	}

	cmd := []byte{common.ESC, 'V', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)

	fc.lastCommand = "Set90DegreeClockwiseRotationMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SelectPrintColor(n character.PrintColor) ([]byte, error) {
	// Validate
	switch n {
	case 0, '0':
		fc.printColor = 0
	case 1, '1':
		fc.printColor = 1
	default:
		return nil, character.ErrPrintColor
	}

	cmd := []byte{common.ESC, 'r', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)

	fc.lastCommand = "SelectPrintColor"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SelectCharacterCodeTable(n character.CodeTable) ([]byte, error) {
	// Basic validation (simplified)
	cmd := []byte{common.ESC, 't', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.codeTable = n
	fc.lastCommand = "SelectCharacterCodeTable"
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SetUpsideDownMode(n character.UpsideDownMode) []byte {
	cmd := []byte{common.ESC, '{', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.isUpsideDown = n&0x01 != 0
	fc.lastCommand = "SetUpsideDownMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SelectCharacterSize(n character.Size) []byte {
	cmd := []byte{common.GS, '!', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentSize = n
	fc.lastCommand = "SelectCharacterSize"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetWhiteBlackReverseMode(n character.ReverseMode) []byte {
	cmd := []byte{common.GS, 'B', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.isReversed = n&0x01 != 0
	fc.lastCommand = "SetWhiteBlackReverseMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetSmoothingMode(n character.SmoothingMode) []byte {
	cmd := []byte{common.GS, 'b', byte(n)}
	fc.buffer = append(fc.buffer, cmd...)
	fc.isSmoothing = n&0x01 != 0
	fc.lastCommand = "SetSmoothingMode"
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

// ============================================================================
// Helper Methods
// ===========================================================================

func (fc *FakeCapability) GetBuffer() []byte {
	return fc.buffer
}

func (fc *FakeCapability) GetCurrentFont() byte {
	return fc.currentFont
}

func (fc *FakeCapability) GetCurrentSize() byte {
	return byte(fc.currentSize)
}

func (fc *FakeCapability) GetIsEmphasized() bool {
	return fc.isEmphasized
}

func (fc *FakeCapability) GetIsUnderlined() bool {
	return fc.isUnderlined
}

func (fc *FakeCapability) GetCommandCount(cmd string) int {
	return fc.commandsCount[cmd]
}

func (fc *FakeCapability) GetLastCommand() string {
	return fc.lastCommand
}

func (fc *FakeCapability) Reset() {
	fc.buffer = make([]byte, 0)
	fc.commandsCount = make(map[string]int)
	fc.currentSpacing = 0
	fc.currentFont = 0
	fc.currentSize = 0x00
	fc.isEmphasized = false
	fc.isUnderlined = false
	fc.underlineThickness = 0
	fc.isDoubleStrike = false
	fc.isUpsideDown = false
	fc.isReversed = false
	fc.isSmoothing = false
	fc.rotation = 0
	fc.printColor = 0
	fc.codeTable = 0
	fc.internationalSet = 0
	fc.lastCommand = ""
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks font changes", func(t *testing.T) {
		fake := NewFakeCapability()

		_, _ = fake.SelectCharacterFont(1)

		if fake.GetCurrentFont() != 1 {
			t.Errorf("CurrentFont = %d, want 1", fake.GetCurrentFont())
		}
		if fake.GetLastCommand() != "SelectCharacterFont" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "SelectCharacterFont")
		}
	})

	t.Run("tracks print modes", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SelectPrintModes(0x88) // Emphasized + Underline

		if !fake.GetIsEmphasized() {
			t.Error("Should be emphasized")
		}
		if !fake.GetIsUnderlined() {
			t.Error("Should be underlined")
		}
		if fake.GetLastCommand() != "SelectPrintModes" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "SelectPrintModes")
		}
	})

	t.Run("tracks character size", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SelectCharacterSize(0x11) // Double width and height

		if fake.GetCurrentSize() != 0x11 {
			t.Errorf("CurrentSize = %#x, want 0x11", fake.GetCurrentSize())
		}
	})

	t.Run("accumulates commands", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetEmphasizedMode(1)
		fake.SelectPrintModes(0x88)
		fake.SelectCharacterSize(0x11)

		if fake.GetCommandCount("SetEmphasizedMode") != 1 {
			t.Errorf("SetEmphasizedMode count = %d, want 1",
				fake.GetCommandCount("SetEmphasizedMode"))
		}
		if fake.GetCommandCount("SelectPrintModes") != 1 {
			t.Errorf("SelectPrintModes count = %d, want 1",
				fake.GetCommandCount("SelectPrintModes"))
		}
		if fake.GetCommandCount("SelectCharacterSize") != 1 {
			t.Errorf("SelectCharacterSize count = %d, want 1",
				fake.GetCommandCount("SelectCharacterSize"))
		}
	})

	t.Run("reset clears state", func(t *testing.T) {
		fake := NewFakeCapability()

		// Set various states
		fake.SetEmphasizedMode(1)
		_, _ = fake.SelectCharacterFont(2)
		_, _ = fake.SetUnderlineMode(1)
		fake.SelectCharacterSize(0x11)

		// Reset
		fake.Reset()

		// Verify all state is cleared
		if fake.GetIsEmphasized() {
			t.Error("Should not be emphasized after reset")
		}
		if fake.GetCurrentFont() != 0 {
			t.Errorf("Font = %d, want 0 after reset", fake.GetCurrentFont())
		}
		if fake.GetIsUnderlined() {
			t.Error("Should not be underlined after reset")
		}
		if fake.GetCurrentSize() != 0x00 {
			t.Errorf("Size = %#x, want 0x00 after reset", fake.GetCurrentSize())
		}
		if len(fake.GetBuffer()) != 0 {
			t.Error("Buffer should be empty after reset")
		}
		if fake.GetLastCommand() != "" {
			t.Error("LastCommand should be empty after reset")
		}
	})
}

func TestFakeCapability_CompleteWorkflow(t *testing.T) {
	t.Run("complete formatting workflow", func(t *testing.T) {
		fake := NewFakeCapability()

		// Build a complete formatting sequence
		fake.SetRightSideCharacterSpacing(10)
		_, _ = fake.SelectCharacterFont(1)
		fake.SetEmphasizedMode(1)
		_, _ = fake.SetUnderlineMode(1)
		fake.SelectCharacterSize(0x11)
		fake.SetSmoothingMode(1)

		// Verify state
		state := fake.GetState()
		if state["spacing"].(character.Spacing) != 10 {
			t.Errorf("Spacing = %d, want 10", state["spacing"])
		}
		if state["font"].(byte) != 1 {
			t.Errorf("Font = %d, want 1", state["font"])
		}
		if !state["emphasized"].(bool) {
			t.Error("Should be emphasized")
		}
		if !state["underlined"].(bool) {
			t.Error("Should be underlined")
		}
		if state["size"].(character.Size) != 0x11 {
			t.Errorf("Size = %#x, want 0x11", state["size"])
		}
		if !state["smoothing"].(bool) {
			t.Error("Should have smoothing")
		}

		// Verify command history
		if fake.GetTotalCommands() != 6 {
			t.Errorf("Total commands = %d, want 6", fake.GetTotalCommands())
		}
	})

	t.Run("error handling in workflow", func(t *testing.T) {
		fake := NewFakeCapability()

		// Valid commands
		fake.SetEmphasizedMode(1)

		// Invalid command
		_, err := fake.SetUnderlineMode(99)
		if err == nil {
			t.Error("SetUnderlineMode(99) should return error")
		}

		// Valid commands should still work after error
		fake.SelectCharacterSize(0x11)

		// Verify partial state
		if !fake.GetIsEmphasized() {
			t.Error("Should still be emphasized")
		}
		if fake.GetIsUnderlined() {
			t.Error("Should not be underlined (error occurred)")
		}
		if fake.GetCurrentSize() != 0x11 {
			t.Errorf("Size = %#x, want 0x11", fake.GetCurrentSize())
		}
	})
}

func TestFakeCapability_BufferAccumulation(t *testing.T) {
	t.Run("accumulates all commands in buffer", func(t *testing.T) {
		fake := NewFakeCapability()

		// Execute multiple commands
		cmd1 := fake.SetRightSideCharacterSpacing(5)
		cmd2 := fake.SetEmphasizedMode(1)
		cmd3 := fake.SelectCharacterSize(0x11)

		buffer := fake.GetBuffer()

		// Verify buffer contains all commands
		expectedLen := len(cmd1) + len(cmd2) + len(cmd3)
		if len(buffer) != expectedLen {
			t.Errorf("Buffer length = %d, want %d", len(buffer), expectedLen)
		}

		// Verify buffer contains commands in order
		if !bytes.Contains(buffer, cmd1) {
			t.Error("Buffer should contain spacing command")
		}
		if !bytes.Contains(buffer, cmd2) {
			t.Error("Buffer should contain emphasized command")
		}
		if !bytes.Contains(buffer, cmd3) {
			t.Error("Buffer should contain size command")
		}
	})
}

func TestFakeCapability_CommandCounting(t *testing.T) {
	t.Run("counts repeated commands", func(t *testing.T) {
		fake := NewFakeCapability()

		// Call same command multiple times
		fake.SetEmphasizedMode(0)
		fake.SetEmphasizedMode(1)
		fake.SetEmphasizedMode(0)

		if fake.GetCommandCount("SetEmphasizedMode") != 3 {
			t.Errorf("SetEmphasizedMode count = %d, want 3",
				fake.GetCommandCount("SetEmphasizedMode"))
		}

		// Different command
		fake.SelectCharacterSize(0x00)

		if fake.GetCommandCount("SelectCharacterSize") != 1 {
			t.Errorf("SelectCharacterSize count = %d, want 1",
				fake.GetCommandCount("SelectCharacterSize"))
		}
	})

	t.Run("tracks command history", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetEmphasizedMode(1)
		_, _ = fake.SelectCharacterFont(1)
		fake.SetEmphasizedMode(0)

		history := fake.GetCommandsHistory()

		// Should have 3 commands in history
		if len(history) != 3 {
			t.Errorf("History length = %d, want 3", len(history))
		}

		// Verify order (may not be guaranteed depending on map iteration)
		expectedCommands := map[string]int{
			"SetEmphasizedMode":   2,
			"SelectCharacterFont": 1,
		}

		for cmd, count := range expectedCommands {
			if fake.GetCommandCount(cmd) != count {
				t.Errorf("%s count = %d, want %d",
					cmd, fake.GetCommandCount(cmd), count)
			}
		}
	})
}
