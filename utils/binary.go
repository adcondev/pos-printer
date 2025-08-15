package utils

import (
	"fmt"
)

// LengthLowHigh convierte una longitud en dos bytes little-endian (dL,dH)
// para usar en comandos ESC/Z.
// length debe estar entre 0 y 0xFFFF (65535), de lo contrario devuelve error.
func LengthLowHigh(length int) (dL, dH byte, err error) {
	if length < 0 || length > 0xFFFF {
		return 0, 0, fmt.Errorf("lengthLowHigh: longitud fuera de rango 0–65535, recibida %d", length)
	}
	dL = byte(length & 0xFF)        // byte de menor peso
	dH = byte((length >> 8) & 0xFF) // byte de mayor peso
	return dL, dH, nil
}
