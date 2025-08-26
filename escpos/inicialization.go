package escpos

// TODO: Comandos para inicialización y configuración básica de la impresora
// - Configuración de página de códigos
// - Configuración regional
// - Reinicio de impresora

// InitializePrinter restores the printer to its default state
func (c *Commands) InitializePrinter() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// SelectStandardMode sets the printer to standard mode
func SelectStandardMode() []byte {
	return []byte{ESC, 'S'}
}

// SelectPageMode sets the printer to page mode
func SelectPageMode() []byte {
	return []byte{ESC, 'L'}
}
