// Package main demuestra cómo imprimir usando una impresora POS en Windows.
package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/devices"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/profile"
)

// printHeader imprime el encabezado del documento de prueba
func printHeader(printer *devices.Printer) error {
	// Título centrado y en negrita
	if err := printer.PrintTitle("PRUEBA RÁPIDA DE IMPRESIÓN"); err != nil {
		return err
	}

	// Separador
	return printer.PrintSeparator("=", 32)
}

// printMainContent imprime el contenido principal del documento
func printMainContent(printer *devices.Printer) error {
	if err := printer.PrintLine("Esta es una prueba básica de la impresora completamente desacoplada."); err != nil {
		return err
	}
	return printer.NewLine()
}

// printAdvantages imprime la lista de ventajas de la nueva arquitectura
func printAdvantages(printer *devices.Printer) error {
	// Encabezado en negrita
	if err := printer.PrintHeader("Ventajas de la nueva arquitectura:"); err != nil {
		return err
	}

	// Lista de ventajas
	advantages := []string{
		"- Comandos ESCPOS completos",
		"- Conectores independientes",
		"- Perfiles intercambiables",
		"- Plantillas de ticket",
		"- Procesamiento de imágenes mejorado",
	}

	return printer.PrintLines(advantages)
}

// finishPrinting alimenta papel y corta
func finishPrinting(printer *devices.Printer) error {
	if err := printer.Feed(3); err != nil {
		return err
	}
	return printer.PartialCut()
}

func main() {
	// Configuración de la impresora
	printerName := "80mm RPT004"

	// Crear conector de impresora
	log.Printf("Intentando conectar a la impresora: %s", printerName)
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear el conector: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error al cerrar el conector: %v", err)
		}
	}()

	// Crear perfil de impresora
	prof := profile.CreateProfile80mm()

	// Crear protocolo ESCPOS
	proto := escpos.NewEscposCommands()

	// Crear instancia de impresora
	printer, err := devices.NewPrinter(proto, prof, conn)
	if err != nil {
		log.Panicf("Error al crear la impresora: %v", err)
	}
	defer func() {
		if err := printer.Close(); err != nil {
			log.Printf("Error al cerrar la impresora: %v", err)
		}
	}()

	// Inicializar impresora
	log.Println("Enviando comandos de prueba...")
	if err := printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// Imprimir documento completo
	if err := printHeader(printer); err != nil {
		log.Printf("Error al imprimir encabezado: %v", err)
	}
	if err := printMainContent(printer); err != nil {
		log.Printf("Error al imprimir contenido: %v", err)
	}
	if err := printAdvantages(printer); err != nil {
		log.Printf("Error al imprimir ventajas: %v", err)
	}
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error al finalizar: %v", err)
	}

	log.Println("Impresión completada!")
}
