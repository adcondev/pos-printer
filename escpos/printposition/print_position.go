package printposition

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters
const (
	// HT moves the print position to the next horizontal tab position.
	HT byte = common.HT // 0x09
)

// Justification modes
const (
	JustifyLeft        byte = 0x00 // n = 0
	JustifyCenter      byte = 0x01 // n = 1
	JustifyRight       byte = 0x02 // n = 2
	JustifyLeftASCII   byte = '0'  // n = 48
	JustifyCenterASCII byte = '1'  // n = 49
	JustifyRightASCII  byte = '2'  // n = 50
)

// Print direction modes (Page mode)
const (
	DirectionLeftToRight byte = 0x00 // n = 0 (upper left start)
	DirectionBottomToTop byte = 0x01 // n = 1 (lower left start)
	DirectionRightToLeft byte = 0x02 // n = 2 (lower right start)
	DirectionTopToBottom byte = 0x03 // n = 3 (upper right start)

	DirectionLeftToRightASCII byte = '0' // n = 48
	DirectionBottomToTopASCII byte = '1' // n = 49
	DirectionRightToLeftASCII byte = '2' // n = 50
	DirectionTopToBottomASCII byte = '3' // n = 51
)

// Beginning of line operations
const (
	BeginLineErase byte = 0x00 // n = 0 (erase buffer)
	BeginLinePrint byte = 0x01 // n = 1 (print buffer)

	BeginLineEraseASCII byte = '0' // n = 48
	BeginLinePrintASCII byte = '1' // n = 49
)

