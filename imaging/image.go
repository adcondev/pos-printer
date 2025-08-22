package imaging

import (
	"github.com/adcondev/pos-printer/utils"
	"golang.org/x/image/draw"

	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

// PrintImage representa una imagen preparada para impresión
type PrintImage struct {
	// La imagen original de Go
	Source image.Image

	// Dimensiones efectivas para impresión
	Width  int
	Height int

	// Metadatos opcionales
	DPI       int
	Threshold uint8

	// Datos pre-procesados opcionales
	MonochromeData []byte

	// Imagen procesada con dithering (si se aplicó)
	ProcessedImage image.Image
	DitherMode     DitherMode
}

// NewPrintImage crea una nueva imagen para impresión
func NewPrintImage(img image.Image, dither DitherMode) *PrintImage {
	bounds := img.Bounds()
	return &PrintImage{
		Source:     img,
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
		Threshold:  128,
		DitherMode: dither,
	}
}

// ApplyDithering aplica un algoritmo de dithering a la imagen
func (p *PrintImage) ApplyDithering(mode DitherMode) error {
	processed, err := ProcessImageWithDithering(p.Source, mode, p.Threshold)
	if err != nil {
		return err
	}

	p.ProcessedImage = processed
	p.DitherMode = mode

	// Invalidar datos monocromáticos anteriores
	p.MonochromeData = nil

	return nil
}

// GetEffectiveImage devuelve la imagen a usar (procesada o original)
func (p *PrintImage) GetEffectiveImage() image.Image {
	if p.ProcessedImage != nil {
		return p.ProcessedImage
	}
	return p.Source
}

// TODO: Revisar linter
// GetPixel obtiene el valor de un pixel como blanco (false) o negro (true)
func (p *PrintImage) GetPixel(x, y int) bool {
	// Si tenemos datos monocromáticos, usarlos
	if p.MonochromeData != nil {
		byteIndex := (y*p.Width + x) / 8
		bitIndex := 7 - uint(x%8) // #nosec G115
		return p.MonochromeData[byteIndex]&(1<<bitIndex) != 0
	}

	// Usar la imagen efectiva (procesada o original)
	img := p.GetEffectiveImage()

	// Si la imagen ya es en escala de grises (resultado de dithering)
	if grayImg, ok := img.(*image.Gray); ok {
		return grayImg.GrayAt(x, y).Y < p.Threshold
	}

	// Convertir a escala de grises
	gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
	return gray.Y < p.Threshold
}

// ToMonochrome convierte la imagen a datos monocromáticos
func (p *PrintImage) ToMonochrome() []byte {
	if p.MonochromeData != nil {
		return p.MonochromeData
	}

	// Calcular bytes necesarios
	bytesPerRow := (p.Width + 7) / 8
	data := make([]byte, bytesPerRow*p.Height)

	// Usar la imagen efectiva
	img := p.GetEffectiveImage()

	// Convertir pixel por pixel
	// TODO: Revisar linters
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			// Para imágenes ya procesadas con dithering
			if grayImg, ok := img.(*image.Gray); ok {
				if grayImg.GrayAt(x, y).Y < p.Threshold {
					byteIndex := y*bytesPerRow + x/8
					bitIndex := 7 - uint(x%8) // #nosec G115
					data[byteIndex] |= 1 << bitIndex
				}
			} else {
				// Para imágenes a color
				if p.GetPixel(x, y) {
					byteIndex := y*bytesPerRow + x/8
					bitIndex := 7 - uint(x%8) // #nosec G115
					data[byteIndex] |= 1 << bitIndex
				}
			}
		}
	}

	p.MonochromeData = data
	return data
}

// ResizeToWidth redimensiona una imagen para ajustarla al ancho especificado
func ResizeToWidth(img image.Image, newWidth, maxWidth int) image.Image {
	if newWidth <= 0 {
		log.Println("Ancho máximo debe ser mayor que 0, devolviendo imagen original")
		return img
	}

	if newWidth > maxWidth {
		log.Printf("Ancho máximo %d excedido, redimensionando a %d", newWidth, maxWidth)
		newWidth = maxWidth
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Si la imagen es menor o igual al ancho máximo, no hacer nada
	if width <= newWidth {
		return img
	}

	// Calcular nueva altura proporcional
	newHeight := int(float64(height) * (float64(newWidth) / float64(width)))

	// Crear nueva imagen redimensionada
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Usar algoritmo de escalado de alta calidad
	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

	return newImg
}

func LoadImage(filename string) (image.Image, error) {
	file, err := utils.SafeOpen(filename)
	if err != nil {
		log.Printf("Error al abrir imagen: %v", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error al cerrar imagen: %v", err)
		}
	}(file)

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Printf("Error al decodificar imagen: %v", err)
	}

	return img, nil
}
