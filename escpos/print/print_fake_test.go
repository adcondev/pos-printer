package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Ensure FakeCapability implements print.Capability
var _ print.Capability = (*FakeCapability)(nil)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates a real printer with state tracking
type FakeCapability struct {
	buffer         []byte
	position       int
	linesFed       int
	formsFed       int
	paperFed       int // Track total units of paper fed
	pageBuffer     []byte
	reverseFed     int
	linesReversed  int
	pagesPrinted   int
	dataCleared    int
	lastCommand    string
	commandHistory []string
}

func (fc *FakeCapability) GetCommandHistory() []string {
	return fc.commandHistory
}

func (fc *FakeCapability) GetState() map[string]interface{} {
	return map[string]interface{}{
		"buffer":        fc.buffer,
		"position":      fc.position,
		"linesFed":      fc.linesFed,
		"formsFed":      fc.formsFed,
		"paperFed":      fc.paperFed,
		"reverseFed":    fc.reverseFed,
		"linesReversed": fc.linesReversed,
		"pagesPrinted":  fc.pagesPrinted,
		"dataCleared":   fc.dataCleared,
	}
}

// NewFakeCapability creates a new fake printer
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:         make([]byte, 0),
		position:       0,
		pageBuffer:     make([]byte, 0),
		commandHistory: make([]string, 0),
	}
}

func (fc *FakeCapability) Text(n string) ([]byte, error) {
	if n == "" {
		return nil, common.ErrEmptyBuffer
	}

	data := print.Formatting([]byte(n))
	fc.buffer = append(fc.buffer, data...)
	fc.position += len(data)
	fc.lastCommand = "Text"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return data, nil
}

func (fc *FakeCapability) PrintAndFeedPaper(n byte) []byte {
	cmd := []byte{common.ESC, 'J', n}
	fc.buffer = append(fc.buffer, cmd...)
	fc.position = 0 // Reset to beginning of line
	fc.paperFed += int(n)
	fc.lastCommand = "PrintAndFeedPaper"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return cmd
}

func (fc *FakeCapability) FormFeed() []byte {
	fc.buffer = append(fc.buffer, print.FF)
	fc.position = 0
	fc.formsFed++
	fc.lastCommand = "FormFeed"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return []byte{print.FF}
}

func (fc *FakeCapability) PrintAndCarriageReturn() []byte {
	fc.buffer = append(fc.buffer, print.CR)
	fc.position = 0
	fc.lastCommand = "PrintAndCarriageReturn"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return []byte{print.CR}
}

func (fc *FakeCapability) PrintAndLineFeed() []byte {
	fc.buffer = append(fc.buffer, print.LF)
	fc.position = 0
	fc.linesFed++
	fc.lastCommand = "PrintAndLineFeed"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return []byte{print.LF}
}

// Helper methods
func (fc *FakeCapability) GetBuffer() []byte {
	return fc.buffer
}

func (fc *FakeCapability) GetPosition() int {
	return fc.position
}

func (fc *FakeCapability) GetLinesFed() int {
	return fc.linesFed
}

func (fc *FakeCapability) GetFormsFed() int {
	return fc.formsFed
}

func (fc *FakeCapability) GetPaperFed() int {
	return fc.paperFed
}

func (fc *FakeCapability) GetLastCommand() string {
	return fc.lastCommand
}

func (fc *FakeCapability) Reset() {
	fc.buffer = make([]byte, 0)
	fc.position = 0
	fc.linesFed = 0
	fc.formsFed = 0
	fc.paperFed = 0
	fc.lastCommand = ""
	fc.pageBuffer = make([]byte, 0)
	fc.reverseFed = 0
	fc.linesReversed = 0
	fc.pagesPrinted = 0
	fc.dataCleared = 0
	fc.commandHistory = make([]string, 0)
}

func (fc *FakeCapability) PrintAndReverseFeed(n byte) ([]byte, error) {
	if n > print.MaxReverseMotionUnits {
		return nil, print.ErrPrintReverseFeed
	}
	cmd := []byte{common.ESC, 'K', n}
	fc.buffer = append(fc.buffer, cmd...)
	fc.reverseFed += int(n)
	fc.lastCommand = "PrintAndReverseFeed"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return cmd, nil
}

func (fc *FakeCapability) PrintAndReverseFeedLines(n byte) ([]byte, error) {
	if n > print.MaxReverseFeedLines {
		return nil, print.ErrPrintReverseFeedLines
	}
	cmd := []byte{common.ESC, 'e', n}
	fc.buffer = append(fc.buffer, cmd...)
	fc.linesReversed += int(n)
	fc.lastCommand = "PrintAndReverseFeedLines"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return cmd, nil
}

func (fc *FakeCapability) CancelData() []byte {
	fc.pageBuffer = make([]byte, 0)
	fc.dataCleared++
	fc.lastCommand = "CancelData"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return []byte{print.CAN}
}

func (fc *FakeCapability) PrintDataInPageMode() []byte {
	fc.pagesPrinted++
	fc.lastCommand = "PrintDataInPageMode"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return []byte{common.ESC, print.FF}
}

