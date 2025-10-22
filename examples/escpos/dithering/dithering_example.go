// Package main demuestra cómo imprimir imágenes con diferentes métodos de dithering usando una impresora POS.
package main

import (
	"image"
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/imaging"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// printDitheringHeader imprime el encabezado del documento de prueba
func printDitheringHeader(printer *pos.EscposPrinter) error {
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := printer.TextLn("PRUEBA DE DITHERING"); err != nil {
		return err
	}
	if err := printer.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}
	return printer.Feed(1)
}

// printImageWithoutDithering imprime una imagen sin aplicar dithering
func printImageWithoutDithering(printer *pos.EscposPrinter, img image.Image) error {
	if err := printer.TextLn("Imagen sin dithering:"); err != nil {
		return err
	}
	if err := printer.PrintImage(img); err != nil {
		return err
	}
	return printer.Feed(1)
}

// printImageWithAtkinson imprime una imagen con dithering Atkinson
func printImageWithAtkinson(printer *pos.EscposPrinter, img image.Image) error {
	if err := printer.TextLn("Imagen con Atkinson:"); err != nil {
		return err
	}

	opts := pos.PrintImageOptions{
		Density:    escpos.DensitySingle,
		DitherMode: imaging.DitherAtkinson,
		Threshold:  128,
		Width:      256,
	}

	return printer.PrintImageWithOptions(img, opts)
}

// finishPrinting imprime texto final, alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) error {
	if err := printer.Feed(1); err != nil {
		return err
	}
	if err := printer.TextLn("Fin del test de imágenes"); err != nil {
		return err
	}
	if err := printer.Feed(3); err != nil {
		return err
	}
	if err := printer.Cut(escpos.PartialCut); err != nil {
		return err
	}
	return printer.Feed(3)
}

func main() {
	// ========== Configuración de la impresora ==========
	printerName := "58mm GP-58N"

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
	prof := profile.CreateProfile58mm()

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

	// ========== Inicializar impresora ==========
	if err := printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// ========== Cargar imagen desde archivo ==========
	img, err := imaging.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Printf("Error al cargar imagen: %v", err)
	}

	// ========== Imprimir documento con diferentes opciones de dithering ==========
	if err := printDitheringHeader(printer); err != nil {
		log.Printf("Error al imprimir encabezado: %v", err)
	}
	if err := printImageWithoutDithering(printer, img); err != nil {
		log.Printf("Error al imprimir imagen sin dithering: %v", err)
	}
	if err := printImageWithAtkinson(printer, img); err != nil {
		log.Printf("Error al imprimir imagen con Atkinson: %v", err)
	}

	// ========== Finalizar impresión ==========
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error al finalizar impresión: %v", err)
	}
}
