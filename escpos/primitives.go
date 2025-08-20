package escpos

// Constantes ESC/POS y parámetros.
const (
	// FF representa el comando Imprimir y Regresar al Modo Estándar en Modo Página.
	//
	// Nombre:
	//   Imprimir y regresar al modo estándar en modo página
	//
	// Formato:
	//   ASCII: FF
	//   Hex: 0C
	//   Decimal: 12
	//
	// Descripción:
	//   Imprime los datos en el búfer de impresión de forma colectiva y regresa al modo estándar.
	//
	// Detalles:
	//   - Los datos del búfer se eliminan después de ser impresos.
	//   - El área de impresión configurada por ESC W se restablece a la configuración predeterminada.
	//   - La impresora no ejecuta el corte de papel.
	//   - Este comando establece la posición de impresión al inicio de la línea.
	//   - Este comando está habilitado únicamente en modo página.
	//
	// Referencia:
	//   ESC FF, ESC L, ESC S
	FF byte = 0x0C

	// CAN representa el comando Cancelar datos de impresión en modo página.
	//
	// Nombre:
	//   Cancelar datos de impresión en modo página
	//
	// Formato:
	//   ASCII: CAN
	//   Hex: 18
	//   Decimal: 24
	//
	// Descripción:
	//   En modo página, elimina todos los datos de impresión en el área imprimible actual.
	//
	// Detalles:
	//   - Este comando está habilitado únicamente en modo página.
	//   - Si existen datos en el área de impresión previamente especificada que también están en el área de impresión actualmente especificada, estos se eliminan.
	//
	// Referencia:
	//   ESC L, ESC W
	CAN byte = 0x18

	// DLE representa el byte de "Data Link Escape" en ESC/POS.
	//
	// Nombre:
	//   Data Link Escape (DLE)
	//
	// Valor:
	//   Hex: 0x10
	//   Decimal: 16
	//
	// Descripción:
	//   DLE se utiliza en protocolos de comunicación, incluyendo ESC/POS, para señalar el inicio
	//   de una secuencia de comandos o para distinguir entre datos y comandos en la transmisión.
	DLE byte = 0x10

	// EOT representa el byte de "End Of Transmission" en ESC/POS.
	//
	// Nombre:
	//   End Of Transmission (EOT)
	//
	// Valor:
	//   Hex: 0x04
	//   Decimal: 4
	//
	// Descripción:
	//   EOT se utiliza para indicar el final de una transmisión. En el contexto de ESC/POS, se emplea
	//   en comandos de transmisión en tiempo real para marcar el final de la comunicación o
	//   para solicitar información de estado.
	EOT byte = 0x04

	// ENQ representa el byte de "Enquiry" en ESC/POS.
	//
	// Nombre:
	//   Enquire (ENQ)
	//
	// Valor:
	//   Hex: 0x05
	//   Decimal: 5
	//
	// Descripción:
	//   ENQ se utiliza para solicitar una respuesta o confirmación del dispositivo. En sistemas ESC/POS,
	//   es comúnmente empleado en comandos de solicitud en tiempo real para recuperar información o
	//   reiniciar estados.
	ENQ byte = 0x05

	// DC4 representa el byte "Device Control 4" en ESC/POS.
	//
	// Nombre:
	//   Device Control 4 (DC4)
	//
	// Valor:
	//   Hex: 0x14
	//   Decimal: 20
	//
	// Descripción:
	//   DC4 es un byte de control utilizado en la comunicación de datos para funciones específicas.
	//   En el contexto de ESC/POS, DC4 forma parte de comandos que generan pulsos en tiempo real,
	//   tales como el comando DLE DC4 n m t.
	DC4 byte = 0x14

	SP  byte = 0x20 // Espacio (carácter de espacio en blanco)
	FS  byte = 0x1C
	GS  byte = 0x1D
	NUL byte = 0x00
)

var (
	// PrintDataPageMode representa el comando para imprimir datos en modo página.
	//
	// Nombre:
	//   Imprimir datos en modo página
	//
	// Formato:
	//   ASCII: ESC FF
	//   Hex: 1B 0C
	//   Decimal: 27 12
	//
	// Descripción:
	//   En modo página, imprime todos los datos almacenados en el búfer dentro del área de impresión de manera colectiva.
	//
	// Detalles:
	//   - Este comando está habilitado únicamente en modo página.
	//   - Después de imprimir, la impresora no borra los datos almacenados en el búfer, los valores configurados para ESC T y ESC W, ni la posición para almacenar datos de caracteres.
	//
	// Referencia:
	//   FF, ESC L, ESC S
	PrintDataPageMode = []byte{ESC, FF}
)

