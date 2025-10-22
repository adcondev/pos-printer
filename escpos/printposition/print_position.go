package printposition

import (
	"errors"
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for print position control.
// ESC/POS is the command system used by thermal receipt printers to control
// print positioning, justification, margins, tab positions, and print area
// configuration in both Standard and Page modes.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters
const (
	// HT moves the print position to the next horizontal tab position.
	HT = 0x09
)

// Justification modes
const (
	JustifyLeft        byte = 0x00
	JustifyCenter      byte = 0x01
	JustifyRight       byte = 0x02
	JustifyLeftASCII   byte = '0'
	JustifyCenterASCII byte = '1'
	JustifyRightASCII  byte = '2'
)

// Print direction modes (Page mode)
const (
	DirectionLeftToRight byte = 0x00 // upper left start
	DirectionBottomToTop byte = 0x01 // lower left start
	DirectionRightToLeft byte = 0x02 // lower right start
	DirectionTopToBottom byte = 0x03 // upper right start

	DirectionLeftToRightASCII byte = '0'
	DirectionBottomToTopASCII byte = '1'
	DirectionRightToLeftASCII byte = '2'
	DirectionTopToBottomASCII byte = '3'
)

// Beginning of line operations
const (
	BeginLineErase byte = 0x00 // erase buffer
	BeginLinePrint byte = 0x01 // print buffer

	BeginLineEraseASCII byte = '0'
	BeginLinePrintASCII byte = '1'
)

// Tab position limits
const (
	MaxTabPositions = 32
	MaxTabValue     = 255
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrJustification represents an invalid justification mode error
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

// NewCommands creates a new instance of print position Commands
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// (Add any private helper functions here if needed)

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateJustification validates if justification mode is valid.
func ValidateJustification(mode byte) error {
	switch mode {
	case JustifyLeft, JustifyCenter, JustifyRight,
		JustifyLeftASCII, JustifyCenterASCII, JustifyRightASCII:
		return nil
	default:
		return ErrJustification
	}
}

// ValidatePrintDirection validates if print direction is valid.
func ValidatePrintDirection(direction byte) error {
	switch direction {
	case DirectionLeftToRight, DirectionBottomToTop, DirectionRightToLeft, DirectionTopToBottom,
		DirectionLeftToRightASCII, DirectionBottomToTopASCII, DirectionRightToLeftASCII, DirectionTopToBottomASCII:
		return nil
	default:
		return ErrPrintDirection
	}
}

// ValidateBeginLineMode validates if begin line mode is valid.
func ValidateBeginLineMode(mode byte) error {
	switch mode {
	case BeginLineErase, BeginLinePrint, BeginLineEraseASCII, BeginLinePrintASCII:
		return nil
	default:
		return ErrBeginLineMode
	}
}

// ValidatePrintArea validates print area dimensions.
func ValidatePrintArea(width, height uint16) error {
	if width == 0 {
		return ErrPrintAreaWidthSize
	}
	if height == 0 {
		return ErrPrintAreaHeightSize
	}
	return nil
}
