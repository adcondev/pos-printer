package barcode_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/common"
)

// Ensure FakeCapability implements barcode.Capability
var _ barcode.Capability = (*FakeCapability)(nil)

// ============================================================================
// Fake Implementation
// ============================================================================

// FakeCapability simulates barcode printing with state tracking
type FakeCapability struct {
	buffer          []byte
	hriPosition     barcode.HRIPosition
	hriFont         barcode.HRIFont
	barcodeHeight   barcode.Height
	barcodeWidth    barcode.Width
	printedBarcodes []BarcodeRecord
	lastCommand     string
	commandHistory  []string
	commandsCount   map[string]int
}

// BarcodeRecord stores information about printed barcodes
type BarcodeRecord struct {
	Symbology barcode.Symbology
	Data      []byte
	CodeSet   barcode.Code128Set
	Height    barcode.Height
	Width     barcode.Width
	HRIPos    barcode.HRIPosition
	HRIFont   barcode.HRIFont
}

// NewFakeCapability creates a new fake barcode capability
func NewFakeCapability() *FakeCapability {
	return &FakeCapability{
		buffer:          make([]byte, 0),
		hriPosition:     barcode.HRINotPrinted,
		hriFont:         barcode.HRIFontA,
		barcodeHeight:   barcode.DefaultHeight,
		barcodeWidth:    barcode.DefaultWidth,
		printedBarcodes: make([]BarcodeRecord, 0),
		commandHistory:  make([]string, 0),
		commandsCount:   make(map[string]int),
	}
}

