// Package main demuestra cómo definir y probar todos los caracteres definidos por el usuario (UDC) de una vez
package main

import (
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// defineAllCharactersAsBlackBlocks define todos los caracteres ASCII imprimibles (0x20-0x7E) como bloques negros en un solo comando
func defineAllCharactersAsBlackBlocks(conn connector.Connector) error {
	const (
		startChar = 0x20 // Espacio
		endChar   = 0x7E // Tilde
		height    = 3    // y = 3 (3 * 8 = 24 puntos de altura)
		width     = 12   // x = 12 puntos de ancho para Fuente A
	)

	numChars := endChar - startChar + 1
	dataPerChar := height * width // 36 bytes por carácter

	// Construir el comando
	cmd := make([]byte, 0, 5+numChars*(1+dataPerChar))

	// Encabezado del comando: ESC & y c1 c2
	cmd = append(cmd, 0x1B)      // ESC
	cmd = append(cmd, 0x26)      // &
	cmd = append(cmd, height)    // y - altura (3 * 8 = 24 puntos)
	cmd = append(cmd, startChar) // c1 - carácter inicial
	cmd = append(cmd, endChar)   // c2 - carácter final

	// Para cada carácter de c1 a c2, agregar ancho y datos
	for char := startChar; char <= endChar; char++ {
		cmd = append(cmd, width)

		// Datos: altura * ancho bytes, todos 0xFF para negro sólido
		for i := 0; i < dataPerChar; i++ {
			cmd = append(cmd, 0xFF)
		}
	}

	log.Printf("Sending UDC definition command (%d bytes total)", len(cmd))
	_, err := conn.Write(cmd)
	return err
}

// printTestHeader imprime el encabezado del test
func printTestHeader(printer *pos.EscposPrinter) error {
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		return fmt.Errorf("error centering: %w", err)
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return fmt.Errorf("error setting bold: %w", err)
	}
	if err := printer.TextLn("USER DEFINED CHAR TEST"); err != nil {
		return fmt.Errorf("error printing title: %w", err)
	}
	if err := printer.TextLn("ALL ASCII 0x20-0x7E"); err != nil {
		return fmt.Errorf("error printing subtitle: %w", err)
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return fmt.Errorf("error turning off bold: %w", err)
	}
	if err := printer.TextLn("================================"); err != nil {
		return fmt.Errorf("error printing line: %w", err)
	}
	return printer.SetJustification(escpos.AlignLeft)
}

// printAllCharactersAsBlocks imprime todos los caracteres como bloques
func printAllCharactersAsBlocks(printer *pos.EscposPrinter, conn connector.Connector) error {
	if err := printer.TextLn(""); err != nil {
		return err
	}
	if err := printer.TextLn("================================"); err != nil {
		return err
	}
	if err := printer.TextLn("ALL CHARS AS BLOCKS:"); err != nil {
		return err
	}

	// Activar caracteres definidos por el usuario
	if err := activateUserDefinedChars(conn, true); err != nil {
		return fmt.Errorf("error activating UDC: %w", err)
	}

	// Imprimir todos los caracteres en filas
	for start := byte(0x20); start <= 0x7E; start += 16 {
		end := start + 15
		if end > 0x7E {
			end = 0x7E
		}

		line := ""
		for char := start; char <= end; char++ {
			line += string(char)
		}
		if err := printer.TextLn(line); err != nil {
			return fmt.Errorf("error printing line: %w", err)
		}
	}

	// Desactivar caracteres definidos por el usuario
	return activateUserDefinedChars(conn, false)
}

// testBatch prueba un lote de caracteres
func testBatch(printer *pos.EscposPrinter, conn connector.Connector, batchName string, start, end byte) error {
	// Imprimir encabezado del lote
	if err := printer.TextLn(""); err != nil {
		return err
	}
	if err := printer.TextLn(fmt.Sprintf("--- %s (0x%02X-0x%02X) ---", batchName, start, end)); err != nil {
		return err
	}

	// Imprimir caracteres originales
	line := "Orig: "
	for char := start; char <= end; char++ {
		line += string(char) + " "
	}
	if err := printer.TextLn(line); err != nil {
		return err
	}

	// Activar caracteres definidos por el usuario
	if err := activateUserDefinedChars(conn, true); err != nil {
		return err
	}

	// Imprimir con caracteres definidos por el usuario (deben ser bloques negros)
	line = "UDef: "
	for char := start; char <= end; char++ {
		line += string(char) + " "
	}
	if err := printer.TextLn(line); err != nil {
		return err
	}

	// Desactivar caracteres definidos por el usuario
	return activateUserDefinedChars(conn, false)
}

// runUDCTest ejecuta la prueba real de caracteres definidos por el usuario
func runUDCTest(printer *pos.EscposPrinter, conn connector.Connector) error {
	// Imprimir encabezado
	if err := printTestHeader(printer); err != nil {
		return err
	}

	// Lotes de prueba
	testBatches := []struct {
		name  string
		start byte
		end   byte
	}{
		{"Special Chars 1", 0x20, 0x2F}, // Espacio a /
		{"Numbers", 0x30, 0x39},         // 0-9
		{"Special Chars 2", 0x3A, 0x40}, // : a @
		{"Uppercase A-P", 0x41, 0x50},   // A-P
		{"Uppercase Q-Z", 0x51, 0x5A},   // Q-Z
		{"Special Chars 3", 0x5B, 0x60}, // [ a `
		{"Lowercase a-p", 0x61, 0x70},   // a-p
		{"Lowercase q-z", 0x71, 0x7A},   // q-z
		{"Special Chars 4", 0x7B, 0x7E}, // { a ~
	}

	for _, batch := range testBatches {
		if err := testBatch(printer, conn, batch.name, batch.start, batch.end); err != nil {
			return fmt.Errorf("error testing batch %s: %w", batch.name, err)
		}
	}

	// Prueba final: Mostrar todos los caracteres como bloques
	return printAllCharactersAsBlocks(printer, conn)
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
	// ========== Configuración ==========
	printerName := "58mm PT-210"

	// ========== Crear conector ==========
	log.Printf("Attempting to connect to printer: %s", printerName)
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error creating connector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("error closing connector: %v", err)
		}
	}(conn)

	// ========== Crear perfil de impresora ==========
	prof := profile.CreatePt210()

	// ========== Crear impresora genérica ==========
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

	// ========== Inicializar impresora ==========
	log.Println("Initializing printer...")
	if err = printer.Initialize(); err != nil {
		log.Printf("Error initializing: %v", err)
	}

	// ========== Definir TODOS los caracteres definidos por el usuario de una vez ==========
	log.Println("Defining all characters (0x20-0x7E) as black blocks...")
	if err := defineAllCharactersAsBlackBlocks(conn); err != nil {
		log.Printf("Error defining all characters: %v", err)
	}

	// ========== Ejecutar el test ==========
	if err := runUDCTest(printer, conn); err != nil {
		log.Printf("Error running test: %v", err)
	}

	// ========== Alimentar y cortar ==========
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error finishing: %v", err)
	}

	log.Println("Test completed!")
}
