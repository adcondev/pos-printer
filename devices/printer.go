// Package devices provides implementations for various POS printer devices.
package devices

import (
	"fmt"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/profile"
)

// Printer represents a POS printer device
type Printer struct {
	Profile    profile.Escpos
	Connection connector.Connector
	Protocol   escpos.Commands
}

// NewPrinter creates a new Printer instance
func NewPrinter(proto *escpos.Commands, prof *profile.Escpos, conn connector.Connector) (*Printer, error) {
	if proto == nil {
		return nil, fmt.Errorf("protocol cannot be nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("connection cannot be nil")
	}
	if prof == nil {
		return nil, fmt.Errorf("profile cannot be nil")
	}
	return &Printer{
		Profile:    *prof,
		Connection: conn,
		Protocol:   *proto,
	}, nil
}

// ============================================================================
// Basic Control Methods
// ============================================================================

// Initialize resets the printer to default settings
func (p *Printer) Initialize() error {
	// TODO: Add profile-specific initialization if needed
	cmd := p.Protocol.Initialize()
	_, err := p.Connection.Write(cmd)
	return err
}

// Close closes the connection to the printer
func (p *Printer) Close() error {
	return p.Connection.Close()
}

// Write sends raw bytes directly to the printer
func (p *Printer) Write(data []byte) error {
	_, err := p.Connection.Write(data)
	return err
}

// ============================================================================
// Text Printing Methods
// ============================================================================

// Print sends text without line feed
func (p *Printer) Print(text string) error {
	cmd, err := p.Protocol.PrintText(text)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PrintLine sends text with line feed
func (p *Printer) PrintLine(text string) error {
	cmd, err := p.Protocol.PrintLine(text)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PrintLines sends multiple lines
func (p *Printer) PrintLines(lines []string) error {
	cmd, err := p.Protocol.PrintLines(lines)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// NewLine sends a line feed
func (p *Printer) NewLine() error {
	return p.Write(p.Protocol.NewLine())
}

// Feed advances paper by n lines
func (p *Printer) Feed(lines byte) error {
	return p.Write(p.Protocol.Feed(lines))
}

// FeedPaper advances paper by n units
func (p *Printer) FeedPaper(units byte) error {
	return p.Write(p.Protocol.FeedPaper(units))
}

// ============================================================================
// Text Formatting Methods
// ============================================================================

// Bold enables or disables bold text
func (p *Printer) Bold(enable bool) error {
	return p.Write(p.Protocol.Bold(enable))
}

// Underline sets underline mode (0=off, 1=single, 2=double)
func (p *Printer) Underline(mode byte) error {
	cmd, err := p.Protocol.Underline(mode)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// Align sets text alignment by name
func (p *Printer) Align(alignment string) error {
	cmd, err := p.Protocol.Align(alignment)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// AlignLeft sets left alignment
func (p *Printer) AlignLeft() error {
	cmd, err := p.Protocol.AlignLeft()
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// AlignCenter sets center alignment
func (p *Printer) AlignCenter() error {
	cmd, err := p.Protocol.AlignCenter()
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// AlignRight sets right alignment
func (p *Printer) AlignRight() error {
	cmd, err := p.Protocol.AlignRight()
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// Size sets text size (width and height multipliers 1-8)
func (p *Printer) Size(width, height byte) error {
	cmd, err := p.Protocol.Size(width, height)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// NormalSize resets text to normal size
func (p *Printer) NormalSize() error {
	return p.Write(p.Protocol.NormalSize())
}

// DoubleWidth enables or disables double width
func (p *Printer) DoubleWidth(enable bool) error {
	return p.Write(p.Protocol.DoubleWidth(enable))
}

// DoubleHeight enables or disables double height
func (p *Printer) DoubleHeight(enable bool) error {
	return p.Write(p.Protocol.DoubleHeight(enable))
}

// ============================================================================
// Paper Control Methods
// ============================================================================

// Cut performs a paper cut
func (p *Printer) Cut(partial bool) error {
	cmd, err := p.Protocol.Cut(partial)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// FullCut performs a full paper cut
func (p *Printer) FullCut() error {
	cmd, err := p.Protocol.FullCut()
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PartialCut performs a partial paper cut
func (p *Printer) PartialCut() error {
	cmd := p.Protocol.PartialCut()
	return p.Write(cmd)
}

// FeedAndCut feeds paper then cuts
func (p *Printer) FeedAndCut(lines byte, partial bool) error {
	cmd, err := p.Protocol.FeedAndCut(lines, partial)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// ============================================================================
// Composite/Template Methods
// ============================================================================

// PrintTitle prints text as title (centered, bold, double size)
func (p *Printer) PrintTitle(title string) error {
	// Center alignment
	if err := p.AlignCenter(); err != nil {
		return err
	}

	// BoldText on
	if err := p.Bold(true); err != nil {
		return err
	}

	// Double size
	if err := p.Size(2, 2); err != nil {
		return err
	}

	// Print the title
	if err := p.PrintLine(title); err != nil {
		return err
	}

	// Reset to normal
	if err := p.NormalSize(); err != nil {
		return err
	}

	if err := p.Bold(false); err != nil {
		return err
	}

	return p.AlignLeft()
}

// PrintSeparator prints a line separator
func (p *Printer) PrintSeparator(char string, width int) error {
	separator := ""
	for i := 0; i < width; i++ {
		separator += char
	}
	return p.PrintLine(separator)
}

// PrintHeader is a convenience method for printing headers
func (p *Printer) PrintHeader(text string) error {
	if err := p.Bold(true); err != nil {
		return err
	}
	if err := p.PrintLine(text); err != nil {
		return err
	}
	return p.Bold(false)
}
