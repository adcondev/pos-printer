package escpos

import (
	"bytes"
	"testing"
)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeLineSpacing simulates line spacing with state tracking
type FakeLineSpacing struct {
	buffer    []byte
	lineSpace byte
}

// NewFakeLineSpacing creates a new fake line spacing instance
func NewFakeLineSpacing() *FakeLineSpacing {
	return &FakeLineSpacing{
		buffer:    make([]byte, 0),
		lineSpace: 30, // Default line spacing
	}
}

func (fls *FakeLineSpacing) SetLineSpacing(n byte) []byte {
	cmd := []byte{ESC, '3', n}
	fls.buffer = append(fls.buffer, cmd...)
	fls.lineSpace = n
	return cmd
}

func (fls *FakeLineSpacing) SelectDefaultLineSpacing() []byte {
	cmd := []byte{ESC, '2'}
	fls.buffer = append(fls.buffer, cmd...)
	fls.lineSpace = 30 // Reset to default
	return cmd
}

// GetLineSpace returns current line spacing for verification
func (fls *FakeLineSpacing) GetLineSpace() byte {
	return fls.lineSpace
}

// GetBuffer returns the command buffer for verification
func (fls *FakeLineSpacing) GetBuffer() []byte {
	return fls.buffer
}

// ============================================================================
// Fake Implementation Tests
// ============================================================================

func TestFakeLineSpacing_SetLineSpacing_StateTracking(t *testing.T) {
	fake := NewFakeLineSpacing()

	// Verify initial state
	if fake.GetLineSpace() != 30 {
		t.Errorf("FakeLineSpacing initial line space = %d, want 30", fake.GetLineSpace())
	}

	// Test setting custom spacing
	customSpacing := byte(45)
	result := fake.SetLineSpacing(customSpacing)

	// Verify return value
	expectedCmd := []byte{ESC, '3', customSpacing}
	if !bytes.Equal(result, expectedCmd) {
		t.Errorf("FakeLineSpacing.SetLineSpacing(%d) = %#v, want %#v",
			customSpacing, result, expectedCmd)
	}

	// Verify state change
	if fake.GetLineSpace() != customSpacing {
		t.Errorf("FakeLineSpacing line space after SetLineSpacing(%d) = %d, want %d",
			customSpacing, fake.GetLineSpace(), customSpacing)
	}

	// Verify buffer accumulation
	if !bytes.Contains(fake.GetBuffer(), expectedCmd) {
		t.Error("FakeLineSpacing buffer should contain SetLineSpacing command")
	}
}

func TestFakeLineSpacing_SelectDefaultLineSpacing_StateReset(t *testing.T) {
	fake := NewFakeLineSpacing()

	// Set custom spacing first
	fake.SetLineSpacing(60)

	// Reset to default
	result := fake.SelectDefaultLineSpacing()

	// Verify return value
	expectedCmd := []byte{ESC, '2'}
	if !bytes.Equal(result, expectedCmd) {
		t.Errorf("FakeLineSpacing.SelectDefaultLineSpacing() = %#v, want %#v",
			result, expectedCmd)
	}

	// Verify state reset to default
	if fake.GetLineSpace() != 30 {
		t.Errorf("FakeLineSpacing line space after SelectDefaultLineSpacing() = %d, want 30",
			fake.GetLineSpace())
	}
}
