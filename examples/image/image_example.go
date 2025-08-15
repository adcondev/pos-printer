package main

import (
	"log"

	"github.com/AdConDev/pos-printer"
	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/profile"
	"github.com/AdConDev/pos-printer/protocol/escpos"
	"github.com/AdConDev/pos-printer/types"
	// En el futuro: "github.com/AdConDev/pos-printer/protocol/zpl"
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
	// Crear protocolo ESC/POS
	proto := escpos.NewESCPOSProtocol()

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	prof := profile.CreateProfile80mm()

	// Crear impresora
	printer, err := posprinter.NewGenericPrinter(proto, conn, prof)
	if err != nil {
		log.Fatal(err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("dithering: error al cerrar impresora")
		}
	}(printer)

	// Cargar imagen
	img, err := imaging.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Fatalf("Error al cargar imagen: %v", err)
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
	opts := posprinter.PrintImageOptions{
		Density:    types.DensitySingle,
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
	err = printer.Text("Fin del test de imágenes")
	if err != nil {
		return
	}
	err = printer.Feed(3)
	if err != nil {
		return
	}

	err = printer.Cut(types.CutFeed, 3)
	if err != nil {
		return
	}
}

// Ejemplo de cómo sería con otro protocolo
// func useZPL(conn connector.Connector) {
// TODO: Cuando implementes ZPL
/*
	protocol := zpl.NewZPLProtocol()
	printer, err := posprinter.NewGenericPrinter(protocol, conn)
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
