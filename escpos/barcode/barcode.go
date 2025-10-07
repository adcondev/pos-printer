package barcode

import (
	"errors"
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// HRI (Human Readable Interpretation) position modes
type HRIPosition byte

const (
	HRINotPrinted      HRIPosition = 0x00 // n = 0
	HRIAbove           HRIPosition = 0x01 // n = 1
	HRIBelow           HRIPosition = 0x02 // n = 2
	HRIBoth            HRIPosition = 0x03 // n = 3
	HRINotPrintedASCII HRIPosition = '0'  // n = 48
	HRIAboveASCII      HRIPosition = '1'  // n = 49
	HRIBelowASCII      HRIPosition = '2'  // n = 50
	HRIBothASCII       HRIPosition = '3'  // n = 51
)

// HRI Font types
type HRIFont byte

const (
	HRIFontA        HRIFont = 0x00 // n = 0
	HRIFontB        HRIFont = 0x01 // n = 1
	HRIFontC        HRIFont = 0x02 // n = 2
	HRIFontD        HRIFont = 0x03 // n = 3
	HRIFontE        HRIFont = 0x04 // n = 4
	HRIFontAASCII   HRIFont = '0'  // n = 48
	HRIFontBASCII   HRIFont = '1'  // n = 49
	HRIFontCASCII   HRIFont = '2'  // n = 50
	HRIFontDASCII   HRIFont = '3'  // n = 51
	HRIFontEASCII   HRIFont = '4'  // n = 52
	HRISpecialFontA HRIFont = 97   // n = 97
	HRISpecialFontB HRIFont = 98   // n = 98
)

// Barcode height range
type Height byte

const (
	MinHeight     Height = 1
	MaxHeight     Height = 255
	DefaultHeight Height = 162 // Model dependent default
)

// Barcode width/module size
type Width byte

const (
	MinWidth     Width = 2
	MaxWidth     Width = 6
	DefaultWidth Width = 3
	// Extended width values (model dependent)
	ExtendedMinWidth Width = 68
	ExtendedMaxWidth Width = 76
)

// Barcode symbology types
type Symbology byte

// Function A symbologies (NUL-terminated)
const (
	UPCA    Symbology = 0 // UPC-A (11-12 digits)
	UPCE    Symbology = 1 // UPC-E (6-8, 11-12 digits)
	JAN13   Symbology = 2 // JAN13/EAN13 (12-13 digits)
	JAN8    Symbology = 3 // JAN8/EAN8 (7-8 digits)
	CODE39  Symbology = 4 // CODE39 (variable length)
	ITF     Symbology = 5 // Interleaved 2 of 5 (even digits)
	CODABAR Symbology = 6 // CODABAR/NW-7 (variable length)
)

// Function B symbologies (length-prefixed)
const (
	UPCAB           Symbology = 65 // UPC-A (11-12 digits)
	UPCEB           Symbology = 66 // UPC-E (6-8, 11-12 digits)
	EAN13           Symbology = 67 // EAN13 (12-13 digits)
	EAN8            Symbology = 68 // EAN8 (7-8 digits)
	CODE39B         Symbology = 69 // CODE39 (1-255 chars)
	ITFB            Symbology = 70 // ITF (2-254 even digits)
	CODABARB        Symbology = 71 // CODABAR (2-255 chars)
	CODE93          Symbology = 72 // CODE93 (1-255 chars)
	CODE128         Symbology = 73 // CODE128 (2-255 bytes)
	GS1128          Symbology = 74 // GS1-128 (2-255 bytes)
	GS1DataBarOmni  Symbology = 75 // GS1 DataBar Omnidirectional (13 digits)
	GS1DataBarTrunc Symbology = 76 // GS1 DataBar Truncated (13 digits)
	GS1DataBarLim   Symbology = 77 // GS1 DataBar Limited (13 digits)
	GS1DataBarExp   Symbology = 78 // GS1 DataBar Expanded (2-255 chars)
	CODE128Auto     Symbology = 79 // CODE128 Auto (1-255 bytes)
)

// CODE128 code sets
type Code128Set byte

const (
	Code128SetA Code128Set = 65 // Code set A (ASCII 0-95)
	Code128SetB Code128Set = 66 // Code set B (ASCII 32-127)
	Code128SetC Code128Set = 67 // Code set C (00-99 numeric pairs)
)

// Special characters
const (
	Code128Prefix byte = '{' // 0x7B - CODE128 prefix
	DataBarPrefix byte = '{' // 0x7B - GS1 DataBar Expanded prefix
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrHRIPosition      = errors.New("invalid HRI position (try 0-3 or '0'..'3')")
	ErrHRIFont          = errors.New("invalid HRI font (try 0-4, '0'..'4', 97, or 98)")
	ErrHeight           = errors.New("invalid barcode height (try 1-255)")
	ErrWidth            = errors.New("invalid barcode width (try 2-6 or 68-76 for extended)")
	ErrSymbology        = errors.New("invalid barcode symbology")
	ErrDataTooShort     = errors.New("barcode data too short")
	ErrDataTooLong      = errors.New("barcode data too long (max 255 bytes)")
	ErrOddITFLength     = errors.New("ITF barcode requires even number of digits")
	ErrCode128Set       = errors.New("invalid CODE128 code set (try 65-67)")
	ErrCode128NoCodeSet = errors.New("CODE128 requires code set specification")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for barcode commands
type Capability interface {
	// HRI (Human Readable Interpretation) settings
	SelectHRICharacterPosition(position HRIPosition) ([]byte, error)
	SelectFontForHRI(font HRIFont) ([]byte, error)

	// Barcode dimensions
	SetBarcodeHeight(height Height) ([]byte, error)
	SetBarcodeWidth(width Width) ([]byte, error)

	// Barcode printing
	PrintBarcode(symbology Symbology, data []byte) ([]byte, error)
	PrintBarcodeWithCodeSet(symbology Symbology, codeSet Code128Set, data []byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for barcode commands
type Commands struct{}

func NewCommands() *Commands {
	return &Commands{}
}

// SelectHRICharacterPosition selects the print position of HRI (Human Readable Interpretation) characters.
//
// Format:
//
//	ASCII: GS H n
//	Hex:   0x1D 0x48 n
//	Decimal: 29 72 n
//
// Range:
//
//	n = 0–3, 48–51
//
// Default:
//
//	n = 0 (Not printed)
//
// Description:
//
//	Selects the print position of HRI characters when printing a barcode:
//	  0 or 48 -> Not printed
//	  1 or 49 -> Above the barcode
//	  2 or 50 -> Below the barcode
//	  3 or 51 -> Both above and below the barcode
//
// Notes:
//   - HRI characters are printed using the font specified by GS f.
//   - The setting persists until ESC @ (initialize), printer reset, or power-off.
//
// Byte sequence:
//
//	GS H n -> 0x1D, 0x48, n
func (c *Commands) SelectHRICharacterPosition(n HRIPosition) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, 3, '0', '1', '2', '3':
		// Valid values
	default:
		return nil, ErrHRIPosition
	}
	return []byte{common.GS, 'H', byte(n)}, nil
}

// SelectFontForHRI selects the font used to print HRI (Human Readable Interpretation) characters.
//
// Format:
//
//	ASCII: GS f n
//	Hex:   0x1D 0x66 n
//	Decimal: 29 102 n
//
// Range:
//
//	n: model-dependent. Common supported values:
//	  0–4, 48–52, 97, 98
//
// Default:
//
//	n = 0
//
// Description:
//
//	Selects the font for HRI characters printed with barcodes:
//	  0 or 48  -> Font A
//	  1 or 49  -> Font B
//	  2 or 50  -> Font C
//	  3 or 51  -> Font D
//	  4 or 52  -> Font E
//	  97       -> Special font A (model dependent)
//	  98       -> Special font B (model dependent)
//
// Notes:
//   - The chosen font applies only to HRI characters.
//   - HRI characters are printed at the position set by GS H.
//   - Built-in font availability and metrics vary by model.
//
// Byte sequence:
//
//	GS f n -> 0x1D, 0x66, n
func (c *Commands) SelectFontForHRI(n HRIFont) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, 3, 4:
		// Numeric values
	case '0', '1', '2', '3', '4':
		// ASCII values
	case 97, 98:
		// Special fonts
	default:
		return nil, ErrHRIFont
	}
	return []byte{common.GS, 'f', byte(n)}, nil
}

