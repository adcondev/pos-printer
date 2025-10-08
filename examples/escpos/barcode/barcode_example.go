package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// printHeader imprime el encabezado del documento
func printHeader(printer *pos.EscposPrinter) {
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		log.Printf("Error setting alignment: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		log.Printf("Error setting emphasis: %v", err)
	}
	if err := printer.TextLn("BARCODE TEST"); err != nil {
		log.Printf("Error printing title: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		log.Printf("Error turning off emphasis: %v", err)
	}
	if err := printer.TextLn("================================"); err != nil {
		log.Printf("Error printing line: %v", err)
	}
	if err := printer.SetJustification(escpos.AlignLeft); err != nil {
		log.Printf("Error setting alignment: %v", err)
	}
}

// printBarcodeExample imprime un código de barras con su etiqueta
func printBarcodeExample(printer *pos.EscposPrinter, label string, symbology barcode.Symbology, data string) {
	if err := printer.TextLn(label); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.PrintBarcode(symbology, data); err != nil {
		log.Printf("Error printing %s: %v", label, err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
}

// printCode128Example imprime un código CODE128 con posición HRI personalizada
func printCode128Example(printer *pos.EscposPrinter) {
	if err := printer.TextLn("CODE128 (HRI Below):"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetBarcodeHRIPosition(barcode.HRIBelow); err != nil {
		log.Printf("Error setting HRI position: %v", err)
	}
	if err := printer.PrintBarcode(barcode.CODE128, "Hello123"); err != nil {
		log.Printf("Error printing CODE128: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
}

// printCustomSizeBarcode imprime un código de barras con dimensiones personalizadas
func printCustomSizeBarcode(printer *pos.EscposPrinter) {
	if err := printer.TextLn("Custom Size Barcode:"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetBarcodeHeight(100); err != nil {
		log.Printf("Error setting barcode height: %v", err)
	}
	if err := printer.SetBarcodeWidth(3); err != nil {
		log.Printf("Error setting barcode width: %v", err)
	}
	if err := printer.PrintBarcode(barcode.EAN8, "12345670"); err != nil {
		log.Printf("Error printing EAN8: %v", err)
	}
}

// finishPrinting alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) {
	if err := printer.Feed(3); err != nil {
		log.Printf("Error feeding paper: %v", err)
	}
	if err := printer.Cut(escpos.PartialCut); err != nil {
		log.Printf("Error cutting: %v", err)
	}
}

func main() {
	// ========== Configuración de la impresora ==========
	printerName := "58mm PT-210"

	// ========== Crear conector de impresora ==========
	log.Printf("Attempting to connect to printer: %s", printerName)
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error creating connector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connector: %v", err)
		}
	}(conn)

	// ========== Crear perfil de impresora ==========
	prof := profile.CreatePt210()

	// ========== Crear instancia de impresora ==========
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Panicf("Error creating printer: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error closing printer: %v", err)
		}
	}(printer)

	// ========== Inicializar impresora ==========
	log.Println("Initializing printer...")
	if err = printer.Initialize(); err != nil {
		log.Panicf("Error initializing: %v", err)
	}

	// ========== Imprimir encabezado ==========
	printHeader(printer)

	// ========== Probar diferentes tipos de códigos de barras ==========
	printBarcodeExample(printer, "UPC-A:", barcode.UPCA, "12345678901")
	printBarcodeExample(printer, "EAN-13:", barcode.EAN13, "4901234567890")
	printBarcodeExample(printer, "CODE39:", barcode.CODE39, "ABC-1234")

	// ========== Código CODE128 con posición HRI personalizada ==========
	printCode128Example(printer)

	// ========== Probar código de barras con dimensiones personalizadas ==========
	printCustomSizeBarcode(printer)

	// ========== Alimentar papel y cortar ==========
	finishPrinting(printer)

	log.Println("Barcode test completed!")
}
