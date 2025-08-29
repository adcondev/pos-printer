package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates a real printer with state tracking
type FakeCapability struct {
	buffer      []byte
	position    int
	linesFed    int
	formsFed    int
	paperFed    int // Track total units of paper fed
	lastCommand string
}

// NewFakeCapability creates a new fake printer
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:   make([]byte, 0),
		position: 0,
	}
}

// Ensure FakeCapability implements print.Capability
var _ print.Capability = (*FakeCapability)(nil)

func (f *FakeCapability) Text(n string) ([]byte, error) {
	if n == "" {
		return nil, common.ErrEmptyBuffer
	}

	data := print.Formatting([]byte(n))
	f.buffer = append(f.buffer, data...)
	f.position += len(data)
	f.lastCommand = "Text"

	return data, nil
}

func (f *FakeCapability) PrintAndFeedPaper(n byte) []byte {
	cmd := []byte{common.ESC, 'J', n}
	f.buffer = append(f.buffer, cmd...)
	f.position = 0 // Reset to beginning of line
	f.paperFed += int(n)
	f.lastCommand = "PrintAndFeedPaper"
	return cmd
}

func (f *FakeCapability) FormFeed() []byte {
	f.buffer = append(f.buffer, print.FF)
	f.position = 0
	f.formsFed++
	f.lastCommand = "FormFeed"
	return []byte{print.FF}
}

func (f *FakeCapability) PrintAndCarriageReturn() []byte {
	f.buffer = append(f.buffer, print.CR)
	f.position = 0
	f.lastCommand = "PrintAndCarriageReturn"
	return []byte{print.CR}
}

func (f *FakeCapability) PrintAndLineFeed() []byte {
	f.buffer = append(f.buffer, print.LF)
	f.position = 0
	f.linesFed++
	f.lastCommand = "PrintAndLineFeed"
	return []byte{print.LF}
}

// Helper methods
func (f *FakeCapability) GetBuffer() []byte {
	return f.buffer
}

func (f *FakeCapability) GetPosition() int {
	return f.position
}

func (f *FakeCapability) GetLinesFed() int {
	return f.linesFed
}

func (f *FakeCapability) GetFormsFed() int {
	return f.formsFed
}

func (f *FakeCapability) GetPaperFed() int {
	return f.paperFed
}

func (f *FakeCapability) GetLastCommand() string {
	return f.lastCommand
}

func (f *FakeCapability) Reset() {
	f.buffer = make([]byte, 0)
	f.position = 0
	f.linesFed = 0
	f.formsFed = 0
	f.paperFed = 0
	f.lastCommand = ""
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
