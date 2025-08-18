package escpos

import (
	"fmt"
	"log"

	"github.com/AdConDev/pos-printer/encoding"
	"github.com/AdConDev/pos-printer/protocol/escpos/types"
)

// TODO: Comandos para manejo de codificación de caracteres
// - Código de página
// - Caracteres internacionales
// - Caracteres especiales

type CodePage byte

const (
	// Tabla de códigos comunes en ESC/POS
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
)

const (
	WCP1252 CodePage = iota + 16 // WCP1252 Windows-1252
	CP866                        // CP866 Cyrillic #2
	CP852                        // CP852 Latin2
	CP858                        // CP858 Multilingual + Euro
	IranII                       // IranII (CP864)
	Latvian                      // Latvian (Windows-1257)
)

// TODO: Revisar si es mejor implementar el map que un rango de valores.

func (cp CodePage) IsValid() bool {
	return cp <= Latvian || (cp >= WCP1252 && cp <= Latvian)
}

func (p *Commands) SelectCharacterTable(table types.CharacterSet) []byte {
	charTable := CodePage(encoding.Registry[table].EscPos)
	// Validar que table esté en un rango válido
	if !charTable.IsValid() {
		// Log de advertencia si está fuera de rango
		log.Printf("advertencia: tabla de caracteres %d fuera de rango, usando 0 por defecto", table)
		charTable = 0 // Default a 0 si está fuera de rango
	}
	// ESC t n - Select character code table
	cmd := []byte{ESC, 't', byte(charTable)}

	return cmd
}

func (p *Commands) CancelKanjiMode() []byte {
	return []byte{FS, '.'}
}

func (p *Commands) SelectKanjiMode() []byte {
	return []byte{FS, '&'}
}

// TODO: Hacer trabajo similar al de las codepages.

func SelectInternationalCharacterSet(n byte) []byte {
	cmd := []byte{ESC, 'R', n}
	return cmd
}

func CancelUserDefinedCharacters(n types.UserDefinedChar) ([]byte, error) {
	if n < 32 || n > 126 {
		return nil, fmt.Errorf("n debe estar en el rango de 32 a 126, recibido: %d", n)
	}
	defChar := byte(n)
	cmd := []byte{ESC, '?', defChar}
	return cmd, nil
}

func DefineUserDefinedCharacters(y, c1, c2 byte, data ...[]byte) []byte {
	cmd := []byte{ESC, '&', y, c1, c2}
	for _, d := range data {
		cmd = append(cmd, d...)
	}
	return cmd
}
