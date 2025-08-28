package lineSpacing_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/lineSpacing"
)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates line spacing with state tracking
type FakeCapability struct {
	buffer         []byte
	currentSpacing byte
	defaultSpacing byte
	timesChanged   int
	lastCommand    string
}

// NewFakeCapability creates a new fake line spacing
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:         make([]byte, 0),
		currentSpacing: 30, // Default
		defaultSpacing: 30,
	}
}

// Ensure FakeCapability implements lineSpacing.Capability
var _ lineSpacing.Capability = (*FakeCapability)(nil)

func (f *FakeCapability) SetLineSpacing(n byte) []byte {
	cmd := []byte{common.ESC, '3', n}
	f.buffer = append(f.buffer, cmd...)
	f.currentSpacing = n
	f.timesChanged++
	f.lastCommand = "SetLineSpacing"
	return cmd
}

func (f *FakeCapability) SelectDefaultLineSpacing() []byte {
	cmd := []byte{common.ESC, '2'}
	f.buffer = append(f.buffer, cmd...)
	f.currentSpacing = f.defaultSpacing
	f.lastCommand = "SelectDefaultLineSpacing"
	return cmd
}

// Helper methods
func (f *FakeCapability) GetCurrentSpacing() byte {
	return f.currentSpacing
}

func (f *FakeCapability) GetBuffer() []byte {
	return f.buffer
}

func (f *FakeCapability) GetTimesChanged() int {
	return f.timesChanged
}

func (f *FakeCapability) GetLastCommand() string {
	return f.lastCommand
}

func (f *FakeCapability) Reset() {
	f.buffer = make([]byte, 0)
	f.currentSpacing = f.defaultSpacing
	f.timesChanged = 0
	f.lastCommand = ""
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks spacing changes", func(t *testing.T) {
		fake := NewFakeCapability()

		result := fake.SetLineSpacing(45)

		expected := []byte{common.ESC, '3', 45}
		if !bytes.Equal(result, expected) {
			t.Errorf("SetLineSpacing(45) = %#v, want %#v", result, expected)
		}
		if fake.GetCurrentSpacing() != 45 {
			t.Errorf("CurrentSpacing = %d, want 45", fake.GetCurrentSpacing())
		}
		if fake.GetTimesChanged() != 1 {
			t.Errorf("TimesChanged = %d, want 1", fake.GetTimesChanged())
		}
		if fake.GetLastCommand() != "SetLineSpacing" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "SetLineSpacing")
		}
	})

	t.Run("resets to default", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetLineSpacing(60)
		result := fake.SelectDefaultLineSpacing()

		expected := []byte{common.ESC, '2'}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectDefaultLineSpacing() = %#v, want %#v", result, expected)
		}
		if fake.GetCurrentSpacing() != 30 {
			t.Errorf("CurrentSpacing = %d, want 30 (default)", fake.GetCurrentSpacing())
		}
		if fake.GetLastCommand() != "SelectDefaultLineSpacing" {
			t.Errorf("LastCommand = %q, want %q", fake.GetLastCommand(), "SelectDefaultLineSpacing")
		}
	})

	t.Run("accumulates buffer", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetLineSpacing(20)
		fake.SetLineSpacing(40)
		fake.SelectDefaultLineSpacing()

		buffer := fake.GetBuffer()

		// Should contain all commands in sequence
		expectedSequence := [][]byte{
			{common.ESC, '3', 20},
			{common.ESC, '3', 40},
			{common.ESC, '2'},
		}

		for _, expected := range expectedSequence {
			if !bytes.Contains(buffer, expected) {
				t.Errorf("Buffer should contain %#v", expected)
			}
		}

		if fake.GetTimesChanged() != 2 {
			t.Errorf("TimesChanged = %d, want 2", fake.GetTimesChanged())
		}
	})

	t.Run("reset clears state", func(t *testing.T) {
		fake := NewFakeCapability()

		fake.SetLineSpacing(50)
		fake.Reset()

		if len(fake.GetBuffer()) != 0 {
			t.Error("Buffer should be empty after reset")
		}
		if fake.GetCurrentSpacing() != 30 {
			t.Error("CurrentSpacing should be default after reset")
		}
		if fake.GetTimesChanged() != 0 {
			t.Error("TimesChanged should be 0 after reset")
		}
	})
}
