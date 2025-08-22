package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/encoding"
)

// TODO: Comandos para manejo de codificación de caracteres
// - Código de página
// - Caracteres internacionales
// - Caracteres especiales

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

// SelectCharacterTable selecciona el conjunto de caracteres para la impresora
func (c *Commands) SelectCharacterTable(table encoding.CharacterSet) ([]byte, error) {
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
	cmd := []byte{ESC, 't', charSet}

	return cmd, nil
}

// CancelKanjiMode deactivates Kanji mode
func (c *Commands) CancelKanjiMode() []byte {
	return []byte{FS, '.'}
}

// SelectKanjiMode activates Kanji mode
func (c *Commands) SelectKanjiMode() []byte {
	return []byte{FS, '&'}
}

// FIXME: Hacer trabajo similar al de las codepages.

// SelectInternationalCharacterSet define los conjuntos de caracteres internacionales
func SelectInternationalCharacterSet(n byte) []byte {
	cmd := []byte{ESC, 'R', n}
	return cmd
}

// CancelUserDefinedCharacters cancela la definición de un carácter definido por el usuario
func CancelUserDefinedCharacters(n UserDefinedChar) ([]byte, error) {
	if n < 32 || n > 126 {
		return nil, fmt.Errorf("n debe estar en el rango de 32 a 126, recibido: %d", n)
	}
	defChar := byte(n)
	cmd := []byte{ESC, '?', defChar}
	return cmd, nil
}

// DefineUserDefinedCharacters define uno o más caracteres definidos por el usuario
func DefineUserDefinedCharacters(y, c1, c2 byte, data ...[]byte) []byte {
	cmd := []byte{ESC, '&', y, c1, c2}
	for _, d := range data {
		cmd = append(cmd, d...)
	}
	return cmd
}
