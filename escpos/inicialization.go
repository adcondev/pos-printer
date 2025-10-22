// ============================================================================
// Initialization helpers
// ============================================================================
// Comandos básicos de inicialización y cambio de modo.

package escpos

import "github.com/adcondev/pos-printer/escpos/shared"

// InitializePrinter restores the printer to its default state
func (c *Protocol) InitializePrinter() []byte {
	// ESC @ - Reset printer
	return []byte{shared.ESC, '@'}
}

// SelectStandardMode sets the printer to standard mode
func SelectStandardMode() []byte {
	return []byte{shared.ESC, 'S'}
}

// SelectPageMode sets the printer to page mode
func SelectPageMode() []byte {
	return []byte{shared.ESC, 'L'}
}
