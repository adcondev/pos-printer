package character

import (
	"fmt"
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
