package escpos

import (
	"fmt"

	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/types"
	"github.com/AdConDev/pos-printer/utils"
)

// TODO: Comandos para impresión de gráficos e imágenes
// - Modos de imagen
// - Compresión de imagen

var densityMap = map[types.Density]byte{
	types.DensitySingle:    0, // Modo normal (200 DPI vertical y horizontal)
	types.DensityDouble:    1, // Modo de doble ancho (200 DPI vertical y 100 DPI horizontal)
	types.DensityTriple:    2, // Modo de doble altura (100 DPI vertical y 200 DPI horizontal)
	types.DensityQuadruple: 3, // Modo cuádruple (100 DPI vertical y horizontal)
}

var bitImageMap = map[types.BitImageMode]byte{
	types.Mode8DotSingleDen:  0,
	types.Mode8DotDoubleDen:  1,
	types.Mode24DotSingleDen: 32,
	types.Mode24DotDoubleDen: 33,
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

// PrintRasterBitImage representa el comando para imprimir una imagen de bits en modo raster.
//
// Nombre:
//
//	Imprimir imagen de bits en modo raster
//
// Formato:
//
//	ASCII: GS v 0 m xL xH yL yH d1...dk
//	Hex: 1D 76 30 m xL xH yL yH d1...dk
//	Decimal: 29 118 48 m xL xH yL yH d1...dk
//
// Rango:
//   - m: 0 ≤ m ≤ 3, 48 ≤ m ≤ 51
//   - xL: 0 ≤ xL ≤ 255
//   - xH: 0 ≤ xH ≤ 255
//   - yL: 0 ≤ yL ≤ 255
//   - d: 0 ≤ d ≤ 255
//   - k = (xL + xH × 256) × (yL + yH × 256) (k ≥ 0)
//
// Descripción:
//
//	Selecciona el modo de imagen de bits raster. El valor de m selecciona el modo, como se indica a continuación:
//	- m = 0, 48: Modo normal (200 DPI vertical y horizontal).
//	- m = 1, 49: Modo de doble ancho (200 DPI vertical y 100 DPI horizontal).
//	- m = 2, 50: Modo de doble altura (100 DPI vertical y 200 DPI horizontal).
//	- m = 3, 51: Modo cuádruple (100 DPI vertical y horizontal).
//
// Detalles:
//   - En modo estándar, este comando es efectivo solo cuando no hay datos en el búfer de impresión.
//   - Este comando no tiene efecto en todos los modos de impresión (tamaño de caracteres, enfatizado, doble impacto, invertido, subrayado, impresión en blanco/negro, etc.) para imágenes de bits raster.
//   - Si el ancho del área de impresión configurado con GS L y GS W es menor que el ancho mínimo, el área de impresión se extiende al ancho mínimo solo en la línea en cuestión. El ancho mínimo es:
//   - 1 punto en modos normal (m=0, 48) y doble altura (m=2, 50).
//   - 2 puntos en modos doble ancho (m=1, 49) y cuádruple (m=3, 51).
//   - Los datos fuera del área de impresión se leen y se descartan punto por punto.
//   - La posición para imprimir caracteres posteriores en imágenes de bits raster se especifica mediante:
//   - HT (Tabulación Horizontal).
//   - ESC $ (Establecer posición de impresión absoluta).
//   - ESC \ (Establecer posición de impresión relativa).
//   - GS L (Establecer margen izquierdo).
//   - Si la posición para imprimir caracteres posteriores no es múltiplo de 8, la velocidad de impresión puede disminuir.
//   - La configuración de ESC a (Seleccionar justificación) también es efectiva para imágenes de bits raster.
//   - Cuando este comando se recibe durante la definición de macro, la impresora termina la definición de macro y comienza a ejecutar este comando. La definición de este comando debe ser borrada.
//   - El valor de d indica los datos de imagen de bits. Configurar un bit en 1 imprime un punto, mientras que configurarlo en 0 no imprime un punto.
func (p *Commands) PrintRasterBitImage(img *imaging.PrintImage, density types.Density) ([]byte, error) {
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

// SelectBitImageMode representa el comando para seleccionar el modo de imagen de bits.
//
// Nombre:
//
//	Seleccionar modo de imagen de bits
//
// Formato:
//
//	ASCII: ESC * m nL nH d1...dk
//	Hex:   1B 2A m nL nH d1...dk
//	Decimal: 27 42 m nL nH d1...dk
//
// Rango:
//   - m: Puede ser 0, 1, 32 o 33.
//   - nL: 0 ≤ nL ≤ 255
//   - nH: 0 ≤ nH ≤ 3
//   - d (datos): 0 ≤ d ≤ 255
//
// Descripción:
//
//	Selecciona un modo de imagen de bits utilizando el parámetro m para determinar
//	el número de puntos especificados por nL y nH, de la siguiente manera:
//
//	Modos según m:
//	  * m = 0: Modo de 8 puntos, densidad simple (8-dot single-density).
//	    - Dirección vertical: 8 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 67 DPI.
//	  * m = 1: Modo de 8 puntos, doble densidad (8-dot double-density).
//	    - Dirección vertical: 8 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 67 DPI y 100 DPI (dependiente del contexto).
//	  * m = 32: Modo de 24 puntos, densidad simple (24-dot single-density).
//	    - Dirección vertical: 24 puntos.
//	    - Dirección horizontal: El número de puntos es nL + nH × 256.
//	    - Densidad: 200 DPI.
//	  * m = 33: Modo de 24 puntos, doble densidad (24-dot double-density).
//	    - Dirección vertical: 24 puntos.
//	    - Dirección horizontal: El número de puntos es (nL + nH × 256) × 3.
//	    - Densidad: 200 DPI.
//
// Detalles:
//   - Si el valor de m no se encuentra dentro de los rangos especificados, nL y los datos
//     siguientes se procesarán como datos normales.
//   - Los parámetros nL y nH indican el número de puntos de la imagen de bits en la dirección horizontal,
//     calculándose como nL + nH × 256.
//   - Si los datos de la imagen de bits exceden el número de puntos que se pueden imprimir en una línea,
//     los datos excedentes serán ignorados.
//   - d representa los datos de la imagen de bits; se debe establecer el bit correspondiente a 1 para imprimir un punto,
//     o a 0 para no imprimirlo.
//   - Si el ancho del área de impresión establecido por GS L y GS W es menor que el requerido por los datos enviados
//     con el comando ESC *, se realizará lo siguiente en la línea en cuestión (sin exceder el área máxima de impresión):
//     ① Se extiende hacia la derecha el ancho del área de impresión para acomodar los datos.
//     ② Si el paso ① no proporciona el ancho suficiente, se reduce el margen izquierdo para acomodar los datos.
//   - Después de imprimir la imagen de bits, la impresora regresa al modo de procesamiento de datos normal.
//   - Este comando no se ve afectado por los modos de impresión (enfatizado, doble impacto, subrayado,
//     tamaño de caracteres o impresión en blanco/negro), excepto en el modo de impresión invertida (upside-down).
//
// Referencia:
//
//	GS L, GS W, ESC \, GS P
func SelectBitImageMode(m types.BitImageMode, nL, nH byte, data []byte) ([]byte, error) {
	mode, ok := bitImageMap[m]
	if !ok {
		return nil, fmt.Errorf("invalid bit image mode: %v", m)
	}
	cmd := []byte{ESC, '*', mode, nL, nH}
	return append(cmd, data...), nil
}
