package print

import (
	"errors"
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters used in printing commands
const (
	// LF (Line Feed)
	LF byte = 0x0A // Hex: 0x0A, Decimal: 10
	// CR (Carriage Return)
	CR byte = 0x0D // Hex: 0x0D, Decimal: 13
	// FF (Form Feed)
	FF byte = 0x0C // Hex: 0x0C, Decimal: 12
	// CAN (Cancel)
	CAN byte = 0x18 // Hex: 0x18, Decimal: 24
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
	// ErrInvalidEmptyText indicates that the provided text is empty
	ErrInvalidEmptyText = common.ErrEmptyBuffer
	// ErrInvalidTextTooLarge indicates that the provided text exceeds buffer limits
	ErrInvalidTextTooLarge = common.ErrBufferOverflow
	// ErrInvalidReverseUnits invalid number of motion units for reverse print
	ErrInvalidReverseUnits = fmt.Errorf("invalid reverse feed units (try 0-%d)", MaxReverseMotionUnits)
	// ErrInvalidReverseLines invalid number of lines for reverse print
	ErrInvalidReverseLines = fmt.Errorf("invalid reverse feed lines (try 0-%d)", MaxReverseFeedLines)
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for print commands
type Capability interface {
	// Text operations
	Text(text string) ([]byte, error)

	// Basic print commands
	PrintAndLineFeed() []byte
	PrintAndCarriageReturn() []byte
	FormFeed() []byte

	// Paper feed operations
	PrintAndFeedPaper(units byte) []byte
	PrintAndFeedLines(lines byte) []byte
	PrintAndReverseFeed(units byte) ([]byte, error)
	PrintAndReverseFeedLines(lines byte) ([]byte, error)

	// Page mode specific
	PrintDataInPageMode() []byte
	CancelData() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for print commands
type Commands struct{}

func NewCommands() *Commands {
	return &Commands{}
}

// Formatting replaces specific characters in the byte slice with their ESC/POS equivalents.
func Formatting(data []byte) []byte {
	formatted := make([]byte, len(data))
	copy(formatted, data)

	for i := range formatted {
		switch formatted[i] {
		case '\n':
			formatted[i] = LF
		case '\r':
			formatted[i] = CR
		case '\t':
			formatted[i] = common.HT
		}
	}
	return formatted
}

// Text formats and sends a string for printing.
//
// Description:
//
//	Converts a string to bytes and applies ESC/POS formatting.
//
// Notes:
//   - Replaces '\n' with LF (0x0A)
//   - Replaces '\r' with CR (0x0D)
//   - Replaces '\t' with HT (0x09)
//   - Validates buffer size according to printer limitations
func (c *Commands) Text(n string) ([]byte, error) {
	if err := common.IsBufLenOk([]byte(n)); err != nil {
		switch {
		case errors.Is(err, common.ErrEmptyBuffer):
			return nil, ErrInvalidEmptyText
		case errors.Is(err, common.ErrBufferOverflow):
			return nil, ErrInvalidTextTooLarge
		default:
			return nil, err
		}
	}

	return Formatting([]byte(n)), nil
}

// FormFeed executes form feed operation (behavior varies by mode).
//
// Format:
//
//	ASCII: FF
//	Hex:   0x0C
//	Decimal: 12
//
// Description:
//
//	In Standard mode:
//	  - Indicates "Printing is completed" for the current job
//	  - Feeds to cutting position if reserved by GS V commands
//
//	In Page mode:
//	  - Prints all data in the print buffer collectively
//	  - Switches from Page mode to Standard mode
//	  - Clears the print area after printing
//	  - Resets ESC W values to defaults
//
// Notes:
//   - Mode-specific behavior; check current mode before use
//   - ESC T value is maintained when returning from Page mode
//
// Byte sequence:
//
//	FF -> 0x0C
func (c *Commands) FormFeed() []byte {
	return []byte{FF}
}

// PrintAndCarriageReturn prints the data in the print buffer and performs a
// carriage return.
//
// Format:
//
//	ASCII: CR
//	Hex:   0x0D
//	Decimal: 13
//
// Description:
//
//	Executes one of the following operations depending on the print head type
//	and auto line feed setting:
//	  - When auto line feed is enabled: executes printing and one line feed (same as LF)
//	  - When auto line feed is disabled: behavior depends on print head type
//	  - In Standard mode: prints data and moves to beginning of line
//	  - In Page mode: only moves print position without printing
//
// Notes:
//   - With a serial interface, the command performs as if auto line feed is disabled.
//   - Auto line feed can be configured by DIP switch or memory switch (GS ( E <Function 3>).
//   - After printing, the print position is moved to the left side of the printable area.
//
// Byte sequence:
//
//	CR -> 0x0D
func (c *Commands) PrintAndCarriageReturn() []byte {
	return []byte{CR}
}

// PrintAndLineFeed prints the data in the print buffer and feeds one line.
//
// Format:
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
//     printable area and the printer enters the "Beginning of the line" status.
//   - When this command is processed in Page mode, only the print position
//     moves and the printer does not perform actual printing.
//
// Byte sequence:
//
//	LF -> 0x0A
func (c *Commands) PrintAndLineFeed() []byte {
	return []byte{LF}
}

// PrintAndFeedPaper prints the data in the print buffer and feeds the paper.
//
// Format:
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
//   - Maximum paper feed amount depends on the printer model.
//   - After printing, the print position moves to the beginning of the line.
//   - In Page mode, only the print position moves without actual printing.
//
// Byte sequence:
//
//	ESC J n -> 0x1B, 0x4A, n
func (c *Commands) PrintAndFeedPaper(n byte) []byte {
	return []byte{common.ESC, 'J', n}
}

// PrintAndFeedLines prints the data in the print buffer and feeds n lines.
//
// Format:
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
//   - Paper feed per line based on line spacing (ESC 2 or ESC 3).
//   - Maximum feed depends on printer model.
//   - After printing, print position moves to beginning of line.
//   - In Page mode, only print position moves without actual printing.
//
// Byte sequence:
//
//	ESC d n -> 0x1B, 0x64, n
func (c *Commands) PrintAndFeedLines(n byte) []byte {
	return []byte{common.ESC, 'd', n}
}

// PrintAndReverseFeed prints the data in the print buffer and feeds paper in reverse.
//
// Format:
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
//   - Motion unit used depends on mode and ESC T setting.
//   - If n exceeds model maximum, reverse feed is not executed but printing occurs.
//   - Some printers perform small forward feed after reverse feed.
//   - In Page mode, only print position moves without actual printing.
//
// Byte sequence:
//
//	ESC K n -> 0x1B, 0x4B, n
func (c *Commands) PrintAndReverseFeed(n byte) ([]byte, error) {
	if n > MaxReverseMotionUnits {
		return nil, ErrInvalidReverseUnits
	}
	return []byte{common.ESC, 'K', n}, nil
}

// PrintAndReverseFeedLines prints the data in the print buffer and feeds n lines in reverse.
//
// Format:
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
//	Prints the data in the print buffer and feeds n lines in the reverse direction.
//
// Notes:
//   - Paper feed per line based on line spacing (ESC 2 or ESC 3).
//   - If n exceeds model maximum, reverse feed is not executed but printing occurs.
//   - Some printers perform small forward feed after reverse feed.
//   - In Page mode, only print position moves without actual printing.
//
// Byte sequence:
//
//	ESC e n -> 0x1B, 0x65, n
func (c *Commands) PrintAndReverseFeedLines(n byte) ([]byte, error) {
	if n > MaxReverseFeedLines {
		return nil, ErrInvalidReverseLines
	}
	return []byte{common.ESC, 'e', n}, nil
}

// PrintDataInPageMode prints the data in the print buffer collectively in Page mode.
//
// Format:
//
//	ASCII: ESC FF
//	Hex:   0x1B 0x0C
//	Decimal: 27 12
//
// Description:
//
//	In Page mode, prints all data currently buffered in the print area collectively.
//
// Notes:
//   - Enabled only in Page mode (selected by ESC L).
//   - After printing, buffer data, print position, and settings remain intact.
//   - Commonly used for reprinting the same Page-mode data multiple times.
//   - Returns to Standard mode when FF, ESC S, or ESC @ is issued.
//
// Byte sequence:
//
//	ESC FF -> 0x1B, 0x0C
func (c *Commands) PrintDataInPageMode() []byte {
	return []byte{common.ESC, FF}
}

// CancelData deletes all print data in the current print area (Page mode only).
//
// Format:
//
//	ASCII: CAN
//	Hex:   0x18
//	Decimal: 24
//
// Description:
//
//	In Page mode, deletes all the print data in the current print area.
//
// Notes:
//   - Enabled only in Page mode (selected by ESC L).
//   - Also deletes overlapping data from previously specified print areas.
//   - Has no effect in Standard mode.
//
// Byte sequence:
//
//	CAN -> 0x18
func (c *Commands) CancelData() []byte {
	return []byte{CAN}
}
