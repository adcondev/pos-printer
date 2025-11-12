package graphics

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
)

const minBorderWidth = 4

// TODO: Check if /internals fits better for custom WriteCloser

// WriteCloser wraps bytes.Buffer to implement io.WriteCloser
type WriteCloser struct {
	*bytes.Buffer
}

// Close is a no-op for WriteCloser
func (wc *WriteCloser) Close() error {
	// TODO: Implement any necessary cleanup if needed
	return nil
}

// NewWriteCloser creates a new WriteCloser instance
func NewWriteCloser() *WriteCloser {
	return &WriteCloser{
		Buffer: new(bytes.Buffer),
	}
}

// QROptions contiene opciones para generar QR (nativo o imagen)
type QROptions struct {
	// === Opciones comunes (funcionan en nativo e imagen) ===
	Model           posqr.Model           // Model1, Model2, MicroQR
	moduleSize      posqr.ModuleSize      // Calculado en base a PixelWidth
	ErrorCorrection posqr.ErrorCorrection // L, M, Q, H

	// === Opciones solo para QR como imagen ===
	// TODO: Try to set to half Dots Per Line
	PixelWidth int // Ancho en píxeles

	// === Opciones útiles para impresora monocromática ===
	LogoPath      string // Ruta al archivo del logo
	LogoSizeMulti int    // Multiplicador del tamaño del logo (2-6)
	CircleShape   bool   // Usar bloques circulares
	HalftonePath  string // Ruta a imagen para efecto semitono
}

// GetModuleSize retorna el tamaño del módulo calculado
func (qro *QROptions) GetModuleSize() posqr.ModuleSize {
	return qro.moduleSize
}

// SetModuleSize calcula y establece el tamaño del módulo basado en PixelWidth y el tamaño de la cuadrícula del QR
func (qro *QROptions) SetModuleSize(data string) (*qrcode.QRCode, error) {
	// Validación de datos vacíos
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	// Establecer valores por defecto si no están configurados
	if qro.PixelWidth <= 0 {
		qro.PixelWidth = 288
		log.Printf("QR: using default pixel width %d", qro.PixelWidth)
	}
	if qro.ErrorCorrection < posqr.LevelL || qro.ErrorCorrection > posqr.LevelH {
		qro.ErrorCorrection = posqr.LevelM
		log.Printf("QR: using default error correction level M")
	}

	// Crear QR code
	qrc, err := qrcode.NewWith(data, WithErrorLevel(qro.ErrorCorrection))
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	gridSize := qrc.Dimension()
	if gridSize == 0 {
		return nil, fmt.Errorf("invalid QR code grid size")
	}

	// Tamaño del módulo con mejor precisión
	moduleSize := qro.PixelWidth / gridSize

	// Aplicar límites
	switch {
	case moduleSize < int(posqr.DefaultModuleSize):
		log.Printf("QR: calculated module size %d too small, using minimum %d",
			moduleSize, posqr.DefaultModuleSize)
		qro.moduleSize = posqr.DefaultModuleSize
	case moduleSize > int(posqr.MaxModuleSize):
		log.Printf("QR: calculated module size %d too large, using maximum %d",
			moduleSize, posqr.MaxModuleSize)
		qro.moduleSize = posqr.MaxModuleSize
	default:
		qro.moduleSize = posqr.ModuleSize(moduleSize)
	}

	log.Printf("QR: grid=%dx%d, pixel_width=%d, module_size=%d",
		gridSize, gridSize, qro.PixelWidth, qro.moduleSize)

	return qrc, nil
}

// DefaultQROptions retorna opciones por defecto optimizadas para impresoras térmicas
func DefaultQROptions() *QROptions {
	return &QROptions{
		Model:           posqr.Model2,
		ErrorCorrection: posqr.LevelQ,
		PixelWidth:      288, // Buen tamaño para 58mm
		LogoSizeMulti:   3,
		CircleShape:     false, // Los cuadrados son mejores para térmica
	}
}

