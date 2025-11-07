package character

import (
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for character encoding conversion.
// ESC/POS is the command system used by thermal receipt printers to control
// character encoding systems (UTF-8, legacy 1-byte) and font priority for
// multi-language support including CJK fonts.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// EncodeSystem represents character encoding system types (1-byte or UTF-8)
type EncodeSystem byte

const (
	// OneByte represents legacy 1-byte encoding
	OneByte EncodeSystem = 1
	// UTF8 represents UTF-8 encoding
	UTF8 EncodeSystem = 2
	// OneByteASCII represents legacy 1-byte encoding (ASCII form)
	OneByteASCII EncodeSystem = '1'
	// UTF8Ascii represents UTF-8 encoding (ASCII form)
	UTF8Ascii EncodeSystem = '2'
)

// FontPriority represents the font priority level for character rendering
type FontPriority byte

const (
	// First represents first priority font
	First FontPriority = 0
	// Second represents second priority font
	Second FontPriority = 1
)

// FontFunction represents the font function type for multi-language support
type FontFunction byte

const (
	// AnkSansSerif represents AnkSansSerif font (Sans serif)
	AnkSansSerif FontFunction = 0
	// JapaneseGothic represents Japanese font (Gothic)
	JapaneseGothic FontFunction = 11
	// SimplifiedChineseMincho represents Simplified Chinese (Mincho)
	SimplifiedChineseMincho FontFunction = 20
	// TraditionalChineseMincho represents Traditional Chinese (Mincho)
	TraditionalChineseMincho FontFunction = 30
	// KoreanGothic represents Korean font (Gothic)
	KoreanGothic FontFunction = 41
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrEncoding represents an invalid encoding method error
var (
	ErrEncoding     = fmt.Errorf("invalid encoding method(1-2 or '1'-'2')")
	ErrFontPriority = fmt.Errorf("invalid font priority(0-1)")
	ErrFontType     = fmt.Errorf("invalid font type(0, 11, 20, 30, 41)")
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

// NewCodeConversionCommands creates a new instance of CodeConversionCommands
func NewCodeConversionCommands() *CodeConversionCommands {
	return &CodeConversionCommands{}
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateEncodeSystem validates if encode system is valid
func ValidateEncodeSystem(encoding EncodeSystem) error {
	switch encoding {
	case OneByte, UTF8, OneByteASCII, UTF8Ascii:
		return nil
	default:
		return ErrEncoding
	}
}

// ValidateFontPriority validates if font priority is valid
func ValidateFontPriority(priority FontPriority) error {
	if priority > Second {
		return ErrFontPriority
	}
	return nil
}

// ValidateFontFunction validates if font function is valid
func ValidateFontFunction(font FontFunction) error {
	switch font {
	case AnkSansSerif, JapaneseGothic, SimplifiedChineseMincho,
		TraditionalChineseMincho, KoreanGothic:
		return nil
	default:
		return ErrFontType
	}
}
