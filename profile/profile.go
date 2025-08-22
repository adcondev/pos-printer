package profile

import (
	"github.com/adcondev/pos-printer/encoding"
)

// Escpos define todas las características físicas y capacidades de una impresora
type Escpos struct {
	// Información básica
	Model       string
	Vendor      string
	Description string

	// Características físicas
	PaperWidth  float64 // en mm (58mm, 80mm, etc.)
	PaperHeight float64 // en mm (0 para rollo continuo)
	DPI         int     // Dots Per Inch (ej. 203, 300)
	DotsPerLine int     // Puntos por línea (ej. 384, 576)
	PrintWidth  int     // Ancho de impresión en mm (ej. 48)

	// Capacidades
	SupportsGraphics bool // Soporta gráficos (imágenes)
	SupportsBarcode  bool // Soporta códigos de barra nativos
	SupportsQR       bool // Soporta códigos QR nativos
	SupportsCutter   bool // Tiene cortador automático
	SupportsDrawer   bool // Soporta cajón de dinero

	QRMaxVersion byte // Máxima versión soportada

	// Juegos de caracteres
	CharacterSets  []encoding.CharacterSet // Códigos de página soportados
	DefaultCharSet encoding.CharacterSet   // Código de página por defecto0

	// Configuración avanzada (opcional)
	ImageThreshold int // Umbral para conversión B/N (0-255)

	// Fuentes
	Fonts map[string]int // Lista de fuentes soportadas, nombre -> ancho (en puntos)

	// Extensible para características específicas
	// Usar un mapa genérico permite agregar características sin cambiar la estructura
	ExtendedFeatures map[string]interface{}
}

// ModelInfo devuelve una representación de string del modelo
func (p *Escpos) ModelInfo() string {
	return p.Vendor + " " + p.Model
}

// GetCharWidth calcula el ancho físico de un caracter en milímetros
func (p *Escpos) GetCharWidth(font string) int {
	return p.DotsPerLine / p.Fonts[font]
}

// CreatePt210 crea un perfil para impresora térmica de 58mm PT-58N
func CreatePt210() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm PT-210"
	p.Vendor = "GOOJPRT"
	p.CharacterSets = []encoding.CharacterSet{
		encoding.CP437,
		encoding.Katakana,
		encoding.CP850,
		encoding.CP860,
		encoding.CP863,
		encoding.CP865,
		encoding.WestEurope,
		encoding.WCP1252,
		encoding.CP866,
		encoding.CP852,
		encoding.CP858,
	}

	p.DefaultCharSet = 0 // CP858 para español
	p.QRMaxVersion = 19  // Máxima versión QR soportada
	p.SupportsQR = true  // Soporta QR nativo
	return p
}

func CreateProfGP_58N() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm GP-58N"
	p.CharacterSets = []encoding.CharacterSet{
		encoding.CP437,
		encoding.Katakana,
		encoding.CP850,
		encoding.CP860,
		encoding.CP863,
		encoding.CP865,
		encoding.WestEurope,
		encoding.Greek,
		encoding.Hebrew,
		// encoding.CP755, // No soportado directamente
		encoding.Iran,
		encoding.WCP1252,
		encoding.CP866,
		encoding.CP852,
		encoding.CP858,
		encoding.IranII,
		encoding.Latvian,
	}

	p.DefaultCharSet = 19 // CP858 para español
	return p
}

// CreateProfile58mm crea un perfil para impresora térmica de 58mm común
func CreateProfile58mm() *Escpos {
	return &Escpos{
		Model:       "Generic 58mm",
		Vendor:      "Generic",
		Description: "Impresora térmica genérica de 58mm",

		PaperWidth:  58,
		DPI:         203,
		DotsPerLine: 384, // Típico para 58mm a 203 DPI
		PrintWidth:  48,  // Ancho de impresión efectivo

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       false, // Muchas impresoras baratas no soportan QR nativo
		SupportsCutter:   false,
		SupportsDrawer:   false,

		CharacterSets: []encoding.CharacterSet{
			encoding.CP437,
			encoding.Katakana,
			encoding.CP850,
			encoding.CP860,
			encoding.CP863,
			encoding.CP865,
			encoding.WestEurope,
			encoding.Greek,
			encoding.Hebrew,
			// encoding.CP755, // No soportado directamente
			encoding.Iran,
			encoding.WCP1252,
			encoding.CP866,
			encoding.CP852,
			encoding.CP858,
			encoding.IranII,
			encoding.Latvian,
		}, // Más juegos de caracteres
		DefaultCharSet: 19, // CP858

		Fonts: map[string]int{
			"FontA": 12, // Ancho de 12 puntos
			"FontB": 9,  // Ancho de 9 puntos
		},

		ExtendedFeatures: make(map[string]interface{}),
	}
}

func CreateProfEC_PM_80250() *Escpos {
	p := CreateProfile80mm()
	p.Model = "80mm EC-PM-80250"
	p.CharacterSets = []encoding.CharacterSet{
		encoding.CP437,
		encoding.Katakana,
		encoding.CP850,
		encoding.CP860,
		encoding.CP863,
		encoding.CP865,
		encoding.WestEurope,
		encoding.Greek,
		encoding.Hebrew,
		// encoding.CP755, // No soportado directamente
		encoding.Iran,
		encoding.WCP1252,
		encoding.CP866,
		encoding.CP852,
		encoding.CP858,
		encoding.IranII,
		encoding.Latvian,
	}

	p.DefaultCharSet = encoding.CP437 // CP858 para español
	return p
}

// CreateProfile80mm crea un perfil para impresora térmica de 80mm común
func CreateProfile80mm() *Escpos {
	return &Escpos{
		Model:       "Generic 80mm",
		Vendor:      "Generic",
		Description: "Impresora térmica genérica de 80mm",

		PaperWidth:  80,
		DPI:         203,
		DotsPerLine: 576, // Típico para 80mm a 203 DPI

		SupportsGraphics: true,
		SupportsBarcode:  true,
		SupportsQR:       true, // Las 80mm suelen tener más funciones
		SupportsCutter:   true,
		SupportsDrawer:   true,

		CharacterSets: []encoding.CharacterSet{
			encoding.CP437,
			encoding.Katakana,
			encoding.CP850,
			encoding.CP860,
			encoding.CP863,
			encoding.CP865,
			encoding.WestEurope,
			encoding.Greek,
			encoding.Hebrew,
			// encoding.CP755, // No soportado directamente
			encoding.Iran,
			encoding.WCP1252,
			encoding.CP866,
			encoding.CP852,
			encoding.CP858,
			encoding.IranII,
			encoding.Latvian,
		},

		// Más juegos de caracteres
		DefaultCharSet: encoding.CP437, // CP858

		ImageThreshold: 128,

		ExtendedFeatures: make(map[string]interface{}),
	}
}

func (p *Escpos) HasImageSupport() bool {
	return p.SupportsGraphics
}
