package printposition_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/printposition"
)

// Ensure FakeCapability implements printposition.Capability
var _ printposition.Capability = (*FakeCapability)(nil)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates print position with state tracking
type FakeCapability struct {
	buffer              []byte
	currentPosition     int
	currentVertPosition int
	leftMargin          int
	printAreaWidth      int
	printAreaX          int
	printAreaY          int
	printAreaHeight     int
	justification       byte
	printDirection      byte
	tabPositions        []int
	isPageMode          bool
	lastCommand         string
	commandHistory      []string
	commandsCount       map[string]int
}

// NewFakeCapability creates a new fake print position capability
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:         make([]byte, 0),
		commandHistory: make([]string, 0),
		commandsCount:  make(map[string]int),
		tabPositions:   []int{8, 16, 24, 32, 40, 48, 56, 64}, // Default tabs
		printAreaWidth: 576,                                  // Default 80mm
	}
}

func (fc *FakeCapability) GetCommandHistory() []string {
	return fc.commandHistory
}

func (fc *FakeCapability) GetState() map[string]interface{} {
	return map[string]interface{}{
		"currentPosition":     fc.currentPosition,
		"currentVertPosition": fc.currentVertPosition,
		"leftMargin":          fc.leftMargin,
		"printAreaWidth":      fc.printAreaWidth,
		"justification":       fc.justification,
		"printDirection":      fc.printDirection,
		"isPageMode":          fc.isPageMode,
		"tabCount":            len(fc.tabPositions),
	}
}

