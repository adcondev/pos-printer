// Package main demuestra cómo utilizar diferentes codificaciones de caracteres
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

// PrinterConfig almacena la configuración para una impresora específica
type PrinterConfig struct {
	Name     string
	CharSets []encoding.CharacterSet
}

// printTestHeader imprime el encabezado del test de codificación
func printTestHeader(p *pos.EscposPrinter, printerName string) error {
	if err := p.SetJustification(escpos.AlignCenter); err != nil {
		return err
	}
	if err := p.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := p.TextLn(fmt.Sprintf("TEST CODIFICACIÓN - %s", printerName)); err != nil {
		return err
	}
	if err := p.SetEmphasis(escpos.EmphasizedOff); err != nil {
		return err
	}
	return p.Feed(1)
}

// printCharsetHeader imprime el encabezado para un charset específico
func printCharsetHeader(p *pos.EscposPrinter, charset encoding.CharacterSet) error {
	if err := p.SetJustification(escpos.AlignLeft); err != nil {
		return err
	}
	if err := p.SetEmphasis(escpos.EmphasizedOn); err != nil {
		return err
	}
	if err := p.TextLn(fmt.Sprintf("=== Charset %d (%s) ===",
		charset, encoding.Registry[charset].Name)); err != nil {
		return err
	}
	return p.SetEmphasis(escpos.EmphasizedOff)
}

// setupCharset configura el charset en la impresora
func setupCharset(p *pos.EscposPrinter, charset encoding.CharacterSet) error {
	if err := p.CancelKanjiMode(); err != nil {
		return err
	}
	return p.SetCharacterSet(charset)
}

// printTestTexts imprime todos los textos de prueba
func printTestTexts(p *pos.EscposPrinter, testTexts []string) error {
	for _, text := range testTexts {
		if err := p.TextLn(text); err != nil {
			return err
		}
	}
	return p.Feed(1)
}

// testCharset prueba un charset específico en la impresora
func testCharset(p *pos.EscposPrinter, charset encoding.CharacterSet, testTexts []string) {
	// Verificar que el charset esté en el Registry
	if _, exists := encoding.Registry[charset]; !exists {
		return
	}

	// Imprimir encabezado del charset
	if err := printCharsetHeader(p, charset); err != nil {
		log.Printf("Error imprimiendo encabezado de charset: %v", err)
		return
	}

	// Configurar el charset
	if err := setupCharset(p, charset); err != nil {
		if err := p.TextLn(fmt.Sprintf("Error: %v", err)); err != nil {
			log.Printf("Error imprimiendo mensaje de error: %v", err)
		}
		return
	}

	// Imprimir textos de prueba
	if err := printTestTexts(p, testTexts); err != nil {
		log.Printf("Error imprimiendo textos: %v", err)
	}
}

// finishPrinting alimenta papel y corta
func finishPrinting(p *pos.EscposPrinter) error {
	if err := p.Feed(1); err != nil {
		return err
	}
	if err := p.Cut(escpos.PartialCut); err != nil {
		return err
	}
	return p.Feed(1)
}

// testPrinter prueba una impresora con los textos proporcionados
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

	// Crear instancia de impresora
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

	// Imprimir encabezado del test
	if err := printTestHeader(p, printer.Name); err != nil {
		log.Printf("Error imprimiendo encabezado: %v", err)
		return
	}

	// Probar cada charset soportado
	for _, charset := range printer.CharSets {
		testCharset(p, charset, testTexts)
	}

	// Finalizar impresión
	if err := finishPrinting(p); err != nil {
		log.Printf("Error finalizando impresión: %v", err)
	}
}

func main() {
	// ========== Configuración de impresoras para probar ==========
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
	}

	// ========== Texto de prueba con caracteres especiales en español para tickets de venta ==========
	testTexts := []string{
		"Acentos: áéíóú ÁÉÍÓÚ",
		"Eñe: ñ Ñ",
		"Diéresis: ü Ü",
		"Moneda: $ ¢",
		"Símbolos: ¡ ¿",
	}

	// ========== Probar cada impresora ==========
	for _, printer := range printers {
		testPrinter(printer, testTexts)
	}
}
