package common

import "errors"

// Buffer limits
var (
	// MinBuf es el tamaño mínimo del buffer
	MinBuf = 1
	// MaxBuf es el tamaño máximo del buffer
	MaxBuf = 65535
)

var (
	// ErrLengthOutOfRange length is out of range (0-65535)
	ErrLengthOutOfRange = errors.New("length is out of range (0-65535)")
	// ErrBufferOverflow buffer is too large
	ErrBufferOverflow = errors.New("can't print overflowed buffer (protocol max 64KB; model may be lower)")
	// ErrEmptyBuffer buffer is empty
	ErrEmptyBuffer = errors.New("can't print an empty buffer")
)

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
func LengthLowHigh(length uint16) (dL, dH byte) {
	dL = byte(length & 0xFF)        // byte de menor peso
	dH = byte((length >> 8) & 0xFF) // byte de mayor peso
	return dL, dH
}
