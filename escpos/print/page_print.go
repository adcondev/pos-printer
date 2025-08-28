package print

import "github.com/adcondev/pos-printer/escpos/common"

const (
	// CAN
	//
	// Format:
	//   ASCII: CAN
	//   Hex:   0x18
	//   Decimal: 24
	//
	// Description:
	//   In Page mode, this command deletes all the print data in the current
	//   print area.
	//
	// Notes:
	//   - This command is enabled only in Page mode. Page mode can be selected
	//     by ESC L.
	//   - If data set in a previously specified print area overlaps the currently
	//     specified print area, that data is deleted as well.
	//   - This command has no effect in Standard mode.
	//
	// Value:
	//   The CAN control code is a single byte 0x18 (decimal 24).
	CAN byte = 0x18
)

// ============================================================================
// Interface Definitions
// ============================================================================

// PageModeCapability defines the interface for page mode capabilities
type PageModeCapability interface {
	ReverseCapability
	PageCapability
}

// ReverseCapability defines the interface for reverse printing capabilities
type ReverseCapability interface {
	PrintAndReverseFeed(n byte) ([]byte, error)
	PrintAndReverseFeedLines(n byte) ([]byte, error)
}

// PageCapability defines the interface for page operations capabilities
type PageCapability interface {
	CancelData() []byte
	PrintDataInPageMode() []byte
	PrintAndFeedLines(n byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// PagePrint groups page mode printing commands
type PagePrint struct{}

// PrintDataInPageMode prints the data in the print buffer collectively when
// the printer is in Page mode.
//
// Text:
//
//	ASCII: ESC FF
//	Hex:   0x1B 0x0C
//	Decimal: 27 12
//
// Description:
//
//	In Page mode, this command prints all the data currently buffered in the
//	print area collectively.
//
// Notes:
//   - This command is enabled only in Page mode. Page mode can be selected by
//     ESC L.
//   - After printing with ESC FF, the printer does NOT clear the buffered
//     data, the print position, or values set by other commands — the buffer
//     remains intact for repeated printing.
//   - The printer returns to Standard mode when FF (in Page mode), ESC S, or
//     ESC @ is issued. If the return to Standard mode is caused by ESC @,
//     all settings are canceled.
//   - This command is commonly used for printing the same Page-mode data
//     multiple times without rebuilding the page buffer.
//
// Byte sequence:
//
//	ESC FF -> 0x1B, 0x0C
func (pp *PagePrint) PrintDataInPageMode() []byte {
	return []byte{common.ESC, FF}
}

// PrintAndReverseFeed prints the data in the print buffer and feeds the
// paper in the reverse direction by n × (vertical or horizontal motion unit).
//
// Text:
//
//	ASCII: ESC K n
//	Hex:   0x1B 0x4B n
//	Decimal: 27 75 n
//
// Range:
//
//	n = 0–48
//
// Default:
//
//	None
//
// Description:
//
//	Prints the data in the print buffer and feeds the paper n × (vertical or
//	horizontal motion unit) in the reverse direction.
//
// Notes:
//   - In Standard mode the vertical motion unit is used.
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T.
//   - When the starting position is set to the upper-left or lower-right of
//     the print area using ESC T, the vertical motion unit is used.
//   - When the starting position is set to the upper-right or lower-left of
//     the print area using ESC T, the horizontal motion unit is used.
//   - When this command is processed in Page mode, only the print position
//     moves; the printer does not perform actual printing.
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - The maximum paper reverse-feed amount depends on the printer model.
//     If n is specified greater than the model maximum, the reverse feed is
//     not executed although the printing is executed.
//   - This command is used to temporarily feed a specific length without
//     changing the line spacing set by other commands.
//   - Some printers perform a small forward paper feed after a reverse feed
//     due to mechanical restrictions.
//
// Byte sequence:
//
//	ESC K n -> 0x1B, 0x4B, n
func (pp *PagePrint) PrintAndReverseFeed(n byte) ([]byte, error) {
	if n > common.MaxReverseMotionUnits {
		return nil, common.ErrPrintReverseFeed
	}
	return []byte{common.ESC, 'K', n}, nil
}

// PrintAndReverseFeedLines prints the data in the print buffer and feeds n
// lines in the reverse direction.
//
// Text:
//
//	ASCII: ESC e n
//	Hex:   0x1B 0x65 n
//	Decimal: 27 101 n
//
// Range:
//
//	n = 0–2
//
// Default:
//
//	None
//
// Description:
//
//	Prints the data in the print buffer and feeds n lines in the reverse
//	direction.
//
// Notes:
//   - The amount of paper fed per line is based on the value set using the
//     line spacing command (ESC 2 or ESC 3).
//   - The maximum paper reverse-feed amount depends on the printer model.
//     If n is specified greater than the model maximum, the reverse feed is
//     not executed although the printing is executed.
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing.
//   - This command is used to temporarily feed a specific number of lines
//     without changing the line spacing set by other commands.
//   - Some printers perform a small forward paper feed after a reverse feed
//     due to mechanical restrictions.
//
// Byte sequence:
//
//	ESC e n -> 0x1B, 0x65, n
func (pp *PagePrint) PrintAndReverseFeedLines(n byte) ([]byte, error) {
	if n > common.MaxReverseFeedLines {
		return nil, common.ErrPrintReverseFeedLines
	}
	return []byte{common.ESC, 'e', n}, nil
}

// PrintAndFeedLines prints the data in the print buffer and feeds n lines.
//
// Text:
//
//	ASCII: ESC d n
//	Hex:   0x1B 0x64 n
//	Decimal: 27 100 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	None
//
// Description:
//
//	Prints the data in the print buffer and feeds n lines.
//
// Notes:
//   - The amount of paper fed per line is based on the value set using the
//     line spacing command (ESC 2 or ESC 3).
//   - The maximum paper feed amount depends on the printer model. If n is
//     specified greater than the model maximum, the printer executes the
//     maximum paper feed amount.
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing.
//   - This command is used to temporarily feed a specific number of lines
//     without changing the line spacing set by other commands.
//
// Byte sequence:
//
//	ESC d n -> 0x1B, 0x64, n
func (pp *PagePrint) PrintAndFeedLines(n byte) ([]byte, error) {
	return []byte{common.ESC, 'd', n}, nil
}

func (pp *PagePrint) CancelData() []byte {
	return []byte{CAN}
}
