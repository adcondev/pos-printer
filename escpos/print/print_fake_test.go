package print_test

import (
	"bytes"
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

func (f *FakeCapability) GetLastCommand() string {
	return f.lastCommand
}

func (f *FakeCapability) Reset() {
	f.buffer = make([]byte, 0)
	f.position = 0
	f.linesFed = 0
	f.formsFed = 0
	f.lastCommand = ""
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks text and position", func(t *testing.T) {
		fake := NewFakeCapability()

		_, err := fake.Text("Hello")
		if err != nil {
			t.Fatalf("Text() unexpected error: %v", err)
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
	})

	t.Run("tracks form feeds", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.FormFeed()

		if fake.GetFormsFed() != 1 {
			t.Errorf("FormsFed = %d, want 1", fake.GetFormsFed())
		}
		if fake.GetLastCommand() != "FormFeed" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "FormFeed")
		}
	})

	t.Run("simulates print sequence", func(t *testing.T) {
		fake := NewFakeCapability()

		// Simulate printing a document
		_, _ = fake.Text("Title")
		fake.PrintAndLineFeed()
		_, _ = fake.Text("Body")
		fake.FormFeed()

		buffer := fake.GetBuffer()

		// Verify sequence
		if !bytes.Contains(buffer, []byte("Title")) {
			t.Error("Buffer should contain 'Title'")
		}
		if !bytes.Contains(buffer, []byte("Body")) {
			t.Error("Buffer should contain 'Body'")
		}
		if bytes.Count(buffer, []byte{print.LF}) != 1 {
			t.Error("Buffer should contain exactly 1 LF")
		}
		if bytes.Count(buffer, []byte{print.FF}) != 1 {
			t.Error("Buffer should contain exactly 1 FF")
		}
		if fake.GetLinesFed() != 1 {
			t.Errorf("LinesFed = %d, want 1", fake.GetLinesFed())
		}
		if fake.GetFormsFed() != 1 {
			t.Errorf("FormsFed = %d, want 1", fake.GetFormsFed())
		}
	})

	t.Run("reset clears state", func(t *testing.T) {
		fake := NewFakeCapability()

		_, _ = fake.Text("data")
		fake.PrintAndLineFeed()
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
	})
}
