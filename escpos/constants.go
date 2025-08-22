package escpos

// Constantes ESC/POS y parámetros.
const (
	// CAN representa el comando Cancelar datos de impresión en modo página.
	CAN byte = 0x18
	// DLE representa el byte de "Data Link Escape" en ESC/POS.
	DLE byte = 0x10
	// EOT representa el byte de "End Of Transmission" en ESC/POS.
	EOT byte = 0x04
	// ENQ representa el byte de "Enquiry" en ESC/POS.
	ENQ byte = 0x05
	// DC4 representa el byte "Device Control 4" en ESC/POS.
	DC4 byte = 0x14
	// SP
	SP byte = 0x20 // Espacio (carácter de espacio en blanco)
	// FS
	FS byte = 0x1C
	// GS
	GS byte = 0x1D
	// NUL
	NUL byte = 0x00
	// ESC Comando de escape para impresoras
	ESC byte = 0x1B
)
