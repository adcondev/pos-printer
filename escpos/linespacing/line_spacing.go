package linespacing

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Spacing represents the line spacing spacing in motion units
type Spacing byte

// Line spacing limits and defaults
const (
	MinSpacing    Spacing = 0   // Minimum line spacing spacing
	MaxSpacing    Spacing = 255 // Maximum line spacing spacing
	NormalSpacing Spacing = 30  // Default line spacing (model dependent, typically 30-80 dots)
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

func NewCommands() *Commands {
	return &Commands{}
}
