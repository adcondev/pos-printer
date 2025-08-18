package escpos

// TODO: Comandos para inicialización y configuración básica de la impresora
// - Configuración de página de códigos
// - Configuración regional
// - Reinicio de impresora
// - Selección de modo

func (p *Commands) InitializePrinter() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// TODO: Pensar en comando compuesto para el final Feed(1) y Cut(1)

// Close genera comandos de cierre (si los hay)
func (p *Commands) Close() []byte {
	// ESC/POS no tiene un comando específico de cierre
	// pero podrías incluir un reset o feed final o ambos
	return []byte{}
}

func SelectStandardMode() []byte {
	return []byte{ESC, 'S'}
}

func SelectPageMode() []byte {
	return []byte{ESC, 'L'}
}
