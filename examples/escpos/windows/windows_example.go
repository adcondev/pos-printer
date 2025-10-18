// Package main demuestra cómo imprimir usando una impresora POS en Windows.
package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// printHeader imprime el encabezado del documento de prueba
func printHeader(printer *pos.EscposPrinter) error {
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := printer.TextLn("PRUEBA RÁPIDA DE IMPRESIÓN"); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}
	if err := printer.TextLn("================================"); err != nil {
		return err
	}
	return printer.SetJustification(escpos.AlignLeft)
}

// printMainContent imprime el contenido principal del documento
func printMainContent(printer *pos.EscposPrinter) error {
	if err := printer.TextLn("Esta es una prueba básica de la impresora completamente desacoplada."); err != nil {
		return err
	}
	return printer.TextLn("")
}

// printAdvantages imprime la lista de ventajas de la nueva arquitectura
func printAdvantages(printer *pos.EscposPrinter) error {
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := printer.TextLn("Ventajas de la nueva arquitectura:"); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}

	// Lista de ventajas
	advantages := []string{
		"- Protocolos intercambiables",
		"- Conectores independientes",
		"- Perfiles intercambiables",
		"- Plantillas de ticket",
		"- Procesamiento de imágenes mejorado",
	}

	for _, advantage := range advantages {
		if err := printer.TextLn(advantage); err != nil {
			return err
		}
	}
	return nil
}

// finishPrinting alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) error {
	if err := printer.Feed(2); err != nil {
		return err
	}
	if _, err := printer.Connector.Write([]byte{common.GS, 'V', 66, 3}); err != nil {
		return err
	}
	return nil
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
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar el conector: %v", err)
		}
	}(conn)

	// Crear perfil de impresora
	prof := profile.CreateProfile80mm()

	// Crear instancia de impresora genérica
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error al crear la impresora: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar la impresora: %v", err)
		}
	}(printer)

	// Inicializar impresora
	log.Println("Enviando comandos de prueba...")
	if err = printer.Initialize(); err != nil {
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
