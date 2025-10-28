package escpos

import (
	"fmt"
	"strings"

	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/bitimage"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/mechanismcontrol"
	"github.com/adcondev/pos-printer/escpos/print"
	"github.com/adcondev/pos-printer/escpos/printposition"
	"github.com/adcondev/pos-printer/escpos/shared"
)

// Commands implements the ESCPOS Commands
type Commands struct {
	Barcode          barcode.Capability
	Bitimage         bitimage.Capability
	Character        character.Capability
	LineSpacing      linespacing.Capability
	MechanismControl mechanismcontrol.Capability
	Print            print.Capability
	PrintPosition    printposition.Capability
}

// Raw sends raw data without processing
func (c *Commands) Raw(data []byte) ([]byte, error) {
	if err := shared.IsBufLenOk(data); err != nil {
		return nil, err
	}
	return data, nil
}

// NewEscposCommands creates a new instance of the ESC/POS protocol
func NewEscposCommands() *Commands {
	return &Commands{
		Barcode:          barcode.NewCommands(),
		Bitimage:         bitimage.NewCommands(),
		Character:        character.NewCommands(),
		LineSpacing:      linespacing.NewCommands(),
		MechanismControl: mechanismcontrol.NewCommands(),
		Print:            print.NewCommands(),
		PrintPosition:    printposition.NewCommands(),
	}
}

// Initialize restores the printer to its default state
func (c *Commands) Initialize() []byte {
	// ESC @ - Reset printer
	return []byte{shared.ESC, '@'}
}

// SetCharacterFont selecciona la fuente
func (c *Commands) SetCharacterFont(font character.FontType) ([]byte, error) {
	return c.Character.SelectCharacterFont(font)
}

// Tab mueve a la siguiente posición de tabulación
func (c *Commands) Tab() []byte {
	return c.PrintPosition.HorizontalTab()
}

// SetLeftMargin establece el margen izquierdo
func (c *Commands) SetLeftMargin(margin uint16) []byte {
	return c.PrintPosition.SetLeftMargin(margin)
}

// SetPrintWidth establece el ancho del área de impresión
func (c *Commands) SetPrintWidth(width uint16) []byte {
	return c.PrintPosition.SetPrintAreaWidth(width)
}

// ============================================================================
// High-Level Print Methods
// ============================================================================

// PrintText sends text to the printer without a line feed.
// Internally uses: Character.Text()
func (c *Commands) PrintText(text string) ([]byte, error) {
	return c.Print.Text(text)
}

// PrintLine sends text to the printer followed by a line feed.
func (c *Commands) PrintLine(text string) ([]byte, error) {
	var result []byte

	if text != "" {
		data, err := c.Print.Text(text)
		if err != nil {
			return nil, fmt.Errorf("print line text: %w", err)
		}
		result = append(result, data...)
	}

	result = append(result, c.Print.PrintAndLineFeed()...)
	return result, nil
}

// PrintLines sends multiple lines of text, each with a line feed.
func (c *Commands) PrintLines(lines []string) ([]byte, error) {
	var result []byte

	for i, line := range lines {
		data, err := c.PrintLine(line)
		if err != nil {
			return nil, fmt.Errorf("print line %d: %w", i, err)
		}
		result = append(result, data...)
	}

	return result, nil
}

// Feed advances the paper by the specified number of lines.
func (c *Commands) Feed(lines byte) []byte {
	return c.Print.PrintAndFeedLines(lines)
}

// FeedPaper advances the paper by the specified units.
func (c *Commands) FeedPaper(units byte) []byte {
	return c.Print.PrintAndFeedPaper(units)
}

// NewLine sends a line feed command.
func (c *Commands) NewLine() []byte {
	return c.Print.PrintAndLineFeed()
}

// CarriageReturn sends a carriage return command.
func (c *Commands) CarriageReturn() []byte {
	return c.Print.PrintAndCarriageReturn()
}

// FormFeed sends a form feed command.
func (c *Commands) FormFeed() []byte {
	return c.Print.FormFeed()
}

// ReverseFeed reverses the paper by the specified number of lines.
func (c *Commands) ReverseFeed(lines byte) ([]byte, error) {
	return c.Print.PrintAndReverseFeed(lines)
}

// ============================================================================
// Text Formatting Methods
// ============================================================================

// Bold enables or disables bold text.
func (c *Commands) Bold(enable bool) []byte {
	if enable {
		return c.Character.SetEmphasizedMode(1)
	}
	return c.Character.SetEmphasizedMode(0)
}

