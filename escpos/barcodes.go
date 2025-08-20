package escpos

import (
	"fmt"
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

func (p *Commands) SetBarcodeHeight(height BarcodeHeight) ([]byte, error) {
	if height == 0 {
		return nil, fmt.Errorf("barcode height cannot be zero")
	}

	return []byte{GS, 'h', byte(height)}, nil
}

func (p *Commands) SetBarcodeWidth(width BarcodeWidth) ([]byte, error) {
	bcWidth, ok := barcodeWidthMap[width]
	if !ok {
		return nil, fmt.Errorf("no barcode width found for width %v", width)
	}

	return []byte{GS, 'w', bcWidth}, nil
}

func (p *Commands) SelectBarcodeTextPosition(position TextPositionBarcode) ([]byte, error) {
	pos, ok := textPosBarcodeMap[position]
	if !ok {
		return nil, fmt.Errorf("unknown position: %d", position)
	}
	return []byte{GS, 'H', pos}, nil
}

func (p *Commands) Barcode(content string, barType BarcodeType) ([]byte, error) {

	return []byte{}, nil
}

func (p *Commands) SelectFontBarcode(font Font) ([]byte, error) {
	bcFont, ok := fontMap[font]
	if !ok {
		return nil, fmt.Errorf("no barcode font found for font %v", font)
	}

	return []byte{GS, 'f', bcFont}, nil
}
