// Package composer provides a high-level interface for composing ESC/POS commands
// for POS printers by combining various capabilities.
package composer

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/bitimage"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/mechanismcontrol"
	"github.com/adcondev/pos-printer/escpos/print"
	"github.com/adcondev/pos-printer/escpos/printposition"
	"github.com/adcondev/pos-printer/escpos/shared"
)

// Composer implements the ESCPOS Commands
type Composer struct {
	barcode          barcode.Capability
	bitImage         bitimage.Capability
	character        character.Capability
	lineSpacing      linespacing.Capability
	mechanismControl mechanismcontrol.Capability
	print            print.Capability
	printPosition    printposition.Capability
	// TODO: Implement other capabilities
	// PrintingPaper        printingpaper.Capability
	// PaperSensor          papersensor.Capability
	// PanelButton          panelbutton.Capability
	// Status               status.Capability
	// TwoDimensionalCode   twodimensionalcode.Capability
	// MacroFunctions       macrofunctions.Capability
	// Kanji 		        kanji.Capability
	// Miscellaneous 	    miscellaneous.Capability
	// Customize 	        customize.Capability
	// CounterPrinting      counterprinting.Capability
}

// NewComposer creates a new instance of the ESC/POS protocol
func NewComposer() *Composer {
	return &Composer{
		barcode:          barcode.NewCommands(),
		bitImage:         bitimage.NewCommands(),
		character:        character.NewCommands(),
		lineSpacing:      linespacing.NewCommands(),
		mechanismControl: mechanismcontrol.NewCommands(),
		print:            print.NewCommands(),
		printPosition:    printposition.NewCommands(),
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
//   - Clears the data in the print buffer and resets printer modes to those at power-on.
//   - Macro definitions are NOT cleared.
//   - Offline response selection is NOT cleared.
//   - Contents of user NV memory are NOT cleared.
//   - NV graphics (NV bit image) and NV user memory are NOT cleared.
//   - The maintenance counter value is NOT affected by this command.
//   - Software setting values are NOT cleared.
//   - DIP switch settings are NOT re-read.
//   - The data in the receiver buffer is NOT cleared.
//   - In Page mode: deletes data in print areas, initializes all settings, and selects Standard mode.
//   - Cancels many active settings (print mode, line feed, etc.) and moves the print position to the left side
//     of the printable area; printer status becomes "Beginning of the line".
//   - Certain ESC = behavior is preserved/adjusted as described by the printer (see model notes).
//   - Use with care when expecting persistent RAM/NV behavior — only RAM settings are reset to power-on defaults.
//
// Errors:
//
//	This function is safe and does not return errors
func InitializePrinter() []byte {
	return []byte{shared.ESC, '@'}
}

// ============================================================================
// Minimal Print Methods
// ============================================================================

// Init restores the printer to its default state
func (c *Composer) Init() []byte {
	// TODO: Implement any additional initialization logic if needed based on defaults and profile
	return InitializePrinter()
}

// FontA selects Font A
func (c *Composer) FontA() []byte {
	cmd, _ := c.character.SelectCharacterFont(character.FontA)
	return cmd
}

// LeftMargin sets the left margin
func (c *Composer) LeftMargin(margin uint16) []byte {
	return c.printPosition.SetLeftMargin(margin)
}

// PrintWidth sets the print area width
func (c *Composer) PrintWidth(width uint16) []byte {
	return c.printPosition.SetPrintAreaWidth(width)
}

// Print sends text to the printer without a line feed.
func (c *Composer) Print(text string) ([]byte, error) {
	cmd, err := c.print.Text(text)
	if err != nil {
		return nil, fmt.Errorf("print: text: %w", err)
	}
	return cmd, nil
}

// PrintLn sends text to the printer followed by a line feed.
func (c *Composer) PrintLn(text string) ([]byte, error) {
	cmd, err := c.print.Text(text)
	if err != nil {
		return nil, fmt.Errorf("println: text: %w", err)
	}
	cmd = append(cmd, c.print.PrintAndLineFeed()...)
	return cmd, nil
}

// Feed feeds n lines.
func (c *Composer) Feed(n byte) []byte {
	return c.print.PrintAndFeedLines(n)
}

// BoldText enables or disables bold mode.
func (c *Composer) BoldText() []byte {
	return c.character.SetEmphasizedMode(character.OnEm)
}

// RegularText disables bold mode.
func (c *Composer) RegularText() []byte {
	return c.character.SetEmphasizedMode(character.OffEm)
}

// SetAlign sets the text alignment.
func (c *Composer) SetAlign(mode printposition.Justification) ([]byte, error) {
	cmd, err := c.printPosition.SelectJustification(mode)
	if err != nil {
		return nil, fmt.Errorf("set align: select justification: %w", err)
	}
	return cmd, nil
}

// CenterAlign centers the text.
func (c *Composer) CenterAlign() []byte {
	cmd, _ := c.printPosition.SelectJustification(printposition.Center)
	return cmd
}

// LeftAlign left-aligns the text.
func (c *Composer) LeftAlign() []byte {
	cmd, _ := c.printPosition.SelectJustification(printposition.Left)
	return cmd
}

// RightAlign right-aligns the text.
func (c *Composer) RightAlign() []byte {
	cmd, _ := c.printPosition.SelectJustification(printposition.Right)
	return cmd
}

// RegularTextSize sets the smallest(regular) text size.
func (c *Composer) RegularTextSize() []byte {
	size, _ := character.NewSize(1, 1)
	return c.character.SelectCharacterSize(size)
}

// DoubleSizeText sets double size text.
func (c *Composer) DoubleSizeText() []byte {
	size, _ := character.NewSize(2, 2)
	return c.character.SelectCharacterSize(size)
}