// RequestPrinterRecovery representa el comando para realizar una petición en tiempo real a la impresora.
//
// Nombre:
//
//	Petición en tiempo real a la impresora
//
// Formato:
//
//	ASCII: DLE ENQ n
//	Hex:   10 05 n
//	Decimal: 16 5 n
//
// Rango:
//
//	1 ≤ n ≤ 2
//
// Descripción:
//
//	Responde a una solicitud de la computadora host. El valor de n especifica la solicitud de la siguiente manera:
//	  n = 1: Recuperar de un error y reiniciar la impresión desde la línea en la que ocurrió el error.
//	  n = 2: Recuperar de un error después de borrar los búferes de recepción e impresión.
//
// Detalles:
//   - Este comando es efectivo únicamente cuando ocurre un error en la autocortadora.
//   - La impresora comienza a procesar los datos al recibir este comando.
//   - Se ejecuta incluso si la impresora está fuera de línea, el búfer de recepción está lleno o se presenta un estado de error en modelos de interfaz serial.
//   - En modelos de interfaz paralela, el comando no se ejecuta cuando la impresora está ocupada; sin embargo, se ejecuta cuando la impresora está fuera de línea o hay error, si el DIP switch 2-1 está activado.
//   - El estado también se transmite cada vez que se recibe la secuencia de datos <10>H <05>H <n> (1 ≤ n ≤ 2).
//     Ejemplo: En el comando "ESC * m nL nH dk", se tiene: d1 = <10>H, d2 = <05>H, d3 = <01>H.
//   - Este comando no debe incluirse dentro de otra secuencia de comandos que consista en dos o más bytes.
//   - DLE ENQ 2 permite que la impresora se recupere de un error luego de borrar los datos en los búferes de recepción e impresión, conservando las configuraciones (por ejemplo, las establecidas por ESC !, ESC 3, etc.) vigentes al ocurrir el error.
//   - La impresora se puede inicializar completamente utilizando este comando junto con ESC @.
//   - Cuando la impresora se deshabilita con ESC = (seleccionar dispositivo periférico), las funciones de recuperación de error (DLE ENQ 1 y DLE ENQ 2) quedan habilitadas y las demás funciones se desactivan.
//
// Referencia:
//
//	DLE EOT
func RequestPrinterRecovery(n byte) []byte {
	cmd := []byte{DLE, ENQ, n}
	return cmd
}

// SetRightSideCharacterSpacing representa el comando para configurar el espaciado a la derecha de los caracteres.
//
// Nombre:
//
//	Configurar espaciado a la derecha de los caracteres
//
// Formato:
//
//	ASCII: ESC SP n
//	Hex: 1B 20 n
//	Decimal: 27 32 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Configura el espaciado de caracteres en el lado derecho del carácter a [n unidades de movimiento horizontal o vertical].
//
// Detalles:
//   - El espaciado a la derecha de los caracteres en modo de doble ancho es el doble del valor normal. Cuando los caracteres se agrandan, el espaciado a la derecha es n veces el valor normal.
//   - Este comando no afecta la configuración de caracteres kanji.
//   - Este comando establece valores de manera independiente en cada modo (modo estándar y modo página).
//   - Las unidades de movimiento horizontal y vertical se especifican mediante el comando GS P. Cambiar la unidad de movimiento horizontal o vertical no afecta el espaciado actual del lado derecho.
//   - El comando GS P puede cambiar la unidad de movimiento horizontal (y vertical). Sin embargo, el valor no puede ser menor que la cantidad mínima de movimiento horizontal y debe estar en unidades pares de la cantidad mínima de movimiento horizontal.
//   - En modo estándar, se utiliza la unidad de movimiento horizontal.
//   - En modo página, la unidad de movimiento horizontal o vertical depende de la posición inicial del área imprimible configurada mediante ESC T:
//     1. Cuando la posición inicial se configura en la esquina superior izquierda o inferior derecha del área imprimible utilizando ESC T, se utiliza la unidad de movimiento horizontal (x).
//     2. Cuando la posición inicial se configura en la esquina superior derecha o inferior izquierda del área imprimible utilizando ESC T, se utiliza la unidad de movimiento vertical (y).
//   - El espaciado máximo del lado derecho es de 255/180 pulgadas. Cualquier configuración que exceda el máximo se convierte automáticamente al valor máximo.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	GS P
func SetRightSideCharacterSpacing(n byte) []byte {
	return []byte{ESC, SP, n}
}

