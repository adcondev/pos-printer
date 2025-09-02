package print

import (
	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters used in printing commands
const (
	// LF prints the data in the print buffer and feeds one line,
	// based on the current line spacing.
	//
	// Format:
	//   ASCII: LF
	//   Hex:   0x0A
	//   Decimal: 10
	//
	// Description:
	//   Prints the data in the print buffer and feeds one line, according to the
	//   currently selected line spacing.
	//
	// Notes:
	//   - The amount of paper fed per line is based on the value set using the
	//     line spacing command (ESC 2 or ESC 3).
	//   - After printing, the print position is moved to the left side of the
	//     printable area and the printer enters the "Beginning of the line"
	//     status.
	//   - When this command is processed in Page mode, only the print position
	//     moves and the printer does not perform actual printing.
	//
	// Value is the single-byte LF (line feed) control code.
	LF byte = 0x0A
	// CR executes a carriage-return operation and (depending
	// on the printer type and auto-line-feed setting) may print and cause a line feed.
	//
	// Format:
	//   ASCII: CR
	//   Hex:   0x0D
	//   Decimal: 13
	//
	// Description:
	//   Executes one of the following operations depending on the print head and
	//   auto-line-feed state:
	//
	//   - Horizontal-alignment heads (Line thermal head or Shuttle head):
	//     * When auto line feed is enabled: Executes printing and one line feed
	//       as if LF was issued.
	//     * When auto line feed is disabled: The command is ignored.
	//
	//   - Vertical-alignment heads (Serial dot head):
	//     * Executes printing and one line feed as if LF was issued.
	//
	//   - In Standard mode:
	//     * Prints the data in the print buffer (when applicable) and moves the
	//       print position to the beginning of the print line (left side of the
	//       printable area).
	//   - In Page mode:
	//     * Moves the print position to the beginning of the print line. When
	//       processed in Page mode, the printer does not perform actual printing.
	//
	// Notes:
	//   - With a serial interface, this command behaves as if auto line feed is
	//     disabled.
	//   - Enabling or disabling auto line feed may be controlled by a DIP
	//     switch or the memory switch. The memory switch can be changed via
	//     GS ( E <Function 3>.
	//   - After printing, the print position is moved to the left side of the
	//     printable area and the printer enters the "Beginning of the line"
	//     status (when printing occurs).
	//
	// Value:
	//   The CR control code is a single byte 0x0D (decimal 13).
	CR byte = 0x0D
	// FF (Form Feed) — behaviour in Standard mode and Page mode.
	//
	// Format:
	//   ASCII: FF
	//   Hex:   0x0C
	//   Decimal: 12
	//
	// Summary:
	//   FF is the single-byte Form Feed control code. Its behaviour depends on
	//   the current printer mode.
	//
	// FF (in Page mode)
	// Name:
	//   Print and return to Standard mode (in Page mode)
	//
	// Description:
	//   In Page mode, FF prints all the data in the print buffer collectively
	//   and switches the printer from Page mode to Standard mode.
	//
	// Notes:
	//   - This command is enabled only in Page mode. Page mode can be selected
	//     by ESC L.
	//   - The data in the print area is deleted after being printed.
	//   - This command resets the values set by ESC W to their defaults.
	//   - The value set by ESC T is maintained.
	//   - After printing, the printer returns to Standard mode, the print
	//     position moves to the left side of the printable area, and the printer
	//     enters the "Beginning of the line" status.
	//
	// FF (in Standard mode)
	// Name:
	//   End job (in Standard mode)
	//
	// Description:
	//   In Standard mode, FF indicates that "Printing is completed" for the
	//   current job. This signals that a series of printing actions has been
	//   completed and subsequent data will be printed as a new, separate job.
	//
	// Notes:
	//   - This command is enabled only in Standard mode. Standard mode can be
	//     selected by ESC S.
	//   - When the cutting position has been reserved by GS V <Function C> or
	//     GS ( V <Function 51>, issuing FF will feed the paper to the cutting
	//     position and perform the cut.
	//
	// Value:
	//   The FF control code is a single byte 0x0C (decimal 12).
	FF byte = 0x0C
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for printer capabilities
type Capability interface {
	Text(text string) ([]byte, error)
	PrintAndFeedPaper(units byte) []byte
	FormFeed() []byte
	PrintAndCarriageReturn() []byte
	PrintAndLineFeed() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Formatting replaces specific characters in the byte slice with their ESC/POS equivalents.
func Formatting(data []byte) []byte {
	for i := range data {
		switch data[i] {
		case '\n':
			data[i] = LF
		case '\t':
			data[i] = common.HT
		case '\r':
			data[i] = CR
		default:
			// Do nothing, keep the character as is
		}
	}
	return data
}

// Commands groups printing-related commands
type Commands struct {
	Page PageModeCapability
}

func NewCommands() *Commands {
	return &Commands{
		Page: NewPagePrint(),
	}
}

// Text formats and sends a string for printing
func (pc *Commands) Text(n string) ([]byte, error) {
	if err := common.IsBufLenOk([]byte(n)); err != nil {
		return nil, err
	}
	return Formatting([]byte(n)), nil
}

// PrintAndFeedPaper prints the data in the print buffer and feeds the paper
// by input × (vertical or horizontal motion unit).
//
// Text:
//
//	ASCII: ESC J input
//	Hex:   0x1B 0x4A input
//	Decimal: 27 74 input
//
// Range:
//
//	input = 0–255
//
// Default:
//
//	None
//
// Description:
//
//	Prints the data in the print buffer and feeds the paper [input × (vertical or
//	horizontal motion unit)].
//
// Notes:
//   - In Standard mode the vertical motion unit is used.
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T.
//   - When the starting position is set to the upper-left or lower-right of
//     the print area using ESC T, the vertical motion unit is used.
//   - When the starting position is set to the upper-right or lower-left of
//     the print area using ESC T, the horizontal motion unit is used.
//   - The maximum paper feed amount depends on the printer model. If input is
//     specified greater than the model maximum, the printer executes the
//     maximum paper feed amount.
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing.
//   - This command is used to temporarily feed a specific length without
//     changing the line spacing set by other commands.
//
// Byte sequence:
//
//	ESC J input -> 0x1B, 0x4A, input
func (pc *Commands) PrintAndFeedPaper(n byte) []byte {
	return []byte{common.ESC, 'J', n}
}

// FormFeed (Form Feed) — behaviour in Standard mode and Page mode.
//
// Formatting:
//
//	ASCII: FF
//	Hex:   0x0C
//	Decimal: 12
//
// Summary:
//
//	FF is the single-byte Form Feed control code. Its behaviour depends on
//	the current printer mode.
//
// FF (in Standard mode)
// Name:
//
//	End job
//
// Description:
//
//	In Standard mode, FF indicates that "Printing is completed" for the
//	current job. This signals that a series of printing actions has been
//	completed and subsequent data will be printed as a new, separate job.
//
// Notes:
//   - This command is enabled only in Standard mode. Standard mode can be
//     selected by ESC S. See the Page-mode description for the behaviour of
//     FF when the printer is in Page mode.
//   - When the cutting position has been reserved by GS V <Function C> or
//     GS ( V <Function 51>, issuing FF will feed the paper to the cutting
//     position and perform the cut.
//
// FF (in Page mode)
// Name:
//
//	Print and return to Standard mode
//
// Description:
//
//	In Page mode, FF prints all the data in the print buffer collectively
//	and switches the printer from Page mode to Standard mode.
//
// Notes:
//   - This command is enabled only in Page mode. Page mode can be selected
//     by ESC L. See the Standard-mode description for the behaviour of FF
//     when the printer is in Standard mode.
//   - The data in the print area is deleted after being printed.
//   - This command resets the values set by ESC W to their defaults.
//   - The value set by ESC T is maintained.
//   - After printing, the printer returns to Standard mode, the print
//     position moves to the left side of the printable area, and the printer
//     enters the "Beginning of the line" status.
//
// Value:
//
//	The FF control code is a single byte 0x0C (decimal 12).
func (pc *Commands) FormFeed() []byte {
	return []byte{FF}
}

// PrintAndCarriageReturn prints the data in the print buffer and performs a
// carriage return (moves the print position to the beginning of the print
// line). Behavior depends on the printer head type, auto line feed setting,
// and current mode (Standard or Page).
//
// Text:
//
//	ASCII: CR
//	Hex:   0x0D
//	Decimal: 13
//
// Description:
//
//	Executes one of the following operations depending on the situation:
//	  - Text head alignment
//	  - When auto line feed is enabled: executes printing and one line feed
//	    (same as LF).
//	  - When auto line feed is disabled: the command is ignored for horizontal
//	    alignment on certain heads (see below).
//	  - Horizontal alignment (applies to Line thermal head or Shuttle head).
//	  - Vertical alignment (applies to Serial dot head).
//	Mode-specific behaviour:
//	  - In Standard mode: prints the data in the print buffer and moves the
//	    print position to the beginning of the print line.
//	  - In Page mode: only moves the print position to the beginning of the
//	    print line; the printer does not perform actual printing.
//
// Notes:
//   - With a serial interface, the command performs as if auto line feed is
//     disabled.
//   - Auto line feed can be configured by DIP switch or memory switch. The
//     memory switch can be changed with GS ( E <Function 3>.
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - When processed in Page mode, only the print position moves and the
//     printer does not perform actual printing.
//
// Value is the single-byte CR (carriage return) control code.
func (pc *Commands) PrintAndCarriageReturn() []byte {
	return []byte{CR}
}

// PrintAndLineFeed prints the data in the print buffer and feeds one line,
// based on the current line spacing.
//
// Text:
//
//	ASCII: LF
//	Hex:   0x0A
//	Decimal: 10
//
// Description:
//
//	Prints the data in the print buffer and feeds one line, based on the
//	current line spacing.
//
// Notes:
//   - The amount of paper fed per line is based on the value set using the
//     line spacing command (ESC 2 or ESC 3).
//   - After printing, the print position is moved to the left side of the
//     printable area and the printer enters the "Beginning of the line"
//     status.
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing.
//
// Value is the single-byte LF (line feed) control code.
func (pc *Commands) PrintAndLineFeed() []byte {
	return []byte{LF}
}
