package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Type / Constant mappings
// ============================================================================
// Mapas que relacionan tipos del paquete `types` con valores ESC/POS.
// Usamos mapas para convertir enums a bytes de comando.

// cashPin mapea el pin del cajón de efectivo a su valor ESC/POS correspondiente.
var cashPin = map[CashDrawerPin]byte{
	Pin2: 0,
	Pin5: 1,
}

// timeMap mapea el tiempo del cajón de efectivo a su valor ESC/POS correspondiente.
var timeMap = map[CashDrawerTimePulse]byte{
	Pulse100ms: 1,
	Pulse200ms: 2,
	Pulse300ms: 3,
	Pulse400ms: 4,
	Pulse500ms: 5,
	Pulse600ms: 6,
	Pulse700ms: 7,
	Pulse800ms: 8,
}

// ============================================================================
// Public API (implementation)
// ============================================================================
// Funciones públicas que generan secuencias de bytes para interactuar con el
// cajón de efectivo.

// Pulse envía un pulso al pin especificado del cajón de efectivo.
func (c *Protocol) Pulse(_ int, _ int, _ int) []byte {
	// TODO: Implementar ESC p m t1 t2
	return []byte{}
}

// GenerateRealTimePulse genera el comando para enviar un pulso al pin especificado del cajón de efectivo.
//
// Parámetros:
//   - m: pin del cajón (Pin2 o Pin5)
//   - t: duración del pulso (Pulse100ms..Pulse800ms)
//
// Devuelve una secuencia DLE DC4 que el printer interpreta como pulso real-time.
func GenerateRealTimePulse(m CashDrawerPin, t CashDrawerTimePulse) ([]byte, error) {
	drawerPin, ok := cashPin[m]
	if !ok {
		// m no es un pin soportado según el mapa cashPin
		return nil, fmt.Errorf("pin de cajón de efectivo no soportado: %v", m)
	}
	pulseTime, ok := timeMap[t]
	if !ok {
		// t no es un tiempo soportado según el mapa timeMap
		return nil, fmt.Errorf("tiempo de pulso no soportado: %v", t)
	}
	// Comando DLE DC4 1 m t
	cmd := []byte{shared.DLE, shared.DC4, 1, drawerPin, pulseTime}
	return cmd, nil
}
