package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para control del cajón de efectivo
// - Apertura de cajón
// - Estado de cajón

// cashPin mapea el pin del cajón de efectivo a su valor ESC/POS correspondiente.
var cashPin = map[types.CashDrawerPin]byte{
	types.Pin2: 0,
	types.Pin5: 1,
}

// timeMap mapea el tiempo del cajón de efectivo a su valor ESC/POS correspondiente.
var timeMap = map[types.CashDrawerTimePulse]byte{
	types.Pulse100ms: 1,
	types.Pulse200ms: 2,
	types.Pulse300ms: 3,
	types.Pulse400ms: 4,
	types.Pulse500ms: 5,
	types.Pulse600ms: 6,
	types.Pulse700ms: 7,
	types.Pulse800ms: 8,
}

func (p *Commands) Pulse(pin int, onMS int, offMS int) []byte {
	// TODO: Implementar ESC p m t1 t2
	return []byte{}
}

// GenerateRealTimePulse representa el comando DLE DC4 n m t para generar un pulso en tiempo real.
//
// Nombre:
//
//	Generar pulso en tiempo real
//
// Formato:
//
//	ASCII: DLE DC4 n m t
//	Hex:   10 14 n m t
//	Decimal: 16 20 n m t
//
// Rango:
//   - n: n = 1
//   - m: m ∈ {0, 1}
//     m = 0: Pin 2 del conector de expulsión del cajón.
//     m = 1: Pin 5 del conector de expulsión del cajón.
//   - t: 1 ≤ t ≤ 8
//
// Descripción:
//
//	Genera el pulso especificado por t en el pin de conector indicado por m.
//	El tiempo de encendido del pulso es t × 100 ms y el tiempo de apagado es también t × 100 ms.
//
// Detalles:
//   - Si la impresora se encuentra en estado de error al procesar este comando, éste se ignora.
//   - Si el pulso se está enviando al pin especificado mientras se ejecuta ESC p o DEL DC4, el comando se ignora.
//   - La impresora ejecuta este comando en cuanto lo recibe.
//   - Con un modelo de interfaz serial, este comando se ejecuta incluso si la impresora está fuera de línea, el búfer de recepción está lleno o hay un estado de error.
//   - Con un modelo de interfaz paralela, el comando no se ejecuta cuando la impresora está ocupada; sin embargo, se ejecuta cuando está fuera de línea o hay error si el DIP switch 2-1 está activado.
//   - Si los datos de impresión incluyen cadenas de caracteres idénticas a este comando, la impresora realizará la misma operación especificada por este comando. Se debe tener en cuenta este comportamiento.
//   - Este comando no debe usarse dentro de la secuencia de datos de otro comando que consista en 2 o más bytes.
//   - Es efectivo incluso cuando la impresora está deshabilitada mediante ESC = (selección de dispositivo periférico).
//
// Referencia:
//
//	ESC p
func GenerateRealTimePulse(m types.CashDrawerPin, t types.CashDrawerTimePulse) ([]byte, error) {
	drawerPin, ok := cashPin[m]
	if !ok {
		return nil, fmt.Errorf("pin de cajón de efectivo no soportado: %v", m)
	}
	pulseTime, ok := timeMap[t]
	if !ok {
		return nil, fmt.Errorf("tiempo de pulso no soportado: %v", t)
	}
	cmd := []byte{DLE, DC4, 1, drawerPin, pulseTime}
	return cmd, nil
}
