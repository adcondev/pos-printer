package graphics

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
)

const (
	// minBorderWidth es el quiet zone mínimo recomendado por el estándar QR (4 módulos)
	minBorderWidth = 4

	// maxPixelWidth es el ancho máximo soportado para impresoras térmicas de 80mm
	maxPixelWidth = 576

	// minPixelWidth es el mínimo para QR Version 1 (21x21) con módulos de 3px + borders
	// Cálculo: (21 + 2*4) * 3 = 87px, pero usamos 63 como mínimo práctico
	minPixelWidth = 87

	minGridSize = 21 // QR Version 1 (21x21 modules)

	// Logo size multiplier limits
	minSizeMulti     = 1
	defaultSizeMulti = 3
	maxSizeMulti     = 5
)

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
	ErrorCorrection posqr.ErrorCorrection // L, M, Q, H

	// === Opciones solo para QR como imagen ===
	PixelWidth int              // Ancho en píxeles
	moduleSize posqr.ModuleSize // Calculado en base a PixelWidth

	// === Opciones útiles para impresora monocromática ===
	LogoPath      string // Ruta al archivo del logo
	LogoSizeMulti int    // Multiplicador del tamaño del logo (1-10)
	CircleShape   bool   // Usar bloques circulares
	HalftonePath  string // Ruta a imagen para efecto semitono
}

// GetModuleSize retorna el tamaño del módulo calculado
func (qro *QROptions) GetModuleSize() posqr.ModuleSize {
	return qro.moduleSize
}