// SelectPrintModes representa el comando para seleccionar modos de impresión.
//
// Nombre:
//
//	Seleccionar modos de impresión
//
// Formato:
//
//	ASCII: ESC ! n
//	Hex: 1B 21 n
//	Decimal: 27 33 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona modos de impresión utilizando n de la siguiente manera:
//
//	Bit Off/On | Hex | Decimal | Función
//	------------------------------------
//	0 Off      | 00  | 0       | Fuente de carácter A (12 × 24).
//	  On       | 01  | 1       | Fuente de carácter B (9 × 17).
//	1          | -   | -       | No definido.
//	2          | -   | -       | No definido.
//	3 Off      | 00  | 0       | Modo enfatizado no seleccionado.
//	  On       | 08  | 8       | Modo enfatizado seleccionado.
//	4 Off      | 00  | 0       | Modo de doble altura no seleccionado.
//	  On       | 10  | 16      | Modo de doble altura seleccionado.
//	5 Off      | 00  | 0       | Modo de doble ancho no seleccionado.
//	  On       | 20  | 32      | Modo de doble ancho seleccionado.
//	6          | -   | -       | No definido.
//	7 Off      | 00  | 0       | Modo subrayado no seleccionado.
//	  On       | 80  | 128     | Modo subrayado seleccionado.
//
// Detalles:
//   - Cuando se seleccionan ambos modos, doble altura y doble ancho, se imprimen caracteres de tamaño cuádruple.
//   - La impresora puede subrayar todos los caracteres, pero no puede subrayar espacios establecidos por HT o caracteres rotados 90° en sentido horario.
//   - El grosor del subrayado es el seleccionado por ESC , independientemente del tamaño del carácter.
//   - Cuando algunos caracteres en una línea tienen doble altura o más, todos los caracteres en la línea se alinean en la línea base.
//   - ESC E también puede activar o desactivar el modo enfatizado. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - ESC —también puede activar o desactivar el modo subrayado. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - GS ! también puede seleccionar el tamaño de los caracteres. Sin embargo, la configuración del último comando recibido es la efectiva.
//   - El modo enfatizado es efectivo para caracteres alfanuméricos y Kanji. Todos los modos de impresión, excepto el modo enfatizado, son efectivos solo para caracteres alfanuméricos.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC -, ESC E, GS !
func SelectPrintModes(n byte) []byte {
	return []byte{ESC, '!', n}
}

// SetAbsolutePrintPosition representa el comando para establecer la posición de impresión absoluta.
//
// Nombre:
//
//	Establecer posición de impresión absoluta
//
// Formato:
//
//	ASCII: ESC $ nL nH
//	Hex:   1B 24 nL nH
//	Decimal: 27 36 nL nH
//
// Rango:
//
//	0 ≤ nL ≤ 255
//	0 ≤ nH ≤ 255
//
// Descripción:
//
//	Establece la distancia desde el comienzo de la línea hasta la posición en la que se imprimirán los caracteres posteriores.
//
// Detalles:
//   - La distancia se calcula como: (nL + nH × 256) × (unidad de movimiento vertical u horizontal) en pulgadas.
//   - Las configuraciones fuera del área imprimible especificada se ignoran.
//   - La unidad de movimiento horizontal y vertical se especifica mediante el comando GS P. Este comando puede modificar dichas unidades,
//     aunque el valor no puede ser menor que la cantidad mínima de movimiento horizontal y debe ser un múltiplo par de dicha cantidad mínima.
//   - En modo estándar se utiliza la unidad de movimiento horizontal (x).
//   - En modo página, la unidad de movimiento puede ser horizontal o vertical, dependiendo de la posición inicial del área imprimible configurada con ESC T:
//     1. Si la posición inicial se establece en la esquina superior izquierda o inferior derecha, se usa la unidad horizontal (x).
//     2. Si la posición inicial se establece en la esquina superior derecha o inferior izquierda, se usa la unidad vertical (y).
//
// Referencia:
//
//	ESC \, GS $, GS P
func SetAbsolutePrintPosition(nL, nH byte) []byte {
	return []byte{ESC, '$', nL, nH}
}

// SelectOrCancelUserDefinedCharset representa el comando para seleccionar o cancelar el conjunto de caracteres definido por el usuario.
//
// Nombre:
//
//	Seleccionar o cancelar conjunto de caracteres definido por el usuario
//
// Formato:
//
//	ASCII: ESC % n
//	Hex:   1B 25 n
//	Decimal: 27 37 n
//
// Rango:
//
//	0 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona o cancela el conjunto de caracteres definido por el usuario.
//
// Detalles:
//   - Cuando el bit menos significativo (LSB) de n es 0, se cancela el conjunto de caracteres definido por el usuario.
//   - Cuando el bit menos significativo de n es 1, se selecciona el conjunto de caracteres definido por el usuario.
//   - Al cancelar el conjunto definido por el usuario, se selecciona automáticamente el conjunto de caracteres interno.
//
// Valor por Defecto:
//
//	n = 0
//
// Referencia:
//
//	ESC &, ESC ?
func SelectOrCancelUserDefinedCharset(n byte) []byte {
	return []byte{ESC, '%', n}
}
