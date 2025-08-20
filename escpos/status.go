package escpos

import (
	"fmt"
)

// TODO: Comandos para obtener estado de la impresora
// - Autodiagnóstico

var realTimeStatusMap = map[RealTimeStatus]byte{
	PrinterStatus:     1,
	OfflineStatus:     2,
	ErrorStatus:       3,
	PaperSensorStatus: 4,
}

func (p *Commands) TransmitRealTimeStatus(n RealTimeStatus) ([]byte, error) {
	status, ok := realTimeStatusMap[n]
	if !ok {
		return nil, fmt.Errorf("estado en tiempo real inválido: %d", n)
	}
	cmd := []byte{DLE, EOT, status}
	return cmd, nil
}