// SelectHRICharacterPosition simulates setting HRI position
func (f *FakeCapability) SelectHRICharacterPosition(position barcode.HRIPosition) ([]byte, error) {
	// Validate position
	switch position {
	case 0, 1, 2, 3, '0', '1', '2', '3':
		// Valid values
	default:
		return nil, barcode.ErrHRIPosition
	}

	f.hriPosition = position
	f.lastCommand = "SelectHRICharacterPosition"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	cmd := []byte{common.GS, 'H', byte(position)}
	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// SelectFontForHRI simulates setting HRI font
func (f *FakeCapability) SelectFontForHRI(font barcode.HRIFont) ([]byte, error) {
	// Validate font
	switch font {
	case 0, 1, 2, 3, 4, '0', '1', '2', '3', '4', 97, 98:
		// Valid values
	default:
		return nil, barcode.ErrHRIFont
	}

	f.hriFont = font
	f.lastCommand = "SelectFontForHRI"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	cmd := []byte{common.GS, 'f', byte(font)}
	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// SetBarcodeHeight simulates setting barcode height
func (f *FakeCapability) SetBarcodeHeight(height barcode.Height) ([]byte, error) {
	if height < barcode.MinHeight || height > barcode.MaxHeight {
		return nil, barcode.ErrHeight
	}

	f.barcodeHeight = height
	f.lastCommand = "SetBarcodeHeight"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	cmd := []byte{common.GS, 'h', byte(height)}
	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// SetBarcodeWidth simulates setting barcode width
func (f *FakeCapability) SetBarcodeWidth(width barcode.Width) ([]byte, error) {
	// Validate standard and extended ranges
	if (width < barcode.MinWidth || width > barcode.MaxWidth) && (width < barcode.ExtendedMinWidth || width > barcode.ExtendedMaxWidth) {
		return nil, barcode.ErrWidth
	}

	f.barcodeWidth = width
	f.lastCommand = "SetBarcodeWidth"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	cmd := []byte{common.GS, 'w', byte(width)}
	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// PrintBarcode simulates printing a barcode
func (f *FakeCapability) PrintBarcode(symbology barcode.Symbology, data []byte) ([]byte, error) {
	// Validate data
	if len(data) == 0 {
		return nil, barcode.ErrDataTooShort
	}
	if len(data) > 255 {
		return nil, barcode.ErrDataTooLong
	}

	// Special validations
	switch symbology {
	case barcode.ITF, barcode.ITFB:
		if len(data)%2 != 0 {
			return nil, barcode.ErrOddITFLength
		}
	case barcode.CODE128, barcode.GS1128:
		return nil, barcode.ErrCode128NoCodeSet
	}

	// Record the barcode
	record := BarcodeRecord{
		Symbology: symbology,
		Data:      append([]byte(nil), data...), // Copy data
		Height:    f.barcodeHeight,
		Width:     f.barcodeWidth,
		HRIPos:    f.hriPosition,
		HRIFont:   f.hriFont,
	}
	f.printedBarcodes = append(f.printedBarcodes, record)

	f.lastCommand = "PrintBarcode"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	// Build command
	var cmd []byte
	if symbology <= barcode.CODABAR {
		// Function A (NUL-terminated)
		cmd = []byte{common.GS, 'k', byte(symbology)}
		cmd = append(cmd, data...)
		cmd = append(cmd, common.NUL)
	} else if symbology >= barcode.UPCAB && symbology <= barcode.CODE128Auto {
		// Function B (length-prefixed)
		cmd = []byte{common.GS, 'k', byte(symbology), byte(len(data))}
		cmd = append(cmd, data...)
	} else {
		return nil, barcode.ErrSymbology
	}

	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// PrintBarcodeWithCodeSet simulates printing a CODE128 barcode with code set
func (f *FakeCapability) PrintBarcodeWithCodeSet(symbology barcode.Symbology, codeSet barcode.Code128Set, data []byte) ([]byte, error) {
	// Validate symbology
	if symbology != barcode.CODE128 && symbology != barcode.GS1128 {
		return nil, barcode.ErrSymbology
	}

	// Validate code set
	if codeSet < barcode.Code128SetA || codeSet > barcode.Code128SetC {
		return nil, barcode.ErrCode128Set
	}

	// Validate data
	if len(data) == 0 {
		return nil, barcode.ErrDataTooShort
	}
	if len(data) > 253 { // 255 - 2 for code set prefix
		return nil, barcode.ErrDataTooLong
	}

	// Record the barcode
	record := BarcodeRecord{
		Symbology: symbology,
		Data:      append([]byte(nil), data...), // Copy data
		CodeSet:   codeSet,
		Height:    f.barcodeHeight,
		Width:     f.barcodeWidth,
		HRIPos:    f.hriPosition,
		HRIFont:   f.hriFont,
	}
	f.printedBarcodes = append(f.printedBarcodes, record)

	f.lastCommand = "PrintBarcodeWithCodeSet"
	f.commandHistory = append(f.commandHistory, f.lastCommand)
	f.commandsCount[f.lastCommand]++

	// Build command with code set prefix
	prefixedData := make([]byte, 0, len(data)+2)
	prefixedData = append(prefixedData, '{', byte(codeSet))
	prefixedData = append(prefixedData, data...)

	cmd := []byte{common.GS, 'k', byte(symbology), byte(len(prefixedData))}
	cmd = append(cmd, prefixedData...)

	f.buffer = append(f.buffer, cmd...)
	return cmd, nil
}

// GetBuffer returns the accumulated command buffer
func (f *FakeCapability) GetBuffer() []byte {
	return f.buffer
}

// GetPrintedBarcodes returns all printed barcode records
func (f *FakeCapability) GetPrintedBarcodes() []BarcodeRecord {
	return f.printedBarcodes
}

// GetLastCommand returns the name of the last command executed
func (f *FakeCapability) GetLastCommand() string {
	return f.lastCommand
}

// GetCommandHistory returns the complete command history
func (f *FakeCapability) GetCommandHistory() []string {
	return f.commandHistory
}

// GetCommandCount returns how many times a command was called
func (f *FakeCapability) GetCommandCount(command string) int {
	return f.commandsCount[command]
}

// Reset clears all state
func (f *FakeCapability) Reset() {
	f.buffer = make([]byte, 0)
	f.hriPosition = barcode.HRINotPrinted
	f.hriFont = barcode.HRIFontA
	f.barcodeHeight = barcode.DefaultHeight
	f.barcodeWidth = barcode.DefaultWidth
	f.printedBarcodes = make([]BarcodeRecord, 0)
	f.lastCommand = ""
	f.commandHistory = make([]string, 0)
	f.commandsCount = make(map[string]int)
}

// ============================================================================
// Fake Tests
// ============================================================================

func TestFakeCapability_StateTracking(t *testing.T) {
	t.Run("tracks HRI configuration state", func(t *testing.T) {
		fake := NewFakeCapability()

		// Set HRI position
		_, err := fake.SelectHRICharacterPosition(barcode.HRIBelow)
		if err != nil {
			t.Fatalf("SelectHRICharacterPosition: %v", err)
		}
		if fake.hriPosition != barcode.HRIBelow {
			t.Errorf("HRI position = %v, want HRIBelow", fake.hriPosition)
		}

		// Set HRI font
		_, err = fake.SelectFontForHRI(barcode.HRIFontC)
		if err != nil {
			t.Fatalf("SelectFontForHRI: %v", err)
		}
		if fake.hriFont != barcode.HRIFontC {
			t.Errorf("HRI font = %v, want HRIFontC", fake.hriFont)
		}
	})

	t.Run("tracks barcode dimensions", func(t *testing.T) {
		fake := NewFakeCapability()

		// Set height
		_, err := fake.SetBarcodeHeight(120)
		if err != nil {
			t.Fatalf("SetBarcodeHeight: %v", err)
		}
		if fake.barcodeHeight != 120 {
			t.Errorf("Barcode height = %d, want 120", fake.barcodeHeight)
		}

		// Set width
		_, err = fake.SetBarcodeWidth(4)
		if err != nil {
			t.Fatalf("SetBarcodeWidth: %v", err)
		}
		if fake.barcodeWidth != 4 {
			t.Errorf("Barcode width = %d, want 4", fake.barcodeWidth)
		}
	})

	t.Run("records printed barcodes", func(t *testing.T) {
		fake := NewFakeCapability()

		// Configure settings
		_, _ = fake.SelectHRICharacterPosition(barcode.HRIBelow)
		_, _ = fake.SetBarcodeHeight(80)
		_, _ = fake.SetBarcodeWidth(3)

		// Print barcodes
		_, err := fake.PrintBarcode(barcode.CODE39, []byte("TEST123"))
		if err != nil {
			t.Fatalf("PrintBarcode: %v", err)
		}

		_, err = fake.PrintBarcodeWithCodeSet(barcode.CODE128, barcode.Code128SetB, []byte("Hello"))
		if err != nil {
			t.Fatalf("PrintBarcodeWithCodeSet: %v", err)
		}

		// Verify records
		records := fake.GetPrintedBarcodes()
		if len(records) != 2 {
			t.Fatalf("Printed barcodes = %d, want 2", len(records))
		}

		// Check first barcode
		first := records[0]
		if first.Symbology != barcode.CODE39 {
			t.Errorf("First symbology = %v, want CODE39", first.Symbology)
		}
		if !bytes.Equal(first.Data, []byte("TEST123")) {
			t.Errorf("First data = %q, want TEST123", first.Data)
		}
		if first.Height != 80 {
			t.Errorf("First height = %d, want 80", first.Height)
		}
		if first.Width != 3 {
			t.Errorf("First width = %d, want 3", first.Width)
		}

		// Check second barcode
		second := records[1]
		if second.Symbology != barcode.CODE128 {
			t.Errorf("Second symbology = %v, want CODE128", second.Symbology)
		}
		if second.CodeSet != barcode.Code128SetB {
			t.Errorf("Second code set = %v, want Code128SetB", second.CodeSet)
		}
		if !bytes.Equal(second.Data, []byte("Hello")) {
			t.Errorf("Second data = %q, want Hello", second.Data)
		}
	})

	t.Run("accumulates buffer correctly", func(t *testing.T) {
		fake := NewFakeCapability()

		// Execute multiple commands
		_, _ = fake.SelectHRICharacterPosition(barcode.HRIBelow)
		_, _ = fake.SelectFontForHRI(barcode.HRIFontB)
		_, _ = fake.SetBarcodeHeight(100)
		_, _ = fake.SetBarcodeWidth(3)
		_, _ = fake.PrintBarcode(barcode.UPCA, []byte("12345678901"))

		buffer := fake.GetBuffer()

		// Verify buffer contains all commands in order
		expectedStart := []byte{
			common.GS, 'H', 2, // HRI position
			common.GS, 'f', 1, // HRI font
			common.GS, 'h', 100, // Height
			common.GS, 'w', 3, // Width
		}
		if !bytes.HasPrefix(buffer, expectedStart) {
			t.Errorf("Buffer doesn't start with expected configuration commands")
		}

		// Verify barcode data is present
		if !bytes.Contains(buffer, []byte("12345678901")) {
			t.Error("Buffer should contain barcode data")
		}
	})

	t.Run("tracks command history", func(t *testing.T) {
		fake := NewFakeCapability()

		// Execute commands in specific order
		_, _ = fake.SelectHRICharacterPosition(barcode.HRIBelow)
		_, _ = fake.SetBarcodeHeight(80)
		_, _ = fake.PrintBarcode(barcode.CODE39, []byte("ABC"))
		_, _ = fake.SetBarcodeHeight(100)
		_, _ = fake.PrintBarcode(barcode.CODE39, []byte("DEF"))

		history := fake.GetCommandHistory()
		expected := []string{
			"SelectHRICharacterPosition",
			"SetBarcodeHeight",
			"PrintBarcode",
			"SetBarcodeHeight",
			"PrintBarcode",
		}

		if len(history) != len(expected) {
			t.Fatalf("History length = %d, want %d", len(history), len(expected))
		}

		for i, cmd := range expected {
			if history[i] != cmd {
				t.Errorf("History[%d] = %s, want %s", i, history[i], cmd)
			}
		}
	})

	t.Run("counts command calls", func(t *testing.T) {
		fake := NewFakeCapability()

		// Execute commands multiple times
		_, _ = fake.SetBarcodeHeight(50)
		_, _ = fake.SetBarcodeHeight(100)
		_, _ = fake.SetBarcodeHeight(150)
		_, _ = fake.PrintBarcode(barcode.CODE39, []byte("A"))
		_, _ = fake.PrintBarcode(barcode.CODE39, []byte("B"))

		if count := fake.GetCommandCount("SetBarcodeHeight"); count != 3 {
			t.Errorf("SetBarcodeHeight count = %d, want 3", count)
		}
		if count := fake.GetCommandCount("PrintBarcode"); count != 2 {
			t.Errorf("PrintBarcode count = %d, want 2", count)
		}
		if count := fake.GetCommandCount("SetBarcodeWidth"); count != 0 {
			t.Errorf("SetBarcodeWidth count = %d, want 0", count)
		}
	})

	t.Run("validates errors correctly", func(t *testing.T) {
		fake := NewFakeCapability()

		// Invalid HRI position
		_, err := fake.SelectHRICharacterPosition(99)
		if err == nil {
			t.Error("Should return error for invalid HRI position")
		}

		// Invalid HRI font
		_, err = fake.SelectFontForHRI(200)
		if err == nil {
			t.Error("Should return error for invalid HRI font")
		}

		// Invalid height
		_, err = fake.SetBarcodeHeight(0)
		if err == nil {
			t.Error("Should return error for invalid height")
		}

		// Invalid width
		_, err = fake.SetBarcodeWidth(100)
		if err == nil {
			t.Error("Should return error for invalid width")
		}

		// ITF with odd length
		_, err = fake.PrintBarcode(barcode.ITF, []byte("12345"))
		if err == nil {
			t.Error("Should return error for ITF with odd length")
		}

		// CODE128 without code set
		_, err = fake.PrintBarcode(barcode.CODE128, []byte("Hello"))
		if err == nil {
			t.Error("Should return error for CODE128 without code set")
		}
	})

	t.Run("reset clears all state", func(t *testing.T) {
		fake := NewFakeCapability()

		// Set up state
		_, _ = fake.SelectHRICharacterPosition(barcode.HRIBelow)
		_, _ = fake.SetBarcodeHeight(100)
		_, _ = fake.PrintBarcode(barcode.CODE39, []byte("TEST"))

		// Reset
		fake.Reset()

		// Verify defaults restored
		if fake.hriPosition != barcode.HRINotPrinted {
			t.Error("HRI position should be reset to default")
		}
		if fake.hriFont != barcode.HRIFontA {
			t.Error("HRI font should be reset to default")
		}
		if fake.barcodeHeight != barcode.DefaultHeight {
			t.Error("Height should be reset to default")
		}
		if fake.barcodeWidth != barcode.DefaultWidth {
			t.Error("Width should be reset to default")
		}
		if len(fake.GetBuffer()) != 0 {
			t.Error("Buffer should be empty after reset")
		}
		if len(fake.GetPrintedBarcodes()) != 0 {
			t.Error("Printed barcodes should be empty after reset")
		}
		if len(fake.GetCommandHistory()) != 0 {
			t.Error("Command history should be empty after reset")
		}
		if fake.GetCommandCount("PrintBarcode") != 0 {
			t.Error("Command counts should be reset")
		}
	})
}

func TestFakeCapability_ComplexScenarios(t *testing.T) {
	t.Run("simulates complete barcode workflow", func(t *testing.T) {
		fake := NewFakeCapability()

		// Simulate printing multiple barcodes with different settings
		workflows := []struct {
			name      string
			setup     func()
			symbology barcode.Symbology
			data      []byte
		}{
			{
				name: "UPC with HRI below",
				setup: func() {
					_, _ = fake.SelectHRICharacterPosition(barcode.HRIBelow)
					_, _ = fake.SetBarcodeHeight(50)
				},
				symbology: barcode.UPCA,
				data:      []byte("12345678901"),
			},
			{
				name: "CODE39 with HRI above",
				setup: func() {
					_, _ = fake.SelectHRICharacterPosition(barcode.HRIAbove)
					_, _ = fake.SetBarcodeHeight(80)
				},
				symbology: barcode.CODE39,
				data:      []byte("*TEST*"),
			},
			{
				name: "ITF without HRI",
				setup: func() {
					_, _ = fake.SelectHRICharacterPosition(barcode.HRINotPrinted)
					_, _ = fake.SetBarcodeHeight(60)
				},
				symbology: barcode.ITF,
				data:      []byte("123456"),
			},
		}

		for i, wf := range workflows {
			wf.setup()
			_, err := fake.PrintBarcode(wf.symbology, wf.data)
			if err != nil {
				t.Errorf("%s: PrintBarcode error: %v", wf.name, err)
				continue
			}

			// Verify the record
			records := fake.GetPrintedBarcodes()
			if len(records) <= i {
				t.Errorf("%s: missing barcode record", wf.name)
				continue
			}

			record := records[i]
			if record.Symbology != wf.symbology {
				t.Errorf("%s: symbology = %v, want %v", wf.name, record.Symbology, wf.symbology)
			}
			if !bytes.Equal(record.Data, wf.data) {
				t.Errorf("%s: data = %q, want %q", wf.name, record.Data, wf.data)
			}
		}

		// Verify all barcodes were recorded
		if len(fake.GetPrintedBarcodes()) != len(workflows) {
			t.Errorf("Recorded barcodes = %d, want %d", len(fake.GetPrintedBarcodes()), len(workflows))
		}
	})

	t.Run("handles CODE128 with different code sets", func(t *testing.T) {
		fake := NewFakeCapability()

		codeSets := []struct {
			codeSet barcode.Code128Set
			data    []byte
		}{
			{barcode.Code128SetA, []byte("UPPERCASE")},
			{barcode.Code128SetB, []byte("MixedCase123")},
			{barcode.Code128SetC, []byte("123456")},
		}

		for _, cs := range codeSets {
			_, err := fake.PrintBarcodeWithCodeSet(barcode.CODE128, cs.codeSet, cs.data)
			if err != nil {
				t.Errorf("CODE128 Set %v: %v", cs.codeSet, err)
			}
		}

		// Verify all were recorded correctly
		records := fake.GetPrintedBarcodes()
		if len(records) != len(codeSets) {
			t.Fatalf("Printed barcodes = %d, want %d", len(records), len(codeSets))
		}

		for i, cs := range codeSets {
			if records[i].CodeSet != cs.codeSet {
				t.Errorf("Record[%d] code set = %v, want %v", i, records[i].CodeSet, cs.codeSet)
			}
		}

		// Verify buffer contains code set prefixes
		buffer := fake.GetBuffer()
		for _, cs := range codeSets {
			prefix := []byte{'{', byte(cs.codeSet)}
			if !bytes.Contains(buffer, prefix) {
				t.Errorf("Buffer should contain code set %v prefix", cs.codeSet)
			}
		}
	})
}
