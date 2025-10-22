package character

import (
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for user-defined characters.
// ESC/POS is the command system used by thermal receipt printers to control
// custom glyph definitions, allowing creation of custom dot-matrix patterns
// for special characters not available in standard character sets.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	// UserDefinedOff represents user-defined character mode off (LSB = 0)
	UserDefinedOff byte = 0x00
	// UserDefinedOn represents user-defined character mode on (LSB = 1)
	UserDefinedOn byte = 0x01

	// UserDefinedOffASCII represents user-defined character mode off (ASCII variant)
	UserDefinedOffASCII byte = '0'
	// UserDefinedOnASCII represents user-defined character mode on (ASCII variant)
	UserDefinedOnASCII byte = '1'

	// UserDefinedMinCode represents minimum character code for user-defined characters
	UserDefinedMinCode byte = 32
	// UserDefinedMaxCode represents maximum character code for user-defined characters
	UserDefinedMaxCode byte = 126
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrCharacterCode represents an invalid character code error
var (
	ErrCharacterCode = fmt.Errorf("invalid character code(try %d-%d)", UserDefinedMinCode, UserDefinedMaxCode)
	ErrYValue        = fmt.Errorf("invalid y value(try y >= 1)")
	ErrCodeRange     = fmt.Errorf("invalid code range(try c2 >= c1 and c2 <= %d)", UserDefinedMaxCode)
	ErrDefinition    = fmt.Errorf("invalid definition count(try matching number of codes in range c1-c2)")
	ErrDataLength    = fmt.Errorf("invalid data length(try using exactly y*width bytes)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance checks
var _ UserDefinedCapability = (*UserDefinedCommands)(nil)

// UserDefinedCapability defines user-defined character operations
type UserDefinedCapability interface {
	SelectUserDefinedCharacterSet(charSet byte) []byte
	DefineUserDefinedCharacters(height, startCode, endCode byte, definitions []UserDefinedChar) ([]byte, error)
	CancelUserDefinedCharacter(charCode byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// UserDefinedCommands implements UserDefinedCapability
type UserDefinedCommands struct{}

// UserDefinedChar represents one glyph definition for a single character code
type UserDefinedChar struct {
	Width byte   // Width in dots (xi)
	Data  []byte // Raw column data, length must equal y * Width
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateCharacterCode validates if character code is within valid range
func ValidateCharacterCode(code byte) error {
	if code < UserDefinedMinCode || code > UserDefinedMaxCode {
		return ErrCharacterCode
	}
	return nil
}

// ValidateCodeRange validates if code range is valid
func ValidateCodeRange(c1, c2 byte) error {
	if err := ValidateCharacterCode(c1); err != nil {
		return err
	}
	if c2 < c1 || c2 > UserDefinedMaxCode {
		return ErrCodeRange
	}
	return nil
}
