package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// TODO: Comandos para obtener estado de la impresora
// - Autodiagnóstico

var realTimeStatusMap = map[RealTimeStatus]byte{
	PrinterStatus:     1,
	OfflineStatus:     2,
	ErrorStatus:       3,
	PaperSensorStatus: 4,
}

// TransmitRealTimeStatus asks the printer to transmit its real-time status
func (c *Commands) TransmitRealTimeStatus(n RealTimeStatus) ([]byte, error) {
	status, ok := realTimeStatusMap[n]
	if !ok {
		return nil, fmt.Errorf("estado en tiempo real inválido: %d", n)
	}
	cmd := []byte{common.DLE, common.EOT, status}
	return cmd, nil
}
