package escpos

// Commands implements the ESC/POS Protocol
type Commands struct {
	Print     PrinterCapability
	LineSpace LineSpacingCapability
}

// Raw sends raw data without processing
func (c *Commands) Raw(n string) ([]byte, error) {
	if err := isBufOk([]byte(n)); err != nil {
		return nil, err
	}
	return []byte(n), nil
}

// NewEscposProtocol creates a new instance of the ESC/POS protocol
// Using Escpos (all caps) for consistency with the protocol name
func NewEscposProtocol() *Commands {
	return &Commands{
		Print: &PrintCommands{
			Page: &PagePrint{},
		},
		LineSpace: &LineSpacingCommands{},
	}
}
