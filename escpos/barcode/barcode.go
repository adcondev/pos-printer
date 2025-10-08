package barcode

import (
	"errors"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Tipos y Constantes
// ============================================================================

// Definición de tipos

// HRIPosition representa las posiciones de impresión de los caracteres HRI
type HRIPosition byte

// HRIFont representa los tipos de fuente para los caracteres HRI
type HRIFont byte

// Height representa la altura del código de barras (en puntos)
type Height byte

// Width representa el ancho del módulo horizontal del código de barras
type Width byte

// Symbology representa los tipos de simbología del código de barras
type Symbology byte

// Code128Set representa los conjuntos de códigos para CODE128 (A/B/C)
type Code128Set byte

// Constantes

const (
	// Modos de posición HRI (numéricos)
	HRINotPrinted HRIPosition = 0x00 // n = 0
	HRIAbove      HRIPosition = 0x01 // n = 1
	HRIBelow      HRIPosition = 0x02 // n = 2
	HRIBoth       HRIPosition = 0x03 // n = 3

	// Modos de posición HRI (ASCII)
	HRINotPrintedASCII HRIPosition = '0' // n = 48
	HRIAboveASCII      HRIPosition = '1' // n = 49
	HRIBelowASCII      HRIPosition = '2' // n = 50
	HRIBothASCII       HRIPosition = '3' // n = 51
)

const (
	// Tipos de fuente HRI (numéricos)
	HRIFontA HRIFont = 0x00 // n = 0
	HRIFontB HRIFont = 0x01 // n = 1
	HRIFontC HRIFont = 0x02 // n = 2
	HRIFontD HRIFont = 0x03 // n = 3
	HRIFontE HRIFont = 0x04 // n = 4

	// Tipos de fuente HRI (ASCII)
	HRIFontAASCII HRIFont = '0' // n = 48
	HRIFontBASCII HRIFont = '1' // n = 49
	HRIFontCASCII HRIFont = '2' // n = 50
	HRIFontDASCII HRIFont = '3' // n = 51
	HRIFontEASCII HRIFont = '4' // n = 52

	// Fuentes especiales
	HRISpecialFontA HRIFont = 97 // n = 97
	HRISpecialFontB HRIFont = 98 // n = 98
)

const (
	MinHeight     Height = 1   // Altura mínima
	MaxHeight     Height = 255 // Altura máxima
	DefaultHeight Height = 162 // Valor por defecto (según modelo)
)

const (
	// Valores estándar de ancho
	MinWidth     Width = 2 // Ancho mínimo
	MaxWidth     Width = 6 // Ancho máximo
	DefaultWidth Width = 3 // Ancho por defecto

	// Valores extendidos (dependen del modelo)
	ExtendedMinWidth Width = 68 // Mínimo extendido
	ExtendedMaxWidth Width = 76 // Máximo extendido
)

const (
	// Simbologías Function A (terminadas en NUL)
	UPCA    Symbology = 0 // UPC-A (11-12 dígitos)
	UPCE    Symbology = 1 // UPC-E (6-8, 11-12 dígitos)
	JAN13   Symbology = 2 // JAN13/EAN13 (12-13 dígitos)
	JAN8    Symbology = 3 // JAN8/EAN8 (7-8 dígitos)
	CODE39  Symbology = 4 // CODE39 (longitud variable)
	ITF     Symbology = 5 // Interleaved 2 of 5 (pares de dígitos)
	CODABAR Symbology = 6 // CODABAR/NW-7 (longitud variable)
)

const (
	// Simbologías Function B (prefijo de longitud)
	UPCAB           Symbology = 65 // UPC-A (11-12 dígitos)
	UPCEB           Symbology = 66 // UPC-E (6-8, 11-12 dígitos)
	EAN13           Symbology = 67 // EAN13 (12-13 dígitos)
	EAN8            Symbology = 68 // EAN8 (7-8 dígitos)
	CODE39B         Symbology = 69 // CODE39 (1-255 chars)
	ITFB            Symbology = 70 // ITF (2-254 pares)
	CODABARB        Symbology = 71 // CODABAR (2-255 chars)
	CODE93          Symbology = 72 // CODE93 (1-255 chars)
	CODE128         Symbology = 73 // CODE128 (2-255 bytes)
	GS1128          Symbology = 74 // GS1-128 (2-255 bytes)
	GS1DataBarOmni  Symbology = 75 // GS1 DataBar Omnidirectional (13 dígitos)
	GS1DataBarTrunc Symbology = 76 // GS1 DataBar Truncated (13 dígitos)
	GS1DataBarLim   Symbology = 77 // GS1 DataBar Limited (13 dígitos)
	GS1DataBarExp   Symbology = 78 // GS1 DataBar Expanded (2-255 chars)
	CODE128Auto     Symbology = 79 // CODE128 Auto (1-255 bytes)
)

const (
	Code128SetA Code128Set = 65 // Set A (ASCII 0-95)
	Code128SetB Code128Set = 66 // Set B (ASCII 32-127)
	Code128SetC Code128Set = 67 // Set C (pares numéricos 00-99)
)

// ============================================================================
// Variables de error
// ============================================================================

var (
	ErrHRIPosition      = errors.New("invalid HRI position (try 0-3 or '0'..'3')")
	ErrHRIFont          = errors.New("invalid HRI font (try 0-4, '0'..'4', 97, or 98")
	ErrHeight           = errors.New("invalid barcode height (try 1-255)")
	ErrWidth            = errors.New("invalid barcode width (try 2-6 or 68-76 for extended)")
	ErrSymbology        = errors.New("invalid barcode symbology")
	ErrDataTooShort     = errors.New("barcode data too short")
	ErrDataTooLong      = errors.New("barcode data too long (max 255 bytes)")
	ErrOddITFLength     = errors.New("ITF barcode requires even number of digits")
	ErrCode128Set       = errors.New("invalid CODE128 code set (try 65-67)")
	ErrCode128NoCodeSet = errors.New("CODE128 requires code set specification")
)

// ============================================================================
// Implementación principal
// ============================================================================

// Commands implementa la interfaz Capability para comandos de código de barras
type Commands struct{}

// NewCommands crea una nueva instancia de Commands
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Definición de interfaces
// ============================================================================

// Comprobación en tiempo de compilación de que Commands implementa Capability
var _ Capability = (*Commands)(nil)

// Capability agrupa las capacidades relacionadas con códigos de barras
type Capability interface {
	// Ajustes HRI
	SelectHRICharacterPosition(position HRIPosition) ([]byte, error)
	SelectFontForHRI(font HRIFont) ([]byte, error)

	// Dimensiones
	SetBarcodeHeight(height Height) ([]byte, error)
	SetBarcodeWidth(width Width) ([]byte, error)

	// Impresión
	PrintBarcode(symbology Symbology, data []byte) ([]byte, error)
	PrintBarcodeWithCodeSet(symbology Symbology, codeSet Code128Set, data []byte) ([]byte, error)
}

// ============================================================================
// Funciones auxiliares
// ============================================================================

// buildFunctionA construye comando Function A (terminado en NUL)
func (c *Commands) buildFunctionA(symbology Symbology, data []byte) ([]byte, error) {
	// Validaciones básicas para simbologías Function A
	if symbology == ITF {
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Construir comando: GS k m datos... NUL
	cmd := []byte{common.GS, 'k', byte(symbology)}
	cmd = append(cmd, data...)
	cmd = append(cmd, common.NUL)
	return cmd, nil
}

// buildFunctionB construye comando Function B (prefijo de longitud)
func (c *Commands) buildFunctionB(symbology Symbology, data []byte) ([]byte, error) {
	// Validar longitud de datos (máximo 255 para longitud de un byte)
	if len(data) > 255 {
		return nil, ErrDataTooLong
	}

	// Validación especial para ciertas simbologías
	switch symbology {
	case CODE128, GS1128:
		// Verificar si los datos tienen el prefijo de conjunto de código requerido
		if len(data) < 2 || data[0] != '{' ||
			data[1] < byte(Code128SetA) || data[1] > byte(Code128SetC) {
			return nil, ErrCode128NoCodeSet
		}
	case ITFB:
		// ITF requiere número par de dígitos
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Construir comando: GS k m n datos...
	cmd := []byte{common.GS, 'k', byte(symbology), byte(len(data))}
	cmd = append(cmd, data...)
	return cmd, nil
}

// ============================================================================
// Funciones de Utilidad para Validación
// ============================================================================

// ValidateNumericData evalúa si todos los bytes son dígitos numéricos
func ValidateNumericData(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// ValidateCode39Data evalúa si todos los bytes son válidos para CODE39
func ValidateCode39Data(data []byte) bool {
	for _, b := range data {
		switch {
		case b >= '0' && b <= '9':
		case b >= 'A' && b <= 'Z':
		case b == ' ' || b == '$' || b == '%' || b == '*' ||
			b == '+' || b == '-' || b == '.' || b == '/':
		default:
			return false
		}
	}
	return true
}

// ValidateCodabarData evalua si los datos de CODABAR son válidos
func ValidateCodabarData(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	// Verificar carácter de inicio
	start := data[0]
	if (start < 'A' || start > 'D') && (start < 'a' || start > 'd') {
		return false
	}
	// Verificar carácter de fin
	stop := data[len(data)-1]
	if (stop < 'A' || stop > 'D') && (stop < 'a' || stop > 'd') {
		return false
	}
	return true
}
