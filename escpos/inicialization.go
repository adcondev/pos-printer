package escpos

// TODO: Comandos para inicialización y configuración básica de la impresora
// - Configuración de página de códigos
// - Configuración regional
// - Reinicio de impresora

func (p *Commands) InitializePrinter() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// Close genera comandos de cierre (si los hay)
func (p *Commands) Close() []byte {
	// ESC/POS no tiene un comando específico de cierre
	// pero podrías incluir un reset o feed final o ambos
	return []byte{NUL}
}

func SelectStandardMode() []byte {
	return []byte{ESC, 'S'}
}

func SelectPageMode() []byte {
	return []byte{ESC, 'L'}
}
