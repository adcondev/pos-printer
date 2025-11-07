/*
Package printer provides a comprehensive ESC/POS printer control library for Go.

This package implements the ESC/POS command protocol used by most thermal receipt printers,
offering a high-level API for printer control, text formatting, graphics printing, and more.

Installation

	go get github.com/adcondev/pos-printer@v3.0.0

Basic Usage

	import (
		"github.com/adcondev/pos-printer/pkg/composer"
		"github.com/adcondev/pos-printer/pkg/connection"
		"github.com/adcondev/pos-printer/pkg/profile"
		"github.com/adcondev/pos-printer/pkg/service"
	)

	func main() {
		// Create printer components
		proto := composer.NewEscpos()
		prof := profile.CreatePt210()
		conn, _ := connection.NewWindowsPrintConnector("POS-58")

		// Initialize printer
		printer, _ := service.NewPrinter(proto, prof, conn)
		defer printer.Close()

		// Print text
		printer.Initialize()
		printer.PrintLine("Hello, World!")
		printer.FullCut(3)
	}

# Architecture

The library follows a modular architecture with clear separation of concerns:

  - Composer: High-level ESC/POS command composition
  - Commands: Individual command implementations (barcode, character, etc.)
  - Connection: Platform-specific printer connections
  - Graphics: Advanced image processing and bitmap conversion
  - Profile: Printer-specific configurations
  - Service: Main printer service interface

Supported Features

  - Text printing with multiple fonts and sizes
  - Barcode printing (UPC-A/E, EAN-13/8, CODE39, ITF, CODABAR, CODE128, QR Code)
  - Image/bitmap printing with dithering algorithms
  - Paper cutting (full/partial)
  - Character formatting (bold, underline, etc.)
  - International character sets and code pages
  - Custom print positioning and alignment

Platform Support

  - Windows: Native Print Spooler API support
  - Linux: CUPS or direct USB/Serial connection
  - macOS: CUPS support

Version History

  - v3.0.0: Major refactor with modular architecture
  - v2.0.0: Initial stable release with core features
  - v1.8.0: Added Task automation
  - v1.7.0: Bit image support
  - v1.6.0: Mechanism control commands

For detailed documentation and examples, visit: https://github.com/adcondev/pos-printer
*/
package printer
