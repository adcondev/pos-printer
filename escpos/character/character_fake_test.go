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
	currentSpacing     byte
	currentFont        byte
	currentSize        byte
	isEmphasized       bool
	isUnderlined       bool
	underlineThickness byte
	isDoubleStrike     bool
	isUpsideDown       bool
	isReversed         bool
	isSmoothing        bool
	rotation           byte
	printColor         byte
	codeTable          byte
	internationalSet   byte
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
func (f *FakeCapability) GetCommandsHistory() []string {
	history := make([]string, 0)
	for cmd := range f.commandsCount {
		for i := 0; i < f.commandsCount[cmd]; i++ {
			history = append(history, cmd)
		}
	}
	return history
}

func (f *FakeCapability) GetTotalCommands() int {
	total := 0
	for _, count := range f.commandsCount {
		total += count
	}
	return total
}

// Standardize state getters
func (f *FakeCapability) GetState() map[string]interface{} {
	return map[string]interface{}{
		"spacing":       f.currentSpacing,
		"font":          f.currentFont,
		"size":          f.currentSize,
		"emphasized":    f.isEmphasized,
		"underlined":    f.isUnderlined,
		"doubleStrike":  f.isDoubleStrike,
		"upsideDown":    f.isUpsideDown,
		"reversed":      f.isReversed,
		"smoothing":     f.isSmoothing,
		"rotation":      f.rotation,
		"printColor":    f.printColor,
		"codeTable":     f.codeTable,
		"international": f.internationalSet,
	}
}

