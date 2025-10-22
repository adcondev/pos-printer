package sharedcommands

const (
	// DLE represent the byte de "Data Link Escape" en ESC/POS.
	DLE byte = 0x10
	// EOT represent the byte de "End of Transmission" en ESC/POS.
	EOT byte = 0x04
	// ENQ represent the byte de "Enquiry" en ESC/POS.
	ENQ byte = 0x05
	// DC4 represent the byte de "Device Control 4" en ESC/POS.
	DC4 byte = 0x14
	// SP represent the byte de "Space" en ESC/POS.
	SP byte = 0x20 // Espacio (carácter de espacio en blanco)
	// FS represent the byte de "File Separator" en ESC/POS.
	FS byte = 0x1C
	// NUL represents the byte de "Null" en ESC/POS.
	NUL byte = 0x00
	// ESC represents the byte de "Escape" en ESC/POS.
	ESC byte = 0x1B
	// GS represents the byte de "Group Separator" en ESC/POS.
	GS byte = 0x1D
)