func (fc *FakeCapability) PrintAndFeedLines(n byte) []byte {
	cmd := []byte{common.ESC, 'd', n}
	fc.buffer = append(fc.buffer, cmd...)
	fc.lastCommand = "PrintAndFeedLines"
	fc.commandHistory = append(fc.commandHistory, fc.lastCommand)
	return cmd
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks text and position", func(t *testing.T) {
		fake := NewFakeCapability()

		data, err := fake.Text("Hello")
		if err != nil {
			t.Fatalf("Text() unexpected error: %v", err)
		}

		expected := []byte("Hello")
		if !bytes.Equal(data, expected) {
			t.Errorf("Text() = %#v, want %#v", data, expected)
		}
		if fake.GetPosition() != 5 {
			t.Errorf("Position = %d, want 5", fake.GetPosition())
		}
		if fake.GetLastCommand() != "Text" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "Text")
		}
		if !bytes.Contains(fake.GetBuffer(), []byte("Hello")) {
			t.Error("Buffer should contain 'Hello'")
		}
	})

	t.Run("handles empty text error", func(t *testing.T) {
		fake := NewFakeCapability()

		_, err := fake.Text("")
		if !errors.Is(err, common.ErrEmptyBuffer) {
			t.Errorf("Text(\"\") error = %v, want %v", err, common.ErrEmptyBuffer)
		}
	})

	t.Run("tracks line feeds", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.PrintAndLineFeed()
		fake.PrintAndLineFeed()

		if fake.GetLinesFed() != 2 {
			t.Errorf("LinesFed = %d, want 2", fake.GetLinesFed())
		}
		if fake.GetPosition() != 0 {
			t.Errorf("Position = %d, want 0 (reset after line feed)", fake.GetPosition())
		}
		if fake.GetLastCommand() != "PrintAndLineFeed" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "PrintAndLineFeed")
		}
	})

	t.Run("tracks paper feed", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.PrintAndFeedPaper(50)
		fake.PrintAndFeedPaper(30)

		if fake.GetPaperFed() != 80 {
			t.Errorf("PaperFed = %d, want 80", fake.GetPaperFed())
		}
		if fake.GetLastCommand() != "PrintAndFeedPaper" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "PrintAndFeedPaper")
		}
	})

	t.Run("tracks form feeds", func(t *testing.T) {
		fake := NewFakeCapability()

		result := fake.FormFeed()

		expected := []byte{print.FF}
		if !bytes.Equal(result, expected) {
			t.Errorf("FormFeed() = %#v, want %#v", result, expected)
		}
		if fake.GetFormsFed() != 1 {
			t.Errorf("FormsFed = %d, want 1", fake.GetFormsFed())
		}
		if fake.GetLastCommand() != "FormFeed" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "FormFeed")
		}
	})

	t.Run("simulates complete print sequence", func(t *testing.T) {
		fake := NewFakeCapability()

		// Simulate printing a receipt
		_, _ = fake.Text("Store Name")
		fake.PrintAndLineFeed()
		_, _ = fake.Text("Item 1")
		fake.PrintAndLineFeed()
		_, _ = fake.Text("Total: $10.00")
		fake.PrintAndFeedPaper(100)
		fake.FormFeed()

		buffer := fake.GetBuffer()

		// Verify complete sequence
		if !bytes.Contains(buffer, []byte("Store Name")) {
			t.Error("Buffer should contain 'Store Name'")
		}
		if !bytes.Contains(buffer, []byte("Item 1")) {
			t.Error("Buffer should contain 'Item 1'")
		}
		if !bytes.Contains(buffer, []byte("Total: $10.00")) {
			t.Error("Buffer should contain 'Total: $10.00'")
		}
		if bytes.Count(buffer, []byte{print.LF}) != 2 {
			t.Errorf("Buffer should contain exactly 2 LF, got %d", bytes.Count(buffer, []byte{print.LF}))
		}
		if bytes.Count(buffer, []byte{print.FF}) != 1 {
			t.Errorf("Buffer should contain exactly 1 FF, got %d", bytes.Count(buffer, []byte{print.FF}))
		}
		if fake.GetLinesFed() != 2 {
			t.Errorf("LinesFed = %d, want 2", fake.GetLinesFed())
		}
		if fake.GetFormsFed() != 1 {
			t.Errorf("FormsFed = %d, want 1", fake.GetFormsFed())
		}
		if fake.GetPaperFed() != 100 {
			t.Errorf("PaperFed = %d, want 100", fake.GetPaperFed())
		}
	})

	t.Run("reset clears all state", func(t *testing.T) {
		fake := NewFakeCapability()

		// Add some state
		_, _ = fake.Text("data")
		fake.PrintAndLineFeed()
		fake.PrintAndFeedPaper(50)
		fake.FormFeed()

		// Reset everything
		fake.Reset()

		if len(fake.GetBuffer()) != 0 {
			t.Error("Buffer should be empty after reset")
		}
		if fake.GetPosition() != 0 {
			t.Error("Position should be 0 after reset")
		}
		if fake.GetLinesFed() != 0 {
			t.Error("LinesFed should be 0 after reset")
		}
		if fake.GetFormsFed() != 0 {
			t.Error("FormsFed should be 0 after reset")
		}
		if fake.GetPaperFed() != 0 {
			t.Error("PaperFed should be 0 after reset")
		}
		if fake.GetLastCommand() != "" {
			t.Error("LastCommand should be empty after reset")
		}
	})
}
