package imaging

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// DitherMode representa los diferentes algoritmos de dithering disponibles
type DitherMode int

const (
	// DitherNone indicates no dithering, just thresholding
	DitherNone DitherMode = iota
	// DitherFloydSteinberg indicates Floyd-Steinberg dithering algorithm
	DitherFloydSteinberg
	// DitherAtkinson indicates Atkinson dithering algorithm
	DitherAtkinson
)

// DitherProcessor es una interfaz para procesar imágenes con dithering
type DitherProcessor interface {
	// Apply aplica el algoritmo de dithering a una imagen
	// threshold es el umbral para conversión a B/N (0-255)
	Apply(img image.Image, threshold uint8) image.Image
}

// === Implementaciones de algoritmos de dithering ===

// FloydSteinbergDithering implementa el algoritmo Floyd-Steinberg
type FloydSteinbergDithering struct{}

// Apply applies the Floyd-Steinberg dithering algorithm to the source image
func (f *FloydSteinbergDithering) Apply(src image.Image, _ uint8) image.Image {
	bounds := src.Bounds()
	// Crear imagen en escala de grises para trabajar
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, src, bounds.Min, draw.Src)

	// Crear imagen de salida en blanco y negro
	dst := image.NewGray(bounds)

	// Aplicar Floyd-Steinberg
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := gray.GrayAt(x, y).Y
			newPixel := uint8(0)

			dst.SetGray(x, y, color.Gray{Y: newPixel})

			// Calcular error
			err := oldPixel - newPixel

			// Distribuir error a píxeles vecinos
			distributeError := func(dx, dy int, factor float32) {
				nx, ny := x+dx, y+dy
				if nx >= bounds.Min.X && nx < bounds.Max.X &&
					ny >= bounds.Min.Y && ny < bounds.Max.Y {
					oldVal := gray.GrayAt(nx, ny).Y
					newVal := oldVal + uint8(float32(err)*factor)

					gray.SetGray(nx, ny, color.Gray{Y: newVal}) //nolint:gosec
				}
			}

			// Matriz de Floyd-Steinberg
			//     X   7/16
			// 3/16 5/16 1/16
			distributeError(1, 0, 7.0/16.0)  // derecha
			distributeError(-1, 1, 3.0/16.0) // abajo-izquierda
			distributeError(0, 1, 5.0/16.0)  // abajo
			distributeError(1, 1, 1.0/16.0)  // abajo-derecha
		}
	}

	return dst
}

// AtkinsonDithering implementa el algoritmo Atkinson
type AtkinsonDithering struct{}

// Apply applies the Atkinson dithering algorithm to the source image
func (a *AtkinsonDithering) Apply(src image.Image, _ uint8) image.Image {
	bounds := src.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, src, bounds.Min, draw.Src)

	dst := image.NewGray(bounds)

	// Aplicar Atkinson
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := gray.GrayAt(x, y).Y
			newPixel := uint8(0)

			dst.SetGray(x, y, color.Gray{Y: newPixel})

			// Calcular error (Atkinson difunde solo 3/4 del error)
			err := oldPixel - newPixel
			diffusedError := err * 3 / 4

			// Distribuir error
			distributeError := func(dx, dy int) {
				nx, ny := x+dx, y+dy
				if nx >= bounds.Min.X && nx < bounds.Max.X &&
					ny >= bounds.Min.Y && ny < bounds.Max.Y {
					oldVal := gray.GrayAt(nx, ny).Y
					// Cada vecino recibe 1/8 del error original
					newVal := oldVal + diffusedError/8

					gray.SetGray(nx, ny, color.Gray{Y: newVal})
				}
			}

			// Patrón de Atkinson
			//     X   1   2
			// 1   1   1
			//     1
			distributeError(1, 0)  // derecha
			distributeError(2, 0)  // derecha x2
			distributeError(-1, 1) // abajo-izquierda
			distributeError(0, 1)  // abajo
			distributeError(1, 1)  // abajo-derecha
			distributeError(0, 2)  // abajo x2
		}
	}

	return dst
}

// ThresholdDithering implementa umbralización simple (sin dithering real)
type ThresholdDithering struct{}

// Apply applies simple thresholding to the source image
func (t *ThresholdDithering) Apply(src image.Image, threshold uint8) image.Image {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Convertir a escala de grises
			gray := color.GrayModel.Convert(src.At(x, y)).(color.Gray)

			// Aplicar umbral
			if gray.Y > threshold {
				dst.SetGray(x, y, color.Gray{Y: 255})
			} else {
				dst.SetGray(x, y, color.Gray{})
			}
		}
	}

	return dst
}

// GetDitherProcessor devuelve el procesador de dithering según el modo
func GetDitherProcessor(mode DitherMode) (DitherProcessor, error) {
	switch mode {
	case DitherNone:
		return &ThresholdDithering{}, nil
	case DitherFloydSteinberg:
		return &FloydSteinbergDithering{}, nil
	case DitherAtkinson:
		return &AtkinsonDithering{}, nil
	default:
		return nil, fmt.Errorf("unsupported dither mode: %d", mode)
	}
}

// ProcessImageWithDithering aplica dithering a una imagen
// Esta es la función principal que usarán los protocolos
func ProcessImageWithDithering(img image.Image, mode DitherMode, threshold uint8) (image.Image, error) {
	if img == nil {
		return nil, fmt.Errorf("imaging cannot be nil")
	}

	processor, err := GetDitherProcessor(mode)
	if err != nil {
		return nil, err
	}

	return processor.Apply(img, threshold), nil
}
