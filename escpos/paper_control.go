package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and constants
// ============================================================================
// Map que convierte CutPaper (enum del paquete) a los bytes ESC/POS
var cutMap = map[CutPaper]byte{
	FullCut:    '0',
	PartialCut: '1',
}

// ============================================================================
// Public API (implementation)
// ============================================================================

// Cut genera comando de corte
func (c *Protocol) Cut(mode CutPaper) ([]byte, error) {
	cut, ok := cutMap[mode]
	if !ok {
		return nil, fmt.Errorf("invalid cut mode: %v", mode)
	}

	return []byte{shared.GS, 'V', cut}, nil
}
