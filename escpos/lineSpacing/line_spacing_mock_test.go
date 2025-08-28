package lineSpacing_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/lineSpacing"
)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a test double for lineSpacing.Capability interface
type MockCapability struct {
	SetLineSpacingCalled bool
	SetLineSpacingInput  byte
	SetLineSpacingReturn []byte

	SelectDefaultCalled bool
	SelectDefaultReturn []byte
}

// Ensure MockCapability implements lineSpacing.Capability
var _ lineSpacing.Capability = (*MockCapability)(nil)

func (m *MockCapability) SetLineSpacing(n byte) []byte {
	m.SetLineSpacingCalled = true
	m.SetLineSpacingInput = n

	if m.SetLineSpacingReturn != nil {
		return m.SetLineSpacingReturn
	}
	return []byte{common.ESC, '3', n}
}

func (m *MockCapability) SelectDefaultLineSpacing() []byte {
	m.SelectDefaultCalled = true

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
		mock := &MockCapability{
			SetLineSpacingReturn: []byte{0xFF, 0xFF, 0xFF},
		}

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
		mock := &MockCapability{}

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
		mock := &MockCapability{}

		result := mock.SetLineSpacing(30)

		expected := []byte{common.ESC, '3', 30}
		if !bytes.Equal(result, expected) {
			t.Errorf("SetLineSpacing() = %#v, want %#v", result, expected)
		}
	})
}
