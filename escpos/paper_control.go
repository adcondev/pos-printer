package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and constants
// ============================================================================

// cutMap mapea los modos de corte a sus valores ESC/POS correspondientes.
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
