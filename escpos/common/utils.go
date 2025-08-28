package common

// IsBufOk validates if the buffer size is within acceptable limits.
func IsBufOk(buf []byte) error {
	if len(buf) < MinBuf {
		return ErrEmptyBuffer
	}
	if len(buf) > MaxBuf {
		return ErrBufferOverflow
	}
	return nil
}

// Format replaces specific characters in the byte slice with their ESC/POS equivalents.
func Format(data []byte) []byte {
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

// LengthLowHigh convierte una longitud en dos bytes little-endian (dL,dH) para usar en comandos ESCPOS.
func LengthLowHigh(length int) (dL, dH byte, err error) {
	if length < 0 || length > 0xFFFF {
		return 0, 0, ErrLengthOutOfRange
	}
	dL = byte(length & 0xFF)        // byte de menor peso
	dH = byte((length >> 8) & 0xFF) // byte de mayor peso
	return dL, dH, nil
}
