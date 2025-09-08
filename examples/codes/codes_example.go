// Package main demonstrates how to utils different character encodings
package main

import (
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/encoding"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

// PrinterConfig holds configuration for a printer to utils
type PrinterConfig struct {
	Name     string
	CharSets []encoding.CharacterSet
}

func main() {
	// Configuración de impresoras para probar
	printers := []PrinterConfig{
		{
			Name:     "80mm EC-PM-80250 x",
			CharSets: []encoding.CharacterSet{encoding.WCP1252, encoding.CP858},
		},
		{
			Name: "58mm PT-210",
			CharSets: []encoding.CharacterSet{
				encoding.CP437,
				encoding.Katakana,
				encoding.CP850,
				encoding.CP860,
				encoding.CP863,
				encoding.CP865,
				encoding.WestEurope,
				encoding.Greek,
				encoding.Hebrew,
				// encoding.CP755, // No soportado directamente
				encoding.Iran,
				encoding.WCP1252,
				encoding.CP866,
				encoding.CP852,
				encoding.CP858,
				encoding.IranII,
				encoding.Latvian,
			},
		},
		{
			Name:     "58mm GP-58N x",
			CharSets: []encoding.CharacterSet{encoding.WCP1252, encoding.CP858},
		},
		// Agregar tu tercera impresora aquí
	}

	// Texto de prueba con caracteres especiales en español para tickets de venta
	testTexts := []string{
		"Acentos: áéíóú ÁÉÍÓÚ",
		"Eñe: ñ Ñ",
		"Diéresis: ü Ü",
		"Moneda: $ ¢",
		"Símbolos: ¡ ¿",
	}

	// Probar cada impresora
	for _, printer := range printers {
		testPrinter(printer, testTexts)
	}
}

// testPrinter tests a single printer with the provided utils texts
func testPrinter(printer PrinterConfig, testTexts []string) {
	fmt.Printf("\n=== Probando %s ===\n", printer.Name)

	// Conectar a la impresora
	conn, err := connector.NewWindowsPrintConnector(printer.Name)
	if err != nil {
		log.Printf("Error conectando a %s: %v", printer.Name, err)
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar el conector de %s: %v", printer.Name, err)
		}
	}()

	// Crear perfil personalizado
	prof := profile.CreateProfile80mm()
	prof.CharacterSets = printer.CharSets
	prof.Model = printer.Name

	p, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error creando impresora: %v", err)
		return
	}
	defer func() {
		err := p.Close()
		if err != nil {
			log.Printf("Error al cerrar la impresora %s: %v", printer.Name, err)
		}
	}()

	// Imprimir encabezado
	if err := p.SetJustification(escpos.AlignCenter); err != nil {
		log.Printf("Error estableciendo alineación centrada: %v", err)
		return
	}

	if err := p.SetEmphasis(escpos.EmphasizedOn); err != nil {
		log.Printf("Error activando negrita: %v", err)
		return
	}

	err = p.TextLn(fmt.Sprintf("TEST CODIFICACIÓN - %s", printer.Name))
	if err != nil {
		log.Printf("Error imprimiendo encabezado: %v", err)
		return
	}

	if err := p.SetEmphasis(escpos.EmphasizedOff); err != nil {
		log.Printf("Error desactivando negrita: %v", err)
		return
	}

	if err := p.Feed(1); err != nil {
		log.Printf("Error alimentando papel: %v", err)
		return
	}

	// Probar cada charset soportado
	for _, charset := range printer.CharSets {
		testCharset(p, charset, testTexts)
	}

	// Cortar
	if err := p.Feed(1); err != nil {
		log.Printf("Error alimentando papel: %v", err)
		return
	}

	if err := p.Cut(escpos.PartialCut); err != nil {
		log.Printf("Error cortando papel: %v", err)
		return
	}

	if err := p.Feed(1); err != nil {
		log.Printf("Error alimentando papel: %v", err)
		return
	}
}

// testCharset tests a specific charset on the printer
func testCharset(p *pos.EscposPrinter, charset encoding.CharacterSet, testTexts []string) {
	// Verificar que el charset esté en nuestro Registry
	if _, exists := encoding.Registry[charset]; !exists {
		return
	}

	err := p.SetJustification(escpos.AlignLeft)
	if err != nil {
		log.Printf("Error estableciendo alineación izquierda: %v", err)
		return
	}

	if err := p.SetEmphasis(escpos.EmphasizedOn); err != nil {
		log.Printf("Error activando negrita: %v", err)
		return
	}

	err = p.TextLn(fmt.Sprintf("=== Charset %d (%s) ===",
		charset, encoding.Registry[charset].Name))
	if err != nil {
		log.Printf("Error imprimiendo encabezado de charset: %v", err)
		return
	}

	if err := p.SetEmphasis(escpos.EmphasizedOff); err != nil {
		log.Printf("Error desactivando negrita: %v", err)
		return
	}

	// Cancelar modo Kanji
	if err := p.CancelKanjiMode(); err != nil {
		log.Printf("Error cancelando modo Kanji: %v", err)
		return
	}

	// Cambiar al charset
	if err := p.SetCharacterSet(charset); err != nil {
		err := p.TextLn(fmt.Sprintf("Error: %v", err))
		if err != nil {
			log.Printf("Error imprimiendo mensaje de error: %v", err)
			return
		}
		return
	}

	// Imprimir textos de prueba
	for _, text := range testTexts {
		if err := p.TextLn(text); err != nil {
			err := p.TextLn(fmt.Sprintf("Error imprimiendo: %v", err))
			if err != nil {
				log.Printf("Error imprimiendo texto: %v", err)
				return
			}
		}
	}

	if err := p.Feed(1); err != nil {
		log.Printf("Error alimentando papel: %v", err)
		return
	}
}
