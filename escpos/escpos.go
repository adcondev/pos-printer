package escpos

import (
	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/bitimage"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/mechanismcontrol"
	"github.com/adcondev/pos-printer/escpos/print"
	"github.com/adcondev/pos-printer/escpos/printposition"
	"github.com/adcondev/pos-printer/escpos/shared"
)

// Protocol implements the ESCPOS Commands
type Protocol struct {
	Barcode          barcode.Capability
	Bitimage         bitimage.Capability
	Character        character.Capability
	LineSpacing      linespacing.Capability
	MechanismControl mechanismcontrol.Capability
	Print            print.Capability
	PrintPosition    printposition.Capability
}

// Raw sends raw data without processing
func (c *Protocol) Raw(data []byte) ([]byte, error) {
	if err := shared.IsBufLenOk(data); err != nil {
		return nil, err
	}
	return data, nil
}

// NewEscposProtocol creates a new instance of the ESC/POS protocol
func NewEscposProtocol() *Protocol {
	return &Protocol{
		Barcode:          barcode.NewCommands(),
		Bitimage:         bitimage.NewCommands(),
		Character:        character.NewCommands(),
		LineSpacing:      linespacing.NewCommands(),
		MechanismControl: mechanismcontrol.NewCommands(),
		Print:            print.NewCommands(),
		PrintPosition:    printposition.NewCommands(),
	}
}
