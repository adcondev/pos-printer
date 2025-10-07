package barcode_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/utils/test"
)

// Ensure MockCapability implements barcode.Capability
var _ barcode.Capability = (*MockCapability)(nil)

// ============================================================================
// Mock Implementation
// ============================================================================

// MockCapability provides a test double for barcode.Capability interface
type MockCapability struct {
	// SelectHRICharacterPosition tracking
	SelectHRICharacterPositionCalled bool
	SelectHRICharacterPositionInput  barcode.HRIPosition
	SelectHRICharacterPositionReturn []byte
	SelectHRICharacterPositionError  error

	// SelectFontForHRI tracking
	SelectFontForHRICalled bool
	SelectFontForHRIInput  barcode.HRIFont
	SelectFontForHRIReturn []byte
	SelectFontForHRIError  error

	// SetBarcodeHeight tracking
	SetBarcodeHeightCalled bool
	SetBarcodeHeightInput  barcode.Height
	SetBarcodeHeightReturn []byte
	SetBarcodeHeightError  error

	// SetBarcodeWidth tracking
	SetBarcodeWidthCalled bool
	SetBarcodeWidthInput  barcode.Width
	SetBarcodeWidthReturn []byte
	SetBarcodeWidthError  error

	// PrintBarcode tracking
	PrintBarcodeCalled    bool
	PrintBarcodeSymbology barcode.Symbology
	PrintBarcodeData      []byte
	PrintBarcodeReturn    []byte
	PrintBarcodeError     error

	// PrintBarcodeWithCodeSet tracking
	PrintBarcodeWithCodeSetCalled    bool
	PrintBarcodeWithCodeSetSymbology barcode.Symbology
	PrintBarcodeWithCodeSetCodeSet   barcode.Code128Set
	PrintBarcodeWithCodeSetData      []byte
	PrintBarcodeWithCodeSetReturn    []byte
	PrintBarcodeWithCodeSetError     error

	// Call counting
	CallCount map[string]int
}

// NewMockCapability creates a new mock barcode capability
func NewMockCapability() *MockCapability {
	return &MockCapability{
		CallCount: make(map[string]int),
	}
}

// Reset clears all mock state
func (m *MockCapability) Reset() {
	*m = *NewMockCapability()
}

// GetCallCount returns the number of times a method was called
func (m *MockCapability) GetCallCount(method string) int {
	return m.CallCount[method]
}

func (m *MockCapability) SelectHRICharacterPosition(position barcode.HRIPosition) ([]byte, error) {
	m.SelectHRICharacterPositionCalled = true
	m.SelectHRICharacterPositionInput = position
	m.CallCount["SelectHRICharacterPosition"]++

	if m.SelectHRICharacterPositionError != nil {
		return nil, m.SelectHRICharacterPositionError
	}
	if m.SelectHRICharacterPositionReturn != nil {
		return m.SelectHRICharacterPositionReturn, nil
	}
	return []byte{common.GS, 'H', byte(position)}, nil
}

func (m *MockCapability) SelectFontForHRI(font barcode.HRIFont) ([]byte, error) {
	m.SelectFontForHRICalled = true
	m.SelectFontForHRIInput = font
	m.CallCount["SelectFontForHRI"]++

	if m.SelectFontForHRIError != nil {
		return nil, m.SelectFontForHRIError
	}
	if m.SelectFontForHRIReturn != nil {
		return m.SelectFontForHRIReturn, nil
	}
	return []byte{common.GS, 'f', byte(font)}, nil
}

func (m *MockCapability) SetBarcodeHeight(height barcode.Height) ([]byte, error) {
	m.SetBarcodeHeightCalled = true
	m.SetBarcodeHeightInput = height
	m.CallCount["SetBarcodeHeight"]++

	if m.SetBarcodeHeightError != nil {
		return nil, m.SetBarcodeHeightError
	}
	if m.SetBarcodeHeightReturn != nil {
		return m.SetBarcodeHeightReturn, nil
	}
	return []byte{common.GS, 'h', byte(height)}, nil
}

