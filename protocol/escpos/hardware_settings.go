package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para configuración de hardware específico
// - Densidad de impresión
// - Velocidad de impresión
// - Modo de ahorro de energía
// - Control de periféricos

var printerEnaMap = map[types.PrinterEnabled]byte{
	types.EnaOff: 0,
	types.EnaOn:  1,
}

// SetPeripheralDevice representa el comando ESC = n para seleccionar el dispositivo periférico.
//
// Nombre:
//
//	Seleccionar dispositivo periférico
//
// Formato:
//
//	ASCII: ESC = n
//	Hex:   1B 3D n
//	Decimal: 27 61 n
//
// Rango:
//
//	1 ≤ n ≤ 255
//
// Descripción:
//
//	Selecciona el dispositivo al que la computadora host envía datos, utilizando el valor de n de la siguiente manera:
//	  - Si el bit correspondiente está apagado (n = 0x00), la impresora se desactiva (Printer disabled).
//	  - Si el bit correspondiente está encendido (n = 0x01), la impresora se activa (Printer enabled).
//	  - Los valores de n de 1 a 7 están indefinidos.
//
// Detalles:
//   - Cuando la impresora está deshabilitada, ignora todos los datos entrantes, excepto los comandos de recuperación de errores
//     (DLE EOT, DLE ENQ, DLE DC4), hasta que se active nuevamente utilizando este comando.
//
// Valor por Defecto:
//
//	n = 1
func SetPeripheralDevice(n types.PrinterEnabled) ([]byte, error) {
	enabled, ok := printerEnaMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid printer enabled value: %v", n)
	}
	cmd := []byte{ESC, '=', enabled}
	return cmd, nil
}
