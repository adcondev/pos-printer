package escpos

// Interface compliance check
var _ LineSpacingCapability = (*LineSpacingCommands)(nil)

// LineSpacingCapability defines the interface for line spacing commands in ESC/POS printers.
type LineSpacingCapability interface {
	SetLineSpacing(n byte) []byte
	SelectDefaultLineSpacing() []byte
}

// LineSpacingCommands implements the LineSpacingCapability interface for ESC/POS printers.
type LineSpacingCommands struct{}

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
//   - The maximum line spacing is 1016 mm (40 inches); actual maximum may be
//     smaller depending on the printer model. If the specified amount
//     exceeds the model maximum, the line spacing is automatically set to
//     the maximum supported value.
//   - In Standard mode the vertical motion unit is used.
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T.
//   - When the starting position is set to the upper-left or lower-right of
//     the print area using ESC T, the vertical motion unit is used.
//   - When the starting position is set to the upper-right or lower-left of
//     the print area using ESC T, the horizontal motion unit is used.
//   - Line spacing can be set independently in Standard mode and in Page
//     mode; this command affects the spacing for the currently selected
//     mode (Standard or Page).
//   - If the motion unit is changed after the line spacing is set, the
//     numeric line spacing value does not change (the unit of measure for
//     that numeric value changes).
//   - The selected line spacing remains in effect until one of the
//     following occurs: ESC 2 is executed, ESC @ is executed, the printer
//     is reset, or power is turned off.
//
// Byte sequence:
//
//	ESC 3 n -> 0x1B, 0x33, n
func (lsc *LineSpacingCommands) SetLineSpacing(n byte) []byte {
	return []byte{ESC, '3', n}
}

// SelectDefaultLineSpacing sets the line spacing to the printer's "default
// line spacing".
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
//   - Line spacing can be set independently in Standard mode and in Page mode.
//     In Standard mode this command sets the line spacing used by Standard
//     mode; in Page mode it sets the line spacing used by Page mode.
//   - The selected line spacing remains in effect until one of the following
//     occurs: ESC 3 is executed, ESC @ is executed, the printer is reset, or
//     power is turned off.
//
// Byte sequence:
//
//	ESC 2 -> 0x1B, 0x32
func (lsc *LineSpacingCommands) SelectDefaultLineSpacing() []byte {
	return []byte{ESC, '2'}
}
