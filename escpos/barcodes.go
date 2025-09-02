package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// TODO: Comandos para impresión de códigos de barras
// - HRI (Human Readable Interpretation)

var textPosBarcodeMap = map[TextPositionBarcode]byte{
	NonePosBarcode:  0,
	AbovePosBarcode: 1,
	BelowPosBarcode: 2,
	BothPosBarcode:  3,
}

var barcodeWidthMap = map[BarcodeWidth]byte{
	ExtraSmallWidth: 2,
	SmallWidth:      3,
	MediumWidth:     4,
	LargeWidth:      5,
	ExtraLargeWidth: 6,
}

// SetBarcodeHeight establece la altura del código de barras en puntos
func (c *Protocol) SetBarcodeHeight(height BarcodeHeight) ([]byte, error) {
	if height == 0 {
		return nil, fmt.Errorf("barcode height cannot be zero")
	}

	return []byte{common.GS, 'h', byte(height)}, nil
}

// SetBarcodeWidth establece el ancho del código de barras
func (c *Protocol) SetBarcodeWidth(width BarcodeWidth) ([]byte, error) {
	bcWidth, ok := barcodeWidthMap[width]
	if !ok {
		return nil, fmt.Errorf("no barcode width found for width %v", width)
	}

	return []byte{common.GS, 'w', bcWidth}, nil
}

// SelectBarcodeTextPosition establece la posición del texto del código de barras
func (c *Protocol) SelectBarcodeTextPosition(position TextPositionBarcode) ([]byte, error) {
	pos, ok := textPosBarcodeMap[position]
	if !ok {
		return nil, fmt.Errorf("unknown position: %d", position)
	}
	return []byte{common.GS, 'H', pos}, nil
}

// Barcode imprime un código de barras con el contenido y tipo especificados
func (c *Protocol) Barcode(_ string, _ BarcodeType) ([]byte, error) {

	return []byte{}, nil
}

// SelectFontBarcode selecciona la fuente para el código de barras
func (c *Protocol) SelectFontBarcode(font Font) ([]byte, error) {
	bcFont, ok := fontMap[font]
	if !ok {
		return nil, fmt.Errorf("no barcode font found for font %v", font)
	}

	return []byte{common.GS, 'f', bcFont}, nil
}
