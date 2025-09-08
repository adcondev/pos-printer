package printposition_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/printposition"
)

// Ensure MockCapability implements printposition.Capability
var _ printposition.Capability = (*MockCapability)(nil)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a utils double for printposition.Capability interface
type MockCapability struct {
	// SetAbsolutePrintPosition tracking
	SetAbsolutePrintPositionCalled bool
	SetAbsolutePrintPositionInput  uint16
	SetAbsolutePrintPositionReturn []byte

	// SetRelativePrintPosition tracking
	SetRelativePrintPositionCalled bool
	SetRelativePrintPositionInput  int16
	SetRelativePrintPositionReturn []byte

	// HorizontalTab tracking
	HorizontalTabCalled bool
	HorizontalTabReturn []byte

	// SetHorizontalTabPositions tracking
	SetHorizontalTabPositionsCalled bool
	SetHorizontalTabPositionsInput  []byte
	SetHorizontalTabPositionsReturn []byte
	SetHorizontalTabPositionsError  error

	// SelectJustification tracking
	SelectJustificationCalled bool
	SelectJustificationInput  byte
	SelectJustificationReturn []byte
	SelectJustificationError  error

	// SetLeftMargin tracking
	SetLeftMarginCalled bool
	SetLeftMarginInput  uint16
	SetLeftMarginReturn []byte

	// SetPrintAreaWidth tracking
	SetPrintAreaWidthCalled bool
	SetPrintAreaWidthInput  uint16
	SetPrintAreaWidthReturn []byte

	// SetPrintPositionBeginningLine tracking
	SetPrintPositionBeginningLineCalled bool
	SetPrintPositionBeginningLineInput  byte
	SetPrintPositionBeginningLineReturn []byte
	SetPrintPositionBeginningLineError  error

	// SelectPrintDirectionPageMode tracking
	SelectPrintDirectionPageModeCalled bool
	SelectPrintDirectionPageModeInput  byte
	SelectPrintDirectionPageModeReturn []byte
	SelectPrintDirectionPageModeError  error

	// SetPrintAreaPageMode tracking
	SetPrintAreaPageModeCalled bool
	SetPrintAreaPageModeX      uint16
	SetPrintAreaPageModeY      uint16
	SetPrintAreaPageModeWidth  uint16
	SetPrintAreaPageModeHeight uint16
	SetPrintAreaPageModeReturn []byte
	SetPrintAreaPageModeError  error

	// SetAbsoluteVerticalPrintPosition tracking
	SetAbsoluteVerticalPrintPositionCalled bool
	SetAbsoluteVerticalPrintPositionInput  uint16
	SetAbsoluteVerticalPrintPositionReturn []byte

	// SetRelativeVerticalPrintPosition tracking
	SetRelativeVerticalPrintPositionCalled bool
	SetRelativeVerticalPrintPositionInput  int16
	SetRelativeVerticalPrintPositionReturn []byte

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

func (m *MockCapability) SetAbsolutePrintPosition(position uint16) []byte {
	m.SetAbsolutePrintPositionCalled = true
	m.SetAbsolutePrintPositionInput = position
	m.CallCount["SetAbsolutePrintPosition"]++

	if m.SetAbsolutePrintPositionReturn != nil {
		return m.SetAbsolutePrintPositionReturn
	}
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	return []byte{common.ESC, '$', nL, nH}
}

func (m *MockCapability) SetRelativePrintPosition(distance int16) []byte {
	m.SetRelativePrintPositionCalled = true
	m.SetRelativePrintPositionInput = distance
	m.CallCount["SetRelativePrintPosition"]++

	if m.SetRelativePrintPositionReturn != nil {
		return m.SetRelativePrintPositionReturn
	}
	// Two's complement for negative values is needed
	value := uint16(distance) // nolint: gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	return []byte{common.ESC, '\\', nL, nH}
}

func (m *MockCapability) HorizontalTab() []byte {
	m.HorizontalTabCalled = true
	m.CallCount["HorizontalTab"]++

	if m.HorizontalTabReturn != nil {
		return m.HorizontalTabReturn
	}
	return []byte{printposition.HT}
}

func (m *MockCapability) SetHorizontalTabPositions(positions []byte) ([]byte, error) {
	m.SetHorizontalTabPositionsCalled = true
	m.SetHorizontalTabPositionsInput = positions
	m.CallCount["SetHorizontalTabPositions"]++

	if m.SetHorizontalTabPositionsError != nil {
		return nil, m.SetHorizontalTabPositionsError
	}
	if m.SetHorizontalTabPositionsReturn != nil {
		return m.SetHorizontalTabPositionsReturn, nil
	}
	cmd := []byte{common.ESC, 'D'}
	cmd = append(cmd, positions...)
	cmd = append(cmd, common.NUL)
	return cmd, nil
}

func (m *MockCapability) SelectJustification(mode byte) ([]byte, error) {
	m.SelectJustificationCalled = true
	m.SelectJustificationInput = mode
	m.CallCount["SelectJustification"]++

	if m.SelectJustificationError != nil {
		return nil, m.SelectJustificationError
	}
	if m.SelectJustificationReturn != nil {
		return m.SelectJustificationReturn, nil
	}
	return []byte{common.ESC, 'a', mode}, nil
}

func (m *MockCapability) SetLeftMargin(margin uint16) []byte {
	m.SetLeftMarginCalled = true
	m.SetLeftMarginInput = margin
	m.CallCount["SetLeftMargin"]++

	if m.SetLeftMarginReturn != nil {
		return m.SetLeftMarginReturn
	}
	nL := byte(margin & 0xFF)
	nH := byte((margin >> 8) & 0xFF)
	return []byte{common.GS, 'L', nL, nH}
}

func (m *MockCapability) SetPrintAreaWidth(width uint16) []byte {
	m.SetPrintAreaWidthCalled = true
	m.SetPrintAreaWidthInput = width
	m.CallCount["SetPrintAreaWidth"]++

	if m.SetPrintAreaWidthReturn != nil {
		return m.SetPrintAreaWidthReturn
	}
	nL := byte(width & 0xFF)
	nH := byte((width >> 8) & 0xFF)
	return []byte{common.GS, 'W', nL, nH}
}

func (m *MockCapability) SetPrintPositionBeginningLine(mode byte) ([]byte, error) {
	m.SetPrintPositionBeginningLineCalled = true
	m.SetPrintPositionBeginningLineInput = mode
	m.CallCount["SetPrintPositionBeginningLine"]++

	if m.SetPrintPositionBeginningLineError != nil {
		return nil, m.SetPrintPositionBeginningLineError
	}
	if m.SetPrintPositionBeginningLineReturn != nil {
		return m.SetPrintPositionBeginningLineReturn, nil
	}
	return []byte{common.GS, 'T', mode}, nil
}

func (m *MockCapability) SelectPrintDirectionPageMode(direction byte) ([]byte, error) {
	m.SelectPrintDirectionPageModeCalled = true
	m.SelectPrintDirectionPageModeInput = direction
	m.CallCount["SelectPrintDirectionPageMode"]++

	if m.SelectPrintDirectionPageModeError != nil {
		return nil, m.SelectPrintDirectionPageModeError
	}
	if m.SelectPrintDirectionPageModeReturn != nil {
		return m.SelectPrintDirectionPageModeReturn, nil
	}
	return []byte{common.ESC, 'T', direction}, nil
}

func (m *MockCapability) SetPrintAreaPageMode(x, y, width, height uint16) ([]byte, error) {
	m.SetPrintAreaPageModeCalled = true
	m.SetPrintAreaPageModeX = x
	m.SetPrintAreaPageModeY = y
	m.SetPrintAreaPageModeWidth = width
	m.SetPrintAreaPageModeHeight = height
	m.CallCount["SetPrintAreaPageMode"]++

	if m.SetPrintAreaPageModeError != nil {
		return nil, m.SetPrintAreaPageModeError
	}
	if m.SetPrintAreaPageModeReturn != nil {
		return m.SetPrintAreaPageModeReturn, nil
	}
	xL := byte(x & 0xFF)
	xH := byte((x >> 8) & 0xFF)
	yL := byte(y & 0xFF)
	yH := byte((y >> 8) & 0xFF)
	dxL := byte(width & 0xFF)
	dxH := byte((width >> 8) & 0xFF)
	dyL := byte(height & 0xFF)
	dyH := byte((height >> 8) & 0xFF)
	return []byte{common.ESC, 'W', xL, xH, yL, yH, dxL, dxH, dyL, dyH}, nil
}

func (m *MockCapability) SetAbsoluteVerticalPrintPosition(position uint16) []byte {
	m.SetAbsoluteVerticalPrintPositionCalled = true
	m.SetAbsoluteVerticalPrintPositionInput = position
	m.CallCount["SetAbsoluteVerticalPrintPosition"]++

	if m.SetAbsoluteVerticalPrintPositionReturn != nil {
		return m.SetAbsoluteVerticalPrintPositionReturn
	}
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	return []byte{common.GS, '$', nL, nH}
}

func (m *MockCapability) SetRelativeVerticalPrintPosition(distance int16) []byte {
	m.SetRelativeVerticalPrintPositionCalled = true
	m.SetRelativeVerticalPrintPositionInput = distance
	m.CallCount["SetRelativeVerticalPrintPosition"]++

	if m.SetRelativeVerticalPrintPositionReturn != nil {
		return m.SetRelativeVerticalPrintPositionReturn
	}
	// Two's complement for negative values is needed
	value := uint16(distance) // nolint: gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	return []byte{common.GS, '\\', nL, nH}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks HorizontalTab calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.HorizontalTabReturn = []byte{0xFF}

		result := mock.HorizontalTab()

		if !mock.HorizontalTabCalled {
			t.Error("HorizontalTab() should be marked as called")
		}
		if !bytes.Equal(result, []byte{0xFF}) {
			t.Errorf("HorizontalTab() = %#v, want custom return", result)
		}
		if mock.GetCallCount("HorizontalTab") != 1 {
			t.Errorf("HorizontalTab call count = %d, want 1", mock.GetCallCount("HorizontalTab"))
		}
	})

	t.Run("tracks SetAbsolutePrintPosition calls", func(t *testing.T) {
		mock := NewMockCapability()

		result := mock.SetAbsolutePrintPosition(500)

		if !mock.SetAbsolutePrintPositionCalled {
			t.Error("SetAbsolutePrintPosition() should be marked as called")
		}
		if mock.SetAbsolutePrintPositionInput != 500 {
			t.Errorf("SetAbsolutePrintPosition() input = %d, want 500",
				mock.SetAbsolutePrintPositionInput)
		}
		expected := []byte{common.ESC, '$', 0xF4, 0x01} // 500 in little-endian
		if !bytes.Equal(result, expected) {
			t.Errorf("SetAbsolutePrintPosition() = %#v, want %#v", result, expected)
		}
	})

	t.Run("simulates SelectJustification error", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SelectJustificationError = printposition.ErrInvalidJustification

		_, err := mock.SelectJustification(99)

		if !mock.SelectJustificationCalled {
			t.Error("SelectJustification() should be marked as called")
		}
		if !errors.Is(err, printposition.ErrInvalidJustification) {
			t.Errorf("SelectJustification() error = %v, want %v",
				err, printposition.ErrInvalidJustification)
		}
	})

	t.Run("tracks SetPrintAreaPageMode calls", func(t *testing.T) {
		mock := NewMockCapability()

		result, _ := mock.SetPrintAreaPageMode(10, 20, 100, 200)

		if !mock.SetPrintAreaPageModeCalled {
			t.Error("SetPrintAreaPageMode() should be marked as called")
		}
		if mock.SetPrintAreaPageModeX != 10 || mock.SetPrintAreaPageModeY != 20 ||
			mock.SetPrintAreaPageModeWidth != 100 || mock.SetPrintAreaPageModeHeight != 200 {
			t.Errorf("SetPrintAreaPageMode() inputs = (%d,%d,%d,%d), want (10,20,100,200)",
				mock.SetPrintAreaPageModeX, mock.SetPrintAreaPageModeY,
				mock.SetPrintAreaPageModeWidth, mock.SetPrintAreaPageModeHeight)
		}
		expected := []byte{common.ESC, 'W', 10, 0, 20, 0, 100, 0, 200, 0}
		if !bytes.Equal(result, expected) {
			t.Errorf("SetPrintAreaPageMode() = %#v, want %#v", result, expected)
		}
	})

	t.Run("simulates SetPrintAreaPageMode width error", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SetPrintAreaPageModeError = printposition.ErrInvalidPrintAreaWidthSize

		_, err := mock.SetPrintAreaPageMode(10, 20, 0, 200)

		if !mock.SetPrintAreaPageModeCalled {
			t.Error("SetPrintAreaPageMode() should be marked as called")
		}
		if !errors.Is(err, printposition.ErrInvalidPrintAreaWidthSize) {
			t.Errorf("SetPrintAreaPageMode() error = %v, want %v",
				err, printposition.ErrInvalidPrintAreaWidthSize)
		}
	})
}
