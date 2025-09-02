package character

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
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

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	CharColorNone byte = '0' // -> None (non-printing dots)
	CharColor1    byte = '1' // -> Color 1 (default)
	CharColor2    byte = '2' // -> Color 2
	CharColor3    byte = '3' // -> Color 3

	DefaultCharColor byte = CharColor1
)

// SelectCharacterColor selects the character color.
//
// Format:
//
//	ASCII: GS ( N pL pH fn m
//	Hex:   0x1D 0x28 0x4E 0x02 0x00 0x30 m
//	Decimal: 29 40 78 2 0 48 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 48–51 (model-dependent)
//
// Default:
//
//	m = 49 (Color 1)
//
// Description:
//
//	Selects the character color:
//	  m = 48 -> None (non-printing dots)
//	  m = 49 -> Color 1 (default)
//	  m = 50 -> Color 2
//	  m = 51 -> Color 3
//
// Notes:
//   - Applies to alphanumeric, Katakana, multilingual, user-defined, and
//     user-defined Kanji characters; does not affect graphics, bit images,
//     barcodes, or 2D codes.
//   - m = 48 treats glyph dots as non-printing (useful with background/shadow).
//   - Underline prints in the selected character color.
//   - Setting persists until ESC @, printer reset, or power-off.
//
// Byte sequence:
//
//	GS ( N 02 00 30 m -> 0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, m
func (ef *EffectsCommands) SelectCharacterColor(m byte) ([]byte, error) {
	// Validate allowed values
	switch m {
	case '0', '1', '2', '3':
		// Valid values
	default:
		return nil, ErrInvalidCharacterColor
	}
	return []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x30, m}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	BackgroundColorNone byte = '0' // -> None (no background dots)
	BackgroundColor1    byte = '1' // -> Color 1
	BackgroundColor2    byte = '2' // -> Color 2
	BackgroundColor3    byte = '3' // -> Color 3

	DefaultBackgroundColor byte = BackgroundColorNone
)

// SelectBackgroundColor selects the background color.
//
// Format:
//
//	ASCII: GS ( N pL pH fn m
//	Hex:   0x1D 0x28 0x4E 0x02 0x00 0x31 m
//	Decimal: 29 40 78 2 0 49 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 48–51 (model-dependent)
//
// Default:
//
//	m = 48 (None)
//
// Description:
//
//	Selects the background color:
//	  m = 48 -> None (no background dots printed)
//	  m = 49 -> Color 1
//	  m = 50 -> Color 2
//	  m = 51 -> Color 3
//
// Notes:
//   - Background color does not affect spaces skipped by HT, ESC $, ESC \,
//     line spacing, or reverse print background.
//   - Inter-character spacing (ESC SP, FS S) prints in this background color.
//   - Settings persist until ESC @, printer reset, or power-off.
//
// Byte sequence:
//
//	GS ( N 02 00 31 m -> 0x1D, 0x28, 0x4E, 0x02, 0x00, 0x31, m
func (ef *EffectsCommands) SelectBackgroundColor(m byte) ([]byte, error) {
	// Validate allowed values
	switch m {
	case '0', '1', '2', '3':
		// Valid values
	default:
		return nil, ErrInvalidBackgroundColor
	}
	return []byte{0x1D, 0x28, 0x4E, 0x02, 0x00, 0x31, m}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
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

	// Defaults
	DefaultShadowMode  byte = ShadowModeOffByte
	DefaultShadowColor byte = ShadowColorNone
)

// SetCharacterShadowMode turns shading (shadow) mode on or off.
//
// Format:
//
//	ASCII: GS ( N pL pH fn m a
//	Hex:   0x1D 0x28 0x4E pL pH 0x32 m a
//	Decimal: 29 40 78 pL pH 50 m a
//
// Range:
//
//	pL = 0x03, pH = 0x00
//	m = 0, 1, 48, 49
//	a = 48–51
//
// Default:
//
//	m = 0, a = 48
//
// Description:
//
//	Turns shadow mode on or off and sets shadow color:
//	  m: Shadow mode (0 or 48 = OFF, 1 or 49 = ON)
//	  a: Shadow color:
//	     48 -> None (not printed)
//	     49 -> Color 1
//	     50 -> Color 2
//	     51 -> Color 3
//
// Notes:
//   - Parameter a (shadowColor) MUST be supplied even when shadow mode is OFF.
//   - Shadow mode prints a shadow in the specified shadow color.
//   - Underline shadow is not printed.
//   - Reverse print does not alter shadow color.
//   - Settings persist until ESC @, printer reset, or power-off.
//
// Byte sequence:
//
//	GS ( N 03 00 32 m a -> 0x1D, 0x28, 0x4E, 0x03, 0x00, 0x32, m, a
func (ef *EffectsCommands) SetCharacterShadowMode(m byte, a byte) ([]byte, error) {
	// Validate allowed values
	switch m {
	case 0, 1, '0', '1':
		// Valid values
	default:
		return nil, ErrInvalidShadowMode
	}
	switch a {
	case '0', '1', '2', '3':
		// Valid values
	default:
		return nil, ErrInvalidShadowColor
	}
	return []byte{common.GS, '(', 'N', 0x03, 0x00, 0x32, m, a}, nil
}
