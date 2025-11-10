// Package composer provides a high-level interface for composing ESC/POS commands
// for POS printers by combining various capabilities.
package composer

import (
	"fmt"

	"github.com/adcondev/pos-printer/pkg/commands/barcode"
	"github.com/adcondev/pos-printer/pkg/commands/bitimage"
	"github.com/adcondev/pos-printer/pkg/commands/character"
	"github.com/adcondev/pos-printer/pkg/commands/common"
	"github.com/adcondev/pos-printer/pkg/commands/linespacing"
	"github.com/adcondev/pos-printer/pkg/commands/mechanismcontrol"
	"github.com/adcondev/pos-printer/pkg/commands/print"
	"github.com/adcondev/pos-printer/pkg/commands/printposition"
	"github.com/adcondev/pos-printer/pkg/commands/qrcode"
)

// EscposProtocol implements the ESCPOS Commands
type EscposProtocol struct {
	Barcode          barcode.Capability
	BitImage         bitimage.Capability
	Character        character.Capability
	LineSpacing      linespacing.Capability
	MechanismControl mechanismcontrol.Capability
	Print            print.Capability
	PrintPosition    printposition.Capability
	QRCode           qrcode.Capability
	// TODO: Implement other capabilities
	// PrintingPaper    printingpaper.Capability
	// PaperSensor      papersensor.Capability
	// PanelButton      panelbutton.Capability
	// Status           status.Capability
	// MacroFunctions   macrofunctions.Capability
	// Kanji 		    kanji.Capability
	// Miscellaneous 	miscellaneous.Capability
	// Customize 	    customize.Capability
	// CounterPrinting  counterprinting.Capability
	// PDF417           pdf417.Capability
	// MaxiCode         maxicode.Capability
	// DataBar          databar.Capability
	// CompositeSym     compositesym.Capability
	// AztecCode        azteccode.Capability
	// DataMatrix       datamatrix.Capability
}

// NewEscpos creates a new instance of the ESC/POS protocol
func NewEscpos() *EscposProtocol {
	return &EscposProtocol{
		Barcode:          barcode.NewCommands(),
		BitImage:         bitimage.NewCommands(),
		Character:        character.NewCommands(),
		LineSpacing:      linespacing.NewCommands(),
		MechanismControl: mechanismcontrol.NewCommands(),
		Print:            print.NewCommands(),
		PrintPosition:    printposition.NewCommands(),
		QRCode:           qrcode.NewCommands(),
	}
}

// TODO: Implement other methods to access capabilities related to initialization and state management

// InitializePrinter provides a reset of the printer to its power-on state for RAM settings.
//
// Format:
//
//	ASCII:   ESC @
//	Hex:     0x1B 0x40
//	Decimal: 27 64
//
// Range:
//
//	None
//
// Default:
//
//	None
//
// Parameters:
//
//	None
//
// Notes:
//   - Clears the data in the Print buffer and resets printer modes to those at power-on.
//   - Macro definitions are NOT cleared.
//   - Offline response selection is NOT cleared.
//   - Contents of user NV memory are NOT cleared.
//   - NV graphics (NV bit image) and NV user memory are NOT cleared.
//   - The maintenance counter value is NOT affected by this command.
//   - Software setting values are NOT cleared.
//   - DIP switch settings are NOT re-read.
//   - The data in the receiver buffer is NOT cleared.
//   - In Page mode: deletes data in Print areas, initializes all settings, and selects Standard mode.
//   - Cancels many active settings (Print mode, line feed, etc.) and moves the Print position to the left side
//     of the printable area; printer status becomes "Beginning of the line".
//   - Certain ESC = behavior is preserved/adjusted as described by the printer (see model notes).
//   - Use with care when expecting persistent RAM/NV behavior — only RAM settings are reset to power-on defaults.
//
// Errors:
//
//	This function is safe and does not return errors
func (c *EscposProtocol) InitializePrinter() []byte {
	return []byte{common.ESC, '@'}
}

// ============================================================================
// Minimal Print Methods
// ============================================================================

// LeftMargin sets the left margin
func (c *EscposProtocol) LeftMargin(margin uint16) []byte {
	return c.PrintPosition.SetLeftMargin(margin)
}

// PrintWidth sets the Print area width
func (c *EscposProtocol) PrintWidth(width uint16) []byte {
	return c.PrintPosition.SetPrintAreaWidth(width)
}

// PrintLn sends text to the printer followed by a line feed.
func (c *EscposProtocol) PrintLn(text string) ([]byte, error) {
	cmd, err := c.Print.Text(text)
	if err != nil {
		return nil, fmt.Errorf("println: text: %w", err)
	}
	cmd = append(cmd, c.Print.PrintAndLineFeed()...)
	return cmd, nil
}

// RegularText disables bold mode.
func (c *EscposProtocol) RegularText() []byte {
	return c.Character.SetEmphasizedMode(character.OffEm)
}

// SetAlign sets the text alignment.
func (c *EscposProtocol) SetAlign(mode printposition.Justification) ([]byte, error) {
	cmd, err := c.PrintPosition.SelectJustification(mode)
	if err != nil {
		return nil, fmt.Errorf("set align: select justification: %w", err)
	}
	return cmd, nil
}

// CenterAlign centers the text.
func (c *EscposProtocol) CenterAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Center)
	return cmd
}

// LeftAlign left-aligns the text.
func (c *EscposProtocol) LeftAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Left)
	return cmd
}

// RightAlign right-aligns the text.
func (c *EscposProtocol) RightAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Right)
	return cmd
}

// RegularTextSize sets the smallest(regular) text size.
func (c *EscposProtocol) RegularTextSize() []byte {
	size, _ := character.NewSize(1, 1)
	return c.Character.SelectCharacterSize(size)
}

// DoubleSizeText sets double size text.
func (c *EscposProtocol) DoubleSizeText() []byte {
	size, _ := character.NewSize(2, 2)
	return c.Character.SelectCharacterSize(size)
}
