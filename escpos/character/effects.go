package character

import (
	"fmt"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	// Character color values (ASCII '0'..'3')
	CharColorNone byte = '0' // -> None (non-printing dots)
	CharColor1    byte = '1' // -> Color 1 (default)
	CharColor2    byte = '2' // -> Color 2
	CharColor3    byte = '3' // -> Color 3

	// Background color values (ASCII '0'..'3')
	BackgroundColorNone byte = '0' // -> None (no background dots)
	BackgroundColor1    byte = '1' // -> Color 1
	BackgroundColor2    byte = '2' // -> Color 2
	BackgroundColor3    byte = '3' // -> Color 3

	// Command bytes
	GSParenN         byte = 0x1D
	LeftParenN       byte = 0x28
	LetterN          byte = 0x4E
	GSParenNFnShadow byte = 0x32 // fn = 50 (0x32)

	// pL/pH for this function (fn + 2 parameters)
	GSParenNShadow_pL byte = 0x03
	GSParenNShadow_pH byte = 0x00

	// Shadow mode values (accepts numeric and ASCII forms)
	ShadowModeOffByte  byte = 0x00 // numeric 0
	ShadowModeOnByte   byte = 0x01 // numeric 1
	ShadowModeOffASCII byte = '0'
	ShadowModeOnASCII  byte = '1'

	// Shadow color values (ASCII '0'..'3')
	ShadowColorNone byte = '0' // -> None (not printed)
	ShadowColor1    byte = '1' // -> Color 1
	ShadowColor2    byte = '2' // -> Color 2
	ShadowColor3    byte = '3' // -> Color 3
)

// ============================================================================
// Error Definitions
// ============================================================================

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

func NewEffectsCommands() *EffectsCommands {
	return &EffectsCommands{}
}
