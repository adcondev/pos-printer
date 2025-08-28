package escpos

import (
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/lineSpacing"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Commands implements the ESC/POS Protocol
type Commands struct {
	Print     print.Capability
	LineSpace lineSpacing.Capability
}

// Raw sends raw data without processing
func (c *Commands) Raw(n string) ([]byte, error) {
	if err := common.IsBufOk([]byte(n)); err != nil {
		return nil, err
	}
	return []byte(n), nil
}

// NewEscposProtocol creates a new instance of the ESC/POS protocol
// Using Escpos (all caps) for consistency with the protocol name
func NewEscposProtocol() *Commands {
	return &Commands{
		Print: &print.Commands{
			Page: &print.PagePrint{},
		},
		LineSpace: &lineSpacing.Commands{},
	}
}
