package main

import (
	"log"

	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/escpos"
	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/pos"
	"github.com/AdConDev/pos-printer/profile"
)

func main() {
	// === Crear conector ===
	// Seleccionar la impresora según tu configuración
	// printerName := "80mm EC-PM-80250"
	printerName := "58mm GP-58N"

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	// prof := profile.CreateProfile80mm()
	prof := profile.CreateProfile58mm() // Usar perfil de 58mm

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
	printer, err := pos.NewEscposPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Fatalf("Error al crear impresora: %v", err)
	}
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// === Inicializar impresora ===
	if err := printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// === Cargar imagen ===
	img, err := imaging.LoadImage("./img/perro.jpeg")
	if err != nil {
		log.Fatalf("Error al cargar imagen: %v", err)
	}

	// === Imprimir título ===
	if err := printer.SetJustification(escpos.AlignCenter); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphOn); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("PRUEBA DE DITHERING"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(escpos.EmphOff); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Opción 1: Imprimir sin dithering ===
	if err := printer.TextLn("Imagen sin dithering:"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.PrintImage(img); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}

	opts := pos.PrintImageOptions{
		Density:    escpos.DensitySingle,
		DitherMode: imaging.DitherFloydSteinberg,
		Threshold:  128,
		Width:      256, // 0 = usar ancho original de imagen. La imagen podría salir más ancha que el papel
	}

	// === Opción 2: Imprimir con Floyd-Steinberg ===
	/*
		if err := printer.Ln("Imagen con Floyd-Steinberg:"); err != nil {
			log.Printf("Error: %v", err)
		}

		if err := printer.PrintImageWithOptions(img, opts); err != nil {
			log.Printf("Error: %v", err)
		}
		if err := printer.Feed(2); err != nil {
			log.Printf("Error: %v", err)
		}
	*/

	// === Opción 3: Imprimir con Atkinson ===
	if err := printer.TextLn("Imagen con Atkinson:"); err != nil {
		log.Printf("Error: %v", err)
	}
	opts.DitherMode = imaging.DitherAtkinson
	if err := printer.PrintImageWithOptions(img, opts); err != nil {
		log.Printf("Error: %v", err)
	}

	// === Finalizar impresión ===
	if err = printer.Feed(1); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.TextLn("Fin del test de imágenes"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Feed(3); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Cut(escpos.PartialCut); err != nil {
		log.Printf("Error: %v", err)
	}
	if err = printer.Feed(3); err != nil {
		log.Printf("Error: %v", err)
	}
}
