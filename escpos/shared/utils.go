// Package shared provides shared utilities and constants for ESC/POS command packages.
//
// This package contains shared buffer validation functions, byte manipulation utilities,
// and shared constants used across multiple ESC/POS command implementations.
package shared

import "errors"

// ============================================================================
// Context
// ============================================================================
// This package provides shared utilities and constants for ESC/POS command packages.
// It includes buffer validation functions, little-endian byte conversion utilities,
// and shared error definitions used across multiple ESC/POS command implementations.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

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

// ============================================================================
// Validation and Helper Functions
// ============================================================================

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

// ToLittleEndian convierte una longitud en dos bytes little-endian (dL,dH) para usar en comandos ESCPOS.
func ToLittleEndian(number uint16) (nL, nH byte) {
	nL = byte(number & 0xFF)        // byte de menor peso
	nH = byte((number >> 8) & 0xFF) // byte de mayor peso
	return nL, nH
}
