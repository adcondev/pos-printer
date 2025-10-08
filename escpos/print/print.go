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
	// ErrEmptyText indicates that the provided text is empty
	ErrEmptyText = common.ErrEmptyBuffer
	// ErrTextTooLarge indicates that the provided text exceeds buffer limits
	ErrTextTooLarge = common.ErrBufferOverflow
	// ErrReverseUnits invalid number of motion units for reverse print
	ErrReverseUnits = fmt.Errorf("invalid reverse feed units (try 0-%d)", MaxReverseMotionUnits)
	// ErrReverseLines invalid number of lines for reverse print
	ErrReverseLines = fmt.Errorf("invalid reverse feed lines (try 0-%d)", MaxReverseFeedLines)
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
			return nil, ErrEmptyText
		case errors.Is(err, common.ErrBufferOverflow):
			return nil, ErrTextTooLarge
		default:
			return nil, err
		}
	}

	return Formatting([]byte(n)), nil
}
