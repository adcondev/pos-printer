package character

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

type EncodeSystem byte

const (
	OneByte      EncodeSystem = 1   // legacy 1-byte encoding
	UTF8         EncodeSystem = 2   // UTF-8 encoding
	OneByteAscii EncodeSystem = '1' // (ASCII form)
	UTF8Ascii    EncodeSystem = '2' // (ASCII form)
)

type FontPriority byte

const (
	First  FontPriority = 0 // m = 0 (1st priority)
	Second FontPriority = 1 // m = 1 (2nd priority)
)

type FontFunction byte

const (
	AnkSansSerif             FontFunction = 0  // AnkSansSerif font (Sans serif)
	JapaneseGothic           FontFunction = 11 // Japanese font (Gothic)
	SimplifiedChineseMincho  FontFunction = 20 // Simplified Chinese (Mincho)
	TraditionalChineseMincho FontFunction = 30 // Traditional Chinese (Mincho)
	KoreanGothic             FontFunction = 41 // Korean font (Gothic)
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrEncoding     = fmt.Errorf("invalid encoding method(1-2 or '1'..'2')")
	ErrFontPriority = fmt.Errorf("invalid font priority(0-1)")
	ErrFontType     = fmt.Errorf("invalid font type(0,11,20,30,41)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance checks
var _ CodeConversionCapability = (*CodeConversionCommands)(nil)

// CodeConversionCapability defines encoding and font priority operations
type CodeConversionCapability interface {
	SelectCharacterEncodeSystem(encoding EncodeSystem) ([]byte, error)
	SetFontPriority(priority FontPriority, fontFunction FontFunction) ([]byte, error)
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
func (c *CodeConversionCommands) SelectCharacterEncodeSystem(m EncodeSystem) ([]byte, error) {
	// Validate allowed values
	switch m {
	case 1, 2, '1', '2':
		// Valid values
	default:
		return nil, ErrEncoding
	}
	return []byte{common.FS, '(', 'C', 0x02, 0x00, 0x30, byte(m)}, nil
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
//	      0  -> AnkSansSerif font (Sans serif)
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
func (c *CodeConversionCommands) SetFontPriority(m FontPriority, a FontFunction) ([]byte, error) {
	// Validate allowed values
	if m > 1 {
		return nil, ErrFontPriority
	}
	switch a {
	case 0, 11, 20, 30, 41:
		// Valid font types
	default:
		return nil, ErrFontType
	}

	return []byte{common.FS, '(', 'C', 0x03, 0x00, 0x3C, byte(m), byte(a)}, nil
}
