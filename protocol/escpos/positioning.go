package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para posicionar texto e imágenes
// - Tabulación
// - Posicionamiento absoluto
// - Posicionamiento relativo

// HT representa el comando Tabulación Horizontal.
//
// Nombre:
//
//	Tabulación horizontal
//
// Formato:
//
//	ASCII: HT
//	Hex: 09
//	Decimal: 9
//
// Descripción:
//
//	Mueve la posición de impresión a la siguiente posición de tabulación horizontal.
//
// Detalles:
//   - Este comando se ignora a menos que se haya configurado la siguiente posición de tabulación horizontal.
//   - Si la siguiente posición de tabulación horizontal excede el área de impresión, la impresora establece la posición de impresión en [Ancho del área de impresión + 1].
//   - Las posiciones de tabulación horizontal se configuran con ESC D.
//   - Si este comando se recibe cuando la posición de impresión está en [Ancho del área de impresión + 1], la impresora ejecuta la impresión del búfer lleno de la línea actual y procesa la tabulación horizontal desde el inicio de la siguiente línea.
//   - La configuración predeterminada de la posición de tabulación horizontal para el rollo de papel es la fuente A (12 × 24) cada 8 caracteres (9°, 17°, 25°, ... columna).
//
// Referencia:
//
//	ESC D
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

// SetHorizontalTabPositions representa el comando ESC D n1...nk NUL para establecer posiciones de tabulación horizontales.
//
// Nombre:
//
//	Establecer posiciones de tabulación horizontales
//
// Formato:
//
//	ASCII: ESC D n1...nk NUL
//	Hex:   1B 44 n1...nk 00
//	Decimal: 27 68 n1...nk 0
//
// Rango:
//   - n: 1 ≤ n ≤ 255 (cada n especifica un número de columna desde el inicio de la línea)
//   - k: 0 ≤ k ≤ 32 (k indica la cantidad total de posiciones de tabulación que se pueden establecer)
//
// Descripción:
//
//	Configura las posiciones de tabulación horizontales. Cada valor n representa la columna en la que se establecerá una posición de tabulación,
//	contada desde el inicio de la línea. Al enviar un código NUL (0) al final, se indica el fin de la secuencia de tabulaciones.
//	Si se envía ESC D NUL, se cancelan todas las posiciones de tabulación horizontales previamente definidas.
//
// Detalles:
//   - Las posiciones de tabulación se almacenan como el producto de [ancho de carácter × n], medido desde el inicio de la línea.
//     El ancho de carácter incluye el espaciado a la derecha, y los caracteres de doble ancho se configuran con el doble del ancho normal.
//   - Este comando borra las configuraciones de tabulación horizontales anteriores.
//   - Al establecer n = 8, la posición de impresión se mueve a la columna 9 mediante el envío del carácter HT.
//   - Se pueden establecer hasta 32 posiciones de tabulación (k = 32). Cualquier dato que exceda de 32 valores se procesa como datos normales.
//   - Se deben enviar los valores de [n] en orden ascendente y finalizar con un código NUL (0).
//   - Si un valor de [n] es menor o igual que el valor anterior (n[k] ≤ n[k-1]), se concluye la configuración y los datos siguientes se interpretan como datos normales.
//   - Las posiciones de tabulación previamente configuradas no cambian, incluso si el ancho del carácter se modifica posteriormente.
//   - El ancho de carácter se memoriza de forma independiente para el modo estándar y el modo página.
//
// Valor por defecto:
//
//	Las posiciones de tabulación por defecto se establecen a intervalos de 8 caracteres (por ejemplo, columnas 9, 17, 25, ...)
//	para la Fuente A (12x24).
//
// Referencia:
//
//	HT
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

// SetLineSpacing representa el comando ESC 3 n para configurar el espaciado de línea.
//
// Nombre:
//
//	Configurar espaciado de línea
//
// Formato:
//
//	ASCII: ESC 3 n
//	Hex:   1B 33 n
//	Decimal: 27 51 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Establece el espaciado de línea a [n × (unidad de movimiento vertical u horizontal)] pulgadas.
//
// Detalles:
//   - El espaciado de línea se puede configurar de manera independiente en el modo estándar y en el modo página.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P. Cambiar la unidad horizontal o vertical
//     no afecta el espaciado de línea actual.
//   - En el modo estándar se utiliza la unidad de movimiento vertical (y).
//   - En el modo página, el comando funciona de la siguiente manera según la posición inicial del área imprimible definida con ESC T:
//     ① Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad vertical (y).
//     ② Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad horizontal (x).
//   - La cantidad máxima de alimentación de papel es de 1016 mm (40 pulgadas). Aunque se configure un valor mayor,
//     la impresora alimenta el papel únicamente hasta 1016 mm (40 pulgadas).
//
// Valor por Defecto:
//
//	Espaciado de línea equivalente a aproximadamente 4,23 mm (1/6 de pulgada).
//
// Referencia:
//
//	ESC 2, GS P
func (p *Commands) SetLineSpacing(n types.LineSpace) []byte {
	return []byte{ESC, '3', byte(n)}
}

// SelectDefaultLineSpacing representa el comando ESC 2 para seleccionar el espaciado de línea por defecto.
//
// Nombre:
//
//	Seleccionar espaciado de línea por defecto
//
// Formato:
//
//	ASCII: ESC 2
//	Hex:   1B 32
//	Decimal: 27 50
//
// Descripción:
//
//	Selecciona un espaciado de línea de 1/6 de pulgada (aproximadamente 4,23 mm).
//
// Detalles:
//   - El espaciado de línea se puede configurar de manera independiente en el modo estándar y en el modo página.
//
// Referencia:
//
//	ESC 3
func SelectDefaultLineSpacing() []byte {
	return []byte{ESC, '2'}
}