// SetBarcodeHeight sets the barcode height.
//
// Format:
//
//	ASCII: GS h n
//	Hex:   0x1D 0x68 n
//	Decimal: 29 104 n
//
// Range:
//
//	n = 1–255
//
// Default:
//
//	n: model dependent (example default: 162)
//
// Description:
//
//	Sets the height of a barcode to n dots.
//
// Notes:
//   - The units for n depend on the printer model.
//   - This setting remains effective until ESC @ (initialize), printer reset, or power-off.
//
// Byte sequence:
//
//	GS h n -> 0x1D, 0x68, n
func (c *Commands) SetBarcodeHeight(height Height) ([]byte, error) {
	if height < MinHeight || height > MaxHeight {
		return nil, fmt.Errorf("%w: %d", ErrHeight, height)
	}
	return []byte{common.GS, 'h', byte(height)}, nil
}

// SetBarcodeWidth sets the horizontal module width for barcodes.
//
// Format:
//
//	ASCII: GS w n
//	Hex:   0x1D 0x77 n
//	Decimal: 29 119 n
//
// Range:
//
//	n = 2–6 (typical numeric values)
//	or model-dependent alternate values 68–76
//
// Default:
//
//	n = 3 (model-dependent)
//
// Description:
//
//	Sets the barcode module width (horizontal size). Units and exact effect
//	depend on the printer model. The setting remains effective until ESC @,
//	printer reset, or power-off.
//
// Notes:
//   - This affects the module width for various barcode types (see printer spec).
//   - The command does not validate model-specific allowed values; caller must
//     supply a value supported by the target printer.
//
// Byte sequence:
//
//	GS w n -> 0x1D, 0x77, n
func (c *Commands) SetBarcodeWidth(width Width) ([]byte, error) {
	// Validate standard and extended ranges
	if (width >= MinWidth && width <= MaxWidth) ||
		(width >= ExtendedMinWidth && width <= ExtendedMaxWidth) {
		return []byte{common.GS, 'w', byte(width)}, nil
	}
	return nil, fmt.Errorf("%w: %d", ErrWidth, width)
}