// Tab position limits
const (
	MaxTabPositions = 32
	MaxTabValue     = 255
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrInvalidJustification  = fmt.Errorf("invalid justification mode (try 0-2 or '0'..'2')")
	ErrInvalidPrintDirection = fmt.Errorf("invalid print direction (try 0-3 or '0'..'3')")
	ErrInvalidBeginLineMode  = fmt.Errorf("invalid begin line mode (try 0-1 or '0'..'1')")
	ErrTooManyTabPositions   = fmt.Errorf("too many tab positions (max %d)", MaxTabPositions)
	ErrInvalidTabPosition    = fmt.Errorf("invalid tab position (must be 1-255 in ascending order)")
	ErrInvalidPrintAreaSize  = fmt.Errorf("invalid print area size (width and height must be >= 1)")
	ErrInvalidPosition       = fmt.Errorf("invalid position value")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for print position commands
type Capability interface {
	// Basic positioning
	SetAbsolutePrintPosition(position uint16) []byte
	SetRelativePrintPosition(distance int16) []byte
	HorizontalTab() []byte
	SetHorizontalTabPositions(positions []byte) ([]byte, error)

	// Justification
	SelectJustification(mode byte) ([]byte, error)

	// Margins and print area
	SetLeftMargin(margin uint16) []byte
	SetPrintAreaWidth(width uint16) []byte
	SetPrintPositionBeginningLine(mode byte) ([]byte, error)

	// Page mode specific
	SelectPrintDirectionPageMode(direction byte) ([]byte, error)
	SetPrintAreaPageMode(x, y, width, height uint16) []byte
	SetAbsoluteVerticalPrintPosition(position uint16) []byte
	SetRelativeVerticalPrintPosition(distance int16) []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for print position commands
type Commands struct{}

func NewCommands() *Commands {
	return &Commands{}
}

// HorizontalTab moves the print position to the next horizontal tab position.
//
// Format:
//
//	ASCII: HT
//	Hex:   0x09
//	Decimal: 9
//
// Description:
//
//	Moves the print position to the next horizontal tab position.
//
// Notes:
//   - Ignored unless the next horizontal tab position has been set (ESC D).
//   - If the next tab position exceeds the print area, print position is set to [Print area width + 1].
//   - If processed when at [Print area width + 1], the printer executes print-buffer-full for the current line and performs horizontal tab processing from the beginning of the next line.
//     In Page mode, printing is not executed but the print position is moved.
//   - The printer will not move to the beginning of the line by executing this command.
//   - When underline mode is on, the underline is not printed under the tab space skipped by this command.
//
// Byte sequence:
//
//	HT -> 0x09
func (c *Commands) HorizontalTab() []byte {
	return []byte{HT}
}

// SetAbsolutePrintPosition sets the absolute print position.
//
// Format:
//
//	ASCII: ESC $ nL nH
//	Hex:   0x1B 0x24 nL nH
//	Decimal: 27 36 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	None
//
// Description:
//
//	Moves the print position to (nL + nH × 256) × (horizontal or vertical motion unit)
//	from the left edge of the print area.
//
// Notes:
//   - The printer ignores any setting that exceeds the print area.
//   - In Standard mode the horizontal motion unit is used.
//   - In Page mode the horizontal or vertical motion unit is used depending on
//     the print direction set by ESC T.
//   - If the starting position is set to upper-left or lower-right using ESC T,
//     the horizontal motion unit is used; for upper-right or lower-left the vertical
//     motion unit is used.
//   - If the motion unit changes after this command, the print position is not changed.
//   - The printer will not move to the beginning of the line by executing this command.
//   - When underline mode is on, the underline is not printed under the space skipped
//     by this command.
//
// Byte sequence:
//
//	ESC $ nL nH -> 0x1B, 0x24, nL, nH
func (c *Commands) SetAbsolutePrintPosition(position uint16) []byte {
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	return []byte{common.ESC, '$', nL, nH}
}

// SetHorizontalTabPositions sets horizontal tab positions.
//
// Format:
//
//	ASCII: ESC D n1 ... nk NUL
//	Hex:   0x1B 0x44 n1 ... nk 0x00
//	Decimal: 27 68 n1 ... nk 0
//
// Range:
//
//	n = 1–255
//	k = 0–32
//
// Default:
//
//	n = 8, 16, 24, 32, ..., 232, 240, 248 (every eight characters for the default font set)
//
// Description:
//
//	Sets horizontal tab positions. Each transmitted n value specifies the number
//	of character widths from the line start to the tab stop. Transmit the tab
//	stops in ascending order and terminate the list with NUL (0x00). Transmitting
//	ESC D NUL clears all horizontal tab positions.
//
// Notes:
//   - The tab position is stored as [character width × n], where character width
//     includes right-side character spacing. Double-width characters count as
//     twice the width.
//   - Character width and font/spacing/enlargement should be set before sending this command.
//   - A maximum of 32 tab positions can be set; data beyond 32 is treated as normal data.
//   - If a transmitted n is less than or equal to the previous value, tab-setting
//     is finished and subsequent bytes are processed as normal data.
//   - Tab settings are preserved until ESC @ (initialize), printer reset, or power-off.
//   - Changing the left margin will shift stored tab positions accordingly.
//   - Horizontal tab positions that exceed the print area are allowed; they become
//     effective or not depending on the current print area width.
//
// Byte sequence:
//
//	ESC D n1 ... nk NUL -> 0x1B, 0x44, n1, ..., nk, 0x00
func (c *Commands) SetHorizontalTabPositions(positions []byte) ([]byte, error) {
	// Check maximum number of positions
	if len(positions) > MaxTabPositions {
		return nil, ErrTooManyTabPositions
	}

	// Validate positions are in ascending order and within range
	prevPos := byte(0)
	for i, pos := range positions {
		if pos == 0 || pos > MaxTabValue {
			return nil, fmt.Errorf("%w: position %d at index %d", ErrInvalidTabPosition, pos, i)
		}
		if pos <= prevPos {
			return nil, fmt.Errorf("%w: position %d at index %d must be greater than %d", ErrInvalidTabPosition, pos, i, prevPos)
		}
		prevPos = pos
	}

	// Build command
	cmd := []byte{common.ESC, 'D'}
	cmd = append(cmd, positions...)
	cmd = append(cmd, common.NUL)
	return cmd, nil
}

// SelectPrintDirectionPageMode selects the print direction and starting position in Page mode.
//
// Format:
//
//	ASCII: ESC T n
//	Hex:   0x1B 0x54 n
//	Decimal: 27 84 n
//
// Range:
//
//	n = 0–3, 48–51
//
// Default:
//
//	n = 0 (Left to right, starting position: upper left)
//
// Description:
//
//	In Page mode, selects print direction and starting position as follows:
//	  n = 0 or 48 -> Print direction: left to right;  Starting position: upper left
//	  n = 1 or 49 -> Print direction: bottom to top;  Starting position: lower left
//	  n = 2 or 50 -> Print direction: right to left;  Starting position: lower right
//	  n = 3 or 51 -> Print direction: top to bottom;  Starting position: upper right
//
// Notes:
//   - Effective only in Page mode; has no effect in Standard mode.
//   - The meaning of horizontal/vertical motion units for other commands depends on the selected starting position (see command reference).
//   - Settings persist until ESC @ (initialize), printer reset, or power-off.
//
// Byte sequence:
//
//	ESC T n -> 0x1B, 0x54, n
func (c *Commands) SelectPrintDirectionPageMode(direction byte) ([]byte, error) {
	// Validate allowed values
	switch direction {
	case 0, 1, 2, 3, '0', '1', '2', '3':
		// Valid values
	default:
		return nil, ErrInvalidPrintDirection
	}
	return []byte{common.ESC, 'T', direction}, nil
}

// SetPrintAreaPageMode sets the print area and logical origin in Page mode.
//
// Format:
//
//	ASCII: ESC W xL xH yL yH dxL dxH dyL dyH
//	Hex:   0x1B 0x57 xL xH yL yH dxL dxH dyL dyH
//	Decimal: 27 87 xL xH yL yH dxL dxH dyL dyH
//
// Description:
//
//	In Page mode, defines the logical origin (horizontal and vertical) and the
//	print area size. The transmitted parameters are interpreted as 16-bit
//	little-endian values: value = (low + high * 256), each measured in the
//	horizontal or vertical motion unit as appropriate.
//
//	  Horizontal logical origin = (xL + xH*256) × (horizontal motion unit)
//	  Vertical logical origin   = (yL + yH*256) × (vertical motion unit)
//	  Print area width          = (dxL + dxH*256) × (horizontal motion unit)
//	  Print area height         = (dyL + dyH*256) × (vertical motion unit)
//
// Notes:
//   - This command only has effect in Page mode (ESC L) and is ignored in Standard mode.
//   - Both print area width and height must be at least 1 (cannot be zero).
//   - Logical origins must lie within the printable area.
//   - If origin + size exceeds the printable area, the size is reduced to fit.
//   - Values are fixed even if motion units change later.
//   - Settings persist until FF (in Page mode), ESC @ (initialize), reset, or power-off.
//   - The absolute origin is the upper-left of the printable area.
//   - For printers supporting GS ( P <Function 48>, the maximum and origin align with that printable area setting.
//
// Byte sequence:
//
//	ESC W xL xH yL yH dxL dxH dyL dyH -> 0x1B, 0x57, xL, xH, yL, yH, dxL, dxH, dyL, dyH
func (c *Commands) SetPrintAreaPageMode(x, y, width, height uint16) []byte {
	xL := byte(x & 0xFF)
	xH := byte((x >> 8) & 0xFF)
	yL := byte(y & 0xFF)
	yH := byte((y >> 8) & 0xFF)
	dxL := byte(width & 0xFF)
	dxH := byte((width >> 8) & 0xFF)
	dyL := byte(height & 0xFF)
	dyH := byte((height >> 8) & 0xFF)

	return []byte{common.ESC, 'W', xL, xH, yL, yH, dxL, dxH, dyL, dyH}
}

// SetRelativePrintPosition moves the print position relative to the current position.
//
// Format:
//
//	ASCII: ESC \ nL nH
//	Hex:   0x1B 0x5C nL nH
//	Decimal: 27 92 nL nH
//
// Range:
//
//	(nL + nH × 256) = -32768 – 32767 (signed 16-bit value; low + high*256 interpreted as int16)
//
// Default:
//
//	None
//
// Description:
//
//	Moves the print position by (nL + nH × 256) × (horizontal or vertical motion unit)
//	from the current position. Positive moves to the right; negative moves to the left.
//
// Notes:
//   - The printer ignores any setting that exceeds the print area.
//   - In Standard mode the horizontal motion unit is used.
//   - In Page mode the horizontal or vertical motion unit is used depending on the print direction set by ESC T.
//   - If the starting position is upper-left or lower-right (ESC T), the horizontal motion unit is used.
//     If the starting position is upper-right or lower-left, the vertical motion unit is used.
//   - Changing motion units after this command does not change the already-set print position.
//   - Underline mode does not print under the space skipped by this command.
//   - In JIS code, '\' corresponds to '¥'.
//
// Byte sequence:
//
//	ESC \ nL nH -> 0x1B, 0x5C, nL, nH
func (c *Commands) SetRelativePrintPosition(distance int16) []byte {
	// Convert signed int16 to unsigned bytes (little-endian)
	// intentional: preserve int16 two's-complement bit pattern for ESC \ command
	value := uint16(distance) // nolint:gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	return []byte{common.ESC, '\\', nL, nH}
}

// SelectJustification selects text justification in Standard mode.
//
// Format:
//
//	ASCII: ESC a n
//	Hex:   0x1B 0x61 n
//	Decimal: 27 97 n
//
// Range:
//
//	n = 0–2, 48–50
//
// Default:
//
//	n = 0 (Left)
//
// Description:
//
//	In Standard mode, aligns all data in one line according to n:
//	  0 or 48 -> Left justification
//	  1 or 49 -> Centered
//	  2 or 50 -> Right justification
//
// Notes:
//   - Effective only in Standard mode and only when processed at the beginning of a line.
//   - Has no effect in Page mode.
//   - Justification is applied within the print area set by GS L and ESC W/GS W.
//   - Affects characters, graphics, barcodes, 2D codes and space areas set by HT, ESC $, ESC \.
//   - Setting persists until ESC @, reset, or power-off.
//
// Byte sequence:
//
//	ESC a n -> 0x1B, 0x61, n
func (c *Commands) SelectJustification(mode byte) ([]byte, error) {
	// Validate allowed values
	switch mode {
	case 0, 1, 2, '0', '1', '2':
		// Valid values
	default:
		return nil, ErrInvalidJustification
	}
	return []byte{common.ESC, 'a', mode}, nil
}

// SetAbsoluteVerticalPrintPosition sets the absolute vertical print position in Page mode.
//
// Format:
//
//	ASCII: GS $ nL nH
//	Hex:   0x1D 0x24 nL nH
//	Decimal: 29 36 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	None
//
// Description:
//
//	In Page mode, moves the vertical print position to (nL + nH × 256) × (vertical or horizontal motion unit)
//	from the starting position set by ESC T.
//
// Notes:
//   - This command is enabled only in Page mode; it is ignored in Standard mode.
//   - The printer ignores any setting that exceeds the print area set by ESC W.
//   - The horizontal or vertical motion unit used depends on the print direction set by ESC T.
//   - If the starting position is upper left or lower right, the vertical motion unit is used.
//     If the starting position is upper right or lower left, the horizontal motion unit is used.
//   - Changing motion units after this command does not change the already-set print position.
//
// Byte sequence:
//
//	GS $ nL nH -> 0x1D, 0x24, nL, nH
func (c *Commands) SetAbsoluteVerticalPrintPosition(position uint16) []byte {
	nL := byte(position & 0xFF)
	nH := byte((position >> 8) & 0xFF)
	return []byte{common.GS, '$', nL, nH}
}

// SetLeftMargin sets the left margin in Standard mode.
//
// Format:
//
//	ASCII: GS L nL nH
//	Hex:   0x1D 0x4C nL nH
//	Decimal: 29 76 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	(nL + nH × 256) = 0
//
// Description:
//
//	In Standard mode, sets the left margin to (nL + nH × 256) × (horizontal motion unit)
//	from the left edge of the printable area.
//
// Notes:
//   - Effective in Standard mode only when processed at the beginning of the line.
//   - Has no effect while in Page mode; if issued in Page mode the value is stored
//     and enabled when returning to Standard mode.
//   - If the setting exceeds the printable area, it is clamped to the printable-area maximum.
//   - If this command and GS W would set the print area width to less than one character,
//     the print area width is extended to accommodate one character.
//   - Uses the horizontal motion unit; changes to the motion unit after setting do not change the margin.
//   - Setting persists until ESC @ (initialize), printer reset, or power-off.
//   - The left margin is measured from the left edge of the printable area; changing the left margin moves that edge.
//
// Byte sequence:
//
//	GS L nL nH -> 0x1D, 0x4C, nL, nH
func (c *Commands) SetLeftMargin(margin uint16) []byte {
	nL := byte(margin & 0xFF)
	nH := byte((margin >> 8) & 0xFF)
	return []byte{common.GS, 'L', nL, nH}
}

// SetPrintPositionBeginningLine moves the print position to the beginning of the print line.
//
// Format:
//
//	ASCII: GS T n
//	Hex:   0x1D 0x54 n
//	Decimal: 29 84 n
//
// Range:
//
//	n = 0, 1, 48, 49
//
// Default:
//
//	None
//
// Description:
//
//	In Standard mode, moves the print position to the beginning (left side) of the printable area
//	after performing the operation specified by n. n controls how the print buffer is processed:
//
//	  n = 0 or 48 -> Erase the data in the print buffer, then move the print position.
//	  n = 1 or 49 -> Print the data in the print buffer, then move the print position (starts a new line based on line spacing).
//
// Notes:
//   - Effective only in Standard mode; ignored in Page mode.
//   - Ignored if the print position is already at the beginning of the line.
//   - If print position is not at the beginning of the line and n = 1 or 49, this behaves the same as LF.
//   - Erase (n = 0 or 48) cancels the current print-buffered data but preserves other settings and buffer contents.
//   - After execution the printer is in the "Beginning of the line" status.
//   - Use this command immediately before other commands that require beginning-of-line to ensure they execute.
//
// Byte sequence:
//
//	GS T n -> 0x1D, 0x54, n
func (c *Commands) SetPrintPositionBeginningLine(mode byte) ([]byte, error) {
	// Validate allowed values
	switch mode {
	case 0, 1, '0', '1':
		// Valid values
	default:
		return nil, ErrInvalidBeginLineMode
	}
	return []byte{common.GS, 'T', mode}, nil
}

// SetPrintAreaWidth sets the print area width in Standard mode.
//
// Format:
//
//	ASCII: GS W nL nH
//	Hex:   0x1D 0x57 nL nH
//	Decimal: 29 87 nL nH
//
// Range:
//
//	(nL + nH × 256) = 0 – 65535
//
// Default:
//
//	Entire printable area (model-dependent; see printer specs)
//
// Description:
//
//	In Standard mode, sets the print area width to (nL + nH × 256) × (horizontal motion unit).
//
// Notes:
//   - This command is effective in Standard mode only when processed at the beginning of a line.
//   - Has no effect in Page mode; if issued in Page mode the value is stored and enabled when returning to Standard mode.
//   - If [left margin + print area width] exceeds the printable area, the print area width is clamped to [printable area - left margin].
//   - If this command together with GS L would set the print area width to less than one character, the width is extended to accommodate one character.
//   - Uses the horizontal motion unit. Changing the motion unit after setting does not change the stored width.
//   - Setting persists until ESC @, printer reset, or power-off.
//
// Model example defaults (n = nL + nH*256):
//
//	80 mm paper -> n = 576
//	58 mm paper -> n = 420
//
// Byte sequence:
//
//	GS W nL nH -> 0x1D, 0x57, nL, nH
func (c *Commands) SetPrintAreaWidth(width uint16) []byte {
	nL := byte(width & 0xFF)
	nH := byte((width >> 8) & 0xFF)
	return []byte{common.GS, 'W', nL, nH}
}

// SetRelativeVerticalPrintPosition moves the vertical print position relative to the current position in Page mode.
//
// Format:
//
//	ASCII: GS \ nL nH
//	Hex:   0x1D 0x5C nL nH
//	Decimal: 29 92 nL nH
//
// Range:
//
//	(nL + nH × 256) = -32768 – 32767  (signed 16-bit; low + high*256 interpreted as int16)
//
// Default:
//
//	None
//
// Description:
//
//	In Page mode, moves the vertical print position by (nL + nH × 256) × (vertical or horizontal motion unit)
//	from the current position. A positive value moves downward; a negative value moves upward.
//
// Notes:
//   - This command is enabled only in Page mode; it is ignored in Standard mode.
//   - The printer ignores any setting that exceeds the print area set by ESC W.
//   - The horizontal or vertical motion unit used depends on the print direction set by ESC T.
//   - If starting position is upper-left or lower-right (ESC T), the vertical motion unit is used.
//   - If starting position is upper-right or lower-left (ESC T), the horizontal motion unit is used.
//   - Changes to the motion units after executing this command do not affect the already-set print position.
//   - In JIS code, '\' corresponds to '¥'.
//
// Byte sequence:
//
//	GS \ nL nH -> 0x1D, 0x5C, nL, nH
func (c *Commands) SetRelativeVerticalPrintPosition(distance int16) []byte {
	// Convert signed int16 to unsigned bytes (little-endian)
	// intentional: preserve int16 two's-complement bit pattern for ESC \ command
	value := uint16(distance) // nolint:gosec
	nL := byte(value & 0xFF)
	nH := byte((value >> 8) & 0xFF)
	return []byte{common.GS, '\\', nL, nH}
}
