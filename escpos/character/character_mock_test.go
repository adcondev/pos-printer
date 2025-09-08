package character_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
)

// Ensure MockCapability implements character.Capability
var _ character.Capability = (*MockCapability)(nil)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a utils double for character.Capability interface
type MockCapability struct {
	// SetRightSideCharacterSpacing tracking
	SetRightSideCharacterSpacingCalled bool
	SetRightSideCharacterSpacingInput  character.Spacing
	SetRightSideCharacterSpacingReturn []byte

	// SelectPrintModes tracking
	SelectPrintModesCalled bool
	SelectPrintModesInput  character.PrintMode
	SelectPrintModesReturn []byte

	// SetUnderlineMode tracking
	SetUnderlineModeCalled bool
	SetUnderlineModeInput  character.UnderlineMode
	SetUnderlineModeReturn []byte
	SetUnderlineModeError  error

	// SetEmphasizedMode tracking
	SetEmphasizedModeCalled bool
	SetEmphasizedModeInput  character.EmphasizedMode
	SetEmphasizedModeReturn []byte

	// SetDoubleStrikeMode tracking
	SetDoubleStrikeModeCalled bool
	SetDoubleStrikeModeInput  character.DoubleStrikeMode
	SetDoubleStrikeModeReturn []byte

	// SelectCharacterFont tracking
	SelectCharacterFontCalled bool
	SelectCharacterFontInput  character.FontType
	SelectCharacterFontReturn []byte
	SelectCharacterFontError  error

	// SelectInternationalCharacterSet tracking
	SelectInternationalCharacterSetCalled bool
	SelectInternationalCharacterSetInput  character.InternationalSet
	SelectInternationalCharacterSetReturn []byte
	SelectInternationalCharacterSetError  error

	// Set90DegreeClockwiseRotationMode tracking
	Set90DegreeClockwiseRotationModeCalled bool
	Set90DegreeClockwiseRotationModeInput  character.RotationMode
	Set90DegreeClockwiseRotationModeReturn []byte
	Set90DegreeClockwiseRotationModeError  error

	// SelectPrintColor tracking
	SelectPrintColorCalled bool
	SelectPrintColorInput  character.PrintColor
	SelectPrintColorReturn []byte
	SelectPrintColorError  error

	// SelectCharacterCodeTable tracking
	SelectCharacterCodeTableCalled bool
	SelectCharacterCodeTableInput  character.CodeTable
	SelectCharacterCodeTableReturn []byte
	SelectCharacterCodeTableError  error

	// SetUpsideDownMode tracking
	SetUpsideDownModeCalled bool
	SetUpsideDownModeInput  character.UpsideDownMode
	SetUpsideDownModeReturn []byte

	// SelectCharacterSize tracking
	SelectCharacterSizeCalled bool
	SelectCharacterSizeInput  character.Size
	SelectCharacterSizeReturn []byte

	// SetWhiteBlackReverseMode tracking
	SetWhiteBlackReverseModeCalled bool
	SetWhiteBlackReverseModeInput  character.ReverseMode
	SetWhiteBlackReverseModeReturn []byte

	// SetSmoothingMode tracking
	SetSmoothingModeCalled bool
	SetSmoothingModeInput  character.SmoothingMode
	SetSmoothingModeReturn []byte

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

func (m *MockCapability) SetRightSideCharacterSpacing(n character.Spacing) []byte {
	m.SetRightSideCharacterSpacingCalled = true
	m.SetRightSideCharacterSpacingInput = n
	m.CallCount["SetRightSideCharacterSpacing"]++

	if m.SetRightSideCharacterSpacingReturn != nil {
		return m.SetRightSideCharacterSpacingReturn
	}
	return []byte{common.ESC, common.SP, byte(n)}
}

func (m *MockCapability) SelectPrintModes(n character.PrintMode) []byte {
	m.SelectPrintModesCalled = true
	m.SelectPrintModesInput = n
	m.CallCount["SelectPrintModes"]++

	if m.SelectPrintModesReturn != nil {
		return m.SelectPrintModesReturn
	}
	return []byte{common.ESC, '!', byte(n)}
}

func (m *MockCapability) SetUnderlineMode(n character.UnderlineMode) ([]byte, error) {
	m.SetUnderlineModeCalled = true
	m.SetUnderlineModeInput = n
	m.CallCount["SetUnderlineMode"]++

	if m.SetUnderlineModeError != nil {
		return nil, m.SetUnderlineModeError
	}
	if m.SetUnderlineModeReturn != nil {
		return m.SetUnderlineModeReturn, nil
	}
	return []byte{common.ESC, '-', byte(n)}, nil
}

func (m *MockCapability) SetEmphasizedMode(n character.EmphasizedMode) []byte {
	m.SetEmphasizedModeCalled = true
	m.SetEmphasizedModeInput = n
	m.CallCount["SetEmphasizedMode"]++

	if m.SetEmphasizedModeReturn != nil {
		return m.SetEmphasizedModeReturn
	}
	return []byte{common.ESC, 'E', byte(n)}
}

func (m *MockCapability) SetDoubleStrikeMode(n character.DoubleStrikeMode) []byte {
	m.SetDoubleStrikeModeCalled = true
	m.SetDoubleStrikeModeInput = n
	m.CallCount["SetDoubleStrikeMode"]++

	if m.SetDoubleStrikeModeReturn != nil {
		return m.SetDoubleStrikeModeReturn
	}
	return []byte{common.ESC, 'G', byte(n)}
}

func (m *MockCapability) SelectCharacterFont(n character.FontType) ([]byte, error) {
	m.SelectCharacterFontCalled = true
	m.SelectCharacterFontInput = n
	m.CallCount["SelectCharacterFont"]++

	if m.SelectCharacterFontError != nil {
		return nil, m.SelectCharacterFontError
	}
	if m.SelectCharacterFontReturn != nil {
		return m.SelectCharacterFontReturn, nil
	}
	return []byte{common.ESC, 'M', byte(n)}, nil
}

func (m *MockCapability) SelectInternationalCharacterSet(n character.InternationalSet) ([]byte, error) {
	m.SelectInternationalCharacterSetCalled = true
	m.SelectInternationalCharacterSetInput = n
	m.CallCount["SelectInternationalCharacterSet"]++

	if m.SelectInternationalCharacterSetError != nil {
		return nil, m.SelectInternationalCharacterSetError
	}
	if m.SelectInternationalCharacterSetReturn != nil {
		return m.SelectInternationalCharacterSetReturn, nil
	}
	return []byte{common.ESC, 'R', byte(n)}, nil
}

func (m *MockCapability) Set90DegreeClockwiseRotationMode(n character.RotationMode) ([]byte, error) {
	m.Set90DegreeClockwiseRotationModeCalled = true
	m.Set90DegreeClockwiseRotationModeInput = n
	m.CallCount["Set90DegreeClockwiseRotationMode"]++

	if m.Set90DegreeClockwiseRotationModeError != nil {
		return nil, m.Set90DegreeClockwiseRotationModeError
	}
	if m.Set90DegreeClockwiseRotationModeReturn != nil {
		return m.Set90DegreeClockwiseRotationModeReturn, nil
	}
	return []byte{common.ESC, 'V', byte(n)}, nil
}

func (m *MockCapability) SelectPrintColor(n character.PrintColor) ([]byte, error) {
	m.SelectPrintColorCalled = true
	m.SelectPrintColorInput = n
	m.CallCount["SelectPrintColor"]++

	if m.SelectPrintColorError != nil {
		return nil, m.SelectPrintColorError
	}
	if m.SelectPrintColorReturn != nil {
		return m.SelectPrintColorReturn, nil
	}
	return []byte{common.ESC, 'r', byte(n)}, nil
}

func (m *MockCapability) SelectCharacterCodeTable(n character.CodeTable) ([]byte, error) {
	m.SelectCharacterCodeTableCalled = true
	m.SelectCharacterCodeTableInput = n
	m.CallCount["SelectCharacterCodeTable"]++

	if m.SelectCharacterCodeTableError != nil {
		return nil, m.SelectCharacterCodeTableError
	}
	if m.SelectCharacterCodeTableReturn != nil {
		return m.SelectCharacterCodeTableReturn, nil
	}
	return []byte{common.ESC, 't', byte(n)}, nil
}

func (m *MockCapability) SetUpsideDownMode(n character.UpsideDownMode) []byte {
	m.SetUpsideDownModeCalled = true
	m.SetUpsideDownModeInput = n
	m.CallCount["SetUpsideDownMode"]++

	if m.SetUpsideDownModeReturn != nil {
		return m.SetUpsideDownModeReturn
	}
	return []byte{common.ESC, '{', byte(n)}
}

func (m *MockCapability) SelectCharacterSize(n character.Size) []byte {
	m.SelectCharacterSizeCalled = true
	m.SelectCharacterSizeInput = n
	m.CallCount["SelectCharacterSize"]++

	if m.SelectCharacterSizeReturn != nil {
		return m.SelectCharacterSizeReturn
	}
	return []byte{common.GS, '!', byte(n)}
}

func (m *MockCapability) SetWhiteBlackReverseMode(n character.ReverseMode) []byte {
	m.SetWhiteBlackReverseModeCalled = true
	m.SetWhiteBlackReverseModeInput = n
	m.CallCount["SetWhiteBlackReverseMode"]++

	if m.SetWhiteBlackReverseModeReturn != nil {
		return m.SetWhiteBlackReverseModeReturn
	}
	return []byte{common.GS, 'B', byte(n)}
}

func (m *MockCapability) SetSmoothingMode(n character.SmoothingMode) []byte {
	m.SetSmoothingModeCalled = true
	m.SetSmoothingModeInput = n
	m.CallCount["SetSmoothingMode"]++

	if m.SetSmoothingModeReturn != nil {
		return m.SetSmoothingModeReturn
	}
	return []byte{common.GS, 'b', byte(n)}
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks SetRightSideCharacterSpacing calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SetRightSideCharacterSpacingReturn = []byte{0xAA, 0xBB, 0xCC}

		result := mock.SetRightSideCharacterSpacing(10)

		if !mock.SetRightSideCharacterSpacingCalled {
			t.Error("SetRightSideCharacterSpacing() should be marked as called")
		}
		if mock.SetRightSideCharacterSpacingInput != 10 {
			t.Errorf("SetRightSideCharacterSpacing() mode = %d, want 10",
				mock.SetRightSideCharacterSpacingInput)
		}
		if !bytes.Equal(result, []byte{0xAA, 0xBB, 0xCC}) {
			t.Errorf("SetRightSideCharacterSpacing() = %#v, want custom return", result)
		}
	})

	t.Run("tracks SelectPrintModes calls", func(t *testing.T) {
		mock := NewMockCapability()

		result := mock.SelectPrintModes(0x88)

		if !mock.SelectPrintModesCalled {
			t.Error("SelectPrintModes() should be marked as called")
		}
		if mock.SelectPrintModesInput != 0x88 {
			t.Errorf("SelectPrintModes() mode = %#x, want 0x88",
				mock.SelectPrintModesInput)
		}
		expected := []byte{common.ESC, '!', 0x88}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectPrintModes() = %#v, want %#v", result, expected)
		}
	})

	t.Run("simulates SetUnderlineMode error", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SetUnderlineModeError = character.ErrUnderlineMode

		_, err := mock.SetUnderlineMode(99)

		if !mock.SetUnderlineModeCalled {
			t.Error("SetUnderlineMode() should be marked as called")
		}
		if !errors.Is(err, character.ErrUnderlineMode) {
			t.Errorf("SetUnderlineMode() error = %v, want %v",
				err, character.ErrUnderlineMode)
		}
	})

	t.Run("tracks SelectCharacterSize calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SelectCharacterSizeReturn = []byte{0x1D, 0x21, 0x11}

		result := mock.SelectCharacterSize(0x11)

		if !mock.SelectCharacterSizeCalled {
			t.Error("SelectCharacterSize() should be marked as called")
		}
		if mock.SelectCharacterSizeInput != 0x11 {
			t.Errorf("SelectCharacterSize() mode = %#x, want 0x11",
				mock.SelectCharacterSizeInput)
		}
		if !bytes.Equal(result, []byte{0x1D, 0x21, 0x11}) {
			t.Errorf("SelectCharacterSize() = %#v, want custom return", result)
		}
	})
}
