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
