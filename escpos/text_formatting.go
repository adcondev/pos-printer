package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and helpers
// ============================================================================
// Mapas usados para formateo de texto. Comentarios traducidos para principiantes.

// TODO: Comandos para dar formato al texto
// - Rotación de texto
// - Espaciado de caracteres

// emphMap mapea el modo enfatizado a su valor ESC/POS correspondiente.
var emphMap = map[EmphasizedMode]byte{
	EmphasizedOff: 0,
	EmphasizedOn:  1,
}

// ulModeMap mapea el modo subrayado a su valor ESC/POS correspondiente.
var ulModeMap = map[UnderlineMode]byte{
	UnderNone:   0,
	UnderSingle: 1,
	UnderDouble: 2,
}

var fontMap = map[Font]byte{
	FontA:    0,
	FontB:    1,
	FontC:    2,
	FontD:    3,
	FontE:    4,
	SpecialA: 'a', // 97
	SpecialB: 'b', // 98
}

// ============================================================================
// Public API (implementation)
// ============================================================================
// Funciones públicas para cambiar formatos de texto.

// SelectCharacterFont sets the character font
func (c *Protocol) SelectCharacterFont(n Font) ([]byte, error) {
	font, ok := fontMap[n]
	if !ok {
		return nil, fmt.Errorf("no font found for font %v", n)
	}

	// ESC M n
	return []byte{shared.ESC, 'M', font}, nil
}

// TurnEmphasizedMode enables or disables emphasized mode
func (c *Protocol) TurnEmphasizedMode(n EmphasizedMode) ([]byte, error) {
	emph, ok := emphMap[n]
	if !ok {
		return nil, fmt.Errorf("no emph mode found")
	}

	return []byte{shared.ESC, 'E', emph}, nil
}

// SetDoubleStrike activa/desactiva doble golpe
func (c *Protocol) SetDoubleStrike(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC G n
	return []byte{shared.ESC, 'G', val}
}

// TurnUnderlineMode enables or disables underline mode
func (c *Protocol) TurnUnderlineMode(n UnderlineMode) ([]byte, error) {
	mode, ok := ulModeMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid underline mode: %d", n)
	}
	// ESC - n
	return []byte{shared.ESC, '-', mode}, nil
}
