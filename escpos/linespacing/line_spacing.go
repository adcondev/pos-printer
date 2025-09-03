package linespacing

import (
	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Line spacing limits and defaults
const (
	MinLineSpacing     = 0    // Minimum line spacing value
	MaxLineSpacing     = 255  // Maximum line spacing value
	DefaultLineSpacing = 30   // Default line spacing (model dependent, typically 30-80 dots)
	MaxLineSpacingMM   = 1016 // Maximum line spacing in mm (40 inches)
)

// ============================================================================
// Error Definitions
// ============================================================================

// No specific errors needed for line spacing as byte range (0-255) covers all valid values

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for line spacing commands in ESC/POS printers.
type Capability interface {
	SetLineSpacing(lines byte) []byte
	SelectDefaultLineSpacing() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for ESC/POS printers.
type Commands struct{}

func NewCommands() *Commands {
	return &Commands{}
}

// SetLineSpacing sets the line spacing to n × (vertical or horizontal motion unit).
//
// Format:
//
//	ASCII: ESC 3 n
//	Hex:   0x1B 0x33 n
//	Decimal: 27 51 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	The amount of line spacing corresponding to the "default line spacing"
//	(equivalent to a value between 30 and 80 dots).
//
// Description:
//
//	Sets the line spacing to n × (vertical or horizontal motion unit).
//
// Notes:
//   - Maximum line spacing is 1016 mm (40 inches); if exceeded, printer uses maximum.
//   - In Standard mode the vertical motion unit is used.
//   - In Page mode the motion unit depends on ESC T setting.
//   - Line spacing can be set independently in Standard and Page modes.
//   - Motion unit changes after setting don't affect the numeric value.
//   - Remains in effect until ESC 2, ESC @, reset, or power off.
//
// Byte sequence:
//
//	ESC 3 n -> 0x1B, 0x33, n
func (lsc *Commands) SetLineSpacing(n byte) []byte {
	return []byte{common.ESC, '3', n}
}

// SelectDefaultLineSpacing sets the line spacing to the printer's default.
//
// Format:
//
//	ASCII: ESC 2
//	Hex:   0x1B 0x32
//	Decimal: 27 50
//
// Description:
//
//	Sets the line spacing to the default line spacing value.
//
// Notes:
//   - Line spacing can be set independently in Standard and Page modes.
//   - Remains in effect until ESC 3, ESC @, reset, or power off.
//
// Byte sequence:
//
//	ESC 2 -> 0x1B, 0x32
func (lsc *Commands) SelectDefaultLineSpacing() []byte {
	return []byte{common.ESC, '2'}
}
