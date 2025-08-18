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

func (p *Commands) TransmitRealTimeStatus(n types.RealTimeStatus) ([]byte, error) {
	status, ok := realTimeStatusMap[n]
	if !ok {
		return nil, fmt.Errorf("estado en tiempo real inválido: %d", n)
	}
	cmd := []byte{DLE, EOT, status}
	return cmd, nil
}