func (m *MockCapability) SetBarcodeWidth(width barcode.Width) ([]byte, error) {
	m.SetBarcodeWidthCalled = true
	m.SetBarcodeWidthInput = width
	m.CallCount["SetBarcodeWidth"]++

	if m.SetBarcodeWidthError != nil {
		return nil, m.SetBarcodeWidthError
	}
	if m.SetBarcodeWidthReturn != nil {
		return m.SetBarcodeWidthReturn, nil
	}
	return []byte{common.GS, 'w', byte(width)}, nil
}

func (m *MockCapability) PrintBarcode(symbology barcode.Symbology, data []byte) ([]byte, error) {
	m.PrintBarcodeCalled = true
	m.PrintBarcodeSymbology = symbology
	m.PrintBarcodeData = data
	m.CallCount["PrintBarcode"]++

	if m.PrintBarcodeError != nil {
		return nil, m.PrintBarcodeError
	}
	if m.PrintBarcodeReturn != nil {
		return m.PrintBarcodeReturn, nil
	}

	// Simple mock implementation
	if symbology <= barcode.CODABAR {
		// Function A
		cmd := []byte{common.GS, 'k', byte(symbology)}
		cmd = append(cmd, data...)
		cmd = append(cmd, common.NUL)
		return cmd, nil
	} else {
		// Function B
		cmd := []byte{common.GS, 'k', byte(symbology), byte(len(data))}
		cmd = append(cmd, data...)
		return cmd, nil
	}
}

func (m *MockCapability) PrintBarcodeWithCodeSet(symbology barcode.Symbology, codeSet barcode.Code128Set, data []byte) ([]byte, error) {
	m.PrintBarcodeWithCodeSetCalled = true
	m.PrintBarcodeWithCodeSetSymbology = symbology
	m.PrintBarcodeWithCodeSetCodeSet = codeSet
	m.PrintBarcodeWithCodeSetData = data
	m.CallCount["PrintBarcodeWithCodeSet"]++

	if m.PrintBarcodeWithCodeSetError != nil {
		return nil, m.PrintBarcodeWithCodeSetError
	}
	if m.PrintBarcodeWithCodeSetReturn != nil {
		return m.PrintBarcodeWithCodeSetReturn, nil
	}

	// Simple mock implementation
	prefixedData := append([]byte{'{', byte(codeSet)}, data...)
	cmd := []byte{common.GS, 'k', byte(symbology), byte(len(prefixedData))}
	cmd = append(cmd, prefixedData...)
	return cmd, nil
}

// ============================================================================
// Mock Tests
// ============================================================================

