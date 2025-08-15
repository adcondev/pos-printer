package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para obtener estado de la impresora
// - Autodiagnóstico

var realTimeStatusMap = map[types.RealTimeStatus]byte{
	types.PrinterStatus:     1,
	types.OfflineStatus:     2,
	types.ErrorStatus:       3,
	types.PaperSensorStatus: 4,
}

// TransmitRealTimeStatus representa el comando para transmitir el estado de la impresora en tiempo real.
//
// Nombre:
//
//	Transmisión de estado en tiempo real
//
// Formato:
//
//	ASCII: DLE EOT n
//	Hex:   10 04 n
//	Decimal: 16 4 n
//
// Rango:
//
//	1 ≤ n ≤ 4
//
// Descripción:
//
//	Transmite en tiempo real el estado seleccionado de la impresora según el valor de n, de la siguiente manera:
//	  n = 1: Transmitir estado de la impresora
//	  n = 2: Transmitir estado fuera de línea
//	  n = 3: Transmitir estado de error
//	  n = 4: Transmitir estado del sensor del papel continuo
//
// Detalles:
//   - El estado se transmite cada vez que se recibe la secuencia de datos <10>H<04>H<n>, con 1 ≤ n ≤ 4.
//     Ejemplo: En "ESC * m nL nH d1...dk", se tiene: d1 = <10>H, d2 = <04>H, d3 = <01>H.
//   - Este comando no debe ser utilizado dentro de la secuencia de datos de otro comando que consista en 2 o más bytes.
//     Ejemplo: Si se intenta transmitir "ESC 3 n" a la impresora, pero DTR (o DSR en la computadora host) cambia a MARK antes de que se transmita n,
//     y luego DLE EOT 3 interrumpe antes de que se reciba n, el código <10>H de DLE EOT 3 se procesa como el código para "ESC 3 <10>H".
//   - Este comando es efectivo incluso si la impresora no ha sido seleccionada mediante ESC = (seleccionar dispositivo periférico).
//   - La impresora transmite el estado actual, donde cada estado se representa con un dato de un byte.
//   - La transmisión se realiza sin confirmar si el host es capaz de recibir los datos.
//   - La impresora ejecuta el comando en cuanto lo recibe, incluso si está fuera de línea, el búfer de recepción está lleno o hay un estado de error
//     en modelos de interfaz serial.
//   - En modelos de interfaz paralela, este comando no se ejecuta cuando la impresora está ocupada, pero sí se ejecuta si está fuera de línea
//     o hay error cuando el DIP switch 2-1 está activado.
//   - Cuando se habilita Auto Status Back (ASB) mediante el comando GS a, se debe distinguir entre el estado transmitido por DLE EOT y el estado ASB.
//
// Parámetros de estado según el valor de n:
//
//	n = 1: Estado de la impresora
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	        |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   0    | 00   |    0    | Señal del cajón (pin 3 del conector) en nivel LOW.
//	        |   1    | 04   |    4    | Señal del cajón (pin 3 del conector) en nivel LOW.
//	   3    |   0    | 00   |    0    | En línea.
//	        |   1    | 08   |    8    | Fuera de línea.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5,6  |   -    |  -   |    -    | Indefinido.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 2: Estado fuera de línea
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   1    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   0    | 00   |    0    | La tapa está cerrada.
//	        |   1    | 04   |    4    | La tapa está abierta.
//	   3    |   0    | 00   |    0    | Papel no se alimenta mediante el botón FEED.
//	        |   1    | 08   |    8    | Papel se alimenta mediante el botón FEED.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5    |   0    | 00   |    0    | No hay tope final de papel.
//	        |   1    | 20   |   32    | Se está deteniendo la impresión.
//	   6    |   0    | 00   |    0    | No hay error.
//	        |   1    | 40   |   64    | Se produce un error.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 3: Estado de error
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   1    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2    |   -    |  -   |    -    | Indefinido.
//	   3    |   0    | 00   |    0    | No hay error en la autocortadora.
//	        |   1    | 08   |    8    | Error en la autocortadora.
//	   4    |   1    | 10   |   16    | No utilizado. Fijado en On.
//	   5    |   0    | 00   |    0    | No hay error irrecuperable.
//	        |   1    | 20   |   32    | Se produce un error irrecuperable.
//	   6    |   0    | 00   |    0    | No hay error auto-recuperable.
//	        |   1    | 40   |   64    | Se produce un error auto-recuperable.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
//	n = 4: Estado del sensor de papel continuo
//	  Bit   | Off/On | Hex  | Decimal | Función
//	  ----- | ------ | ---- | ------- | -------------------------------------------------------------
//	   0    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//	   0    |   1    | 02   |    2    | No utilizado. Fijado en On.
//	   2,3  |   0    | 00   |    0    | Sensor de fin de papel: papel adecuado.
//	        |   1    |   -  |    -    | Fin de papel detectado (sensor activado).
//	   5,6  |   0C   | 10   |   16    | Se detecta fin de papel mediante el sensor.
//	   4    |   0    | 00   |    0    | No utilizado. Fijado en On.
//	        |   0    | 00   |    0    | Sensor de papel: papel presente.
//	   1    |   1    | 60   |   96    | Sensor de fin de papel: fin de papel detectado.
//	   7    |   0    | 00   |    0    | No utilizado. Fijado en Off.
//
// Referencia:
//
//	DLE ENQ, GS a, GS r
func (p *Commands) TransmitRealTimeStatus(n types.RealTimeStatus) ([]byte, error) {
	status, ok := realTimeStatusMap[n]
	if !ok {
		return nil, fmt.Errorf("estado en tiempo real inválido: %d", n)
	}
	cmd := []byte{DLE, EOT, status}
	return cmd, nil
}
