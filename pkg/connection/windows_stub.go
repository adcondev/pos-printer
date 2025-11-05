//go:build !windows

package connection

import (
	"errors"
)

// WindowsPrintConnector es un stub para sistemas no-Windows
type WindowsPrintConnector struct{}

// NewWindowsPrintConnector devuelve un error en sistemas no-Windows
func NewWindowsPrintConnector(_ string) (*WindowsPrintConnector, error) {
	return nil, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// Write implementación para sistemas no-Windows
func (c *WindowsPrintConnector) Write(_ []byte) (int, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// Read implementación para sistemas no-Windows
func (c *WindowsPrintConnector) Read(_ []byte) (int, error) {
	return 0, errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}

// Close implementación para sistemas no-Windows
func (c *WindowsPrintConnector) Close() error {
	return errors.New("WindowsPrintConnector no está disponible en este sistema operativo")
}
