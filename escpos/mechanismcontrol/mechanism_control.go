package mechanismcontrol

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for mechanism control functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// physical printer mechanisms such as paper cutting, print head movement,
// paper feeding, and unidirectional print mode settings.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions and Constants

// UnidirectionalMode represents the unidirectional print mode setting
type UnidirectionalMode byte

const (
	// UnidirectionalOff turns off unidirectional mode (bidirectional on)
	UnidirectionalOff UnidirectionalMode = 0x00
	// UnidirectionalOn turns on unidirectional mode
	UnidirectionalOn UnidirectionalMode = 0x01
)

// CutMode represents the paper cut mode for simple cut commands
type CutMode byte

const (
	// CutModeFull performs a full cut
	CutModeFull CutMode = 0x00
	// CutModePartial performs a partial cut
	CutModePartial CutMode = 0x01
	// CutModeFullASCII performs a full cut (ASCII mode)
	CutModeFullASCII CutMode = '0'
	// CutModePartialASCII performs a partial cut (ASCII mode)
	CutModePartialASCII CutMode = '1'
)

// CutType represents the paper cut type for extended cut commands
type CutType byte

const (
	// CutTypeFull performs a full cut
	CutTypeFull CutType = 0x00
	// CutTypePartial performs a partial cut
	CutTypePartial CutType = 0x01
)

// FeedCutMode represents the feed and cut mode
type FeedCutMode byte

const (
	// FeedCutModeFull performs a full cut after feeding
	FeedCutModeFull FeedCutMode = 65
	// FeedCutModePartial performs a partial cut after feeding
	FeedCutModePartial FeedCutMode = 66
)

// PositionCutMode represents the position cut mode
type PositionCutMode byte

const (
	// PositionCutModeFull performs a full cut at position
	PositionCutModeFull PositionCutMode = 97
	// PositionCutModePartial performs a partial cut at position
	PositionCutModePartial PositionCutMode = 98
)

// FeedCutReturnMode represents the feed, cut and return mode
type FeedCutReturnMode byte

const (
	// FeedCutReturnModeFull performs a full cut with return
	FeedCutReturnModeFull FeedCutReturnMode = 103
	// FeedCutReturnModePartial performs a partial cut with return
	FeedCutReturnModePartial FeedCutReturnMode = 104
)

const (
	// MaxFeedAmount represents the maximum feed amount before cutting
	MaxFeedAmount byte = 255
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	// ErrCutMode indicates an invalid cut mode
	ErrCutMode = errors.New("invalid cut mode (try 0, 1, 48, or 49)")
	// ErrCutType indicates an invalid cut type
	ErrCutType = errors.New("invalid cut type (try 0 or 1)")
	// ErrFeedCutMode indicates an invalid feed and cut mode
	ErrFeedCutMode = errors.New("invalid feed and cut mode (try 65 or 66)")
	// ErrPositionCutMode indicates an invalid position cut mode
	ErrPositionCutMode = errors.New("invalid position cut mode (try 97 or 98)")
	// ErrFeedCutReturnMode indicates an invalid feed, cut and return mode
	ErrFeedCutReturnMode = errors.New("invalid feed, cut and return mode (try 103 or 104)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Compile-time check that Commands implements Capability
var _ Capability = (*Commands)(nil)

// Capability defines the interface for mechanism control commands
type Capability interface {
	// Print head control
	ReturnHome() []byte
	SetUnidirectionalPrintMode(mode UnidirectionalMode) []byte

	// Legacy cut commands (deprecated)
	PartialCut() []byte
	PartialCutThreePoints() []byte

	// Modern cut commands
	PaperCut(cutType CutType) ([]byte, error)
	PaperFeedAndCut(cutType CutType, feedAmount byte) ([]byte, error)
	ReservePaperCut(cutType CutType, feedAmount byte) ([]byte, error)

	// Simple cut commands
	CutPaper(mode CutMode) ([]byte, error)
	FeedAndCutPaper(mode FeedCutMode, feedAmount byte) ([]byte, error)
	SetCutPosition(mode PositionCutMode, position byte) ([]byte, error)
	FeedCutAndReturnPaper(mode FeedCutReturnMode, feedAmount byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for mechanism control
type Commands struct{}

// NewCommands creates a new Commands instance
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateCutMode validates if cut mode is valid
func ValidateCutMode(mode CutMode) error {
	switch mode {
	case CutModeFull, CutModePartial, CutModeFullASCII, CutModePartialASCII:
		return nil
	default:
		return ErrCutMode
	}
}

// ValidateCutType validates if cut type is valid
func ValidateCutType(cutType CutType) error {
	switch cutType {
	case CutTypeFull, CutTypePartial:
		return nil
	default:
		return ErrCutType
	}
}

// ValidateFeedCutMode validates if feed cut mode is valid
func ValidateFeedCutMode(mode FeedCutMode) error {
	switch mode {
	case FeedCutModeFull, FeedCutModePartial:
		return nil
	default:
		return ErrFeedCutMode
	}
}

// ValidatePositionCutMode validates if position cut mode is valid
func ValidatePositionCutMode(mode PositionCutMode) error {
	switch mode {
	case PositionCutModeFull, PositionCutModePartial:
		return nil
	default:
		return ErrPositionCutMode
	}
}

// ValidateFeedCutReturnMode validates if feed cut return mode is valid
func ValidateFeedCutReturnMode(mode FeedCutReturnMode) error {
	switch mode {
	case FeedCutReturnModeFull, FeedCutReturnModePartial:
		return nil
	default:
		return ErrFeedCutReturnMode
	}
}
