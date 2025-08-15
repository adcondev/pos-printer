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
	types.FontA: 0,
	types.FontB: 1,
}

// SelectCharacterFont permite configurar la fuente de los caracteres en la impresora.
//
// Formato:
//
//	ASCII: ESC M n
//	Hex:   1B 4D n
//	Decimal: 27 77 n
//
// Rango:
//
//	n = 0, 1, 48, 49
//
// Descripción:
//
//	Selecciona la fuente de caracteres que será utilizada por la impresora:
//	  - n = 0 o 48: Se selecciona la fuente de caracteres A (12 × 24).
//	  - n = 1 o 49: Se selecciona la fuente de caracteres B (9 × 17).
//
// Detalles:
//   - Este comando determina el tamaño de fuente activo para la impresión de texto.
func (p *Commands) SelectCharacterFont(n types.Font) ([]byte, error) {
	font, ok := fontMap[n]
	if !ok {
		return nil, fmt.Errorf("no font found for font %v", n)
	}

	// ESC M n
	return []byte{ESC, 'M', font}, nil
}

// TurnEmphasizedMode representa el comando ESC E n para activar o desactivar el modo enfatizado.
//
// Nombre:
//
//	Activar/Desactivar modo enfatizado
//
// Formato:
//
//	ASCII: ESC E n
//	Hex:   1B 45 n
//	Decimal: 27 69 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Activa o desactiva el modo enfatizado basado en el valor del bit menos significativo (LSB) de n:
//	  - Si el LSB de n es 0, el modo enfatizado se desactiva.
//	  - Si el LSB de n es 1, el modo enfatizado se activa.
//
// Detalles:
//   - Solo el bit menos significativo (LSB) de n es utilizado para determinar el estado del modo enfatizado.
//   - Este comando y ESC ! activan o desactivan el modo enfatizado de la misma manera. Es importante tener cuidado
//     al utilizar ambos comandos simultáneamente, ya que el último comando recibido prevalecerá.
//
// Valor por defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC !
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

// TurnUnderlineMode representa el comando ESC - n para activar o desactivar el modo subrayado.
//
// Nombre:
//
//	Activar/Desactivar modo subrayado
//
// Formato:
//
//	ASCII: ESC - n
//	Hex:   1B 2D n
//	Decimal: 27 45 n
//
// Rango:
//
//	n puede tomar uno de los siguientes valores:
//	  • 0, 48: Desactiva el modo subrayado.
//	  • 1, 49: Activa el modo subrayado con grosor de 1 punto.
//	  • 2, 50: Activa el modo subrayado con grosor de 2 puntos.
//	Valores válidos: 0 ≤ n ≤ 2, 48 ≤ n ≤ 50
//
// Descripción:
//
//	Activa o desactiva el modo subrayado en la impresora basado en el valor de n recibido.
//
// Detalles:
//   - La impresora puede subrayar todos los caracteres (incluyendo el espaciado derecho), pero no puede subrayar
//     los espacios establecidos por HT.
//   - No es posible subrayar caracteres rotados 90° en sentido horario ni caracteres invertidos de blanco y negro.
//   - Cuando el modo subrayado se desactiva (n = 0 o 48), los datos posteriores no se subrayan y el grosor del subrayado
//     configurado previamente se mantiene sin cambios (el grosor por defecto es de 1 punto).
//   - Cambiar el tamaño de los caracteres no afecta el grosor de subrayado actual.
//   - El modo subrayado también puede activarse o desactivarse mediante ESC !. Sin embargo, solo se considera efectivo
//     el último comando recibido.
//   - Este comando no afecta la configuración de los caracteres Kanji.
//
// Valor por defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC !
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
