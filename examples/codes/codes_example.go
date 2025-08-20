package main

import (
	"fmt"
	"log"

	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/encoding"
	"github.com/AdConDev/pos-printer/escpos"
	"github.com/AdConDev/pos-printer/pos"
	"github.com/AdConDev/pos-printer/profile"
)

func main() {
	// Configuración de impresoras para probar
	printers := []struct {
		Name     string
		CharSets []encoding.CharacterSet // Charsets reportados por el fabricante
	}{
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
		fmt.Printf("\n=== Probando %s ===\n", printer.Name)

		// Conectar a la impresora
		conn, err := connector.NewWindowsPrintConnector(printer.Name)
		if err != nil {
			log.Printf("Error conectando a %s: %v", printer.Name, err)
			continue
		}
		defer func(conn *connector.WindowsPrintConnector) {
			err := conn.Close()
			if err != nil {
				log.Printf("Error al cerrar el conector de %s: %v", printer.Name, err)
			}
		}(conn)

		// Crear perfil personalizado
		prof := profile.CreateProfile80mm()
		prof.CharacterSets = printer.CharSets
		prof.Model = printer.Name

		p, err := pos.NewEscposPrinter(pos.EscposProto, conn, prof)
		if err != nil {
			log.Printf("Error creando impresora: %v", err)
			continue
		}
		defer func(p *pos.EscposPrinter) {
			err := p.Close()
			if err != nil {
				log.Printf("Error al cerrar la impresora %s: %v", printer.Name, err)
			}
		}(p)

		// Imprimir encabezado
		if err := p.SetJustification(escpos.AlignCenter); err != nil {
			log.Printf("Error estableciendo alineación centrada: %v", err)
			continue
		}

		if err := p.SetEmphasis(escpos.EmphOn); err != nil {
			log.Printf("Error activando negrita: %v", err)
			continue
		}

		err = p.TextLn(fmt.Sprintf("TEST CODIFICACIÓN - %s", printer.Name))
		if err != nil {
			log.Printf("Error imprimiendo encabezado: %v", err)
			continue
		}

		if err := p.SetEmphasis(escpos.EmphOff); err != nil {
			log.Printf("Error desactivando negrita: %v", err)
			continue
		}

		if err := p.Feed(1); err != nil {
			log.Printf("Error alimentando papel: %v", err)
			continue
		}

		// Probar cada charset soportado
		for _, charset := range printer.CharSets {
			// Verificar que el charset esté en nuestro Registry
			if _, exists := encoding.Registry[charset]; !exists {
				continue
			}

			err := p.SetJustification(escpos.AlignLeft)
			if err != nil {
				log.Printf("Error estableciendo alineación izquierda: %v", err)
				continue
			}

			if err := p.SetEmphasis(escpos.EmphOn); err != nil {
				log.Printf("Error activando negrita: %v", err)
				continue
			}
			err = p.TextLn(fmt.Sprintf("=== Charset %d (%s) ===",
				charset, encoding.Registry[charset].Name))
			if err != nil {
				log.Printf("Error imprimiendo encabezado de charset: %v", err)
				continue
			}
			if err := p.SetEmphasis(escpos.EmphOff); err != nil {
				log.Printf("Error desactivando negrita: %v", err)
				continue
			}

			// Cancelar modo Kanji
			if err := p.CancelKanjiMode(); err != nil {
				log.Printf("Error cancelando modo Kanji: %v", err)
				continue
			}

			// Cambiar al charset
			if err := p.SetCharacterSet(charset); err != nil {
				err := p.TextLn(fmt.Sprintf("Error: %v", err))
				if err != nil {
					log.Printf("Error imprimiendo mensaje de error: %v", err)
					continue
				}
				continue
			}

			// Imprimir textos de prueba
			for _, text := range testTexts {
				if err := p.TextLn(text); err != nil {
					err := p.TextLn(fmt.Sprintf("Error imprimiendo: %v", err))
					if err != nil {
						log.Printf("Error imprimiendo texto: %v", err)
						continue
					}
				}
			}

			if err := p.Feed(1); err != nil {
				log.Printf("Error alimentando papel: %v", err)
				continue
			}
		}

		// Cortar

		if err := p.Feed(1); err != nil {
			log.Printf("Error alimentando papel: %v", err)
			continue
		}

		if err := p.Cut(escpos.PartialCut); err != nil {
			log.Printf("Error cortando papel: %v", err)
			continue
		}

		if err := p.Feed(1); err != nil {
			log.Printf("Error alimentando papel: %v", err)
			continue
		}
	}
}
