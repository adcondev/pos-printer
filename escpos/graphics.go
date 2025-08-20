package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/utils"
)

// TODO: Comandos para impresión de gráficos e imágenes
// - Modos de imagen
// - Compresión de imagen

var densityMap = map[Density]byte{
	DensitySingle:    0, // Modo normal (200 DPI vertical y horizontal)
	DensityDouble:    1, // Modo de doble ancho (200 DPI vertical y 100 DPI horizontal)
	DensityTriple:    2, // Modo de doble altura (100 DPI vertical y 200 DPI horizontal)
	DensityQuadruple: 3, // Modo cuádruple (100 DPI vertical y horizontal)
}

var bitImageMap = map[BitImageMode]byte{
	Mode8DotSingleDen:  0,
	Mode8DotDoubleDen:  1,
	Mode24DotSingleDen: 32,
	Mode24DotDoubleDen: 33,
}

// ESCImage ahora es más simple, solo guarda referencia a PrintRasterBitImage
type ESCImage struct {
	printImage *imaging.PrintImage

	// Cache de datos procesados
	rasterData []byte
}

// newESCImageFromPrintImage crea una ESCImage desde PrintRasterBitImage
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

// GetWidth devuelve el ancho en píxeles
func (e *ESCImage) GetWidth() int {
	return e.printImage.Width
}

// GetHeight devuelve el alto en píxeles
func (e *ESCImage) GetHeight() int {
	return e.printImage.Height
}

// GetWidthBytes devuelve el ancho en bytes
func (e *ESCImage) GetWidthBytes() int {
	return (e.printImage.Width + 7) / 8
}

// toRasterFormat convierte la imagen al formato raster de ESC/POS
func (e *ESCImage) toRasterFormat() ([]byte, error) {
	// Si ya tenemos los datos en cache, devolverlos
	if e.rasterData != nil {
		return e.rasterData, nil
	}

	// Obtener datos monocromáticos de la imagen
	// PrintRasterBitImage se encarga de aplicar dithering si fue configurado
	e.rasterData = e.printImage.ToMonochrome()

	return e.rasterData, nil
}

func (p *Commands) PrintRasterBitImage(img *imaging.PrintImage, density Density) ([]byte, error) {
	// Crear ESCImage
	escImg, err := newESCImageFromPrintImage(img)
	if err != nil {
		return nil, err
	}

	// Obtener datos rasterizados
	rasterData, err := escImg.toRasterFormat()
	if err != nil {
		return nil, err
	}

	mode, ok := densityMap[density]
	if !ok {
		return nil, fmt.Errorf("densidad no soportada: %v", density)
	}

	// Construir comando GS v 0
	cmd := []byte{GS, 'v', '0', mode}

	// Agregar dimensiones
	wL, wH, err := utils.LengthLowHigh(escImg.GetWidthBytes())
	if err != nil {
		return nil, err
	}
	hL, hH, err := utils.LengthLowHigh(escImg.GetHeight())
	if err != nil {
		return nil, err
	}

	cmd = append(cmd, wL, wH) // Ancho en bytes
	cmd = append(cmd, hL, hH) // Alto en píxeles
	cmd = append(cmd, rasterData...)

	return cmd, nil
}

// GetMaxImageWidth devuelve el ancho máximo de imagen que soporta la impresora
func (p *Commands) GetMaxImageWidth(paperWidth, dpi int) int {
	// Cálculo basado en el ancho del papel y resolución
	// Formula: (ancho_papel_mm / 25.4) * dpi
	if paperWidth > 0 && dpi > 0 {
		return int((float64(paperWidth) / 25.4) * float64(dpi))
	}

	// Valores predeterminados si no hay configuración
	if paperWidth >= 80 {
		return 576 // Para papel de 80mm a 203dpi
	} else {
		return 384 // Para papel de 58mm a 203dpi
	}
}

func SelectBitImageMode(m BitImageMode, nL, nH byte, data []byte) ([]byte, error) {
	mode, ok := bitImageMap[m]
	if !ok {
		return nil, fmt.Errorf("invalid bit image mode: %v", m)
	}
	cmd := []byte{ESC, '*', mode, nL, nH}
	return append(cmd, data...), nil
}
