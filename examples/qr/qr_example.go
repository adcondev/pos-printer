// Package main demonstrates how to print a QR code using a POS printer in Go.
package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

func main() {
	// === Crear conector ===
	// Seleccionar la impresora según tu configuración
	// printerName := "80mm EC-PM-80250"
	printerName := "58mm PT-210"

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	// prof := profile.CreateProfile80mm()
	prof := profile.CreatePt210() // Usar perfil de 58mm

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

	// === Crear impresora genérica ===
	// === Inicializar impresora ===
	printer, err := pos.NewEscposPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error al crear impresora: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// === Imprimir título ===
	if err := printer.SetFont(escpos.FontA); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphOn); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("PRUEBA DE QR á é í ó ú ñ"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphOff); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	// === Imprimir QR Code ===
	if err := printer.PrintQR(
		"https://github.com/adcondev/pos-printer", // Contenido del QR Code
		escpos.Model2, // Modelo de QR Code (Model1, Model2)
		escpos.ECHigh, // Nivel de corrección de errores (Low, Medium, High, Highest)
		8,             // Tamaño del módulo (1-16)
		256,           // Tamaño del QR Code (en pixeles, si el protocolo no soporta QR nativo)
	); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Cut(escpos.PartialCut); err != nil {
		log.Printf("Error al cortar: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	log.Println("Impresión de QR completada exitosamente.")
}
