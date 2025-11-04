// Package service provides implementations for various POS printer service.
package service

import (
	"fmt"

	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/controllers/escpos/character"
	"github.com/adcondev/pos-printer/pkg/controllers/escpos/mechanismcontrol"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/profile"
)

// Printer represents a POS printer device
type Printer struct {
	Profile    profile.Escpos
	Connection connection.Connector
	Protocol   composer.Escpos
}

// NewPrinter creates a new Printer instance
func NewPrinter(proto *composer.Escpos, prof *profile.Escpos, conn connection.Connector) (*Printer, error) {
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
	cmd := p.Protocol.InitializePrinter()
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
	cmd, err := p.Protocol.Print.Text(text)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PrintLine sends text with line feed
func (p *Printer) PrintLine(text string) error {
	cmd, err := p.Protocol.PrintLn(text)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// FeedLines advances paper by n lines
func (p *Printer) FeedLines(lines byte) error {
	return p.Write(p.Protocol.Print.PrintAndFeedLines(lines))
}

// ============================================================================
// Text Formatting Methods
// ============================================================================

// FontA sets the font to Font A
func (p *Printer) FontA() error {
	cmd, _ := p.Protocol.Character.SelectCharacterFont(character.FontA)
	return p.Write(cmd)
}

// FontB sets the font to Font B
func (p *Printer) FontB() error {
	cmd, _ := p.Protocol.Character.SelectCharacterFont(character.FontB)
	return p.Write(cmd)
}

// Bold enables or disables bold text
func (p *Printer) Bold() error {
	return p.Write(p.Protocol.Character.SetEmphasizedMode(character.OnEm))
}

// AlignLeft sets left alignment
func (p *Printer) AlignLeft() error {
	cmd := p.Protocol.LeftAlign()
	return p.Write(cmd)
}

// AlignCenter sets center alignment
func (p *Printer) AlignCenter() error {
	cmd := p.Protocol.CenterAlign()
	return p.Write(cmd)
}

// AlignRight sets right alignment
func (p *Printer) AlignRight() error {
	cmd := p.Protocol.RightAlign()
	return p.Write(cmd)
}

// NormalSize resets text to normal size
func (p *Printer) NormalSize() error {
	return p.Write(p.Protocol.RegularTextSize())
}

// DoubleSize enables or disables double width
func (p *Printer) DoubleSize() error {
	return p.Write(p.Protocol.DoubleSizeText())
}

// ============================================================================
// Paper Control Methods
// ============================================================================

// FullCut performs a full paper cut
func (p *Printer) FullCut(lines byte) error {
	cmd := p.Protocol.FullPaperCut(lines)
	return p.Write(cmd)
}

// PartialFeedAndCut performs a partial paper cut
func (p *Printer) PartialFeedAndCut(lines byte) error {
	cmd, _ := p.Protocol.MechanismControl.FeedAndCutPaper(mechanismcontrol.FeedCutPartial, lines)
	return p.Write(cmd)
}

// ============================================================================
// Image Printing Methods
// ============================================================================

// PrintBitmap prints a monochrome bitmap using raster graphics
func (p *Printer) PrintBitmap(bitmap *graphics.MonochromeBitmap) error {
	if bitmap == nil {
		return fmt.Errorf("bitmap cannot be nil")
	}

	// Usar el comando GS v 0 para imprimir imagen raster
	cmd, err := p.Protocol.BitImage.PrintRasterBitImage(
		0,                              // Modo normal
		uint16(bitmap.GetWidthBytes()), // Ancho en BYTES, not pixels!
		uint16(bitmap.Height),          // Alto en puntos (dots)
		bitmap.GetRasterData(),         // Datos de imagen
	)
	if err != nil {
		return fmt.Errorf("generate raster command: %w", err)
	}

	return p.Write(cmd)
}