func (f *FakeCapability) SetRightSideCharacterSpacing(n byte) []byte {
	cmd := []byte{common.ESC, common.SP, n}
	f.buffer = append(f.buffer, cmd...)
	f.currentSpacing = n
	f.lastCommand = "SetRightSideCharacterSpacing"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SelectPrintModes(n byte) []byte {
	cmd := []byte{common.ESC, '!', n}
	f.buffer = append(f.buffer, cmd...)

	// Parse mode bits
	if n&0x01 != 0 {
		f.currentFont = 1
	} else {
		f.currentFont = 0
	}
	f.isEmphasized = n&0x08 != 0
	f.isUnderlined = n&0x80 != 0

	// Double height/width
	switch {
	case n&0x10 != 0 && n&0x20 != 0:
		f.currentSize = 0x11 // Double both
	case n&0x10 != 0:
		f.currentSize = 0x01 // Double height
	case n&0x20 != 0:
		f.currentSize = 0x10 // Double width
	default:
		f.currentSize = 0x00 // Normal
	}

	f.lastCommand = "SelectPrintModes"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SetUnderlineMode(n byte) ([]byte, error) {
	// Validate
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid
	default:
		return nil, character.ErrInvalidUnderlineMode
	}

	cmd := []byte{common.ESC, '-', n}
	f.buffer = append(f.buffer, cmd...)

	switch n {
	case 0, '0':
		f.isUnderlined = false
		f.underlineThickness = 0
	case 1, '1':
		f.isUnderlined = true
		f.underlineThickness = 1
	case 2, '2':
		f.isUnderlined = true
		f.underlineThickness = 2
	}

	f.lastCommand = "SetUnderlineMode"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) SetEmphasizedMode(n byte) []byte {
	cmd := []byte{common.ESC, 'E', n}
	f.buffer = append(f.buffer, cmd...)
	f.isEmphasized = n&0x01 != 0
	f.lastCommand = "SetEmphasizedMode"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SetDoubleStrikeMode(n byte) []byte {
	cmd := []byte{common.ESC, 'G', n}
	f.buffer = append(f.buffer, cmd...)
	f.isDoubleStrike = n&0x01 != 0
	f.lastCommand = "SetDoubleStrikeMode"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SelectCharacterFont(n byte) ([]byte, error) {
	// Validate
	switch n {
	case 0, 1, 2, 3, 4, '0', '1', '2', '3', '4', 97, 98:
		// Valid
	default:
		return nil, character.ErrInvalidCharacterFont
	}

	cmd := []byte{common.ESC, 'M', n}
	f.buffer = append(f.buffer, cmd...)

	switch n {
	case 0, '0':
		f.currentFont = 0
	case 1, '1':
		f.currentFont = 1
	case 2, '2':
		f.currentFont = 2
	case 3, '3':
		f.currentFont = 3
	case 4, '4':
		f.currentFont = 4
	case 97:
		f.currentFont = 97
	case 98:
		f.currentFont = 98
	}

	f.lastCommand = "SelectCharacterFont"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) SelectInternationalCharacterSet(n byte) ([]byte, error) {
	// Basic validation (simplified)
	if n > 17 && (n < 66 || n > 75) && n != 82 {
		return nil, character.ErrInvalidCharacterSet
	}

	cmd := []byte{common.ESC, 'R', n}
	f.buffer = append(f.buffer, cmd...)
	f.internationalSet = n
	f.lastCommand = "SelectInternationalCharacterSet"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) Set90DegreeClockwiseRotationMode(n byte) ([]byte, error) {
	// Validate
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid
	default:
		return nil, character.ErrInvalidRotationMode
	}

	cmd := []byte{common.ESC, 'V', n}
	f.buffer = append(f.buffer, cmd...)

	switch n {
	case 0, '0':
		f.rotation = 0
	case 1, '1':
		f.rotation = 1
	case 2, '2':
		f.rotation = 2
	}

	f.lastCommand = "Set90DegreeClockwiseRotationMode"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) SelectPrintColor(n byte) ([]byte, error) {
	// Validate
	switch n {
	case 0, 1, '0', '1':
		// Valid
	default:
		return nil, character.ErrInvalidPrintColor
	}

	cmd := []byte{common.ESC, 'r', n}
	f.buffer = append(f.buffer, cmd...)

	switch n {
	case 0, '0':
		f.printColor = 0
	case 1, '1':
		f.printColor = 1
	}

	f.lastCommand = "SelectPrintColor"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) SelectCharacterCodeTable(n byte) ([]byte, error) {
	// Basic validation (simplified)
	cmd := []byte{common.ESC, 't', n}
	f.buffer = append(f.buffer, cmd...)
	f.codeTable = n
	f.lastCommand = "SelectCharacterCodeTable"
	f.commandsCount[f.lastCommand]++
	return cmd, nil
}

func (f *FakeCapability) SetUpsideDownMode(n byte) []byte {
	cmd := []byte{common.ESC, '{', n}
	f.buffer = append(f.buffer, cmd...)
	f.isUpsideDown = n&0x01 != 0
	f.lastCommand = "SetUpsideDownMode"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SelectCharacterSize(n byte) []byte {
	cmd := []byte{common.GS, '!', n}
	f.buffer = append(f.buffer, cmd...)
	f.currentSize = n
	f.lastCommand = "SelectCharacterSize"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SetWhiteBlackReverseMode(n byte) []byte {
	cmd := []byte{common.GS, 'B', n}
	f.buffer = append(f.buffer, cmd...)
	f.isReversed = n&0x01 != 0
	f.lastCommand = "SetWhiteBlackReverseMode"
	f.commandsCount[f.lastCommand]++
	return cmd
}

func (f *FakeCapability) SetSmoothingMode(n byte) []byte {
	cmd := []byte{common.GS, 'b', n}
	f.buffer = append(f.buffer, cmd...)
	f.isSmoothing = n&0x01 != 0
	f.lastCommand = "SetSmoothingMode"
	f.commandsCount[f.lastCommand]++
	return cmd
}

// Helper methods
func (f *FakeCapability) GetBuffer() []byte {
	return f.buffer
}

func (f *FakeCapability) GetCurrentFont() byte {
	return f.currentFont
}

func (f *FakeCapability) GetCurrentSize() byte {
	return f.currentSize
}

func (f *FakeCapability) IsEmphasized() bool {
	return f.isEmphasized
}

func (f *FakeCapability) IsUnderlined() bool {
	return f.isUnderlined
}

func (f *FakeCapability) GetCommandCount(cmd string) int {
	return f.commandsCount[cmd]
}

func (f *FakeCapability) GetLastCommand() string {
	return f.lastCommand
}

func (f *FakeCapability) Reset() {
	f.buffer = make([]byte, 0)
	f.commandsCount = make(map[string]int)
	f.currentSpacing = 0
	f.currentFont = 0
	f.currentSize = 0x00
	f.isEmphasized = false
	f.isUnderlined = false
	f.underlineThickness = 0
	f.isDoubleStrike = false
	f.isUpsideDown = false
	f.isReversed = false
	f.isSmoothing = false
	f.rotation = 0
	f.printColor = 0
	f.codeTable = 0
	f.internationalSet = 0
	f.lastCommand = ""
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

		if !fake.IsEmphasized() {
			t.Error("Should be emphasized")
		}
		if !fake.IsUnderlined() {
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
		if fake.IsEmphasized() {
			t.Error("Should not be emphasized after reset")
		}
		if fake.GetCurrentFont() != 0 {
			t.Errorf("Font = %d, want 0 after reset", fake.GetCurrentFont())
		}
		if fake.IsUnderlined() {
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
		if state["spacing"].(byte) != 10 {
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
		if state["size"].(byte) != 0x11 {
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
		if !fake.IsEmphasized() {
			t.Error("Should still be emphasized")
		}
		if fake.IsUnderlined() {
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
