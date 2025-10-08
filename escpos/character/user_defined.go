package character

import (
	"fmt"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	UserDefinedOff byte = 0x00 // LSB = 0 -> user-defined OFF
	UserDefinedOn  byte = 0x01 // LSB = 1 -> user-defined ON

	// ASCII-digit variants sometimes accepted by implementations

	UserDefinedOffASCII byte = '0'
	UserDefinedOnASCII  byte = '1'

	UserDefinedMinCode byte = 32
	UserDefinedMaxCode byte = 126
)

// ============================================================================
// Error Definitions
// ============================================================================

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