// SetModuleSize calcula y establece el tamaño del módulo basado en PixelWidth y el tamaño de la cuadrícula del QR
func (qro *QROptions) SetModuleSize(data string) (*qrcode.QRCode, error) {
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	if len(data) > posqr.MaxDataLength {
		return nil, fmt.Errorf("QR data too long: %d bytes (maximum %d)",
			len(data), posqr.MaxDataLength)
	}
	if !utf8.ValidString(data) {
		log.Printf("warning: QR data contains invalid UTF-8 characters")
	}

	// Validación de PixelWidth
	if qro.PixelWidth < minPixelWidth {
		log.Printf("warning: pixel_width %d < minimum %d, adjusting to minimum",
			qro.PixelWidth, minPixelWidth)
		qro.PixelWidth = minPixelWidth
	}

	// Establecer valores por defecto si no están configurados
	if qro.PixelWidth == 0 {
		qro.PixelWidth = 288
		log.Printf("QR: using default pixel width %d", qro.PixelWidth)
	}

	if qro.PixelWidth > maxPixelWidth {
		log.Printf("warning: pixel_width %d exceeds maximum %d, clamping",
			qro.PixelWidth, maxPixelWidth)
		qro.PixelWidth = maxPixelWidth
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

	if gridSize < minGridSize {
		return nil, fmt.Errorf("QR grid size %d is too small (minimum %d)",
			gridSize, minGridSize)
	}

	// Tamaño del módulo con mejor precisión
	totalModules := gridSize + (2 * minBorderWidth)
	moduleSize := qro.PixelWidth / totalModules

	log.Printf("QR: grid=%dx%d, border=%d modules, total=%d modules, requested=%dpx",
		gridSize, gridSize, minBorderWidth, totalModules, qro.PixelWidth)

	// Aplicar límites al module size
	switch {
	case moduleSize < int(posqr.DefaultModuleSize):
		log.Printf("QR: calculated module size %d too small, using default %d",
			moduleSize, posqr.DefaultModuleSize)
		qro.moduleSize = posqr.DefaultModuleSize
	case moduleSize > int(posqr.MaxModuleSize):
		log.Printf("QR: calculated module size %d too large, using maximum %d",
			moduleSize, posqr.MaxModuleSize)
		qro.moduleSize = posqr.MaxModuleSize
	default:
		qro.moduleSize = posqr.ModuleSize(moduleSize)
	}

	// Calcular y reportar tamaño final real
	actualWidth := totalModules * int(qro.moduleSize)
	dataWidth := gridSize * int(qro.moduleSize)
	borderSize := (2 * minBorderWidth) * int(qro.moduleSize)

	log.Printf("QR: module_size=%d, data=%dx%dpx, border=%dpx, actual_total=%dx%dpx",
		qro.moduleSize, dataWidth, dataWidth, borderSize, actualWidth, actualWidth)

	// ⚠Advertencia si el tamaño difiere del solicitado
	if actualWidth != qro.PixelWidth {
		diff := actualWidth - qro.PixelWidth
		if diff > 0 {
			log.Printf("warning: actual QR size %dpx exceeds requested %dpx by %dpx (rounding up to module boundary)",
				actualWidth, qro.PixelWidth, diff)
		} else {
			log.Printf("info: actual QR size %dpx is smaller than requested %dpx by %dpx (rounding down to module boundary)",
				actualWidth, qro.PixelWidth, -diff)
		}
	}

	return qrc, nil
}

// DefaultQROptions retorna opciones por defecto optimizadas para impresoras térmicas.
//
// PixelWidth por defecto: 288px total
//   - Para grid típico 25x25: total modules = 33, module size = 8px
//   - Área datos: 25 × 8 = 200px
//   - Quiet zone: 8 × 8 = 64px (32px por lado)
//   - Total: 264px (ajustado a múltiplo de module size)
func DefaultQROptions() *QROptions {
	return &QROptions{
		Model:           posqr.Model2,
		ErrorCorrection: posqr.LevelQ,
		PixelWidth:      288, // Tamaño total incluyendo quiet zone
		LogoSizeMulti:   3,
		CircleShape:     false,
	}
}

// GenerateQRImage genera un QR code como imagen optimizada para impresora térmica
func GenerateQRImage(data string, opts *QROptions) (image.Image, error) {
	if data == "" {
		return nil, fmt.Errorf("QR data cannot be empty")
	}
	if opts == nil {
		opts = DefaultQROptions()
	}

	// Validar archivos antes de generar
	if opts.LogoPath != "" {
		if _, err := os.Stat(opts.LogoPath); os.IsNotExist(err) {
			log.Printf("warning: logo file not found: %s, ignoring", opts.LogoPath)
			opts.LogoPath = "" // Limpiar para evitar error
		}
	}

	if opts.HalftonePath != "" {
		if _, err := os.Stat(opts.HalftonePath); os.IsNotExist(err) {
			log.Printf("warning: halftone file not found: %s, ignoring", opts.HalftonePath)
			opts.HalftonePath = "" // Limpiar para evitar error
		}
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

	imgOpts = append(imgOpts, standard.WithQRWidth(uint8(opts.moduleSize))) // Module PixelWidth
	imgOpts = append(imgOpts, standard.WithBorderWidth(minBorderWidth))     // Silence Zone

	// LogoPath si está habilitado y existe
	if opts.LogoPath != "" {
		// Detectar formato por extensión
		lowerPath := strings.ToLower(opts.LogoPath)
		switch {
		case strings.HasSuffix(lowerPath, ".png"):
			imgOpts = append(imgOpts, standard.WithLogoImageFilePNG(opts.LogoPath))
		case strings.HasSuffix(lowerPath, ".jpg") || strings.HasSuffix(lowerPath, ".jpeg"):
			imgOpts = append(imgOpts, standard.WithLogoImageFileJPEG(opts.LogoPath))
		default:
			log.Printf("warning: unsupported logo format: %s (only PNG/JPEG supported)", opts.LogoPath)
		}

		// Tamaño del logo
		if opts.LogoSizeMulti > 0 {
			if opts.LogoSizeMulti < minSizeMulti || opts.LogoSizeMulti > maxSizeMulti {
				log.Printf("warning: logo_size_multi %d out of range [%d-%d], using default %d",
					opts.LogoSizeMulti, minSizeMulti, maxSizeMulti, defaultSizeMulti)
				opts.LogoSizeMulti = defaultSizeMulti
			}
			imgOpts = append(imgOpts, standard.WithLogoSizeMultiplier(opts.LogoSizeMulti))
		}
		// Zona segura para el logo
		imgOpts = append(imgOpts, standard.WithLogoSafeZone())
	}

	// Can't be used together: HalftonePath and CircleShape
	if opts.HalftonePath != "" {
		log.Printf("qr: using halftone image (circle shape disabled)")
		imgOpts = append(imgOpts, standard.WithHalftone(opts.HalftonePath))
	} else if opts.CircleShape {
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
