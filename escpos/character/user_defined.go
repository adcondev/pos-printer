package character

import "fmt"

// ============================================================================
// Constant and Var Definitions
// ============================================================================

const (
	UserDefinedOff byte = 0x00 // LSB = 0 -> user-defined OFF
	UserDefinedOn  byte = 0x01 // LSB = 1 -> user-defined ON

	// ASCII-digit variants sometimes accepted by implementations

	UserDefinedOffASCII byte = '0'
	UserDefinedOnASCII  byte = '1'

	ESCAmpersand              = 0x26
	DefineUDChars_preambleLen = 5 // ESC & y c1 c2

	UserDefinedMinCode byte = 32
	UserDefinedMaxCode byte = 126
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrInvalidCharacterCode = fmt.Errorf("invalid character code(try 32-126)")
	ErrInvalidYValue        = fmt.Errorf("invalid y value(try y >= 1)")
	ErrInvalidCodeRange     = fmt.Errorf("invalid code range(try c2 >= c1 and c2 <= 126)")
	ErrInvalidDefinition    = fmt.Errorf("invalid definition count(try matching number of codes in range c1-c2)")
	ErrInvalidDataLength    = fmt.Errorf("invalid data length(try using exactly y*width bytes)")
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

// SelectUserDefinedCharacterSet selects or cancels the user-defined character set.
//
// Format:
//
//	ASCII: ESC % n
//	Hex:   0x1B 0x25 n
//	Decimal: 27 37 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Selects or cancels the user-defined character set. When the least-significant
//	bit (LSB) of n is 0, the user-defined character set is canceled. When the
//	LSB of n is 1, the user-defined character set is selected.
//
// Notes:
//   - When the user-defined character set is canceled the resident (built-in)
//     character set is automatically selected.
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off.
//   - This command affects alphanumeric, Kana, multilingual, and user-defined
//     characters as applicable per model.
//
// Byte sequence:
//
//	ESC % n -> 0x1B, 0x25, n
func (udc *UserDefinedCommands) SelectUserDefinedCharacterSet(n byte) []byte {
	return []byte{0x1B, 0x25, n}
}

// DefineUserDefinedCharacters defines user-defined glyph patterns for character codes.
//
// Format:
//
//	ASCII: ESC & y c1 c2 [x1 d1...d(y*x1)]...[xk d1...d(y*xk)]
//	Hex:   0x1B 0x26 y c1 c2 [x1 data...]...
//	Decimal: 27 38 y c1 c2 [x1 data...]...
//
// Parameters:
//
//	y  - Number of bytes in the vertical direction for each column (1 byte).
//	     Each column is described by y bytes (little-endian vertical bit order).
//	c1 - First character code to define (inclusive).
//	c2 - Last character code to define (inclusive).
//	defs - Slice of per-character definitions in order for codes c1..c2.
//	       Each definition is encoded as: 1 byte width xi, followed by y*xi bytes
//	       of column data (column-major).
//
// Range / Notes:
//   - Typical y values are model/font dependent (e.g., 3 for 12x24 or 9x17 fonts).
//   - c1 and c2 typically in the printable range (32..126) depending on model.
//   - For each character i between c1 and c2, defs[i - int(c1)] must contain
//     the width byte and exactly y*width data bytes.
//   - Existing user-defined characters for the specified codes are replaced.
//   - Definitions persist until cleared (ESC ?, ESC @), reset, or power-off.
//   - To use defined glyphs, send ESC % 1 (select user-defined character set).
//
// Byte sequence:
//
//	ESC & y c1 c2 [x1 data...]...[xk data...] -> 0x1B, 0x26, y, c1, c2, ...
func (udc *UserDefinedCommands) DefineUserDefinedCharacters(y, c1, c2 byte, definitions []UserDefinedChar) ([]byte, error) {
	// Validation
	if y == 0 {
		return nil, ErrInvalidYValue
	}
	if c1 < 32 || c1 > 126 {
		return nil, fmt.Errorf("%w: c1=%d", ErrInvalidCharacterCode, c1)
	}
	if c2 < c1 || c2 > 126 {
		return nil, fmt.Errorf("%w: c2=%d", ErrInvalidCodeRange, c2)
	}
	expected := int(c2 - c1 + 1)
	if len(definitions) != expected {
		return nil, fmt.Errorf("%w: got %d, expected %d", ErrInvalidDefinition, len(definitions), expected)
	}

	// Build command
	seq := []byte{0x1B, 0x26, y, c1, c2}
	bytesPerCol := int(y)

	for idx, def := range definitions {
		// Width validation is printer & font dependent
		if def.Width == 0 {
			// Zero width is allowed (blank char)
			seq = append(seq, def.Width)
			continue
		}
		expectedDataLen := bytesPerCol * int(def.Width)
		if len(def.Data) != expectedDataLen {
			return nil, fmt.Errorf("%w: char %d has %d bytes, expected %d",
				ErrInvalidDataLength, int(c1)+idx, len(def.Data), expectedDataLen)
		}
		seq = append(seq, def.Width)
		seq = append(seq, def.Data...)
	}

	return seq, nil
}

// CancelUserDefinedCharacter deletes (cancels) a user-defined character.
//
// Format:
//
//	ASCII: ESC ? n
//	Hex:   0x1B 0x3F n
//	Decimal: 27 63 n
//
// Range:
//
//	n = 32–126
//
// Default:
//
//	None
//
// Description:
//
//	Deletes the user-defined character pattern specified by character code n.
//	After cancellation, the resident (built-in) character for that code is printed.
//
// Notes:
//   - This command can cancel user-defined characters per font. Select the font
//     with ESC ! or ESC M before issuing this command if needed.
//   - Settings take effect immediately; the deleted definition remains cleared
//     until redefined, ESC @ (initialize), power-off, or reset.
//
// Byte sequence:
//
//	ESC ? n -> 0x1B, 0x3F, n
func (udc *UserDefinedCommands) CancelUserDefinedCharacter(n byte) ([]byte, error) {
	if n < 32 || n > 126 {
		return nil, ErrInvalidCharacterCode
	}
	return []byte{0x1B, 0x3F, n}, nil
}
