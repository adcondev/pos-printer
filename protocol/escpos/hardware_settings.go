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

func SetPeripheralDevice(n types.PrinterEnabled) ([]byte, error) {
	enabled, ok := printerEnaMap[n]
	if !ok {
		return nil, fmt.Errorf("invalid printer enabled value: %v", n)
	}
	cmd := []byte{ESC, '=', enabled}
	return cmd, nil
}
