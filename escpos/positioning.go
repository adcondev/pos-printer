package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and helpers
// ============================================================================

// alignMap mapea las alineaciones genéricas a sus valores ESC/POS correspondientes.
var alignMap = map[Alignment]byte{
	AlignLeft:   0, // ESC/POS: 0 = left
	AlignCenter: 1, // ESC/POS: 1 = center
	AlignRight:  2, // ESC/POS: 2 = right
}

// ============================================================================
// Public API (implementation)
// ============================================================================
// Funciones que usan los mapas anteriores para generar comandos ESC/POS.

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (c *Commands) SetJustification(justification Alignment) ([]byte, error) {
	alignment, ok := alignMap[justification]
	if !ok {
		return nil, fmt.Errorf("justificación no soportada: %v", justification)
	}
	// ESC a n
	return []byte{shared.ESC, 'a', alignment}, nil
}
