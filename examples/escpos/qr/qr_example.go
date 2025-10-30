// Package main demuestra cómo imprimir códigos QR usando una impresora POS.
package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// printQRHeader imprime el encabezado del documento QR
func printQRHeader(printer *pos.EscposPrinter) error {
	if err := printer.SetFont(character.FontType(escpos.FontA)); err != nil {
		return err
	}
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := printer.TextLn("PRUEBA DE QR á é í ó ú ñ"); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}
	return printer.Feed(1)
}

// printQRCode imprime el código QR con los parámetros especificados
func printQRCode(printer *pos.EscposPrinter) error {
	return printer.PrintQR(
		"https://github.com/adcondev/pos-printer",
		escpos.Model2,
		escpos.ECHigh,
		8,
		256,
	)
}

// finishPrinting alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) error {
	if err := printer.Feed(1); err != nil {
		return err
	}
	if err := printer.Cut(escpos.PartialCut); err != nil {
		return err
	}
	return printer.Feed(1)
}

func main() {
	// ========== Configuración de la impresora ==========
	printerName := "58mm PT-210"

	// ========== Crear conector de impresora ==========
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear conector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar conector de impresora: %v", err)
		}
	}(conn)

	// ========== Crear perfil de impresora ==========
	prof := profile.CreatePt210()

	// ========== Crear instancia de impresora ==========
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error al crear impresora: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// ========== Imprimir documento QR ==========
	if err := printQRHeader(printer); err != nil {
		log.Printf("Error al imprimir encabezado: %v", err)
	}
	if err := printQRCode(printer); err != nil {
		log.Printf("Error al imprimir código QR: %v", err)
	}
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error al finalizar impresión: %v", err)
	}

	log.Println("Impresión de QR completada exitosamente.")
}
