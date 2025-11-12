# English Version

This document combines the Go package template (`package_template.go`) and the detailed godoc guide (
`godocs_template.md`).

## 1\. Go Package Template (from `package_template.go`)

This file provides a complete template for a Go package module implementing ESC/POS commands.

```go
package functionality

import (
	"errors"
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for [specific functionality].
// ESC/POS is the command system used by thermal receipt printers to control
// [describe what this package controls].

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions
// Define custom types for parameters to ensure type safety
// Example: type Speed byte

// Constants
// Define command-specific constants and limits
// Example: const MaxSpeed Speed = 255

// Error variables
// Consistent Error Handling:
//  - All errors should be defined as var declarations
//  - Error messages should follow pattern: "invalid X (try Y)"
//  - Use fmt.Errorf with %w for error wrapping
// Example: var ErrInvalidSpeed = errors.New("invalid speed (try 0-255)")

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
// var _ Capability = (*Commands)(nil)

// Interface definition
// type Capability interface {
//     // Command methods with clear, descriptive names
// }

// ============================================================================
// Main Implementation
// ============================================================================

// Struct definition
// type Commands struct {
//     // Include any sub-modules if needed
// }

// Constructor
// func NewCommands() *Commands {
//     return &Commands{}
// }

// Command methods (in logical order)
// Each command should include:
// - Detailed documentation header
// - Format specification (ASCII, Hex, Decimal)
// - Range and default values
// - Parameter descriptions
// - Implementation notes
// - Error conditions

// CommandName provides a brief description of what the command does.
//
// Format:
//
//	ASCII:   <ASCII representation> (e.g., ESC a n)
//	Hex:     <Hex representation> (e.g., 0x1B 0x61 n)
//	Decimal: <Decimal representation> (e.g., 27 97 n)
//
// Range:
//
//	<Parameter ranges, if applicable> (e.g., n = 0-2, 48-50)
//
// Default:
//
//	<Default parameter values> (e.g., n = 0)
//
// Parameters:
//
//	<Parameter name>: <Description of parameter>
//	   <Parameter values and their meanings, if applicable>
//	   <Additional parameter values as needed>
//
// Notes:
//   - <Important implementation note 1>
//   - <Important implementation note 2>
//   - <Any additional notes about behavior or considerations>
//   - <Model-specific information if relevant>
//
// Errors:
//
//	<Error conditions, if applicable> (e.g., Returns ErrParameterName if parameter is outside valid range)
//	<If function is safe and doesn't return errors, note: This function is safe and does not return errors>
func (c *Commands) CommandName(parameters ...interface{}) ([]byte, error) {
	// Implementation
	return nil, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// Private helper methods
// Example: buildCommand, parseResponse, etc.

// ============================================================================
// Validation Helper Functions
// ============================================================================

// Public validation functions
// Consistent Validation Pattern:
// func (c *Commands) SomeCommand(param ParamType) ([]byte, error) {
//     // Validate first
//     if err := ValidateParam(param); err != nil {
//         return nil, err
//     }
//     // Then build command
//     return []byte{...}, nil
// }

// Example validation function:
// func ValidateParam(param ParamType) error {
//     if param < MinValue || param > MaxValue {
//         return fmt.Errorf("%w: %d", ErrInvalidParam, param)
//     }
//     return nil
// }
```

## 2\. ESC/POS Godoc Template Guide (from `godocs_template.md`)

Here's a standardized template for ESC/POS command godocs that you can use in the future:

```go
// CommandName provides a brief description of what the command does.
//
// Format:
//   ASCII:   <ASCII representation> (e.g., ESC a n)
//   Hex:     <Hex representation> (e.g., 0x1B 0x61 n)
//   Decimal: <Decimal representation> (e.g., 27 97 n)
//
// Range:
//   <Parameter ranges, if applicable> (e.g., n = 0-2, 48-50)
//
// Default:
//   <Default parameter values> (e.g., n = 0)
//
// Parameters:
//   <Parameter name>: <Description of parameter>
//      <Parameter values and their meanings, if applicable>
//      <Additional parameter values as needed>
//
// Notes:
//   - <Important implementation note 1>
//   - <Important implementation note 2>
//   - <Any additional notes about behavior or considerations>
//   - <Model-specific information if relevant>
//
// Errors:
//   <Error conditions, if applicable> (e.g., Returns ErrParameterName if parameter is outside valid range)
//   <If function is safe and doesn't return errors, note: This function is safe and does not return errors>
func (c *Commands) CommandName(parameters ...interface{}) ([]byte, error) {
// Implementation
return nil, nil
}
```

### Key Sections

1. **Function Name and Brief Description**

    * First line should succinctly describe what the command does

2. **Format Section**

    * ASCII representation (with command letters and parameter placeholders)
    * Hexadecimal representation (with 0x prefix for hex values)
    * Decimal representation (space-separated values)

3. **Range Section**

    * Valid parameter ranges
    * Use "Not applicable" if there are no parameters

4. **Default Section**

    * Default values for parameters
    * Use "None" or "Not applicable" if there's no default

5. **Parameters Section**

    * Name and description of each parameter
    * Detailed explanation of parameter values and their effects
    * Use indentation for value descriptions

6. **Notes Section**

    * Important implementation details and behaviors
    * Begin each note with a hyphen
    * Include compatibility information, dependencies on other commands
    * Mention mode-specific behavior (Standard vs. Page mode)
    * Include persistence of settings

7. **Errors Section**

    * Document possible errors the function may return
    * If the function doesn't return errors, state that explicitly

### Example

```go
// SelectJustification selects text justification in Standard mode.
//
// Format:
//   ASCII:   ESC a n
//   Hex:     0x1B 0x61 n
//   Decimal: 27 97 n
//
// Range:
//   n = 0–2, 48–50
//
// Default:
//   n = 0 (Left)
//
// Parameters:
//   mode: Justification mode:
//         0 or 48 -> Left justification
//         1 or 49 -> CenterAlign
//         2 or 50 -> Right justification
//
// Notes:
//   - In Standard mode, aligns all data in one line according to n
//   - Effective only in Standard mode and only when processed at the beginning of a line
//   - Has no effect in Page mode
//   - Justification is applied within the Print area set by GS L and ESC W/GS W
//   - Affects characters, graphics, barcodes, 2D codes and space areas set by HT, ESC $, ESC \
//   - Setting persists until ESC @, reset, or power-off
//
// Errors:
//   Returns ErrJustification if mode is not a valid value (0-2, 48-50).
func (c *Commands) SelectJustification(mode byte) ([]byte, error) {
// Implementation
return nil, nil
}
```

-----
