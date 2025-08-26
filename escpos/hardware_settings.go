package escpos

import (
	"fmt"
)

// TODO: Comandos para configuración de hardware específico
// - Densidad de impresión
// - Velocidad de impresión
// - Modo de ahorro de energía
// - Control de periféricos

var printerEnaMap = map[PrinterEnabled]byte{
	EnaOff: 0,
	EnaOn:  1,
}

// SetPeripheralDevice configura el estado del dispositivo periférico (como un lector de tarjetas o un escáner de códigos de barras).
func SetPeripheralDevice(n PrinterEnabled) ([]byte, error) {
	enabled, ok := printerEnaMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid printer enabled value: %v", n)
	}
	cmd := []byte{ESC, '=', enabled}
	return cmd, nil
}
