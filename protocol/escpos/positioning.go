package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/protocol/escpos/types"
)

// TODO: Comandos para posicionar texto e imágenes
// - Tabulación
// - Posicionamiento absoluto
// - Posicionamiento relativo

const HT byte = 0x09

func (p *Commands) SetPrintLeftMargin(margin int) []byte {
	// TODO: Implementar usando GS L nL nH
	return []byte{}
}

func (p *Commands) SetPrintWidth(width int) []byte {
	// TODO: Implementar usando GS W nL nH
	return []byte{}
}

// SetJustification convierte el tipo genérico al específico de ESC/POS
func (p *Commands) SetJustification(justification types.Alignment) []byte {
	// Mapear el tipo genérico al valor ESC/POS
	var escposValue byte
	switch justification {
	case types.AlignLeft:
		escposValue = 0 // ESC/POS: 0 = left
	case types.AlignCenter:
		escposValue = 1 // ESC/POS: 1 = center
	case types.AlignRight:
		escposValue = 2 // ESC/POS: 2 = right
	default:
		escposValue = 0 // Default to left
	}

	// ESC a n
	return []byte{ESC, 'a', escposValue}
}

func SetHorizontalTabPositions(n types.TabColumnNumber, k types.TabTotalPosition) ([]byte, error) {
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

func (p *Commands) SetLineSpacing(n types.LineSpace) []byte {
	return []byte{ESC, '3', byte(n)}
}

func SelectDefaultLineSpacing() []byte {
	return []byte{ESC, '2'}
}
