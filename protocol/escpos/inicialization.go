package escpos

// TODO: Comandos para inicialización y configuración básica de la impresora
// - Configuración de página de códigos
// - Configuración regional
// - Reinicio de impresora
// - Selección de modo

// InitializePrinter representa el comando ESC @ para inicializar la impresora.
//
// Nombre:
//
//	Inicializar impresora
//
// Formato:
//
//	ASCII: ESC @
//	Hex:   1B 40
//	Decimal: 27 64
//
// Descripción:
//
//	Borra los datos en el búfer de impresión y restablece el modo de la impresora al que estaba vigente al encenderla.
//
// Detalles:
//   - La configuración de los DIP switches no se vuelve a verificar.
//   - Los datos en el búfer de recepción no se borran.
//   - La definición de macros no se limpia.
//   - Los datos de la imagen de bits en la memoria NV no se borran.
//   - Los datos de la memoria NV del usuario no se borran.
func (p *Commands) InitializePrinter() []byte {
	// ESC @ - Reset printer
	return []byte{ESC, '@'}
}

// TODO: Pensar en comando compuesto para el final Feed(1) y Cut(1)

// Close genera comandos de cierre (si los hay)
func (p *Commands) Close() []byte {
	// ESC/POS no tiene un comando específico de cierre
	// pero podrías incluir un reset o feed final o ambos
	return []byte{}
}

// SelectStandardMode cambia la impresora del modo página al modo estándar.
//
// Formato:
//
//	ASCII: ESC S
//	Hex:   1B 53
//	Decimal: 27 83
//
// Descripción:
//
//	Activa el modo estándar en la impresora. Este comando es válido solo en el modo página.
//
// Detalles:
//   - Los datos almacenados en el búfer en modo página se eliminan al cambiar al modo estándar.
//   - La posición de impresión se establece al comienzo de la línea.
//   - El área de impresión configurada por ESC W se inicializa.
//   - Los siguientes comandos se configuran con valores para el modo estándar, donde los valores pueden establecerse de forma independiente para modo estándar y modo página:
//     ① Configurar el espaciado de caracteres hacia la derecha: ESC SP, FS S
//     ② Seleccionar el espaciado de línea predeterminado: ESC 2, ESC 3
//   - Los siguientes comandos están habilitados solo para configurar en modo estándar:
//     ① Configurar el área de impresión en modo página: ESC W
//     ② Seleccionar la dirección de impresión en modo página: ESC T
//   - Los siguientes comandos se ignoran en modo estándar:
//     ① Configurar posición de impresión vertical absoluta en modo página: GS $
//     ② Configurar posición de impresión vertical relativa en modo página: GS \
//   - El modo estándar se selecciona automáticamente cuando se enciende la impresora, se reinicia o se utiliza el comando ESC @.
//
// Referencia:
//
//	FF, ESC FF, ESC L
func SelectStandardMode() []byte {
	return []byte{ESC, 'S'}
}

// SelectPageMode cambia la impresora del modo estándar al modo página.
//
// Formato:
//
//	ASCII: ESC L
//	Hex:   1B 4C
//	Decimal: 27 76
//
// Descripción:
//
//	Activa el modo página en la impresora. Este comando es válido solo cuando se procesa al comienzo de una línea en modo estándar
//	y no tiene efecto si ya se encuentra en modo página.
//
// Detalles:
//   - Después de completar la impresión utilizando el comando FF o ESC S, la impresora vuelve al modo estándar.
//   - Este comando establece la posición donde se almacenan los datos en el búfer según la posición especificada por ESC T dentro del área de impresión definida por ESC W.
//   - Los siguientes comandos se configuran con valores para el modo página, donde se pueden establecer valores de forma independiente en modo estándar y modo página:
//     ① Configurar el espaciado de caracteres hacia la derecha: ESC SP, FS S
//     ② Seleccionar el espaciado de línea predeterminado: ESC 2, ESC 3
//   - En modo página, solo es posible configurar valores para los siguientes comandos; estos comandos no se ejecutan:
//     ① Activar/desactivar el modo de rotación 90° en sentido horario: ESC V
//     ② Seleccionar justificación: ESC a
//     ③ Activar/desactivar el modo de impresión invertida: ESC {
//     ④ Configurar margen izquierdo: GS L
//     ⑤ Configurar el ancho del área imprimible: GS W
//   - El siguiente comando se ignora en modo página:
//     ① Ejecutar impresión de prueba: GS ( A
//   - Los siguientes comandos no están disponibles en modo página:
//     ① Imprimir imagen NV: FS p
//     ② Definir imagen NV: FS q
//     ③ Escribir en la memoria NV del usuario: FS g 1
//     ④ Imprimir imagen de bits rasterizada: GS v 0
//   - La impresora vuelve al modo estándar cuando se enciende, se reinicia o se utiliza el comando ESC @.
//
// Referencia:
//
//	FF, CAN, ESC FF, ESC S, ESC T, ESC W, GS $, GS \
func SelectPageMode() []byte {
	return []byte{ESC, 'L'}
}
