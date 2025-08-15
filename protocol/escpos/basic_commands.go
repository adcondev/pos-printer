package escpos

import "strings"

const (
	// LF representa el comando Imprimir y Alimentar Línea.
	//
	// Nombre:
	//   Imprimir y alimentar línea
	//
	// Formato:
	//   ASCII: LF
	//   Hex: 0A
	//   Decimal: 10
	//
	// Descripción:
	//   Imprime los datos en el búfer de impresión y alimenta una línea según el espaciado de línea actual.
	//
	// Nota:
	//   Este comando establece la posición de impresión al inicio de la línea.
	//
	// Referencia:
	//   ESC 2, ESC 3
	LF byte = 0x0A

	// ESC representa el carácter de control Escape, utilizado como prefijo
	// para la mayoría de los comandos de control de la impresora.
	//
	// Nombre:
	//   Escape (prefijo de comando)
	//
	// Formato:
	//   ASCII: ESC
	//   Hex: 1B
	//   Decimal: 27
	//
	// Descripción:
	//   ESC es un carácter de escape que precede a muchos comandos ESC/POS.
	//   Indica que el siguiente byte (o secuencia) representa una instrucción
	//   de control para la impresora.
	//
	// Detalles:
	//   - La mayoría de los comandos de formato, espaciado, inicialización y alineación
	//     comienzan con este byte.
	//   - Es obligatorio para interpretar correctamente comandos como ESC @, ESC a n, etc.
	//   - No tiene efecto por sí solo; siempre debe ir seguido de un comando válido.
	//
	// Referencia:
	//   ESC
	ESC byte = 0x1B
)

// Text convierte texto a bytes con encoding apropiado
func (p *Commands) Text(str string) []byte {
	cmd := strings.ReplaceAll(str, "\n", string(LF))
	return []byte(cmd)
}

// TextLn agrega un salto de línea al final
func (p *Commands) TextLn(str string) []byte {
	text := p.Text(str)
	// Agregar LF al final
	return append(text, LF)
}

// TextRaw envía bytes sin procesar
func (p *Commands) TextRaw(str string) []byte {
	return []byte(str)
}

// PrintAndFeedPaper representa el comando ESC J n para imprimir y alimentar el papel.
//
// Nombre:
//
//	Imprimir y alimentar papel
//
// Formato:
//
//	ASCII: ESC J n
//	Hex:   1B 4A n
//	Decimal: 27 74 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Imprime los datos en el búfer de impresión y alimenta el papel en una cantidad equivalente a
//	[n × unidad de movimiento vertical u horizontal] pulgadas.
//
// Detalles:
//   - Después de completar la impresión, este comando establece la posición inicial de impresión al comienzo de la línea.
//   - La cantidad de alimentación de papel configurada por este comando no afecta los valores establecidos por ESC 2 o ESC 3.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P.
//   - El comando GS P puede cambiar las unidades de movimiento vertical y horizontal. Sin embargo, el valor no puede ser menor
//     que la cantidad mínima de movimiento vertical y debe ser un múltiplo par de dicha cantidad mínima.
//   - En modo estándar, la impresora utiliza la unidad de movimiento vertical (y).
//   - En modo página, el comando funciona de la siguiente manera según la posición inicial del área imprimible definida con ESC T:
//     ① Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad vertical (y).
//     ② Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad horizontal (x).
//   - La cantidad máxima de alimentación de papel es de 1016 mm (40 pulgadas). Si el valor configurado excede este límite,
//     se ajustará automáticamente al máximo permitido.
//
// Referencia:
//
//	GS P
func PrintAndFeedPaper(n byte) []byte {
	return []byte{ESC, 'J', n}
}
