package main

import (
	"log"
	"strings"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

func main() {
	// === Configuration ===
	printerName := "58mm PT-210" // Change this to your printer name

	// === Create connector ===
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

	// === Create printer profile ===
	prof := profile.CreatePt210()

	// === Create generic printer ===
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

	// === Initialize printer ===
	log.Println("Initializing printer...")
	if err = printer.Initialize(); err != nil {
		log.Printf("Error initializing: %v", err)
	}

	// === Define ALL user-defined characters at once ===
	log.Println("Defining all ASCII characters as black blocks...")
	if err := defineAllCharactersAsBlackBlocks(conn); err != nil {
		log.Printf("Error defining characters: %v", err)
	}

	// === Print header ===
	if err = printer.SetJustification(escpos.AlignCenter); err != nil {
		log.Printf("Error setting alignment: %v", err)
	}

	if err = printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		log.Printf("Error setting emphasis: %v", err)
	}

	if err = printer.TextLn("CARACTERES ESPECIALES"); err != nil {
		log.Printf("Error printing title: %v", err)
	}

	if err = printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		log.Printf("Error turning off emphasis: %v", err)
	}

	if err = printer.SetJustification(escpos.AlignLeft); err != nil {
		log.Printf("Error setting alignment: %v", err)
	}

	// Map of special characters to ASCII replacements
	specialChars := map[rune]byte{
		'á': '{', 'é': '|', 'í': '}', 'ó': '~', 'ú': '@',
		'ü': '#', 'ñ': '$', 'Á': '%', 'É': '^', 'Í': '&',
		'Ó': '*', 'Ú': '(', 'Ü': ')', 'Ñ': '_',
	}

	// Test sentences with accented characters
	testSentences := []string{
		"¡Hola! ¿Cómo estás?",
		// "El niño jugó fútbol en la montaña",
		// "José y María celebran el año nuevo",
		// "El pingüino pequeño nada rápido",
		// "Él comió jamón añejo con café",
	}

	// Print each test sentence
	for _, sentence := range testSentences {
		if err := printWithSpecialChars(printer, conn, sentence, specialChars); err != nil {
			log.Printf("Error printing: %v", err)
		}
	}

	// === Feed and cut ===
	if err = printer.Feed(2); err != nil {
		log.Printf("Error feeding paper: %v", err)
	}

	if err = printer.Cut(escpos.PartialCut); err != nil {
		log.Printf("Error cutting: %v", err)
	}

	log.Println("Test completed!")
}

// defineAllCharactersAsBlackBlocks defines all printable ASCII characters as black blocks
func defineAllCharactersAsBlackBlocks(conn connector.Connector) error {
	const (
		startChar = 0x20 // Space
		endChar   = 0x7E // Tilde
		height    = 3    // y = 3 (24 dots height)
		width     = 12   // x = 12 dots width
	)

	// Build command header
	cmd := []byte{
		0x1B, 0x26, // ESC &
		height,    // y - height
		startChar, // c1 - starting character
		endChar,   // c2 - ending character
	}

	// For each character, add width and data (solid black blocks)
	for char := startChar; char <= endChar; char++ {
		// Add width
		cmd = append(cmd, width)

		// Add data: all 0xFF for solid black (height * width bytes)
		for i := 0; i < height*width; i++ {
			cmd = append(cmd, 0xFF)
		}
	}

	_, err := conn.Write(cmd)
	return err
}

// printWithSpecialChars prints text with special characters replaced by black boxes
func printWithSpecialChars(printer *pos.EscposPrinter, conn connector.Connector, text string, specialCharsMap map[rune]byte) error {
	segments := splitTextBySpecialChars(text, specialCharsMap)

	for _, segment := range segments {
		if len(segment) == 1 {
			r := []rune(segment)[0]
			if replacement, isSpecial := specialCharsMap[r]; isSpecial {
				// Special character - activate UDC, print replacement, deactivate UDC
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
				// Regular character
				if err := printer.Print(segment); err != nil {
					return err
				}
			}
		} else {
			// Regular text segment
			if err := printer.Print(segment); err != nil {
				return err
			}
		}
	}

	// Add newline
	return printer.TextLn("")
}

// splitTextBySpecialChars splits text into segments where special characters are their own segments
func splitTextBySpecialChars(text string, specialChars map[rune]byte) []string {
	var segments []string
	var currentSegment strings.Builder

	for _, char := range text {
		if _, isSpecial := specialChars[char]; isSpecial {
			// End current segment
			if currentSegment.Len() > 0 {
				segments = append(segments, currentSegment.String())
				currentSegment.Reset()
			}
			// Add special character as its own segment
			segments = append(segments, string(char))
		} else {
			// Add to current segment
			currentSegment.WriteRune(char)
		}
	}

	// Add any remaining text
	if currentSegment.Len() > 0 {
		segments = append(segments, currentSegment.String())
	}

	return segments
}

// activateUserDefinedChars activates or deactivates user-defined character set
func activateUserDefinedChars(conn connector.Connector, enable bool) error {
	cmd := []byte{0x1B, 0x25} // ESC %

	if enable {
		cmd = append(cmd, 0x01) // Enable
	} else {
		cmd = append(cmd, 0x00) // Disable
	}

	_, err := conn.Write(cmd)
	return err
}
