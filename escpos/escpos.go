package escpos

import (
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/linespacing"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Commands implements the ESC/POS Protocol
type Commands struct {
	Print     print.Capability
	LineSpace linespacing.Capability
}

// Raw sends raw data without processing
func (c *Commands) Raw(n []byte) ([]byte, error) {
	if err := common.IsBufLenOk(n); err != nil {
		return nil, err
	}
	return n, nil
}

// NewEscposCommands creates a new instance of the ESC/POS protocol
// Using Escpos (all caps) for consistency with the protocol name
func NewEscposCommands() *Commands {
	return &Commands{
		Print: &print.Commands{
			Page: &print.PagePrint{},
		},
		LineSpace: &linespacing.Commands{},
	}
}
