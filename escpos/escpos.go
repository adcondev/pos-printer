package escpos

import (
	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/print"
	"github.com/adcondev/pos-printer/escpos/printposition"
	"github.com/adcondev/pos-printer/escpos/sharedcommands"
)

// Protocol implements the ESCPOS Commands
type Protocol struct {
	Print         print.Capability
	LineSpace     linespacing.Capability
	Character     character.Capability
	PrintPosition printposition.Capability
	Barcode       barcode.Capability
}

// Raw sends raw data without processing
func (c *Protocol) Raw(data []byte) ([]byte, error) {
	if err := sharedcommands.IsBufLenOk(data); err != nil {
		return nil, err
	}
	return data, nil
}

// NewEscposProtocol creates a new instance of the ESC/POS protocol
func NewEscposProtocol() *Protocol {
	return &Protocol{
		Print:         print.NewCommands(),
		LineSpace:     linespacing.NewCommands(),
		Character:     character.NewCommands(),
		PrintPosition: printposition.NewCommands(),
		Barcode:       barcode.NewCommands(),
	}
}
