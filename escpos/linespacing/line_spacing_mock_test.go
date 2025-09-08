package linespacing_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
)

// Ensure MockCapability implements linespacing.Capability
var _ linespacing.Capability = (*MockCapability)(nil)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a utils double for linespacing.Capability interface
type MockCapability struct {
	SetLineSpacingCalled bool
	SetLineSpacingInput  linespacing.Spacing
	SetLineSpacingReturn []byte

	SelectDefaultCalled bool
	SelectDefaultReturn []byte

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

func (m *MockCapability) SetLineSpacing(n linespacing.Spacing) []byte {
	m.SetLineSpacingCalled = true
	m.SetLineSpacingInput = n
	m.CallCount["SetLineSpacing"]++

	if m.SetLineSpacingReturn != nil {
		return m.SetLineSpacingReturn
	}
	return []byte{common.ESC, '3', byte(n)}
}

func (m *MockCapability) SelectDefaultLineSpacing() []byte {
	m.SelectDefaultCalled = true
	m.CallCount["SelectDefaultLineSpacing"]++

	if m.SelectDefaultReturn != nil {
		return m.SelectDefaultReturn
	}
	return []byte{common.ESC, '2'}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks SetLineSpacing calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SetLineSpacingReturn = []byte{0xFF, 0xFF, 0xFF}

		result := mock.SetLineSpacing(50)

		if !mock.SetLineSpacingCalled {
			t.Error("SetLineSpacing() should be marked as called")
		}
		if mock.SetLineSpacingInput != 50 {
			t.Errorf("SetLineSpacing() input = %d, want 50", mock.SetLineSpacingInput)
		}
		if !bytes.Equal(result, []byte{0xFF, 0xFF, 0xFF}) {
			t.Errorf("SetLineSpacing() = %#v, want custom return", result)
		}
	})

	t.Run("tracks SelectDefaultLineSpacing calls", func(t *testing.T) {
		mock := NewMockCapability()

		result := mock.SelectDefaultLineSpacing()

		if !mock.SelectDefaultCalled {
			t.Error("SelectDefaultLineSpacing() should be marked as called")
		}
		expected := []byte{common.ESC, '2'}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectDefaultLineSpacing() = %#v, want %#v", result, expected)
		}
	})

	t.Run("returns default behavior when no return configured", func(t *testing.T) {
		mock := NewMockCapability()

		result := mock.SetLineSpacing(30)

		expected := []byte{common.ESC, '3', 30}
		if !bytes.Equal(result, expected) {
			t.Errorf("SetLineSpacing() = %#v, want %#v", result, expected)
		}
	})
}
