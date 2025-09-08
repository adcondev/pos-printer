// Package main demonstrates how to print images using a POS printer with different dithering algorithms.
package main

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/imaging"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
	// En el futuro: "github.com/adcondev/pos-printer/protocol/zpl"
)

func main() {
	// === Crear conector ===
	conn, err := connector.NewWindowsPrintConnector("58mm PT-210")
	// conn, err := connector.NewWindowsPrintConnector("80mm EC-PM-80250") // Cambia el nombre según tu impresora
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar conector")
		}
	}(conn)

	// === Opción 1: Usar protocolo ESC/POS ===
	useESCPOS(conn)

	// === Opción 2: Usar protocolo ZPL (cuando esté implementado) ===
	// useZPL(conn)
}

func useESCPOS(conn connector.Connector) {
	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	prof := profile.CreateProfile80mm()

	// Crear impresora
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

	// Cargar imagen
	img, err := imaging.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Printf("Error al cargar imagen: %v", err)
	}

	err = printer.Feed(3)
	if err != nil {
		return
	}
	// Imprimir con densidad normal
	if err = printer.PrintImage(img); err != nil {
		log.Printf("Error imprimiendo imagen: %v", err)
	}

	// === Opción 2: Imprimir con Floyd-Steinberg ===
	log.Println("Imprimiendo con Floyd-Steinberg...")
	opts := pos.PrintImageOptions{
		Density:    escpos.DensitySingle,
		DitherMode: imaging.DitherFloydSteinberg,
		Threshold:  128,
		Width:      256,
	}

	// if err := printer.PrintImageWithOptions(img, opts); err != nil {
	// 	log.Printf("Error: %v", err)
	// }

	// === Opción 3: Imprimir con Atkinson ===
	log.Println("Imprimiendo con Atkinson...")
	opts.DitherMode = imaging.DitherAtkinson
	if err = printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}

	err = printer.Feed(1)
	if err != nil {
		return
	}
	err = printer.Print("Fin del utils de imágenes")
	if err != nil {
		return
	}
	err = printer.Feed(3)
	if err != nil {
		return
	}
	err = printer.Cut(escpos.PartialCut)
	if err != nil {
		return
	}
	err = printer.Feed(3)
	if err != nil {
		return
	}
}

// Ejemplo de cómo sería con otro protocolo
// func useZPL(conn connector.Connector) {
// TODO: Cuando implementes ZPL
/*
	protocol := zpl.NewZPLProtocol()
	printer, err := pos.NewPrinter(protocol, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	img := loadTestImage()

	// ZPL procesará la imagen de manera diferente internamente,
	// pero la API es la misma
	if err := printer.PrintRasterBitImage(img, types.DensitySingle); err != nil {
		log.Printf("Error: %v", err)
	}
*/
// }
