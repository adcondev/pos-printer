package escpos

import (
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Protocol implements the ESCPOS Commands
type Protocol struct {
	Print     print.Capability
	LineSpace linespacing.Capability
	Character character.Capability
}

// Raw sends raw data without processing
func (c *Protocol) Raw(n []byte) ([]byte, error) {
	if err := common.IsBufLenOk(n); err != nil {
		return nil, err
	}
	return n, nil
}

// NewEscposCommands creates a new instance of the ESC/POS protocol
// Using Escpos (all caps) for consistency with the protocol name
func NewEscposCommands() *Protocol {
	return &Protocol{
		Print:     print.NewCommands(),
		LineSpace: linespacing.NewCommands(),
		Character: character.NewCommands(),
	}
}
