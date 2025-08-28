package common

// IsBufLenOk validates if the buffer size is within acceptable limits.
func IsBufLenOk(buf []byte) error {
	if len(buf) < MinBuf {
		return ErrEmptyBuffer
	}
	if len(buf) > MaxBuf {
		return ErrBufferOverflow
	}
	return nil
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
