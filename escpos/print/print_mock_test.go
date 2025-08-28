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
	return []byte(n), nil
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
			TextReturn: []byte("mocked"),
		}

		result, err := mock.Text("input")

		if !mock.TextCalled {
			t.Error("Text() should be marked as called")
		}
		if mock.TextInput != "input" {
			t.Errorf("Text() input = %q, want %q", mock.TextInput, "input")
		}
		if err != nil {
			t.Errorf("Text() unexpected error: %v", err)
		}
		if !bytes.Equal(result, []byte("mocked")) {
			t.Errorf("Text() = %#v, want %#v", result, []byte("mocked"))
		}
	})

	t.Run("simulates errors", func(t *testing.T) {
		mock := &MockCapability{
			TextError: common.ErrEmptyBuffer,
		}

		_, err := mock.Text("")

		if !errors.Is(err, common.ErrEmptyBuffer) {
			t.Errorf("Text() error = %v, want %v", err, common.ErrEmptyBuffer)
		}
	})

	t.Run("tracks PrintAndFeedPaper calls", func(t *testing.T) {
		mock := &MockCapability{}

		result := mock.PrintAndFeedPaper(5)

		if !mock.PrintAndFeedPaperCalled {
			t.Error("PrintAndFeedPaper() should be marked as called")
		}
		if mock.PrintAndFeedPaperInput != 5 {
			t.Errorf("PrintAndFeedPaper() input = %d, want 5", mock.PrintAndFeedPaperInput)
		}
		expected := []byte{common.ESC, 'J', 5}
		if !bytes.Equal(result, expected) {
			t.Errorf("PrintAndFeedPaper() = %#v, want %#v", result, expected)
		}
	})
}
