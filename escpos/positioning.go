package escpos

import (
	"fmt"
)

// TODO: Comandos para posicionar texto e imágenes
// - Tabulación
// - Posicionamiento absoluto
// - Posicionamiento relativo

// HT is the control code for horizontal tab
const HT byte = 0x09

// SetPrintLeftMargin sets the left margin for printing
func (c *Commands) SetPrintLeftMargin(_ byte) []byte {
	// TODO: Implementar usando GS L nL nH
	return []byte{}
}

// SetPrintWidth establece el ancho de impresión
func (c *Commands) SetPrintWidth(_ byte) []byte {
	// TODO: Implementar usando GS W nL nH
	return []byte{}
}

var alignMap = map[Alignment]byte{
	AlignLeft:   0, // ESC/POS: 0 = left
	AlignCenter: 1, // ESC/POS: 1 = center
	AlignRight:  2, // ESC/POS: 2 = right
}

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (c *Commands) SetJustification(justification Alignment) ([]byte, error) {
	alignment, ok := alignMap[justification]
	if !ok {
		return nil, fmt.Errorf("justificación no soportada: %v", justification)
	}
	// ESC a n
	return []byte{ESC, 'a', alignment}, nil
}

// SetHorizontalTabPositions establece las posiciones de tabulación horizontal
func SetHorizontalTabPositions(n TabColumnNumber, k TabTotalPosition) ([]byte, error) {
	if k > 32 {
		return nil, fmt.Errorf("k fuera de rango (0-32): %d", k)
	}
	if len(n) != int(k) {
		return nil, fmt.Errorf("la longitud de n (%d) no coincide con k (%d)", len(n), k)
	}
	for i := 0; i < len(n)-1; i++ {
		if n[i] >= n[i+1] {
			return nil, fmt.Errorf("las posiciones de tabulación deben ser estrictamente ascendentes: %d no es menor que %d", n[i], n[i+1])
		}
	}

	cmd := []byte{ESC, 'D'}
	cmd = append(cmd, n...)
	cmd = append(cmd, NUL) // NUL al final

	return cmd, nil
}
