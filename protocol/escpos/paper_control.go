package escpos

import (
	"github.com/AdConDev/pos-printer/protocol/escpos/types"
)

// TODO: Comandos para control y manejo del papel
// - Corte de papel (parcial/completo)
// - Expulsión de papel

const (
	// Modo de corte de papel

	Cut     byte = 49 // 'A'
	CutFeed byte = 66 // 'B'
)

// Cut genera comando de corte
func (p *Commands) Cut(mode types.CutMode, lines int) []byte {
	// TODO: Implementar validación de lines

	cmd := []byte{GS, 'V'}

	switch mode {
	case types.CutFeed:
		cmd = append(cmd, CutFeed, byte(lines)) // o 48 ('0') según el modelo
	case types.Cut:
		cmd = append(cmd, Cut) // o 49 ('1') según el modelo
	default:
		cmd = append(cmd, 0)
	}

	return cmd
}

// Feed genera comando de alimentación de papel
func (p *Commands) Feed(lines int) []byte {
	// TODO: Validar que lines esté en rango válido
	if lines <= 0 {
		return []byte{}
	}

	// ESC d n - Print and feed n lines
	return []byte{ESC, 'd', byte(lines)}
}

func (p *Commands) FeedReverse(lines int) []byte {
	// TODO: Implementar ESC e n
	return []byte{}
}

func (p *Commands) FeedForm() []byte {
	// TODO: Implementar FF
	return []byte{}
}
