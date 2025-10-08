// Package main demuestra cómo usar caracteres definidos por el usuario (UDC) para imprimir caracteres especiales
package main

import (
	"log"
	"strings"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// defineAllCharactersAsBlackBlocks define todos los caracteres ASCII imprimibles como bloques negros
func defineAllCharactersAsBlackBlocks(conn connector.Connector) error {
	const (
		startChar = 0x20 // Espacio
		endChar   = 0x7E // Tilde
		height    = 3    // y = 3 (24 puntos de altura)
		width     = 12   // x = 12 puntos de ancho
	)

	// Construir encabezado del comando
	cmd := []byte{
		0x1B, 0x26, // ESC &
		height,    // y - altura
		startChar, // c1 - carácter inicial
		endChar,   // c2 - carácter final
	}

	// Para cada carácter, agregar ancho y datos (bloques negros sólidos)
	for char := startChar; char <= endChar; char++ {
		cmd = append(cmd, width)

		// Agregar datos: todo 0xFF para negro sólido (altura * ancho bytes)
		for i := 0; i < height*width; i++ {
			cmd = append(cmd, 0xFF)
		}
	}

	_, err := conn.Write(cmd)
	return err
}

// printHeader imprime el encabezado del documento de prueba
func printHeader(printer *pos.EscposPrinter) error {
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := printer.TextLn("CARACTERES ESPECIALES"); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}
	return printer.SetJustification(escpos.AlignLeft)
}

// printWithSpecialChars imprime texto con caracteres especiales reemplazados por bloques negros
func printWithSpecialChars(printer *pos.EscposPrinter, conn connector.Connector, text string, specialCharsMap map[rune]byte) error {
	segments := splitTextBySpecialChars(text, specialCharsMap)

	for _, segment := range segments {
		if len(segment) == 1 {
			r := []rune(segment)[0]
			if replacement, isSpecial := specialCharsMap[r]; isSpecial {
				// Carácter especial - activar UDC, imprimir reemplazo, desactivar UDC
				if err := activateUserDefinedChars(conn, true); err != nil {
					return err
				}
				if err := printer.Print(string(replacement)); err != nil {
					return err
				}
				if err := activateUserDefinedChars(conn, false); err != nil {
					return err
				}
			} else {
				// Carácter regular
				if err := printer.Print(segment); err != nil {
					return err
				}
			}
		} else {
			// Segmento de texto regular
			if err := printer.Print(segment); err != nil {
				return err
			}
		}
	}

	return printer.TextLn("")
}

// splitTextBySpecialChars divide el texto en segmentos donde los caracteres especiales son segmentos propios
func splitTextBySpecialChars(text string, specialChars map[rune]byte) []string {
	var segments []string
	var currentSegment strings.Builder

	for _, char := range text {
		if _, isSpecial := specialChars[char]; isSpecial {
			// Finalizar segmento actual
			if currentSegment.Len() > 0 {
				segments = append(segments, currentSegment.String())
				currentSegment.Reset()
			}
			// Agregar carácter especial como su propio segmento
			segments = append(segments, string(char))
		} else {
			// Agregar al segmento actual
			currentSegment.WriteRune(char)
		}
	}

	// Agregar cualquier texto restante
	if currentSegment.Len() > 0 {
		segments = append(segments, currentSegment.String())
	}

	return segments
}

// activateUserDefinedChars activa o desactiva el conjunto de caracteres definidos por el usuario
func activateUserDefinedChars(conn connector.Connector, enable bool) error {
	cmd := []byte{0x1B, 0x25} // ESC %

	if enable {
		cmd = append(cmd, 0x01) // Habilitar
	} else {
		cmd = append(cmd, 0x00) // Deshabilitar
	}

	_, err := conn.Write(cmd)
	return err
}

// finishPrinting alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) error {
	if err := printer.Feed(2); err != nil {
		return err
	}
	return printer.Cut(escpos.PartialCut)
}

func main() {
	// Configuración
	printerName := "58mm PT-210"

	// Crear conector
	log.Printf("Attempting to connect to printer: %s", printerName)
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Printf("Error creating connector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("error closing connector: %v", err)
		}
	}(conn)

	// Crear perfil de impresora
	prof := profile.CreatePt210()

	// Crear impresora genérica
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error creating printer: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("error closing connector: %v", err)
		}
	}(printer)

	// Inicializar impresora
	log.Println("Initializing printer...")
	if err = printer.Initialize(); err != nil {
		log.Printf("Error initializing: %v", err)
	}

	// Definir TODOS los caracteres definidos por el usuario de una vez
	log.Println("Defining all ASCII characters as black blocks...")
	if err := defineAllCharactersAsBlackBlocks(conn); err != nil {
		log.Printf("Error defining characters: %v", err)
	}

	// Imprimir encabezado
	if err := printHeader(printer); err != nil {
		log.Printf("Error printing header: %v", err)
	}

	// Mapa de caracteres especiales a reemplazos ASCII
	specialChars := map[rune]byte{
		'á': '{', 'é': '|', 'í': '}', 'ó': '~', 'ú': '@',
		'ü': '#', 'ñ': '$', 'Á': '%', 'É': '^', 'Í': '&',
		'Ó': '*', 'Ú': '(', 'Ü': ')', 'Ñ': '_',
	}

	// Frases de prueba con caracteres acentuados
	testSentences := []string{
		"¡Hola! ¿Cómo estás?",
	}

	// Imprimir cada frase de prueba
	for _, sentence := range testSentences {
		if err := printWithSpecialChars(printer, conn, sentence, specialChars); err != nil {
			log.Printf("Error printing: %v", err)
		}
	}

	// Alimentar y cortar
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error finishing: %v", err)
	}

	log.Println("Test completed!")
}
