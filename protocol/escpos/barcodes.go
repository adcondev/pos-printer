package escpos

import (
	"github.com/AdConDev/pos-printer/protocol/escpos/types"
)

// TODO: Comandos para impresión de códigos de barras
// - HRI (Human Readable Interpretation)

func (p *Commands) SetBarcodeHeight(height int) []byte {
	// TODO: Implementar GS h n
	return []byte{}
}

func (p *Commands) SetBarcodeWidth(width int) []byte {
	// TODO: Implementar GS w n
	return []byte{}
}

func (p *Commands) SetBarcodeTextPosition(position types.BarcodeTextPosition) []byte {
	// TODO: Mapear position a valores ESC/POS y usar GS H n
	return []byte{}
}

func (p *Commands) Barcode(content string, barType types.BarcodeType) ([]byte, error) {
	// TODO: Esta es la más compleja, necesitas:
	// 1. Mapear barType genérico a tipo ESC/POS
	// 2. Validar content según el tipo
	// 3. Generar comando según p.capabilities["barcode_b"]
	return []byte{}, nil
}
