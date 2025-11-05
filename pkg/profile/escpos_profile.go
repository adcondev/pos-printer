package profile

import (
	"github.com/adcondev/pos-printer/pkg/controllers/escpos/character"
	"github.com/adcondev/pos-printer/pkg/graphics"
)

// Escpos define todas las características físicas y capacidades de una impresora
type Escpos struct {
	// Información básica
	Model string // Same name used for printer connection

	// Características físicas
	PaperWidth  float64 // en mm (58mm, 80mm, etc.)
	PaperHeight float64 // en mm (0 para rollo continuo)
	DPI         int     // Dots Per Inch (ej. 203, 300)
	DotsPerLine int     // Puntos por línea (ej. 384, 576)
	PrintWidth  int     // Ancho de impresión en mm (ej. 48, 42 en Font A)

	// Capacidades
	SupportsGraphics bool // Soporta gráficos (imágenes)
	SupportsBarcode  bool // Soporta códigos de barra nativos
	SupportsQR       bool // Soporta códigos QR nativos
	SupportsCutter   bool // Tiene cortador automático
	SupportsDrawer   bool // Soporta cajón de dinero

	QRMaxVersion byte // Máxima versión soportada

	// Code table and encoding configuration
	CodeTable character.CodeTable

	// Configuración avanzada (opcional)
	ImageThreshold int                 // Umbral para conversión B/N (0-255)
	Dithering      graphics.DitherMode // Tipo de dithering por defecto

	// Fuentes
	Fonts map[string]int // Lista de fuentes soportadas, nombre -> ancho (en puntos)
}

// CreatePt210 crea un perfil para impresora térmica de 58mm PT-58N
func CreatePt210() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm PT-210"

	p.CodeTable = character.PC850
	p.QRMaxVersion = 19 // Máxima versión QR soportada
	p.SupportsQR = true // Soporta QR nativo
	return p
}

// CreateGP58N crea un perfil para impresora térmica de 58mm GP-58N
func CreateGP58N() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm GP-58N"

	p.CodeTable = character.PC850
	return p
}

// CreateProfile58mm crea un perfil para impresora térmica de 58mm común
func CreateProfile58mm() *Escpos {
	return &Escpos{
		Model: "Generic 58mm",

		PaperWidth:  58,
		DPI:         203,
		DotsPerLine: 384, // Típico para 58mm a 203 DPI
		PrintWidth:  48,  // Ancho de impresión efectivo

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       false, // Muchas impresoras baratas no soportan QR nativo
		SupportsCutter:   false,
		SupportsDrawer:   false,

		CodeTable: character.PC850,

		Fonts: map[string]int{
			"FontA": 12, // Ancho de 12 puntos
			"FontB": 9,  // Ancho de 9 puntos
		},
	}
}

// CreateECPM80250 crea un perfil para impresora térmica de 80mm EC-PM-80250
func CreateECPM80250() *Escpos {
	p := CreateProfile80mm()
	p.Model = "80mm EC-PM-80250"
	return p
}

// CreateProfile80mm crea un perfil para impresora térmica de 80mm común
func CreateProfile80mm() *Escpos {
	return &Escpos{
		Model: "Generic 80mm",

		PaperWidth:  80,
		DPI:         203,
		DotsPerLine: 576, // Típico para 80mm (72mm) a 203 DPI

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       true, // Las 80mm suelen tener más funciones
		SupportsCutter:   true,
		SupportsDrawer:   true,

		// Más juegos de caracteres
		CodeTable: character.PC850, // CP850

		ImageThreshold: 128,
	}
}
