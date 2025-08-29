package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a test double for print.Capability interface
type MockCapability struct {
	// Text tracking
	TextCalled bool
	TextInput  string
	TextReturn []byte
	TextError  error

	// PrintAndFeedPaper tracking
	PrintAndFeedPaperCalled bool
	PrintAndFeedPaperInput  byte
	PrintAndFeedPaperReturn []byte

	// FormFeed tracking
	FormFeedCalled bool
	FormFeedReturn []byte

	// PrintAndCarriageReturn tracking
	PrintAndCarriageReturnCalled bool
	PrintAndCarriageReturnReturn []byte

	// PrintAndLineFeed tracking
	PrintAndLineFeedCalled bool
	PrintAndLineFeedReturn []byte
}

// Ensure MockCapability implements print.Capability
var _ print.Capability = (*MockCapability)(nil)

func (m *MockCapability) Text(n string) ([]byte, error) {
	m.TextCalled = true
	m.TextInput = n

	if m.TextError != nil {
		return nil, m.TextError
	}
	if m.TextReturn != nil {
		return m.TextReturn, nil
	}
	return print.Formatting([]byte(n)), nil
}

func (m *MockCapability) PrintAndFeedPaper(n byte) []byte {
	m.PrintAndFeedPaperCalled = true
	m.PrintAndFeedPaperInput = n

	if m.PrintAndFeedPaperReturn != nil {
		return m.PrintAndFeedPaperReturn
	}
	return []byte{common.ESC, 'J', n}
}

func (m *MockCapability) FormFeed() []byte {
	m.FormFeedCalled = true
	if m.FormFeedReturn != nil {
		return m.FormFeedReturn
	}
	return []byte{print.FF}
}

func (m *MockCapability) PrintAndCarriageReturn() []byte {
	m.PrintAndCarriageReturnCalled = true
	if m.PrintAndCarriageReturnReturn != nil {
		return m.PrintAndCarriageReturnReturn
	}
	return []byte{print.CR}
}

func (m *MockCapability) PrintAndLineFeed() []byte {
	m.PrintAndLineFeedCalled = true
	if m.PrintAndLineFeedReturn != nil {
		return m.PrintAndLineFeedReturn
	}
	return []byte{print.LF}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks Text calls", func(t *testing.T) {
		mock := &MockCapability{
			TextReturn: []byte{0xFF, 0xFE, 0xFD},
		}

		result, err := mock.Text("test input")

		if !mock.TextCalled {
			t.Error("Text() should be marked as called")
		}
		if mock.TextInput != "test input" {
			t.Errorf("Text() input = %q, want %q", mock.TextInput, "test input")
		}
		if err != nil {
			t.Errorf("Text() unexpected error: %v", err)
		}
		if !bytes.Equal(result, []byte{0xFF, 0xFE, 0xFD}) {
			t.Errorf("Text() = %#v, want custom return", result)
		}
	})

	t.Run("simulates Text error", func(t *testing.T) {
		mock := &MockCapability{
			TextError: common.ErrEmptyBuffer,
		}

		_, err := mock.Text("")

		if !mock.TextCalled {
			t.Error("Text() should be marked as called")
		}
		if !errors.Is(err, common.ErrEmptyBuffer) {
			t.Errorf("Text() error = %v, want %v", err, common.ErrEmptyBuffer)
		}
	})

	t.Run("tracks PrintAndFeedPaper calls", func(t *testing.T) {
		mock := &MockCapability{
			PrintAndFeedPaperReturn: []byte{0xAA, 0xBB, 0xCC},
		}

		result := mock.PrintAndFeedPaper(100)

		if !mock.PrintAndFeedPaperCalled {
			t.Error("PrintAndFeedPaper() should be marked as called")
		}
		if mock.PrintAndFeedPaperInput != 100 {
			t.Errorf("PrintAndFeedPaper() input = %d, want 100", mock.PrintAndFeedPaperInput)
		}
		if !bytes.Equal(result, []byte{0xAA, 0xBB, 0xCC}) {
			t.Errorf("PrintAndFeedPaper() = %#v, want custom return", result)
		}
	})

	t.Run("returns default behavior when no return configured", func(t *testing.T) {
		mock := &MockCapability{}

		result := mock.PrintAndLineFeed()

		if !mock.PrintAndLineFeedCalled {
			t.Error("PrintAndLineFeed() should be marked as called")
		}
		expected := []byte{print.LF}
		if !bytes.Equal(result, expected) {
			t.Errorf("PrintAndLineFeed() = %#v, want %#v", result, expected)
		}
	})
}
