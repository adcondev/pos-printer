package profile

import "github.com/AdConDev/pos-printer/types"

// Profile define todas las características físicas y capacidades de una impresora
type Profile struct {
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
	SupportsColor    bool // Soporta impresión a color

	QRMaxVersion byte // Máxima versión soportada

	// Juegos de caracteres
	CharacterSets    []types.CharacterSet // Códigos de página soportados
	DefaultCharSet   types.CharacterSet   // Código de página por defecto0
	ActiveCharSet    types.CharacterSet   // Código de página actualmente activo
	DefaultKanjiMode bool                 // Modo Kanji por defecto (true/false)

	// Configuración avanzada (opcional)
	FeedLinesAfterCut int // Líneas de avance después de cortar
	ImageThreshold    int // Umbral para conversión B/N (0-255)

	// Fuentes
	Fonts map[string]int // Lista de fuentes soportadas, nombre -> ancho (en puntos)

	// Extensible para características específicas
	// Usar un mapa genérico permite agregar características sin cambiar la estructura
	ExtendedFeatures map[string]interface{}
}

// ModelInfo devuelve una representación de string del modelo
func (p *Profile) ModelInfo() string {
	return p.Vendor + " " + p.Model
}

// GetCharWidth calcula el ancho físico de un caracter en milímetros
func (p *Profile) GetCharWidth(font string) int {
	return p.DotsPerLine / p.Fonts[font]
}

// CreatePt210 crea un perfil para impresora térmica de 58mm PT-58N
func CreatePt210() *Profile {
	p := CreateProfile58mm()
	p.Model = "58mm PT-210"
	p.Vendor = "GOOJPRT"
	p.DefaultKanjiMode = true // Esta impresora inicia con Kanji activo
	p.CharacterSets = []types.CharacterSet{
		types.CP437,
		types.Katakana,
		types.CP850,
		types.CP860,
		types.CP863,
		types.CP865,
		types.WestEurope,
		types.WCP1252,
		types.CP866,
		types.CP852,
		types.CP858,
	}

	p.DefaultCharSet = 0 // CP858 para español
	p.QRMaxVersion = 19  // Máxima versión QR soportada
	p.SupportsQR = true  // Soporta QR nativo
	return p
}

func CreateProfGP_58N() *Profile {
	p := CreateProfile58mm()
	p.Model = "58mm GP-58N"
	p.CharacterSets = []types.CharacterSet{
		types.CP437,
		types.Katakana,
		types.CP850,
		types.CP860,
		types.CP863,
		types.CP865,
		types.WestEurope,
		types.Greek,
		types.Hebrew,
		// types.CP755, // No soportado directamente
		types.Iran,
		types.WCP1252,
		types.CP866,
		types.CP852,
		types.CP858,
		types.IranII,
		types.Latvian,
	}

	p.DefaultCharSet = 19 // CP858 para español
	return p
}

// CreateProfile58mm crea un perfil para impresora térmica de 58mm común
func CreateProfile58mm() *Profile {
	return &Profile{
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
		SupportsColor:    false,

		CharacterSets: []types.CharacterSet{
			types.CP437,
			types.Katakana,
			types.CP850,
			types.CP860,
			types.CP863,
			types.CP865,
			types.WestEurope,
			types.Greek,
			types.Hebrew,
			// types.CP755, // No soportado directamente
			types.Iran,
			types.WCP1252,
			types.CP866,
			types.CP852,
			types.CP858,
			types.IranII,
			types.Latvian,
		}, // Más juegos de caracteres
		DefaultCharSet: 19, // CP858

		Fonts: map[string]int{
			"FontA": 12, // Ancho de 12 puntos
			"FontB": 9,  // Ancho de 9 puntos
		},

		ExtendedFeatures: make(map[string]interface{}),
	}
}

func CreateProfEC_PM_80250() *Profile {
	p := CreateProfile80mm()
	p.Model = "80mm EC-PM-80250"
	p.CharacterSets = []types.CharacterSet{
		types.CP437,
		types.Katakana,
		types.CP850,
		types.CP860,
		types.CP863,
		types.CP865,
		types.WestEurope,
		types.Greek,
		types.Hebrew,
		// types.CP755, // No soportado directamente
		types.Iran,
		types.WCP1252,
		types.CP866,
		types.CP852,
		types.CP858,
		types.IranII,
		types.Latvian,
	}

	p.DefaultCharSet = types.CP437 // CP858 para español
	return p
}

// CreateProfile80mm crea un perfil para impresora térmica de 80mm común
func CreateProfile80mm() *Profile {
	return &Profile{
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
		SupportsColor:    false,

		CharacterSets: []types.CharacterSet{
			types.CP437,
			types.Katakana,
			types.CP850,
			types.CP860,
			types.CP863,
			types.CP865,
			types.WestEurope,
			types.Greek,
			types.Hebrew,
			// types.CP755, // No soportado directamente
			types.Iran,
			types.WCP1252,
			types.CP866,
			types.CP852,
			types.CP858,
			types.IranII,
			types.Latvian,
		},

		// Más juegos de caracteres
		DefaultCharSet: types.CP437, // CP858

		FeedLinesAfterCut: 5,
		ImageThreshold:    128,

		ExtendedFeatures: make(map[string]interface{}),
	}
}

func (p *Profile) HasImageSupport() bool {
	return p.SupportsGraphics
}
