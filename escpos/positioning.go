package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// TODO: Comandos para posicionar texto e imágenes
// - Tabulación
// - Posicionamiento absoluto
// - Posicionamiento relativo

// SetPrintLeftMargin sets the left margin for printing
func (c *Protocol) SetPrintLeftMargin(_ byte) []byte {
	// TODO: Implementar usando GS L nL nH
	return []byte{}
}

// SetPrintWidth establece el ancho de impresión
func (c *Protocol) SetPrintWidth(_ byte) []byte {
	// TODO: Implementar usando GS W nL nH
	return []byte{}
}

var alignMap = map[Alignment]byte{
	AlignLeft:   0, // ESC/POS: 0 = left
	AlignCenter: 1, // ESC/POS: 1 = center
	AlignRight:  2, // ESC/POS: 2 = right
}

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (c *Protocol) SetJustification(justification Alignment) ([]byte, error) {
	alignment, ok := alignMap[justification]
	if !ok {
		return nil, fmt.Errorf("justificación no soportada: %v", justification)
	}
	// ESC a n
	return []byte{common.ESC, 'a', alignment}, nil
}
