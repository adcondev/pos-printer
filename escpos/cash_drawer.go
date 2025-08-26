package escpos

import (
	"fmt"
)

// TODO: Comandos para control del cajón de efectivo
// - Apertura de cajón
// - Estado de cajón

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

// Pulse envía un pulso al pin especificado del cajón de efectivo.
func (c *Commands) Pulse(_ int, _ int, _ int) []byte {
	// TODO: Implementar ESC p m t1 t2
	return []byte{}
}

// GenerateRealTimePulse genera el comando para enviar un pulso al pin especificado del cajón de efectivo.
func GenerateRealTimePulse(m CashDrawerPin, t CashDrawerTimePulse) ([]byte, error) {
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
