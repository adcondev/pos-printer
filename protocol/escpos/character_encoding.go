package escpos

import (
	"fmt"
	"log"

	"github.com/AdConDev/pos-printer/encoding"
	"github.com/AdConDev/pos-printer/types"
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

// CancelKanjiMode cancela el modo de caracteres Kanji.
//
// Formato:
//
//	ASCII: FS .
//	Hex:   1C 2E
//	Decimal: 28 46
//
// Descripción:
//
//	Deshabilita el modo de caracteres Kanji en la impresora.
//
// Referencia:
//
//	FS &, FS C
func (p *Commands) CancelKanjiMode() []byte {
	return []byte{FS, '.'}
}

func (p *Commands) SelectKanjiMode() []byte {
	return []byte{FS, '&'}
}

// TODO: Hacer trabajo similar al de las codepages.

// SelectInternationalCharacterSet configura el conjunto de caracteres internacionales de la impresora.
//
// Formato:
//
//	ASCII: ESC R n
//	Hex:   1B R n
//	Decimal: 27 R n
//
// Rango:
//
//	0 ≤ n ≤ 15
//
// Descripción:
//
//	Selecciona un conjunto de caracteres internacionales basado en el valor de n, según la siguiente tabla:
//	  n   Conjunto de caracteres
//	  0   Estados Unidos
//	  1   Francia
//	  2   Alemania
//	  3   Reino Unido
//	  4   Dinamarca
//	  5   Suecia
//	  6   Italia
//	  7   España
//	  8   Japón
//	  9   Noruega
//	  10  Dinamarca
//	  11  España
//	  12  Latino
//	  13  Chino
//	  14  Corea
//	  15  Eslovenia/Croacia
//
// Predeterminado:
//
//	Para el modelo Simplificado Chino: n = 15; Para modelos distintos del Simplificado Chino: n = 0
//
// Nota:
//
//	Los conjuntos de caracteres para Eslovenia/Croacia y China solo son compatibles con el modelo Simplificado Chino.
func SelectInternationalCharacterSet(n byte) []byte {
	cmd := []byte{ESC, 'R', n}
	return cmd
}

// CancelUserDefinedCharacters representa el comando ESC ? n para cancelar caracteres definidos por el usuario.
//
// Nombre:
//
//	Cancelar caracteres definidos por el usuario
//
// Formato:
//
//	ASCII: ESC ? n
//	Hex:   1B 3F n
//	Decimal: 27 63 n
//
// Rango:
//
//	32 ≤ n ≤ 126
//
// Descripción:
//
//	Cancela los caracteres definidos por el usuario. Después de cancelar, se imprime el patrón correspondiente
//	al carácter interno.
//
// Detalles:
//   - El comando elimina el patrón definido para el código de carácter especificado por n en la fuente seleccionada mediante ESC !.
//   - Si un carácter definido por el usuario no ha sido previamente definido para el código especificado, la impresora ignora este comando.
//
// Referencia:
//
//	ESC &, ESC %
func CancelUserDefinedCharacters(n types.UserDefinedChar) ([]byte, error) {
	if n < 32 || n > 126 {
		return nil, fmt.Errorf("n debe estar en el rango de 32 a 126, recibido: %d", n)
	}
	defChar := byte(n)
	cmd := []byte{ESC, '?', defChar}
	return cmd, nil
}

// DefineUserDefinedCharacters representa el comando para definir caracteres personalizados.
//
// Nombre:
//
//	Definir caracteres personalizados
//
// Formato:
//
//	ASCII: ESC & y c1 c2 [x1 d1...d(y × x1)] ... [xk d1...d(y × xk)]
//	Hex:   1B 26 y c1 c2 [x1 d1...d(y×x1)] ... [xk d1...d(y×xk)]
//
// Rango:
//   - y: Especifica el número de bytes en la dirección vertical (por ejemplo, y = 3).
//   - c1, c2: Rango de códigos de caracteres a definir; 32 (0x20) ≤ c1 ≤ c2 ≤ 126 (0x7E).
//   - x: Número de puntos en la dirección horizontal.
//     Para la Fuente A (12×24): 0 ≤ x ≤ 12.
//     Para la Fuente B (9×17): 0 ≤ x ≤ 9.
//   - d (datos de puntos): 0 ≤ d ≤ 255.
//   - La cantidad de datos para cada carácter es (y × x) bytes.
//
// Descripción:
//
//	Define caracteres personalizados utilizando un conjunto de datos que especifican el patrón de puntos de cada
//	carácter. Se pueden definir múltiples caracteres consecutivos asignados a códigos de carácter desde c1 hasta c2.
//	Para definir un único carácter, se utiliza c1 = c2.
//
// Detalles:
//   - y especifica el número de bytes en la dirección vertical.
//   - c1 indica el código del primer carácter a definir y c2, el código final.
//   - x define el número de puntos (dots) en la dirección horizontal para cada carácter.
//   - El patrón de puntos se configura de izquierda a derecha; los puntos restantes a la derecha se dejan en blanco.
//   - Cada carácter se define con (y × x) bytes, donde cada bit determina si se imprime (bit = 1) o no imprime (bit = 0) un punto.
//   - Es posible definir diferentes patrones de caracteres para cada fuente. Para seleccionar la fuente deseada, utilice el comando ESC !
//   - No es posible definir un carácter personalizado y una imagen de bits descargada simultáneamente. La ejecución de este comando
//     borrará cualquier imagen de bits descargada.
//   - La definición de caracteres personalizados se borra cuando se ejecuta cualquiera de los siguientes comandos:
//     ① ESC @
//     ② ESC ?
//     ③ FS q
//     ④ GS *
//     ⑤ Al reiniciar la impresora o apagar la alimentación.
//   - Para caracteres definidos en la Fuente B (9×17), únicamente es efectivo el bit más significativo del tercer byte en la dirección vertical.
//   - Al cancelar la definición personalizada mediante otros comandos (por ejemplo, ESC %, ESC ?), se selecciona el conjunto de caracteres interno.
//
// Referencia:
//
//	ESC %, ESC ?
func DefineUserDefinedCharacters(y, c1, c2 byte, data ...[]byte) []byte {
	cmd := []byte{ESC, '&', y, c1, c2}
	for _, d := range data {
		cmd = append(cmd, d...)
	}
	return cmd
}
