package escpos

import (
	"fmt"
)

// TODO: Comandos para dar formato al texto
// - Doble ancho/altura
// - Rotación de texto
// - Espaciado de caracteres

// emphMap mapea el modo enfatizado a su valor ESC/POS correspondiente.
var emphMap = map[EmphasizedMode]byte{
	EmphOff: 0,
	EmphOn:  1,
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

func (p *Commands) SelectCharacterFont(n Font) ([]byte, error) {
	font, ok := fontMap[n]
	if !ok {
		return nil, fmt.Errorf("no font found for font %v", n)
	}

	// ESC M n
	return []byte{ESC, 'M', font}, nil
}

func (p *Commands) TurnEmphasizedMode(n EmphasizedMode) ([]byte, error) {
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

func (p *Commands) TurnUnderlineMode(n UnderlineMode) ([]byte, error) {
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