func TestMockCapability_BehaviorTracking(t *testing.T) {
	t.Run("tracks SelectHRICharacterPosition calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SelectHRICharacterPositionReturn = []byte{0xAA, 0xBB, 0xCC}

		result, err := mock.SelectHRICharacterPosition(barcode.HRIBelow)

		if !mock.SelectHRICharacterPositionCalled {
			t.Error("SelectHRICharacterPosition() should be marked as called")
		}
		if mock.SelectHRICharacterPositionInput != barcode.HRIBelow {
			t.Errorf("SelectHRICharacterPosition() input = %v, want HRIBelow",
				mock.SelectHRICharacterPositionInput)
		}
		if err != nil {
			t.Errorf("SelectHRICharacterPosition() unexpected error: %v", err)
		}
		if !bytes.Equal(result, []byte{0xAA, 0xBB, 0xCC}) {
			t.Errorf("SelectHRICharacterPosition() = %#v, want custom return", result)
		}
		if mock.GetCallCount("SelectHRICharacterPosition") != 1 {
			t.Errorf("Call count = %d, want 1", mock.GetCallCount("SelectHRICharacterPosition"))
		}
	})

	t.Run("simulates SetBarcodeHeight error", func(t *testing.T) {
		mock := NewMockCapability()
		mock.SetBarcodeHeightError = barcode.ErrHeight

		_, err := mock.SetBarcodeHeight(0)

		if !mock.SetBarcodeHeightCalled {
			t.Error("SetBarcodeHeight() should be marked as called")
		}
		if !errors.Is(err, barcode.ErrHeight) {
			t.Errorf("SetBarcodeHeight() error = %v, want %v", err, barcode.ErrHeight)
		}
	})

	t.Run("tracks PrintBarcode calls", func(t *testing.T) {
		mock := NewMockCapability()

		result, err := mock.PrintBarcode(barcode.CODE39, []byte("TEST123"))

		if !mock.PrintBarcodeCalled {
			t.Error("PrintBarcode() should be marked as called")
		}
		if mock.PrintBarcodeSymbology != barcode.CODE39 {
			t.Errorf("PrintBarcode() symbology = %v, want CODE39", mock.PrintBarcodeSymbology)
		}
		if !bytes.Equal(mock.PrintBarcodeData, []byte("TEST123")) {
			t.Errorf("PrintBarcode() data = %q, want %q", mock.PrintBarcodeData, "TEST123")
		}
		if err != nil {
			t.Errorf("PrintBarcode() unexpected error: %v", err)
		}
		// Default mock implementation for Function A
		expected := append([]byte{common.GS, 'k', byte(barcode.CODE39)}, append([]byte("TEST123"), common.NUL)...)
		if !bytes.Equal(result, expected) {
			t.Errorf("PrintBarcode() = %#v, want %#v", result, expected)
		}
	})

	t.Run("tracks PrintBarcodeWithCodeSet calls", func(t *testing.T) {
		mock := NewMockCapability()
		mock.PrintBarcodeWithCodeSetReturn = []byte{0x11, 0x22, 0x33}

		result, err := mock.PrintBarcodeWithCodeSet(
			barcode.CODE128,
			barcode.Code128SetB,
			[]byte("Hello"),
		)

		if !mock.PrintBarcodeWithCodeSetCalled {
			t.Error("PrintBarcodeWithCodeSet() should be marked as called")
		}
		if mock.PrintBarcodeWithCodeSetSymbology != barcode.CODE128 {
			t.Errorf("Symbology = %v, want CODE128", mock.PrintBarcodeWithCodeSetSymbology)
		}
		if mock.PrintBarcodeWithCodeSetCodeSet != barcode.Code128SetB {
			t.Errorf("CodeSet = %v, want Code128SetB", mock.PrintBarcodeWithCodeSetCodeSet)
		}
		if !bytes.Equal(mock.PrintBarcodeWithCodeSetData, []byte("Hello")) {
			t.Errorf("Data = %q, want %q", mock.PrintBarcodeWithCodeSetData, "Hello")
		}
		if err != nil {
			t.Errorf("PrintBarcodeWithCodeSet() unexpected error: %v", err)
		}
		if !bytes.Equal(result, []byte{0x11, 0x22, 0x33}) {
			t.Errorf("PrintBarcodeWithCodeSet() = %#v, want custom return", result)
		}
	})

	t.Run("returns default behavior when no return configured", func(t *testing.T) {
		mock := NewMockCapability()

		result, _ := mock.SelectFontForHRI(barcode.HRIFontB)

		expected := []byte{common.GS, 'f', byte(barcode.HRIFontB)}
		if !bytes.Equal(result, expected) {
			t.Errorf("SelectFontForHRI() = %#v, want %#v", result, expected)
		}
	})

	t.Run("reset clears all state", func(t *testing.T) {
		mock := NewMockCapability()

		// Set some state
		_, _ = mock.SelectHRICharacterPosition(barcode.HRIBelow)
		_, _ = mock.SetBarcodeHeight(100)
		_, _ = mock.PrintBarcode(barcode.CODE39, []byte("TEST"))

		// Reset
		mock.Reset()

		// Verify state is cleared
		if mock.SelectHRICharacterPositionCalled {
			t.Error("SelectHRICharacterPositionCalled should be false after reset")
		}
		if mock.SetBarcodeHeightCalled {
			t.Error("SetBarcodeHeightCalled should be false after reset")
		}
		if mock.PrintBarcodeCalled {
			t.Error("PrintBarcodeCalled should be false after reset")
		}
		if mock.GetCallCount("SelectHRICharacterPosition") != 0 {
			t.Error("Call counts should be 0 after reset")
		}
	})
}