// PrintBarcode builds the GS k command byte sequence to print a barcode.
//
// Command summary:
//
//	Function A (m = 0–6):
//	  Format:  GS k m d1...dk NUL
//	  Data end: NUL (0x00) terminator (length byte NOT sent)
//	Function B (m = 65–79):
//	  Format:  GS k m n d1...dn
//	  Data length: single length byte n (1–255), NO terminator
//
// Byte sequence prefix:
//
//	GS k -> 0x1D 0x6B
//
// Parameter m (symbology selector):
//
//	Function A (classic forms):
//	  0  UPC-A           (k = 11 or 12 digits)  (numeric)
//	  1  UPC-E           (k = 6–8, 11, 12)      (numeric; k=7/8/11/12 must start with '0')
//	  2  JAN13 / EAN13   (k = 12 or 13 digits)  (numeric)
//	  3  JAN8  / EAN8    (k = 7 or 8 digits)    (numeric)
//	  4  CODE39          (k >= 1)               (0–9 A–Z space $ % * + - . /)  Start/stop '*' auto if omitted
//	  5  ITF (Interleaved 2 of 5) (k >= 2 even) (numeric; odd final digit ignored)
//	  6  CODABAR (NW-7)  (k >= 2)               (Start/stop A–D/a–d must be present; not auto-added)
//	Function B (extended forms):
//	  65 UPC-A      (n = 11 or 12)
//	  66 UPC-E      (n = 6–8, 11, 12)
//	  67 EAN13      (n = 12 or 13)
//	  68 EAN8       (n = 7 or 8)
//	  69 CODE39     (1–255)
//	  70 ITF        (2–254 even)
//	  71 CODABAR    (2–255)
//	  72 CODE93     (1–255) (start/stop + 2 check chars auto)
//	  73 CODE128    (2–255) (d1= '{' (0x7B)=123, d2= 65–67 => Set A/B/C; check digit auto)
//	  74 GS1-128    (2–255) (FNC1, check digits auto; special SP,(,),* rules)
//	  75 GS1 DataBar Omnidirectional (n=13 digits; AI(01), check digit auto)
//	  76 GS1 DataBar Truncated        (n=13)
//	  77 GS1 DataBar Limited          (n=13; first digit constraint)
//	  78 GS1 DataBar Expanded         (2–255; uses '{'+code for FNC1 / '(' / ')')
//	  79 CODE128 Auto                 (1–255; 0–255 byte data)
//
// Notes:
//   - This function DOES NOT validate symbology-specific content or lengths;
//     caller must supply conforming data.
//   - After printing, printer returns to "beginning of line" state.
//   - Not affected by most text print modes (except upside-down).
//   - In Page mode, data is buffered (rendering per Page mode rules).
//   - Width exceeding print area is ignored/clipped by device.
//
// Byte sequence:
//
//	Function A: GS k m data... NUL -> 0x1D, 0x6B, m, data..., 0x00
//	Function B: GS k m n data...   -> 0x1D, 0x6B, m, n, data...
func (c *Commands) PrintBarcode(symbology Symbology, data []byte) ([]byte, error) {
	// Validate data exists
	if len(data) == 0 {
		return nil, ErrDataTooShort
	}

	// Build command based on symbology type
	if symbology <= CODABAR {
		// Function A (NUL-terminated)
		return c.buildFunctionA(symbology, data)
	} else if symbology >= UPCAB && symbology <= CODE128Auto {
		// Function B (length-prefixed)
		return c.buildFunctionB(symbology, data)
	}

	return nil, ErrSymbology
}

