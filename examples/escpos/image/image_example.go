// Package main demuestra cómo imprimir imágenes usando una impresora POS con diferentes algoritmos de dithering.
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

// printImageDefault imprime imagen con configuración por defecto
func printImageDefault(printer *pos.EscposPrinter, img image.Image) error {
	if err := printer.Feed(3); err != nil {
		return err
	}
	return printer.PrintImage(img)
}

// printImageWithAtkinson imprime imagen con algoritmo de dithering Atkinson
func printImageWithAtkinson(printer *pos.EscposPrinter, img image.Image) error {
	log.Println("Imprimiendo con Atkinson...")

	opts := pos.PrintImageOptions{
		Density:    escpos.DensitySingle,
		DitherMode: imaging.DitherAtkinson,
		Threshold:  128,
		Width:      256,
	}

	return printer.PrintImageWithOptions(img, opts)
}

// finishPrinting imprime mensaje final, alimenta papel y corta
func finishPrinting(printer *pos.EscposPrinter) error {
	if err := printer.Feed(1); err != nil {
		return err
	}
	if err := printer.Print("Fin del test de imágenes"); err != nil {
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

// useESCPOS configura y ejecuta la impresión usando protocolo ESC/POS
func useESCPOS(conn connector.Connector) {
	// Crear perfil de impresora
	prof := profile.CreateProfile80mm()

	// Crear instancia de impresora
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Fatal(err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar impresora")
		}
	}(printer)

	// Cargar imagen desde archivo
	img, err := imaging.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Printf("Error al cargar imagen: %v", err)
	}

	// Imprimir con densidad normal
	if err := printImageDefault(printer, img); err != nil {
		log.Printf("Error imprimiendo imagen: %v", err)
	}

	// Imprimir con algoritmo Atkinson
	if err := printImageWithAtkinson(printer, img); err != nil {
		log.Printf("Error: %v", err)
	}

	// Finalizar impresión
	if err := finishPrinting(printer); err != nil {
		log.Printf("Error al finalizar: %v", err)
	}
}

func main() {
	// ========== Configuración de la impresora ==========
	printerName := "80mm RPT004"

	// ========== Crear conector de impresora ==========
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar conector")
		}
	}(conn)

	// ========== Usar protocolo ESC/POS para imprimir ==========
	useESCPOS(conn)
}
