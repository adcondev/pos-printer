package encoding

import (
	"fmt"

	"github.com/AdConDev/pos-printer/types"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// CharacterSetData representa un conjunto de caracteres con su codificación
type CharacterSetData struct {
	EscPos   int               // Código numérico del charset (ej: 0, 2, 3)
	Name     string            // Nombre descriptivo (ej: "CP437", "CP850")
	Desc     string            // Descripción del charset (opcional)
	Encoding encoding.Encoding // Codificación real de golang.org/x/text
}

// TODO: Generalizar encodings para los diferentes protocolos de impresoras

// Registry contiene todos los character sets disponibles.
// Numeración "típica" (pero no garantizada universalmente)
var Registry = map[types.CharacterSet]*CharacterSetData{
	types.CP437: {
		EscPos:   0,
		Name:     "CP437",
		Desc:     "Inglés/EE. UU. y símbolos gráficos DOS",
		Encoding: charmap.CodePage437,
	},
	types.Katakana: {
		EscPos:   1,
		Name:     "Katakana",
		Desc:     "Japonés",
		Encoding: japanese.ISO2022JP, // CP932 es común para Katakana
	},
	types.CP850: {
		EscPos:   2,
		Name:     "CP850",
		Desc:     "Europa Occidental (Latin-1)",
		Encoding: charmap.CodePage850,
	},
	types.CP860: {
		EscPos:   3,
		Name:     "CP860",
		Desc:     "Portugués (Portugal)",
		Encoding: charmap.CodePage860,
	},
	types.CP863: {
		EscPos:   4,
		Name:     "CP863",
		Desc:     "Francés canadiense",
		Encoding: charmap.CodePage863,
	},
	types.CP865: {
		EscPos:   5,
		Name:     "CP865",
		Desc:     "Nórdico (escandinavo)",
		Encoding: charmap.CodePage865,
	},
	types.WestEurope: {
		EscPos:   6,
		Name:     "ISO8859-1",
		Desc:     "Europa Central y del Este",
		Encoding: charmap.ISO8859_1,
	},
	types.WCP1252: {
		EscPos:   16,
		Name:     "WPC1252",
		Desc:     "Windows Europa Occidental",
		Encoding: charmap.Windows1252,
	},
	types.CP866: {
		EscPos:   17,
		Name:     "CP866",
		Desc:     "Cirílico (Ruso MS-DOS)",
		Encoding: charmap.CodePage866,
	},
	types.CP852: {
		EscPos:   18,
		Name:     "CP852",
		Desc:     "Europa Central (Latin-2)",
		Encoding: charmap.CodePage852,
	},
	types.CP858: {
		EscPos:   19,
		Name:     "CP858",
		Encoding: charmap.CodePage858,
	},
	// Agregar más según necesites
	// IMPORTANTE: No existe un estándar universal obligatorio para la numeración de
	// tablas de codificación (code pages) en impresoras térmicas.
}

// GetEncoder devuelve un encoder para el charset especificado
func GetEncoder(charsetCode types.CharacterSet) (*encoding.Encoder, error) {
	cs, ok := Registry[charsetCode]
	if !ok {
		return nil, fmt.Errorf("charset not supported in registry: %d", charsetCode)
	}

	return cs.Encoding.NewEncoder(), nil
}

// EncodeString codifica un string usando el charset especificado
func EncodeString(str string, charsetCode types.CharacterSet) ([]byte, error) {
	encoder, err := GetEncoder(charsetCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoder for charset %d: %w", charsetCode, err)
	}
	return encoder.Bytes([]byte(str))
}
