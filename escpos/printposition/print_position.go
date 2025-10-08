package printposition

import (
	"errors"
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters
const (
	// HT moves the print position to the next horizontal tab position.
	HT = common.HT // 0x09
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
	ErrJustification       = errors.New("invalid justification mode (try 0-2 or '0'..'2')")
	ErrPrintDirection      = errors.New("invalid print direction (try 0-3 or '0'..'3')")
	ErrBeginLineMode       = errors.New("invalid begin line mode (try 0-1 or '0'..'1')")
	ErrTooManyTabPositions = fmt.Errorf("too many tab positions (max %d)", MaxTabPositions)
	ErrTabPosition         = errors.New("invalid tab position (must be 1-255 in ascending order)")
	ErrPrintAreaWidthSize  = errors.New("invalid print area size (width must be >= 1)")
	ErrPrintAreaHeightSize = errors.New("invalid print area size (height must be >= 1)")
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
	SetPrintAreaPageMode(x, y, width, height uint16) ([]byte, error)
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
