package escpos

import (
	"bytes"
	"testing"
)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockLineSpacingCapability provides a test double for LineSpacingCapability interface
type MockLineSpacingCapability struct {
	SetLineSpacingCalled bool
	SetLineSpacingInput  byte
	SetLineSpacingReturn []byte

	SelectDefaultCalled bool
	SelectDefaultReturn []byte
}

// SetLineSpacing records the call and returns configured response
func (m *MockLineSpacingCapability) SetLineSpacing(n byte) []byte {
	m.SetLineSpacingCalled = true
	m.SetLineSpacingInput = n
	if m.SetLineSpacingReturn != nil {
		return m.SetLineSpacingReturn
	}
	return []byte{ESC, '3', n}
}

// SelectDefaultLineSpacing records the call and returns configured response
func (m *MockLineSpacingCapability) SelectDefaultLineSpacing() []byte {
	m.SelectDefaultCalled = true
	if m.SelectDefaultReturn != nil {
		return m.SelectDefaultReturn
	}
	return []byte{ESC, '2'}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockLineSpacingCapability_SetLineSpacing_BehaviorTracking(t *testing.T) {
	mock := &MockLineSpacingCapability{}

	// Test SetLineSpacing
	result := mock.SetLineSpacing(50)

	if !mock.SetLineSpacingCalled {
		t.Error("MockLineSpacingCapability.SetLineSpacing() should mark SetLineSpacingCalled as true")
	}
	if mock.SetLineSpacingInput != 50 {
		t.Errorf("MockLineSpacingCapability.SetLineSpacing() input = %d, want 50", mock.SetLineSpacingInput)
	}

	expected := []byte{ESC, '3', 50}
	if !bytes.Equal(result, expected) {
		t.Errorf("MockLineSpacingCapability.SetLineSpacing() = %#v, want %#v", result, expected)
	}
}

func TestMockLineSpacingCapability_SelectDefaultLineSpacing_BehaviorTracking(t *testing.T) {
	mock := &MockLineSpacingCapability{}

	// Test SelectDefaultLineSpacing
	result := mock.SelectDefaultLineSpacing()

	if !mock.SelectDefaultCalled {
		t.Error("MockLineSpacingCapability.SelectDefaultLineSpacing() should mark SelectDefaultCalled as true")
	}

	expected := []byte{ESC, '2'}
	if !bytes.Equal(result, expected) {
		t.Errorf("MockLineSpacingCapability.SelectDefaultLineSpacing() = %#v, want %#v", result, expected)
	}
}

func TestMockLineSpacingCapability_Integration_WithCommands(t *testing.T) {
	mock := &MockLineSpacingCapability{}
	cmd := &Commands{
		Print:     &PrintCommands{Page: &PagePrint{}},
		LineSpace: mock,
	}

	// Use the command
	result := cmd.LineSpace.SetLineSpacing(60)

	// Verify the mock was called
	if !mock.SetLineSpacingCalled {
		t.Error("MockLineSpacingCapability.SetLineSpacing() was not called")
	}
	if mock.SetLineSpacingInput != 60 {
		t.Errorf("MockLineSpacingCapability received input %d, want 60", mock.SetLineSpacingInput)
	}

	expected := []byte{ESC, '3', 60}
	if !bytes.Equal(result, expected) {
		t.Errorf("Commands.LineSpace.SetLineSpacing() = %#v, want %#v", result, expected)
	}
}
