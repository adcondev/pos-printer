package print

import (
	"errors"

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

// Reverse motion units and lines
var (
	// MaxReverseMotionUnits is the maximum number of motion units for reverse printing
	MaxReverseMotionUnits byte = 48
	// MaxReverseFeedLines is the maximum number of lines for reverse printing
	MaxReverseFeedLines byte = 2
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	// ErrPrintReverseFeed invalid number of motion units for reverse print
	ErrPrintReverseFeed = errors.New("invalid number of motion units for reverse print (try 0-48)")
	// ErrPrintReverseFeedLines invalid number of lines for reverse print
	ErrPrintReverseFeedLines = errors.New("invalid number of lines for reverse print (try 0-2)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var (
	_ Capability         = (*Commands)(nil)
	_ PageModeCapability = (*PagePrint)(nil)
)

// Capability defines the interface for printer capabilities
type Capability interface {
	Text(text string) ([]byte, error)
	PrintAndFeedPaper(units byte) []byte
	FormFeed() []byte
	PrintAndCarriageReturn() []byte
	PrintAndLineFeed() []byte
}

// PageModeCapability defines the interface for page mode capabilities
type PageModeCapability interface {
	PrintAndReverseFeed(units byte) ([]byte, error)
	PrintAndReverseFeedLines(lines byte) ([]byte, error)
	CancelData() []byte
	PrintDataInPageMode() []byte
	PrintAndFeedLines(lines byte) ([]byte, error)
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
		Page: &PagePrint{},
	}
}

// PagePrint groups page mode printing commands
type PagePrint struct{}

// Text formats and sends a string for printing
func (pc *Commands) Text(n string) ([]byte, error) {
	if err := common.IsBufLenOk([]byte(n)); err != nil {
		return nil, err
	}
	return Formatting([]byte(n)), nil
}

// PrintAndFeedPaper prints the data in the print buffer and feeds the paper
// by n × (vertical or horizontal motion unit).
//
// Text:
//
//	ASCII: ESC J n
//	Hex:   0x1B 0x4A n
//	Decimal: 27 74 n
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
//	Prints the data in the print buffer and feeds the paper [n × (vertical or
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
//   - The maximum paper feed amount depends on the printer model. If n is
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
//	ESC J n -> 0x1B, 0x4A, n
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
	if n > MaxReverseMotionUnits {
		return nil, ErrPrintReverseFeed
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
	if n > MaxReverseFeedLines {
		return nil, ErrPrintReverseFeedLines
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