func (fc *FakeCapability) SetAbsolutePrintPosition(position uint16) []byte {
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	cmd := []byte{common.ESC, '$', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentPosition = int(position)
	fc.lastCommand = "SetAbsolutePrintPosition"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetRelativePrintPosition(distance int16) []byte {
	// Two's complement for negative values is needed
	value := uint16(distance) // nolint:gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	cmd := []byte{common.ESC, '\\', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentPosition += int(distance)
	if fc.currentPosition < 0 {
		fc.currentPosition = 0
	}
	fc.lastCommand = "SetRelativePrintPosition"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) HorizontalTab() []byte {
	cmd := []byte{printposition.HT}
	fc.buffer = append(fc.buffer, cmd...)

	// Move to next tab position
	for _, tab := range fc.tabPositions {
		if tab > fc.currentPosition {
			fc.currentPosition = tab
			break
		}
	}

	fc.lastCommand = "HorizontalTab"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetHorizontalTabPositions(positions []byte) ([]byte, error) {
	// Validate
	if len(positions) > printposition.MaxTabPositions {
		return nil, printposition.ErrTooManyTabPositions
	}

	prevPos := byte(0)
	for _, pos := range positions {
		if pos == 0 || pos > printposition.MaxTabValue {
			return nil, printposition.ErrTabPosition
		}
		if pos <= prevPos {
			return nil, printposition.ErrTabPosition
		}
		prevPos = pos
	}

	cmd := []byte{common.ESC, 'D'}
	cmd = append(cmd, positions...)
	cmd = append(cmd, common.NUL)
	fc.buffer = append(fc.buffer, cmd...)

	// Update tab positions
	fc.tabPositions = make([]int, len(positions))
	for i, pos := range positions {
		fc.tabPositions[i] = int(pos)
	}

	fc.lastCommand = "SetHorizontalTabPositions"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SelectJustification(mode byte) ([]byte, error) {
	// Validate
	switch mode {
	case 0, 1, 2, '0', '1', '2':
		// Valid
	default:
		return nil, printposition.ErrJustification
	}

	cmd := []byte{common.ESC, 'a', mode}
	fc.buffer = append(fc.buffer, cmd...)
	fc.justification = mode
	fc.lastCommand = "SelectJustification"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SetLeftMargin(margin uint16) []byte {
	nL := byte(margin & 0xFF)
	nH := byte((margin >> 8) & 0xFF)
	cmd := []byte{common.GS, 'L', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.leftMargin = int(margin)
	fc.lastCommand = "SetLeftMargin"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetPrintAreaWidth(width uint16) []byte {
	nL := byte(width & 0xFF)
	nH := byte((width >> 8) & 0xFF)
	cmd := []byte{common.GS, 'W', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.printAreaWidth = int(width)
	fc.lastCommand = "SetPrintAreaWidth"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetPrintPositionBeginningLine(mode byte) ([]byte, error) {
	// Validate
	switch mode {
	case 0, 1, '0', '1':
		// Valid
	default:
		return nil, printposition.ErrBeginLineMode
	}

	cmd := []byte{common.GS, 'T', mode}
	fc.buffer = append(fc.buffer, cmd...)

	// Reset position to beginning
	fc.currentPosition = fc.leftMargin

	fc.lastCommand = "SetPrintPositionBeginningLine"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SelectPrintDirectionPageMode(direction byte) ([]byte, error) {
	// Validate
	switch direction {
	case 0, 1, 2, 3, '0', '1', '2', '3':
		// Valid
	default:
		return nil, printposition.ErrPrintDirection
	}

	cmd := []byte{common.ESC, 'T', direction}
	fc.buffer = append(fc.buffer, cmd...)
	fc.printDirection = direction
	fc.lastCommand = "SelectPrintDirectionPageMode"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SetPrintAreaPageMode(x, y, width, height uint16) ([]byte, error) {
	switch {
	case width == 0:
		return nil, printposition.ErrPrintAreaWidthSize
	case height == 0:
		return nil, printposition.ErrPrintAreaHeightSize
	}

	xL := byte(x & 0xFF)
	xH := byte((x >> 8) & 0xFF)
	yL := byte(y & 0xFF)
	yH := byte((y >> 8) & 0xFF)
	dxL := byte(width & 0xFF)
	dxH := byte((width >> 8) & 0xFF)
	dyL := byte(height & 0xFF)
	dyH := byte((height >> 8) & 0xFF)

	cmd := []byte{common.ESC, 'W', xL, xH, yL, yH, dxL, dxH, dyL, dyH}
	fc.buffer = append(fc.buffer, cmd...)

	fc.printAreaX = int(x)
	fc.printAreaY = int(y)
	fc.printAreaWidth = int(width)
	fc.printAreaHeight = int(height)
	fc.isPageMode = true

	fc.lastCommand = "SetPrintAreaPageMode"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd, nil
}

func (fc *FakeCapability) SetAbsoluteVerticalPrintPosition(position uint16) []byte {
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	cmd := []byte{common.GS, '$', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentVertPosition = int(position)
	fc.lastCommand = "SetAbsoluteVerticalPrintPosition"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

func (fc *FakeCapability) SetRelativeVerticalPrintPosition(distance int16) []byte {
	// Two's complement for negative values is needed
	value := uint16(distance) // nolint:gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	cmd := []byte{common.GS, '\\', nL, nH}
	fc.buffer = append(fc.buffer, cmd...)
	fc.currentVertPosition += int(distance)
	if fc.currentVertPosition < 0 {
		fc.currentVertPosition = 0
	}
	fc.lastCommand = "SetRelativeVerticalPrintPosition"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	fc.commandsCount[fc.lastCommand]++
	return cmd
}

// ============================================================================
// Helper Methods
// ============================================================================

func (fc *FakeCapability) GetBuffer() []byte {
	return fc.buffer
}

func (fc *FakeCapability) GetCurrentPosition() int {
	return fc.currentPosition
}

func (fc *FakeCapability) GetCurrentVertPosition() int {
	return fc.currentVertPosition
}

func (fc *FakeCapability) GetJustification() byte {
	return fc.justification
}

func (fc *FakeCapability) GetTabPositions() []int {
	return fc.tabPositions
}

func (fc *FakeCapability) GetCommandCount(cmd string) int {
	return fc.commandsCount[cmd]
}

func (fc *FakeCapability) GetLastCommand() string {
	return fc.lastCommand
}

func (fc *FakeCapability) Reset() {
	*fc = *NewFakeCapability()
	// TODO: Check if these defaults are appropriate or should be configurable
	fc.printAreaWidth = 576
	fc.tabPositions = []int{8, 16, 24, 32, 40, 48, 56, 64}
	fc.isPageMode = false
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks absolute position", func(t *testing.T) {
		fake := NewFakeCapability()

		result := fake.SetAbsolutePrintPosition(300)

		expected := []byte{common.ESC, '$', 0x2C, 0x01}
		if !bytes.Equal(result, expected) {
			t.Errorf("SetAbsolutePrintPosition(300) = %#v, want %#v", result, expected)
		}
		if fake.GetCurrentPosition() != 300 {
			t.Errorf("CurrentPosition = %d, want 300", fake.GetCurrentPosition())
		}
		if fake.GetLastCommand() != "SetAbsolutePrintPosition" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "SetAbsolutePrintPosition")
		}
	})

	t.Run("tracks relative position", func(t *testing.T) {
		fake := NewFakeCapability()
		fake.SetAbsolutePrintPosition(100)

		fake.SetRelativePrintPosition(50)
		if fake.GetCurrentPosition() != 150 {
			t.Errorf("CurrentPosition = %d, want 150", fake.GetCurrentPosition())
		}

		fake.SetRelativePrintPosition(-75)
		if fake.GetCurrentPosition() != 75 {
			t.Errorf("CurrentPosition = %d, want 75", fake.GetCurrentPosition())
		}
	})

	t.Run("tracks justification", func(t *testing.T) {
		fake := NewFakeCapability()

		_, err := fake.SelectJustification(1)
		if err != nil {
			t.Fatalf("SelectJustification(1) unexpected error: %v", err)
		}

		if fake.GetJustification() != 1 {
			t.Errorf("Justification = %d, want 1", fake.GetJustification())
		}
	})

	t.Run("tracks tab positions", func(t *testing.T) {
		fake := NewFakeCapability()

		newTabs := []byte{10, 20, 30, 40}
		_, err := fake.SetHorizontalTabPositions(newTabs)
		if err != nil {
			t.Fatalf("SetHorizontalTabPositions() unexpected error: %v", err)
		}

		tabs := fake.GetTabPositions()
		if len(tabs) != 4 {
			t.Errorf("Tab count = %d, want 4", len(tabs))
		}
		if tabs[0] != 10 || tabs[3] != 40 {
			t.Errorf("Tabs = %v, want [10 20 30 40]", tabs)
		}
	})

	t.Run("horizontal tab moves to next position", func(t *testing.T) {
		fake := NewFakeCapability()
		_, _ = fake.SetHorizontalTabPositions([]byte{10, 20, 30})
		fake.SetAbsolutePrintPosition(5)

		fake.HorizontalTab()
		if fake.GetCurrentPosition() != 10 {
			t.Errorf("Position after tab = %d, want 10", fake.GetCurrentPosition())
		}

		fake.HorizontalTab()
		if fake.GetCurrentPosition() != 20 {
			t.Errorf("Position after second tab = %d, want 20", fake.GetCurrentPosition())
		}
	})
}

func TestFakeCapability_PageMode(t *testing.T) {
	t.Run("tracks page mode settings", func(t *testing.T) {
		fake := NewFakeCapability()

		_, _ = fake.SetPrintAreaPageMode(10, 20, 100, 200)

		state := fake.GetState()
		if state["isPageMode"] != true {
			t.Error("Should be in page mode")
		}
		if fake.printAreaX != 10 || fake.printAreaY != 20 {
			t.Errorf("Print area origin = (%d,%d), want (10,20)",
				fake.printAreaX, fake.printAreaY)
		}
		if fake.printAreaWidth != 100 || fake.printAreaHeight != 200 {
			t.Errorf("Print area size = (%d,%d), want (100,200)",
				fake.printAreaWidth, fake.printAreaHeight)
		}
	})

	t.Run("tracks vertical position", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetAbsoluteVerticalPrintPosition(150)
		if fake.GetCurrentVertPosition() != 150 {
			t.Errorf("VerticalPosition = %d, want 150", fake.GetCurrentVertPosition())
		}

		fake.SetRelativeVerticalPrintPosition(50)
		if fake.GetCurrentVertPosition() != 200 {
			t.Errorf("VerticalPosition = %d, want 200", fake.GetCurrentVertPosition())
		}

		fake.SetRelativeVerticalPrintPosition(-100)
		if fake.GetCurrentVertPosition() != 100 {
			t.Errorf("VerticalPosition = %d, want 100", fake.GetCurrentVertPosition())
		}
	})

	t.Run("tracks print direction", func(t *testing.T) {
		fake := NewFakeCapability()

		_, err := fake.SelectPrintDirectionPageMode(2)
		if err != nil {
			t.Fatalf("SelectPrintDirectionPageMode(2) unexpected error: %v", err)
		}

		if fake.printDirection != 2 {
			t.Errorf("PrintDirection = %d, want 2", fake.printDirection)
		}
	})
}

func TestFakeCapability_CompleteWorkflow(t *testing.T) {
	t.Run("standard mode positioning", func(t *testing.T) {
		fake := NewFakeCapability()

		// Set up print area
		fake.SetLeftMargin(50)
		fake.SetPrintAreaWidth(400)

		// Set justification
		_, _ = fake.SelectJustification(1) // Center

		// Set tabs
		_, _ = fake.SetHorizontalTabPositions([]byte{20, 40, 60})

		// Position operations
		fake.SetAbsolutePrintPosition(100)
		fake.HorizontalTab()
		fake.SetRelativePrintPosition(-10)

		// Verify state
		state := fake.GetState()
		if state["leftMargin"].(int) != 50 {
			t.Errorf("LeftMargin = %d, want 50", state["leftMargin"])
		}
		if state["printAreaWidth"].(int) != 400 {
			t.Errorf("PrintAreaWidth = %d, want 400", state["printAreaWidth"])
		}
		if state["justification"].(byte) != 1 {
			t.Errorf("Justification = %d, want 1", state["justification"])
		}

		// Check command history
		history := fake.GetCommandHistory()
		if len(history) != 7 {
			t.Errorf("Command history length = %d, want 7", len(history))
		}
	})

	t.Run("reset clears all state", func(t *testing.T) {
		fake := NewFakeCapability()

		// Modify state
		fake.SetAbsolutePrintPosition(100)
		fake.SetLeftMargin(50)
		_, _ = fake.SelectJustification(2)

		// Reset
		fake.Reset()

		// Verify reset
		if fake.GetCurrentPosition() != 0 {
			t.Error("Position should be 0 after reset")
		}
		if fake.leftMargin != 0 {
			t.Error("Left margin should be 0 after reset")
		}
		if fake.GetJustification() != 0 {
			t.Error("Justification should be 0 after reset")
		}
		if len(fake.GetBuffer()) != 0 {
			t.Error("Buffer should be empty after reset")
		}
		if len(fake.GetCommandHistory()) != 0 {
			t.Error("Command history should be empty after reset")
		}
	})
}
