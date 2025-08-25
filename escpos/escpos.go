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
	// Sin procesar
	return []byte(n), nil
}

// NewEscposProtocol crea una nueva instancia del protocolo ESC/POS
func NewEscposProtocol() *Commands {
	c := &Commands{
		Print: &PrintCommands{
			Page: &PagePrint{},
		},
		LineSpace: &LineSpacingCommands{},
	}
	return c
}