// PrintBarcodeWithCodeSet prints a CODE128 or GS1-128 barcode with explicit code set.
//
// Description:
//
//	Specialized method for CODE128 (m=73) and GS1-128 (m=74) that require
//	code set specification. The first two bytes of the barcode data must be:
//	  d1 = '{' (0x7B)
//	  d2 = 65-67 (Code set A/B/C)
//
// Notes:
//   - Use this method when you need explicit control over CODE128 code sets.
//   - For automatic code set selection, use PrintBarcode with CODE128Auto (m=79).
//
// Byte sequence:
//
//	GS k m n '{' codeSet data... -> 0x1D, 0x6B, m, n, 0x7B, codeSet, data...
func (c *Commands) PrintBarcodeWithCodeSet(symbology Symbology, codeSet Code128Set, data []byte) ([]byte, error) {
	// Validate symbology supports code sets
	if symbology != CODE128 && symbology != GS1128 {
		return nil, fmt.Errorf("%w: symbology %d does not support code sets", ErrSymbology, symbology)
	}

	// Validate code set
	if codeSet < Code128SetA || codeSet > Code128SetC {
		return nil, ErrCode128Set
	}

	// Build data with code set prefix
	prefixedData := make([]byte, 0, len(data)+2)
	prefixedData = append(prefixedData, Code128Prefix, byte(codeSet))
	prefixedData = append(prefixedData, data...)

	return c.buildFunctionB(symbology, prefixedData)
}

// buildFunctionA builds Function A barcode command (NUL-terminated)
func (c *Commands) buildFunctionA(symbology Symbology, data []byte) ([]byte, error) {
	// Basic validation for Function A symbologies
	switch symbology {
	case ITF:
		// ITF requires even number of digits
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Build command: GS k m data... NUL
	cmd := []byte{common.GS, 'k', byte(symbology)}
	cmd = append(cmd, data...)
	cmd = append(cmd, common.NUL)
	return cmd, nil
}

// buildFunctionB builds Function B barcode command (length-prefixed)
func (c *Commands) buildFunctionB(symbology Symbology, data []byte) ([]byte, error) {
	// Validate data length (max 255 for single byte length)
	if len(data) > 255 {
		return nil, ErrDataTooLong
	}

	// Special validation for certain symbologies
	switch symbology {
	case CODE128, GS1128:
		// Check if data has the required code set prefix
		if len(data) < 2 || data[0] != Code128Prefix ||
			data[1] < byte(Code128SetA) || data[1] > byte(Code128SetC) {
			return nil, ErrCode128NoCodeSet
		}
	case ITFB:
		// ITF requires even number of digits
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Build command: GS k m n data...
	cmd := []byte{common.GS, 'k', byte(symbology), byte(len(data))}
	cmd = append(cmd, data...)
	return cmd, nil
}

// Helper functions for common barcode operations

// ValidateNumericData checks if all bytes are numeric digits
func ValidateNumericData(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// ValidateCode39Data checks if all bytes are valid CODE39 characters
func ValidateCode39Data(data []byte) bool {
	for _, b := range data {
		switch {
		case b >= '0' && b <= '9':
		case b >= 'A' && b <= 'Z':
		case b == ' ' || b == '$' || b == '%' || b == '*' ||
			b == '+' || b == '-' || b == '.' || b == '/':
		default:
			return false
		}
	}
	return true
}

// ValidateCodabarData checks if data has valid CODABAR start/stop characters
func ValidateCodabarData(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	// Check start character
	start := data[0]
	if (start < 'A' || start > 'D') && (start < 'a' || start > 'd') {
		return false
	}
	// Check stop character
	stop := data[len(data)-1]
	if (stop < 'A' || stop > 'D') && (stop < 'a' || stop > 'd') {
		return false
	}
	return true
}
