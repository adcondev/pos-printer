package linespacing

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for line spacing control.
// ESC/POS is the command system used by thermal receipt printers to control
// the vertical spacing between printed lines, allowing customization of
// document density and readability.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Spacing represents the line spacing spacing in motion units
type Spacing byte

// Line spacing limits and defaults
const (
	NormalSpacing Spacing = 30 // Default line spacing (model dependent, typically 30-80 dots)
)

// ============================================================================
// Error Definitions
// ============================================================================

// No specific errors needed for line spacing as byte range (0-255) covers all valid values

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for line spacing commands in ESC/POS printers.
type Capability interface {
	SetLineSpacing(lines Spacing) []byte
	SelectDefaultLineSpacing() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for ESC/POS printers.
type Commands struct{}

// NewCommands creates a new instance of line spacing Commands
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateSpacing validates if the spacing value is within acceptable limits.
// Although technically any byte value (0-255) is valid, this function can be
// extended to add printer-specific or logical limits if needed.
func ValidateSpacing(_ Spacing) error {
	// All byte values are valid for line spacing
	// This function exists for consistency and future extensibility
	return nil
}
