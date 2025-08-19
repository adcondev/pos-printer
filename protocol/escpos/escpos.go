package escpos

// Commands implementa Protocol para ESC/POS
type Commands struct {
}

// NewESCPOSProtocol crea una nueva instancia del protocolo ESC/POS
func NewESCPOSProtocol() *Commands {
	c := &Commands{}
	return c
}
