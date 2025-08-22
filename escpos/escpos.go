package escpos

// Commands implements the ESC/POS Protocol
type Commands struct {
	Print PrintCommands
}

// PrintCommands groups printing-related commands
type PrintCommands struct {
	Page PagePrint
}

// Raw sends raw data without processing
func (c *Commands) Raw(n string) ([]byte, error) {
	if err := isBufOk([]byte(n)); err != nil {
		return nil, err
	}
	// Sin procesar
	return []byte(n), nil
}

// Text formats and sends a string for printing
func (p *PrintCommands) Text(n string) ([]byte, error) {
	if err := isBufOk([]byte(n)); err != nil {
		return nil, err
	}
	return format([]byte(n)), nil
}

// PagePrint define los comandos de modo página de ESC/POS
type PagePrint struct{}

// NewESCPOSProtocol crea una nueva instancia del protocolo ESC/POS
func NewESCPOSProtocol() *Commands {
	c := &Commands{
		Print: PrintCommands{
			Page: PagePrint{},
		},
	}
	return c
}
