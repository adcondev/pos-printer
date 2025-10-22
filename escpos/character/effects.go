package character

import (
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for character visual effects.
// ESC/POS is the command system used by thermal receipt printers to control
// character colors, background colors, and shadow effects for enhanced
// visual appearance in multi-color thermal printing.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	// CharColorNone represents no character color
	CharColorNone byte = '0'
	// CharColor1 represents character color 1
	CharColor1 byte = '1'
	// CharColor2 represents character color 2
	CharColor2 byte = '2'
	// CharColor3 represents character color 3
	CharColor3 byte = '3'

	// BackgroundColorNone represents no background color
	BackgroundColorNone byte = '0'
	// BackgroundColor1 represents background color 1
	BackgroundColor1 byte = '1'
	// BackgroundColor2 represents background color 2
	BackgroundColor2 byte = '2'
	// BackgroundColor3 represents background color 3
	BackgroundColor3 byte = '3'

	// GS represents the GS command byte for special functions
	GS byte = 0x1D
	// LeftParenN represents the left parenthesis byte
	LeftParenN byte = 0x28
	// LetterN represents the letter N byte
	LetterN byte = 0x4E
	// GSParenNFnShadow represents the shadow function code
	GSParenNFnShadow byte = 0x32

	// GSParenNShadowPL represents the pL parameter for shadow function
	GSParenNShadowPL byte = 0x03
	// GSParenNShadowPH represents the pH parameter for shadow function
	GSParenNShadowPH byte = 0x00

	// ShadowModeOffByte represents shadow mode off (numeric)
	ShadowModeOffByte byte = 0x00
	// ShadowModeOnByte represents shadow mode on (numeric)
	ShadowModeOnByte byte = 0x01
	// ShadowModeOffASCII represents shadow mode off (ASCII)
	ShadowModeOffASCII byte = '0'
	// ShadowModeOnASCII represents shadow mode on (ASCII)
	ShadowModeOnASCII byte = '1'

	// ShadowColorNone represents no shadow color
	ShadowColorNone byte = '0'
	// ShadowColor1 represents shadow color 1
	ShadowColor1 byte = '1'
	// ShadowColor2 represents shadow color 2
	ShadowColor2 byte = '2'
	// ShadowColor3 represents shadow color 3
	ShadowColor3 byte = '3'
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrInvalidCharacterColor represents an invalid character color error
var (
	ErrInvalidCharacterColor  = fmt.Errorf("invalid character color('0'..'3')")
	ErrInvalidBackgroundColor = fmt.Errorf("invalid background color('0'..'3')")
	ErrInvalidShadowColor     = fmt.Errorf("invalid shadow color('0'..'3')")
	ErrInvalidShadowMode      = fmt.Errorf("invalid shadow mode(0-1 or '0'-'1')")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance checks
var _ EffectsCapability = (*EffectsCommands)(nil)

// EffectsCapability defines character effects operations
type EffectsCapability interface {
	SelectCharacterColor(color byte) ([]byte, error)
	SelectBackgroundColor(color byte) ([]byte, error)
	SetCharacterShadowMode(shadowMode byte, shadowColor byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// EffectsCommands implements EffectsCapability
type EffectsCommands struct{}

// NewEffectsCommands creates a new instance of EffectsCommands
func NewEffectsCommands() *EffectsCommands {
	return &EffectsCommands{}
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateCharacterColor validates if character color is valid
func ValidateCharacterColor(color byte) error {
	switch color {
	case CharColorNone, CharColor1, CharColor2, CharColor3:
		return nil
	default:
		return ErrInvalidCharacterColor
	}
}

// ValidateBackgroundColor validates if background color is valid
func ValidateBackgroundColor(color byte) error {
	switch color {
	case BackgroundColorNone, BackgroundColor1, BackgroundColor2, BackgroundColor3:
		return nil
	default:
		return ErrInvalidBackgroundColor
	}
}

// ValidateShadowMode validates if shadow mode is valid
func ValidateShadowMode(mode byte) error {
	switch mode {
	case ShadowModeOffByte, ShadowModeOnByte, ShadowModeOffASCII, ShadowModeOnASCII:
		return nil
	default:
		return ErrInvalidShadowMode
	}
}

// ValidateShadowColor validates if shadow color is valid
func ValidateShadowColor(color byte) error {
	switch color {
	case ShadowColorNone, ShadowColor1, ShadowColor2, ShadowColor3:
		return nil
	default:
		return ErrInvalidShadowColor
	}
}
