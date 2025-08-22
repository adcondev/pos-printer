package escpos

// isBufOk validates if the buffer size is within acceptable limits.
func isBufOk(buf []byte) error {
	if len(buf) < MinBuf {
		return errEmptyBuffer
	}
	if len(buf) > MaxBuf {
		return errBufferOverflow
	}
	return nil
}

// format replaces specific characters in the byte slice with their ESC/POS equivalents.
func format(data []byte) []byte {
	for i := range data {
		switch data[i] {
		case '\n':
			data[i] = LF
		case '\t':
			data[i] = HT
		case '\r':
			data[i] = CR
		default:
			// Do nothing, keep the character as is
		}
	}
	return data
}

// lengthLowHigh convierte una longitud en dos bytes little-endian (dL,dH)
// para usar en comandos ESC/Z.
// length debe estar entre 0 y 0xFFFF (65535), de lo contrario devuelve error.
func lengthLowHigh(length int) (dL, dH byte, err error) {
	if length < 0 {
		return 0, 0, errNegativeInt
	}
	dL = byte(length & 0xFF)        // byte de menor peso
	dH = byte((length >> 8) & 0xFF) // byte de mayor peso
	return dL, dH, nil
}
