package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/encoding"
	"github.com/adcondev/pos-printer/escpos/shared"
)

// ============================================================================
// Code page definitions and mappings
// ============================================================================
// Definimos CodePage y un mapa que convierte de nuestro enum a los valores
// numéricos esperados por ESC/POS.

// CodePage define los conjuntos de caracteres estándar
type CodePage byte

const (
	CP437      CodePage = iota // CP437 U.S.A. / Standard Europe
	Katakana                   // Katakana (JIS X 0201)
	CP850                      // CP850 Multilingual
	CP860                      // CP860 Portuguese
	CP863                      // CP863 Canadian French
	CP865                      // CP865 Nordic
	WestEurope                 // WestEurope (ISO-8859-1)
	Greek                      // Greek (ISO-8859-7)
	Hebrew                     // Hebrew (ISO-8859-8)
	CP755                      // CP755 East Europe (not directly supported)
	Iran                       // Iran (CP720 Arabic)
	WCP1252                    // WCP1252 Windows-1252
	CP866                      // CP866 Cyrillic #2
	CP852                      // CP852 Latin2
	CP858                      // CP858 Multilingual + Euro
	IranII                     // IranII (CP864)
	Latvian                    // Latvian (Windows-1257)
)

var codePageMap = map[CodePage]byte{
	CP437:      0,
	Katakana:   1,
	CP850:      2,
	CP860:      3,
	CP863:      4,
	CP865:      5,
	WestEurope: 6,
	Greek:      7,
	Hebrew:     8,
	CP755:      9, // Not directly supported
	Iran:       10,
	WCP1252:    16,
	CP866:      17,
	CP852:      18,
	CP858:      19,
	IranII:     20,
	Latvian:    21,
}

// ============================================================================
// Public API (implementation)
// ============================================================================

// SelectCharacterTable selecciona el conjunto de caracteres para la impresora
func (c *Protocol) SelectCharacterTable(table encoding.CharacterSet) ([]byte, error) {
	encoder, ok := encoding.Registry[table]
	if !ok {
		return nil, fmt.Errorf("error: conjunto de caracteres %s no soportado por implementacion", encoder.Name)
	}
	charTable := CodePage(encoder.EscPos)
	charSet, ok := codePageMap[charTable]
	if !ok {
		// Log de error si no se encuentra el código de página
		return nil, fmt.Errorf("error: código de página %s no soportado por protocolo", encoder.Name)
	}
	// ESC t n - Select character code table
	cmd := []byte{shared.ESC, 't', charSet}

	return cmd, nil
}

// CancelKanjiMode deactivates Kanji mode
func (c *Protocol) CancelKanjiMode() []byte {
	return []byte{shared.FS, '.'}
}

// SelectKanjiMode activates Kanji mode
func (c *Protocol) SelectKanjiMode() []byte {
	return []byte{shared.FS, '&'}
}
