package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Maps and helpers
// ============================================================================
// Mapas que definen códigos de estado en tiempo real para ESC/POS.

var realTimeStatusMap = map[RealTimeStatus]byte{
	PrinterStatus:     1,
	OfflineStatus:     2,
	ErrorStatus:       3,
	PaperSensorStatus: 4,
}

// ============================================================================
// Public API (implementation)
// ============================================================================

// TransmitRealTimeStatus pide a la impresora transmitir su estado en tiempo real
func (c *Protocol) TransmitRealTimeStatus(n RealTimeStatus) ([]byte, error) {
	status, ok := realTimeStatusMap[n]
	if !ok {
		return nil, fmt.Errorf("estado en tiempo real inválido: %d", n)
	}
	cmd := []byte{shared.DLE, shared.EOT, status}
	return cmd, nil
}
