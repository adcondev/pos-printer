package main

import (
	"log"

	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/posprinter"
	"github.com/AdConDev/pos-printer/profile"
	"github.com/AdConDev/pos-printer/types"
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
	printer, err := posprinter.NewGenericPrinter(types.EscposProto, conn, prof)
	if err != nil {
		log.Fatalf("Error al crear impresora: %v", err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// === Imprimir título ===
	if err := printer.SetFont(types.FontA); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(types.EmphOn); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("PRUEBA DE QR á é í ó ú ñ"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(types.EmphOff); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	// === Imprimir QR Code ===
	if err := printer.PrintQR(
		"https://github.com/AdConDev/pos-daemon", // Contenido del QR Code
		types.Model2,                             // Modelo de QR Code (Model1, Model2)
		types.ECHigh,                             // Nivel de corrección de errores (Low, Medium, High, Highest)
		8,                                        // Tamaño del módulo (1-16)
		256,                                      // Tamaño del QR Code (en pixeles, si el protocolo no soporta QR nativo)
	); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Cut(types.CutFeed, 1); err != nil {
		log.Printf("Error al cortar: %v", err)
	}
	log.Println("Impresión de QR completada exitosamente.")
}