// GenerateQRImage genera un QR code como imagen optimizada para impresora térmica
func GenerateQRImage(data string, opts *QROptions) (image.Image, error) {
	if opts == nil {
		opts = DefaultQROptions()
	}

	// El objetivo es hacer el QR tan grande y legible como sea posible, sin pasarse de tu límite de PixelWidth.
	qrc, err := opts.SetModuleSize(data)
	if err != nil {
		opts.moduleSize = posqr.DefaultModuleSize
		log.Printf(
			"warning: could not set module size based on image width and grid size: %v. Using minimum module size %d",
			err,
			opts.moduleSize,
		)
	}

	// Construir opciones de imagen
	imgOpts := buildImageOptions(opts)

	// Generar imagen en memoria
	buf := NewWriteCloser()
	w := standard.NewWithWriter(buf, imgOpts...)
	defer func(w *standard.Writer) {
		err := w.Close()
		if err != nil {
			log.Printf("error closing QR writer: %v", err)
		}
	}(w)

	if err := qrc.Save(w); err != nil {
		return nil, fmt.Errorf("generate QR image: %w", err)
	}

	// Decodificar imagen
	img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("decode QR image: %w", err)
	}

	return img, nil
}

// buildImageOptions construye las opciones mínimas útiles para impresora térmica
func buildImageOptions(opts *QROptions) []standard.ImageOption {
	var imgOpts []standard.ImageOption

	// TODO: Define module size base on PixelWidth and grid size, final user should not see module size only PixelWidth

	imgOpts = append(imgOpts, standard.WithQRWidth(uint8(opts.moduleSize))) // Module PixelWidth
	imgOpts = append(imgOpts, standard.WithBorderWidth(minBorderWidth))     // Silence Zone

	// TODO: Implement HalftonePath option
	// imgOpts = append(imgOpts, standard.WithHalftone("./assets/images/gopher.jpeg")) // Halftone effect

	// LogoPath si está habilitado y existe
	if opts.LogoPath != "" {
		// Detectar formato por extensión
		if strings.HasSuffix(opts.LogoPath, ".png") {
			imgOpts = append(imgOpts, standard.WithLogoImageFilePNG(opts.LogoPath))
		} else if strings.HasSuffix(opts.LogoPath, ".jpg") || strings.HasSuffix(opts.LogoPath, ".jpeg") {
			imgOpts = append(imgOpts, standard.WithLogoImageFileJPEG(opts.LogoPath))
		}

		// Tamaño del logo
		if opts.LogoSizeMulti > 0 {
			imgOpts = append(imgOpts, standard.WithLogoSizeMultiplier(opts.LogoSizeMulti))
		}
		// Zona segura para el logo (contorno blanco)
		imgOpts = append(imgOpts, standard.WithLogoSafeZone())
	}

	// Forma circular (opcional, puede afectar legibilidad)
	if opts.CircleShape {
		imgOpts = append(imgOpts, standard.WithCircleShape())
	}

	return imgOpts
}

// WithErrorLevel convierte el nivel de corrección de errores poster a go-qrcode
func WithErrorLevel(level posqr.ErrorCorrection) qrcode.EncodeOption {
	switch level {
	case posqr.LevelL:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow)
	case posqr.LevelM:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium)
	case posqr.LevelQ:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart)
	case posqr.LevelH:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest)
	default:
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium)
	}
}

// TODO: Integrate halftone to QROptions

// createQRWithHalftone generates a QR code using the WithHalftone option.
func createQRWithHalftone(content string) {
	qr, err := qrcode.New(content)
	if err != nil {
		fmt.Printf("create qrcode failed: %v\n", err)
		return
	}

	// Please replace with the actual path to the halftone image.
	halftonePath := "../assets/example/monna-lisa.png"
	if _, err := os.Stat(halftonePath); os.IsNotExist(err) {
		fmt.Printf("halftone image file %s not found\n", halftonePath)
		return
	}

	options := []standard.ImageOption{
		standard.WithHalftone(halftonePath),
		standard.WithQRWidth(21),
	}
	writer, err := standard.New("../assets/example/qrcode_with_halftone.png", options...)
	if err != nil {
		fmt.Printf("create writer failed: %v\n", err)
		return
	}
	defer func(writer *standard.Writer) {
		err := writer.Close()
		if err != nil {
			log.Printf("error closing writer: %v\n", err)
		}
	}(writer)
	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}
}