func TestMockCapability_ErrorSimulation(t *testing.T) {
	t.Run("simulates multiple error conditions", func(t *testing.T) {
		mock := NewMockCapability()

		// Configure errors
		mock.SelectHRICharacterPositionError = barcode.ErrHRIPosition
		mock.SelectFontForHRIError = barcode.ErrHRIFont
		mock.SetBarcodeHeightError = barcode.ErrHeight
		mock.SetBarcodeWidthError = barcode.ErrWidth
		mock.PrintBarcodeError = barcode.ErrSymbology
		mock.PrintBarcodeWithCodeSetError = barcode.ErrCode128NoCodeSet

		// Test each error
		if _, err := mock.SelectHRICharacterPosition(0); !errors.Is(err, barcode.ErrHRIPosition) {
			t.Error("Should return ErrHRIPosition")
		}
		if _, err := mock.SelectFontForHRI(0); !errors.Is(err, barcode.ErrHRIFont) {
			t.Error("Should return ErrHRIFont")
		}
		if _, err := mock.SetBarcodeHeight(0); !errors.Is(err, barcode.ErrHeight) {
			t.Error("Should return ErrHeight")
		}
		if _, err := mock.SetBarcodeWidth(0); !errors.Is(err, barcode.ErrWidth) {
			t.Error("Should return ErrWidth")
		}
		if _, err := mock.PrintBarcode(0, nil); !errors.Is(err, barcode.ErrSymbology) {
			t.Error("Should return ErrSymbology")
		}
		if _, err := mock.PrintBarcodeWithCodeSet(0, 0, nil); !errors.Is(err, barcode.ErrCode128NoCodeSet) {
			t.Error("Should return ErrCode128NoCodeSet")
		}
	})

	t.Run("allows selective error configuration", func(t *testing.T) {
		mock := NewMockCapability()

		// Only configure one error
		mock.PrintBarcodeError = barcode.ErrDataTooLong

		// Other methods should work normally
		if _, err := mock.SetBarcodeHeight(100); err != nil {
			t.Errorf("SetBarcodeHeight should not error: %v", err)
		}

		// Configured method should error
		if _, err := mock.PrintBarcode(barcode.CODE39, []byte("TEST")); !errors.Is(err, barcode.ErrDataTooLong) {
			t.Error("PrintBarcode should return configured error")
		}
	})
}

func TestMockCapability_CustomReturns(t *testing.T) {
	t.Run("returns custom byte sequences", func(t *testing.T) {
		mock := NewMockCapability()

		// Configure custom returns
		mock.SelectHRICharacterPositionReturn = []byte{0x01, 0x02, 0x03}
		mock.SelectFontForHRIReturn = []byte{0x04, 0x05, 0x06}
		mock.SetBarcodeHeightReturn = []byte{0x07, 0x08, 0x09}
		mock.SetBarcodeWidthReturn = []byte{0x0A, 0x0B, 0x0C}
		mock.PrintBarcodeReturn = []byte{0x0D, 0x0E, 0x0F}
		mock.PrintBarcodeWithCodeSetReturn = []byte{0x10, 0x11, 0x12}

		// Verify each custom return
		result1, _ := mock.SelectHRICharacterPosition(0)
		if !bytes.Equal(result1, []byte{0x01, 0x02, 0x03}) {
			t.Errorf("SelectHRICharacterPosition custom return not used")
		}

		result2, _ := mock.SelectFontForHRI(0)
		if !bytes.Equal(result2, []byte{0x04, 0x05, 0x06}) {
			t.Errorf("SelectFontForHRI custom return not used")
		}

		result3, _ := mock.SetBarcodeHeight(0)
		if !bytes.Equal(result3, []byte{0x07, 0x08, 0x09}) {
			t.Errorf("SetBarcodeHeight custom return not used")
		}

		result4, _ := mock.SetBarcodeWidth(0)
		if !bytes.Equal(result4, []byte{0x0A, 0x0B, 0x0C}) {
			t.Errorf("SetBarcodeWidth custom return not used")
		}

		result5, _ := mock.PrintBarcode(0, nil)
		if !bytes.Equal(result5, []byte{0x0D, 0x0E, 0x0F}) {
			t.Errorf("PrintBarcode custom return not used")
		}

		result6, _ := mock.PrintBarcodeWithCodeSet(0, 0, nil)
		if !bytes.Equal(result6, []byte{0x10, 0x11, 0x12}) {
			t.Errorf("PrintBarcodeWithCodeSet custom return not used")
		}
	})
}

