package escpos

// isBufOk validates if the buffer size is within acceptable limits.
func isBufOk(buf []byte) error {
	if len(buf) < MinBuf {
		return ErrEmptyBuffer
	}
	if len(buf) > MaxBuf {
		return ErrBufferOverflow
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
