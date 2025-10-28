package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/shared"
	"github.com/adcondev/pos-printer/imaging"
)

// ============================================================================
// Type / Constant maps and helpers
// ============================================================================

// densityMap mapea las densidades de impresión a sus valores ESC/POS correspondientes.
var densityMap = map[Density]byte{
	DensitySingle:    0, // Modo normal (200 DPI vertical y horizontal)
	DensityDouble:    1, // Modo de doble ancho (200 DPI vertical y 100 DPI horizontal)
	DensityTriple:    2, // Modo de doble altura (100 DPI vertical y 200 DPI horizontal)
	DensityQuadruple: 3, // Modo cuádruple (100 DPI vertical y horizontal)
}

// ESCImage encapsula una imagen preparada para ESC/POS.
type ESCImage struct {
	printImage *imaging.PrintImage

	// Cache de datos procesados
	rasterData []byte
}

// newESCImageFromPrintImage crea una ESCImage desde PrintImage
func newESCImageFromPrintImage(img *imaging.PrintImage) (*ESCImage, error) {
	if img == nil {
		return nil, fmt.Errorf("print imaging cannot be nil")
	}

	if img.Width <= 0 || img.Height <= 0 {
		return nil, fmt.Errorf("invalid imaging dimensions: %dx%d", img.Width, img.Height)
	}

	return &ESCImage{
		printImage: img,
	}, nil
}

// GetHeight devuelve el alto en píxeles
func (e *ESCImage) GetHeight() int {
	return e.printImage.Height
}

// GetWidthBytes devuelve el ancho en bytes (cada byte = 8 píxeles)
func (e *ESCImage) GetWidthBytes() int {
	return (e.printImage.Width + 7) / 8
}

// toRasterFormat convierte la imagen al formato raster de ESC/POS
func (e *ESCImage) toRasterFormat() []byte {
	// Si ya tenemos los datos en cache, devolverlos
	if e.rasterData != nil {
		return e.rasterData
	}

	// Obtener datos monocromáticos de la imagen
	// PrintRasterBitImage se encarga de aplicar dithering si fue configurado
	e.rasterData = e.printImage.ToMonochrome()

	return e.rasterData
}

// ============================================================================
// Public API (implementation)
// ============================================================================
// Funciones que generan comandos ESC/POS para imágenes y modos de bits.

// PrintRasterBitImage genera los comandos para imprimir una imagen rasterizada
func (c *Commands) PrintRasterBitImage(img *imaging.PrintImage, density Density) ([]byte, error) {
	// Crear ESCImage
	escImg, err := newESCImageFromPrintImage(img)
	if err != nil {
		return nil, err
	}

	// Obtener datos rasterizados
	rasterData := escImg.toRasterFormat()

	mode, ok := densityMap[density]
	if !ok {
		return nil, fmt.Errorf("densidad no soportada: %v", density)
	}

	// Construir comando GS v 0
	cmd := []byte{shared.GS, 'v', '0', mode}

	// Agregar dimensiones
	var widthBytes uint16
	switch {
	case escImg.GetWidthBytes() > 0xFFFF:
		widthBytes = 0xFFFF
	case escImg.GetWidthBytes() < 0x0:
		widthBytes = 0x0
	default:
		// Secure, it has been validated
		widthBytes = uint16(escImg.GetWidthBytes()) // nolint:gosec
	}
	wL, wH := shared.ToLittleEndian(widthBytes)

	var heightBytes uint16
	switch {
	case escImg.GetHeight() > 0xFFFF:
		heightBytes = 0xFFFF
	case escImg.GetHeight() < 0:
		heightBytes = 0
	default:
		// Secure, it has been validated
		heightBytes = uint16(escImg.GetHeight()) // nolint:gosec
	}
	hL, hH := shared.ToLittleEndian(heightBytes)

	cmd = append(cmd, wL, wH) // Ancho en bytes
	cmd = append(cmd, hL, hH) // Alto en píxeles
	cmd = append(cmd, rasterData...)

	return cmd, nil
}
