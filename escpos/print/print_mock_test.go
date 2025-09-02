package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Ensure MockCapability implements print.Capability
var _ print.Capability = (*MockCapability)(nil)

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

	// Add call counting
	CallCount map[string]int
}

// Add constructor
func NewMockCapability() *MockCapability {
	return &MockCapability{
		CallCount: make(map[string]int),
	}
}

// Add Reset method
func (m *MockCapability) Reset() {
	*m = *NewMockCapability()
}

// Add helper methods
func (m *MockCapability) GetCallCount(method string) int {
	return m.CallCount[method]
}

func (m *MockCapability) Text(n string) ([]byte, error) {
	m.TextCalled = true
	m.TextInput = n
	m.CallCount["Text"]++

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
	m.CallCount["PrintAndFeedPaper"]++

	if m.PrintAndFeedPaperReturn != nil {
		return m.PrintAndFeedPaperReturn
	}
	return []byte{common.ESC, 'J', n}
}

func (m *MockCapability) FormFeed() []byte {
	m.FormFeedCalled = true
	m.CallCount["FormFeed"]++

	if m.FormFeedReturn != nil {
		return m.FormFeedReturn
	}
	return []byte{print.FF}
}

func (m *MockCapability) PrintAndCarriageReturn() []byte {
	m.PrintAndCarriageReturnCalled = true
	m.CallCount["PrintAndCarriageReturn"]++

	if m.PrintAndCarriageReturnReturn != nil {
		return m.PrintAndCarriageReturnReturn
	}
	return []byte{print.CR}
}

func (m *MockCapability) PrintAndLineFeed() []byte {
	m.PrintAndLineFeedCalled = true
	m.CallCount["PrintAndLineFeed"]++

	if m.PrintAndLineFeedReturn != nil {
		return m.PrintAndLineFeedReturn
	}
	return []byte{print.LF}
}

type MockPageModeCapability struct {
	// PrintAndReverseFeed tracking
	PrintAndReverseFeedCalled bool
	PrintAndReverseFeedInput  byte
	PrintAndReverseFeedReturn []byte
	PrintAndReverseFeedError  error

	// PrintAndReverseFeedLines tracking
	PrintAndReverseFeedLinesCalled bool
	PrintAndReverseFeedLinesInput  byte
	PrintAndReverseFeedLinesReturn []byte
	PrintAndReverseFeedLinesError  error

	// CancelData tracking
	CancelDataCalled bool
	CancelDataReturn []byte

	// PrintDataInPageMode tracking
	PrintDataInPageModeCalled bool
	PrintDataInPageModeReturn []byte

	// PrintAndFeedLines tracking
	PrintAndFeedLinesCalled bool
	PrintAndFeedLinesInput  byte
	PrintAndFeedLinesReturn []byte
	PrintAndFeedLinesError  error

	// Add call counting
	CallCount map[string]int
}

// Add constructor
func NewMockPageModeCapability() *MockPageModeCapability {
	return &MockPageModeCapability{
		CallCount: make(map[string]int),
	}
}

// Add Reset method
func (m *MockPageModeCapability) Reset() {
	*m = *NewMockPageModeCapability()
}

// Add helper methods
func (m *MockPageModeCapability) GetCallCount(method string) int {
	return m.CallCount[method]

}

func (m *MockPageModeCapability) PrintAndReverseFeed(units byte) ([]byte, error) {
	m.PrintAndReverseFeedCalled = true
	m.PrintAndReverseFeedInput = units
	m.CallCount["PrintAndReverseFeed"]++

	if m.PrintAndReverseFeedError != nil {
		return nil, m.PrintAndReverseFeedError
	}
	if m.PrintAndReverseFeedReturn != nil {
		return m.PrintAndReverseFeedReturn, nil
	}
	if units > print.MaxReverseMotionUnits {
		return nil, print.ErrPrintReverseFeed
	}
	return []byte{common.ESC, 'K', units}, nil
}

func (m *MockPageModeCapability) PrintAndReverseFeedLines(lines byte) ([]byte, error) {
	m.PrintAndReverseFeedLinesCalled = true
	m.PrintAndReverseFeedLinesInput = lines
	m.CallCount["PrintAndReverseFeedLines"]++

	if m.PrintAndReverseFeedLinesError != nil {
		return nil, m.PrintAndReverseFeedLinesError
	}
	if m.PrintAndReverseFeedLinesReturn != nil {
		return m.PrintAndReverseFeedLinesReturn, nil
	}
	if lines > print.MaxReverseFeedLines {
		return nil, print.ErrPrintReverseFeed
	}
	return []byte{common.ESC, 'K', lines}, nil
}

func (m *MockPageModeCapability) CancelData() []byte {
	m.CancelDataCalled = true
	m.CallCount["CancelData"]++

	if m.CancelDataReturn != nil {
		return m.CancelDataReturn
	}
	return []byte{print.CAN}
}

func (m *MockPageModeCapability) PrintDataInPageMode() []byte {
	m.PrintDataInPageModeCalled = true
	m.CallCount["PrintDataInPageMode"]++

	if m.PrintDataInPageModeReturn != nil {
		return m.PrintDataInPageModeReturn
	}
	return []byte{common.ESC, print.FF}
}

func (m *MockPageModeCapability) PrintAndFeedLines(lines byte) ([]byte, error) {
	m.PrintAndFeedLinesCalled = true
	m.PrintAndFeedLinesInput = lines
	m.CallCount["PrintAndFeedLines"]++

	if m.PrintAndFeedLinesError != nil {
		return nil, m.PrintAndFeedLinesError
	}
	if m.PrintAndFeedLinesReturn != nil {
		return m.PrintAndFeedLinesReturn, nil
	}
	return []byte{common.ESC, 'd', lines}, nil
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks Text calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.TextReturn = []byte{0xFF, 0xFE, 0xFD}

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
		mock := NewMockCapability()
		mock.TextError = common.ErrEmptyBuffer

		_, err := mock.Text("")

		if !mock.TextCalled {
			t.Error("Text() should be marked as called")
		}
		if !errors.Is(err, common.ErrEmptyBuffer) {
			t.Errorf("Text() error = %v, want %v", err, common.ErrEmptyBuffer)
		}
	})

	t.Run("tracks PrintAndFeedPaper calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.PrintAndFeedPaperReturn = []byte{0xAA, 0xBB, 0xCC}

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
		mock := NewMockCapability()

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
