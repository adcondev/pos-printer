package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and helpers
// ============================================================================

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

// fontMap mapea las fuentes de caracteres a sus valores ESC/POS correspondientes.
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

// TurnEmphasizedMode enables or disables emphasized mode
func (c *Commands) TurnEmphasizedMode(n EmphasizedMode) ([]byte, error) {
	emph, ok := emphMap[n]
	if !ok {
		return nil, fmt.Errorf("no emph mode found")
	}

	return []byte{shared.ESC, 'E', emph}, nil
}

// SetDoubleStrike activa/desactiva doble golpe
func (c *Commands) SetDoubleStrike(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC G n
	return []byte{shared.ESC, 'G', val}
}

// TurnUnderlineMode enables or disables underline mode
func (c *Commands) TurnUnderlineMode(n UnderlineMode) ([]byte, error) {
	mode, ok := ulModeMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid underline mode: %d", n)
	}
	// ESC - n
	return []byte{shared.ESC, '-', mode}, nil
}
