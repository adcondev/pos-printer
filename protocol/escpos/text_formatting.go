package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para dar formato al texto
// - Doble ancho/altura
// - Rotación de texto
// - Espaciado de caracteres

// emphMap mapea el modo enfatizado a su valor ESC/POS correspondiente.
var emphMap = map[types.EmphasizedMode]byte{
	types.EmphOff: 0,
	types.EmphOn:  1,
}

// ulModeMap mapea el modo subrayado a su valor ESC/POS correspondiente.
var ulModeMap = map[types.UnderlineMode]byte{
	types.UnderNone:   0,
	types.UnderSingle: 1,
	types.UnderDouble: 2,
}

var fontMap = map[types.Font]byte{
	types.FontA:    0,
	types.FontB:    1,
	types.FontC:    2,
	types.FontD:    3,
	types.FontE:    4,
	types.SpecialA: 'a', // 97
	types.SpecialB: 'b', // 98
}

func (p *Commands) SelectCharacterFont(n types.Font) ([]byte, error) {
	font, ok := fontMap[n]
	if !ok {
		return nil, fmt.Errorf("no font found for font %v", n)
	}

	// ESC M n
	return []byte{ESC, 'M', font}, nil
}

func (p *Commands) TurnEmphasizedMode(n types.EmphasizedMode) ([]byte, error) {
	emph, ok := emphMap[n]
	if !ok {
		return nil, fmt.Errorf("no emph mode found")
	}

	return []byte{ESC, 'E', emph}, nil
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *Commands) SetDoubleStrike(on bool) []byte {
	val := byte(0)
	if on {
		val = 1
	}
	// ESC G n
	return []byte{ESC, 'G', val}
}

func (p *Commands) TurnUnderlineMode(n types.UnderlineMode) ([]byte, error) {
	mode, ok := ulModeMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid underline mode: %d", n)
	}
	// ESC - n
	return []byte{ESC, '-', mode}, nil
}

// SetTextSize Implementar
func (p *Commands) SetTextSize(widthMultiplier, heightMultiplier int) []byte {
	// TODO: Implementar usando GS ! n
	// Hint: n = (widthMultiplier-1)<<4 | (heightMultiplier-1)
	return []byte{}
}
