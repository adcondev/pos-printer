package main

import (
	"fmt"
	"log"

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
		log.Fatalf("Error creating connector: %v", err)
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
	log.Println("Defining all characters (0x20-0x7E) as black blocks...")
	if err := defineAllCharactersAsBlackBlocks(conn); err != nil {
		log.Printf("Error defining all characters: %v", err)
	}

	// === Run the test ===
	if err := runUDCTest(printer, conn); err != nil {
		log.Printf("Error running test: %v", err)
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

// defineAllCharactersAsBlackBlocks defines all printable ASCII characters (0x20-0x7E) as black blocks in a single command
func defineAllCharactersAsBlackBlocks(conn connector.Connector) error {
	const (
		startChar = 0x20 // Space
		endChar   = 0x7E // Tilde
		height    = 3    // y = 3 (3 * 8 = 24 dots height)
		width     = 12   // x = 12 dots width for Font A
	)

	numChars := endChar - startChar + 1
	dataPerChar := height * width // 36 bytes per character

	// Build the command
	cmd := make([]byte, 0, 5+numChars*(1+dataPerChar))

	// Command header: ESC & y c1 c2
	cmd = append(cmd, 0x1B)      // ESC
	cmd = append(cmd, 0x26)      // &
	cmd = append(cmd, height)    // y - height (3 * 8 = 24 dots)
	cmd = append(cmd, startChar) // c1 - starting character
	cmd = append(cmd, endChar)   // c2 - ending character

	// For each character from c1 to c2, append width and data
	for char := startChar; char <= endChar; char++ {
		// Width for this character
		cmd = append(cmd, width)

		// Data: height * width bytes, all 0xFF for solid black
		for i := 0; i < dataPerChar; i++ {
			cmd = append(cmd, 0xFF)
		}
	}

	log.Printf("Sending UDC definition command (%d bytes total)", len(cmd))
	_, err := conn.Write(cmd)
	return err
}

// runUDCTest performs the actual test of user-defined characters
func runUDCTest(printer *pos.EscposPrinter, conn connector.Connector) error {
	// === Print header ===
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

	if err := printer.SetJustification(escpos.AlignLeft); err != nil {
		return fmt.Errorf("error aligning left: %w", err)
	}

	// === Test batches ===
	testBatches := []struct {
		name  string
		start byte
		end   byte
	}{
		{"Special Chars 1", 0x20, 0x2F}, // Space to /
		{"Numbers", 0x30, 0x39},         // 0-9
		{"Special Chars 2", 0x3A, 0x40}, // : to @
		{"Uppercase A-P", 0x41, 0x50},   // A-P
		{"Uppercase Q-Z", 0x51, 0x5A},   // Q-Z
		{"Special Chars 3", 0x5B, 0x60}, // [ to `
		{"Lowercase a-p", 0x61, 0x70},   // a-p
		{"Lowercase q-z", 0x71, 0x7A},   // q-z
		{"Special Chars 4", 0x7B, 0x7E}, // { to ~
	}

	for _, batch := range testBatches {
		if err := testBatch(printer, conn, batch.name, batch.start, batch.end); err != nil {
			return fmt.Errorf("error testing batch %s: %w", batch.name, err)
		}
	}

	// === Final test: Show all characters as blocks ===
	if err := printer.TextLn(""); err != nil {
		return err
	}

	if err := printer.TextLn("================================"); err != nil {
		return err
	}

	if err := printer.TextLn("ALL CHARS AS BLOCKS:"); err != nil {
		return err
	}

	// Activate user-defined characters
	if err := activateUserDefinedChars(conn, true); err != nil {
		return fmt.Errorf("error activating UDC: %w", err)
	}

	// Print all characters in rows
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

	// Deactivate user-defined characters
	if err := activateUserDefinedChars(conn, false); err != nil {
		return fmt.Errorf("error deactivating UDC: %w", err)
	}

	return nil
}

// testBatch tests a batch of characters
func testBatch(printer *pos.EscposPrinter, conn connector.Connector, batchName string, start, end byte) error {
	// Print batch header
	if err := printer.TextLn(""); err != nil {
		return err
	}

	if err := printer.TextLn(fmt.Sprintf("--- %s (0x%02X-0x%02X) ---", batchName, start, end)); err != nil {
		return err
	}

	// Print original characters
	line := "Orig: "
	for char := start; char <= end; char++ {
		line += string(char) + " "
	}
	if err := printer.TextLn(line); err != nil {
		return err
	}

	// Activate user-defined characters
	if err := activateUserDefinedChars(conn, true); err != nil {
		return err
	}

	// Print with user-defined characters (should be black blocks)
	line = "UDef: "
	for char := start; char <= end; char++ {
		line += string(char) + " "
	}
	if err := printer.TextLn(line); err != nil {
		return err
	}

	// Deactivate user-defined characters
	if err := activateUserDefinedChars(conn, false); err != nil {
		return err
	}

	return nil
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
