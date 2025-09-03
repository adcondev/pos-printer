package character

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	Encoding1Byte      byte = 1   // legacy 1-byte encoding
	EncodingUTF8       byte = 2   // UTF-8 encoding
	Encoding1ByteASCII byte = '1' // (ASCII form)
	EncodingUTF8ASCII  byte = '2' // (ASCII form)

	DefaultEncoding byte = Encoding1Byte

	FontPriorityFirst  byte = 0 // m = 0 (1st priority)
	FontPrioritySecond byte = 1 // m = 1 (2nd priority)

	FontANK                byte = 0  // ANK font (Sans serif)
	FontJapaneseGothic     byte = 11 // Japanese font (Gothic)
	FontSimplifiedChinese  byte = 20 // Simplified Chinese (Mincho)
	FontTraditionalChinese byte = 30 // Traditional Chinese (Mincho)
	FontKoreanGothic       byte = 41 // Korean font (Gothic)
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrInvalidEncoding     = fmt.Errorf("invalid encoding method(1-2 or '1'..'2')")
	ErrInvalidFontPriority = fmt.Errorf("invalid font priority(0-1)")
	ErrInvalidFontType     = fmt.Errorf("invalid font type(0,11,20,30,41)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance checks
var _ CodeConversionCapability = (*CodeConversionCommands)(nil)

// CodeConversionCapability defines encoding and font priority operations
type CodeConversionCapability interface {
	SelectCharacterEncodeSystem(encoding byte) ([]byte, error)
	SetFontPriority(priority byte, fontType byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// CodeConversionCommands implements CodeConversionCapability
type CodeConversionCommands struct{}

func NewCodeConversionCommands() *CodeConversionCommands {
	return &CodeConversionCommands{}
}

// SelectCharacterEncodeSystem selects the character encoding system.
//
// Format:
//
//	ASCII: FS ( C pL pH fn m
//	Hex:   0x1C 0x28 0x43 0x02 0x00 0x30 m
//	Decimal: 28 40 67 2 0 48 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 1, 2, 49, 50
//
// Default:
//
//	m = 1 (1-byte encoding)
//
// Description:
//
//	Selects the character encoding system:
//	  m = 1 or 49 -> 1-byte (legacy) encoding (model-dependent legacy code pages)
//	  m = 2 or 50 -> UTF-8 (Unicode)
//
// Notes:
//   - When UTF-8 is selected, ESC t (code table selection) is ignored.
//   - Settings persist until ESC @ (initialize), printer reset, or power off.
//   - Availability of specific legacy encodings is model-dependent.
//
// Byte sequence:
//
//	FS ( C 02 00 30 m -> 0x1C, 0x28, 0x43, 0x02, 0x00, 0x30, m
func (c *CodeConversionCommands) SelectCharacterEncodeSystem(m byte) ([]byte, error) {
	// Validate allowed values
	switch m {
	case 1, 2, '1', '2':
		// Valid values
	default:
		return nil, ErrInvalidEncoding
	}
	return []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, m}, nil
}

// SetFontPriority sets the font priority.
//
// Format:
//
//	ASCII: FS ( C pL pH fn m a
//	Hex:   0x1C 0x28 0x43 0x03 0x00 0x3C m a
//	Decimal: 28 40 67 3 0 60 m a
//
// Range:
//
//	pL = 0x03, pH = 0x00
//	m = 0–1
//	a = 0, 11, 20, 30, 41
//
// Default:
//
//	m = 0, a = 0
//
// Description:
//
//	Sets font priority where:
//	  - m: Priority rank (0 = 1st priority, 1 = 2nd priority)
//	  - a: Font type
//	      0  -> ANK font (Sans serif)
//	      11 -> Japanese font (Gothic)
//	      20 -> Simplified Chinese font (Mincho)
//	      30 -> Traditional Chinese font (Mincho)
//	      41 -> Korean font (Gothic)
//
// Notes:
//   - Assigns a font style to a priority slot (1st or 2nd).
//   - If the style already exists in the priority list, promotion/demotion
//     is handled so that the newly specified font becomes the selected
//     priority.
//   - Settings persist until ESC @ (initialize), printer reset, or power-off.
//
// Byte sequence:
//
//	FS ( C 03 00 3C m a -> 0x1C, 0x28, 0x43, 0x03, 0x00, 0x3C, m, a
func (c *CodeConversionCommands) SetFontPriority(m byte, a byte) ([]byte, error) {
	// Validate allowed values
	if m > 1 {
		return nil, ErrInvalidFontPriority
	}
	switch a {
	case 0, 11, 20, 30, 41:
		// Valid font types
	default:
		return nil, ErrInvalidFontType
	}

	return []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, m, a}, nil
}
