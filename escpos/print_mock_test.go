package escpos

import (
	"bytes"
	"errors"
	"testing"
)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockPrinterCapability provides a test double for PrinterCapability interface
type MockPrinterCapability struct {
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

// Text records the call and returns configured response
func (m *MockPrinterCapability) Text(n string) ([]byte, error) {
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

// PrintAndFeedPaper records the call and returns configured response
func (m *MockPrinterCapability) PrintAndFeedPaper(n byte) []byte {
	m.PrintAndFeedPaperCalled = true
	m.PrintAndFeedPaperInput = n

	if m.PrintAndFeedPaperReturn != nil {
		return m.PrintAndFeedPaperReturn
	}
	return []byte{ESC, 'J', n}
}

// FormFeed records the call and returns configured response
func (m *MockPrinterCapability) FormFeed() []byte {
	m.FormFeedCalled = true
	if m.FormFeedReturn != nil {
		return m.FormFeedReturn
	}
	return []byte{FF}
}

// PrintAndCarriageReturn records the call and returns configured response
func (m *MockPrinterCapability) PrintAndCarriageReturn() []byte {
	m.PrintAndCarriageReturnCalled = true
	if m.PrintAndCarriageReturnReturn != nil {
		return m.PrintAndCarriageReturnReturn
	}
	return []byte{CR}
}

// PrintAndLineFeed records the call and returns configured response
func (m *MockPrinterCapability) PrintAndLineFeed() []byte {
	m.PrintAndLineFeedCalled = true
	if m.PrintAndLineFeedReturn != nil {
		return m.PrintAndLineFeedReturn
	}
	return []byte{LF}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockPrinterCapability_Text_BehaviorTracking(t *testing.T) {
	t.Run("tracks method calls and input", func(t *testing.T) {
		mock := &MockPrinterCapability{
			TextReturn: []byte("mocked response"),
		}

		result, err := mock.Text("test input")

		// Verify tracking
		if !mock.TextCalled {
			t.Error("MockPrinterCapability.Text() should mark TextCalled as true")
		}
		if mock.TextInput != "test input" {
			t.Errorf("MockPrinterCapability.Text() input = %q, want %q", mock.TextInput, "test input")
		}
		if err != nil {
			t.Errorf("MockPrinterCapability.Text() unexpected error: %v", err)
		}
		if !bytes.Equal(result, []byte("mocked response")) {
			t.Errorf("MockPrinterCapability.Text() = %#v, want %#v", result, []byte("mocked response"))
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		expectedErr := errors.New("mock error")
		mock := &MockPrinterCapability{
			TextError: expectedErr,
		}

		_, err := mock.Text("any input")

		if err != expectedErr {
			t.Errorf("MockPrinterCapability.Text() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestMockPrinterCapability_Integration_WithCommands(t *testing.T) {
	mock := &MockPrinterCapability{
		TextReturn: []byte("mocked"),
	}

	// Inject mock into Commands
	cmd := &Commands{
		Print: mock,
	}

	// Use the command
	result, err := cmd.Print.Text("test")

	if err != nil {
		t.Errorf("Commands.Print.Text() unexpected error: %v", err)
	}

	// Verify behavior
	if !mock.TextCalled {
		t.Error("MockPrinterCapability.Text() was not called")
	}
	if mock.TextInput != "test" {
		t.Errorf("MockPrinterCapability received input %q, want %q", mock.TextInput, "test")
	}
	if !bytes.Equal(result, []byte("mocked")) {
		t.Errorf("Commands.Print.Text() = %#v, want %#v", result, []byte("mocked"))
	}
}