func TestMockCapability_CallCounting(t *testing.T) {
	t.Run("tracks call counts properly", func(t *testing.T) {
		mock := NewMockCapability()

		// Call methods multiple times
		_, _ = mock.SelectHRICharacterPosition(0)
		_, _ = mock.SelectHRICharacterPosition(1)
		_, _ = mock.SelectHRICharacterPosition(2)

		_, _ = mock.SetBarcodeHeight(100)
		_, _ = mock.SetBarcodeHeight(200)

		_, _ = mock.PrintBarcode(barcode.CODE39, []byte("TEST"))

		// Check counts
		if count := mock.GetCallCount("SelectHRICharacterPosition"); count != 3 {
			t.Errorf("SelectHRICharacterPosition count = %d, want 3", count)
		}
		if count := mock.GetCallCount("SetBarcodeHeight"); count != 2 {
			t.Errorf("SetBarcodeHeight count = %d, want 2", count)
		}
		if count := mock.GetCallCount("PrintBarcode"); count != 1 {
			t.Errorf("PrintBarcode count = %d, want 1", count)
		}
		if count := mock.GetCallCount("SelectFontForHRI"); count != 0 {
			t.Errorf("SelectFontForHRI count = %d, want 0", count)
		}
	})

	t.Run("counts are reset properly", func(t *testing.T) {
		mock := NewMockCapability()

		// Call methods
		_, _ = mock.SelectHRICharacterPosition(0)
		_, _ = mock.SetBarcodeHeight(100)
		_, _ = mock.PrintBarcode(barcode.CODE39, []byte("TEST"))

		// Reset
		mock.Reset()

		// All counts should be zero
		if count := mock.GetCallCount("SelectHRICharacterPosition"); count != 0 {
			t.Errorf("After reset, SelectHRICharacterPosition count = %d, want 0", count)
		}
		if count := mock.GetCallCount("SetBarcodeHeight"); count != 0 {
			t.Errorf("After reset, SetBarcodeHeight count = %d, want 0", count)
		}
		if count := mock.GetCallCount("PrintBarcode"); count != 0 {
			t.Errorf("After reset, PrintBarcode count = %d, want 0", count)
		}
	})
}

func TestMockCapability_FunctionTypeHandling(t *testing.T) {
	t.Run("handles Function A barcodes correctly", func(t *testing.T) {
		mock := NewMockCapability()

		result, _ := mock.PrintBarcode(barcode.CODE39, []byte("CODE39"))

		// Function A format ends with NUL terminator
		test.AssertHasNullTerminator(t, result, "Function A barcode format")

		// Verify last byte specifically
		if result[len(result)-1] != common.NUL {
			t.Error("Function A barcode should end with NUL terminator")
		}
	})

	t.Run("handles Function B barcodes correctly", func(t *testing.T) {
		mock := NewMockCapability()

		data := []byte("CODE128")
		result, _ := mock.PrintBarcode(barcode.CODE128, data)

		// Function B format includes length byte
		if result[3] != byte(len(data)) {
			t.Errorf("Function B barcode length byte = %d, want %d", result[3], len(data))
		}

		// Verify data length is within valid range
		test.AssertValidLength(t, data, 1, 255, "Barcode data length")
	})
}