// Underline sets the underline mode (0=off, 1=single, 2=double).
func (c *Commands) Underline(mode byte) ([]byte, error) {
	return c.Character.SetUnderlineMode(character.UnderlineMode(mode))
}

// Align sets text justification (left, center, right).
func (c *Commands) Align(alignment string) ([]byte, error) {
	var mode byte
	switch strings.ToLower(alignment) {
	case "left":
		mode = 0
	case "center":
		mode = 1
	case "right":
		mode = 2
	default:
		return nil, fmt.Errorf("invalid alignment: %s (use left, center, or right)", alignment)
	}
	return c.PrintPosition.SelectJustification(mode)
}

// AlignLeft sets left justification.
func (c *Commands) AlignLeft() ([]byte, error) {
	return c.PrintPosition.SelectJustification(0)
}

// AlignCenter sets center justification.
func (c *Commands) AlignCenter() ([]byte, error) {
	return c.PrintPosition.SelectJustification(1)
}

// AlignRight sets right justification.
func (c *Commands) AlignRight() ([]byte, error) {
	return c.PrintPosition.SelectJustification(2)
}

// Size sets the character size (width and height multipliers 1-8).
func (c *Commands) Size(width, height byte) ([]byte, error) {
	size, err := character.NewSize(width, height)
	if err != nil {
		return nil, fmt.Errorf("invalid size: %w", err)
	}
	return c.Character.SelectCharacterSize(size), nil
}

// NormalSize resets the character size to normal.
// Internally uses: Character.SelectCharacterSize(1, 1)
func (c *Commands) NormalSize() []byte {
	size, _ := character.NewSize(1, 1)
	return c.Character.SelectCharacterSize(size)
}

// DoubleWidth enables or disables double-width characters.
func (c *Commands) DoubleWidth(enable bool) []byte {
	// This would depend on your Character implementation
	// Using SelectCharacterSize as an example
	if enable {
		size, _ := character.NewSize(2, 1)
		data := c.Character.SelectCharacterSize(size)
		return data
	}
	size, _ := character.NewSize(1, 1)
	data := c.Character.SelectCharacterSize(size)
	return data
}

// DoubleHeight enables or disables double-height characters.
func (c *Commands) DoubleHeight(enable bool) []byte {
	if enable {
		size, _ := character.NewSize(1, 2)
		data := c.Character.SelectCharacterSize(size)
		return data
	}
	size, _ := character.NewSize(1, 1)
	data := c.Character.SelectCharacterSize(size)
	return data
}

// AbsolutePosition sets the print position to an absolute value.
func (c *Commands) AbsolutePosition(position uint16) []byte {
	return c.PrintPosition.SetAbsolutePrintPosition(position)
}

// RelativePosition sets the print position relative to the current position.
func (c *Commands) RelativePosition(position int16) []byte {
	return c.PrintPosition.SetRelativePrintPosition(position)
}

// ============================================================================
// Paper Control Methods
// ============================================================================

// Cut performs a paper cut (full or partial).
func (c *Commands) Cut(partial bool) ([]byte, error) {
	mode := byte(0) // full cut
	if partial {
		mode = 1 // partial cut
	}
	return c.MechanismControl.CutPaper(mechanismcontrol.CutMode(mode))
}

// FullCut performs a full paper cut.
func (c *Commands) FullCut() ([]byte, error) {
	return c.MechanismControl.CutPaper(mechanismcontrol.CutModeFull)
}

// PartialCut performs a partial paper cut.
func (c *Commands) PartialCut() ([]byte, error) {
	return c.MechanismControl.CutPaper(mechanismcontrol.CutModePartial)
}

// FeedAndCut feeds paper then performs a cut.
func (c *Commands) FeedAndCut(feedLines byte, partial bool) ([]byte, error) {
	mode := mechanismcontrol.FeedCutModeFull
	if partial {
		mode = mechanismcontrol.FeedCutModePartial
	}
	return c.MechanismControl.FeedAndCutPaper(mode, feedLines)
}

// Home returns the print position to the home position.
func (c *Commands) Home() []byte {
	return c.MechanismControl.ReturnHome()
}

// UnidirectionalPrint enables or disables unidirectional print mode.
func (c *Commands) UnidirectionalPrint(mode bool) []byte {
	if mode {
		return c.MechanismControl.SetUnidirectionalPrintMode(mechanismcontrol.UnidirectionalOn)
	}
	return c.MechanismControl.SetUnidirectionalPrintMode(mechanismcontrol.UnidirectionalOff)
}
