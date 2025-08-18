package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
)

// TODO: Comandos para impresión de códigos de barras
// - HRI (Human Readable Interpretation)

var textPosBarcodeMap = map[types.TextPositionBarcode]byte{
	types.NonePosBarcode:  0,
	types.AbovePosBarcode: 1,
	types.BelowPosBarcode: 2,
	types.BothPosBarcode:  3,
}

var barcodeWidthMap = map[types.BarcodeWidth]byte{
	types.ExtraSmallWidth: 2,
	types.SmallWidth:      3,
	types.MediumWidth:     4,
	types.LargeWidth:      5,
	types.ExtraLargeWidth: 6,
}

func (p *Commands) SetBarcodeHeight(height types.BarcodeHeight) ([]byte, error) {
	if height == 0 {
		return nil, fmt.Errorf("barcode height cannot be zero")
	}

	return []byte{GS, 'h', byte(height)}, nil
}

func (p *Commands) SetBarcodeWidth(width types.BarcodeWidth) ([]byte, error) {
	bcWidth, ok := barcodeWidthMap[width]
	if !ok {
		return nil, fmt.Errorf("no barcode width found for width %v", width)
	}

	return []byte{GS, 'w', bcWidth}, nil
}

func (p *Commands) SelectBarcodeTextPosition(position types.TextPositionBarcode) ([]byte, error) {
	pos, ok := textPosBarcodeMap[position]
	if !ok {
		return nil, fmt.Errorf("unknown position: %d", position)
	}
	return []byte{GS, 'H', pos}, nil
}

func (p *Commands) Barcode(content string, barType types.BarcodeType) ([]byte, error) {

	return []byte{}, nil
}

func (p *Commands) SelectFontBarcode(font types.Font) ([]byte, error) {
	bcFont, ok := fontMap[font]
	if !ok {
		return nil, fmt.Errorf("no barcode font found for font %v", font)
	}

	return []byte{GS, 'f', bcFont}, nil
}
