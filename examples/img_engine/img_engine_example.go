// Package main demonstrates the new advanced graphics engine for ESC/POS printing
package main

import (
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

func main() {
	// ========== Crear perfil de impresora ==========
	prof := profile.CreateECPM80250()

	// ========== Crear conector de impresora ==========
	conn, err := connection.NewWindowsPrintConnector(prof.Model)
	if err != nil {
		log.Fatalf("Error al crear conector: %v", err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar conector de impresora: %v", err)
		}
	}(conn)

	// ========== Crear protocolo de impresión ==========
	proto := composer.NewEscpos()

	// ========== Crear instancia de impresora ==========
	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		log.Printf("Error al crear impresora: %v", err)
	}
	defer func(printer *service.Printer) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar impresora: %v", err)
		}
	}(printer)

	// Load the image
	imgName := "test_image.jpg"
	imgPath := "./assets/images"
	log.Printf("Loading image from: %s", imgPath)
	img, err := graphics.LoadImageFromFile(imgPath, imgName)
	if err != nil {
		log.Panicf("Failed to load image: %v", err)
	}

	// Log original image dimensions
	bounds := img.Bounds()
	log.Printf("Original image size: %dx%d", bounds.Dx(), bounds.Dy())

	// Configure processing options
	opts := &graphics.Options{
		Width:          384,
		Threshold:      128,
		Mode:           graphics.Atkinson,
		AutoRotate:     false,
		PreserveAspect: true,
	}

	// Create and run the processing pipeline
	pipeline := graphics.NewPipeline(opts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		log.Panicf("Failed to process image: %v", err)
	}

	log.Printf("Processed image size: %dx%d", bitmap.Width, bitmap.Height)
	log.Printf("Raster data size: %d bytes", len(bitmap.GetRasterData()))

	// Generate ESC/POS commands
	err = printJob(printer, bitmap)
	if err != nil {
		log.Panicf("Failed to print job: %v", err)
	}

	log.Printf("Print job sent successfully")

	// Print some statistics
	printStatistics(bitmap)
}

// printStatistics displays processing statistics
func printStatistics(bitmap *graphics.MonochromeBitmap) {
	// Count black pixels
	blackPixels := 0
	totalPixels := bitmap.Width * bitmap.Height

	for y := 0; y < bitmap.Height; y++ {
		for x := 0; x < bitmap.Width; x++ {
			if bitmap.GetPixel(x, y) {
				blackPixels++
			}
		}
	}

	density := float64(blackPixels) * 100.0 / float64(totalPixels)

	fmt.Println("\n=== Processing Statistics ===")
	fmt.Printf("Image dimensions: %d x %d pixels\n", bitmap.Width, bitmap.Height)
	fmt.Printf("Total pixels: %d\n", totalPixels)
	fmt.Printf("Black pixels: %d (%.2f%%)\n", blackPixels, density)
	fmt.Printf("White pixels: %d (%.2f%%)\n", totalPixels-blackPixels, 100.0-density)
	fmt.Printf("Raster bytes: %d\n", len(bitmap.GetRasterData()))
	fmt.Printf("Compression ratio: %.2fx\n", float64(totalPixels/8)/float64(len(bitmap.GetRasterData())))
}

func printJob(printer *service.Printer, bitmap *graphics.MonochromeBitmap) error {
	err := printer.Initialize()
	if err != nil {
		return err
	}

	err = printer.AlignCenter()
	if err != nil {
		return err
	}

	err = printer.PrintLine("SET YOUR")
	if err != nil {
		return err
	}

	err = printer.FeedLines(1)
	if err != nil {
		return err
	}

	// ========== IMPRIMIR IMAGEN ==========
	log.Println("Imprimiendo imagen...")
	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("print bitmap: %w", err)
	}

	err = printer.FeedLines(1)
	if err != nil {
		return err
	}

	err = printer.PrintLine("HEART ABLAZE")
	if err != nil {
		return err
	}

	err = printer.FeedLines(3)
	if err != nil {
		return err
	}
	// Aquí se agregarían los comandos para imprimir la imagen procesada

	err = printer.PartialFeedAndCut(9)
	if err != nil {
		return err
	}

	return nil
}
