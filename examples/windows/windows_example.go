package main

import (
	"log"

	"github.com/AdConDev/pos-printer"
	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/profile"
	"github.com/AdConDev/pos-printer/protocol/escpos"
	"github.com/AdConDev/pos-printer/types"
)

func main() {
	// === Configuración ===
	// printerName := "80mm EC-PM-80250"
	printerName := "58mm PT-210" // Cambia esto al nombre de tu impresora

	// === Crear conector ===
	log.Printf("Intentando conectar a la impresora: %s", printerName)
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		log.Fatalf("Error al crear el conector: %v", err)
	}
	defer func(conn *connector.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error al cerrar el conector: %v", err)
		}
	}(conn)

	// === Crear protocolo ===
	// Aquí es donde eliges el protocolo (ESC/POS, ZPL, etc.)
	proto := escpos.NewESCPOSProtocol()

	// === Crear Perfil de impresora ===
	// Puedes definir un perfil si necesitas configuraciones específicas
	prof := profile.CreateProfile80mm()

	// === Crear impresora genérica ===
	printer, err := posprinter.NewGenericPrinter(proto, conn, prof)
	if err != nil {
		log.Fatalf("Error al crear la impresora: %v", err)
	}
	defer func(printer *posprinter.GenericPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar la impresora: %v", err)
		}
	}(printer)

	// === Prueba básica de impresión ===
	log.Println("Enviando comandos de prueba...")

	// Inicializar
	if err = printer.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}

	// Texto centrado (usando tipos del paquete types)
	if err = printer.SetJustification(types.AlignCenter); err != nil {
		log.Printf("Error al centrar: %v", err)
	}

	// Texto en negrita
	if err = printer.SetEmphasis(types.EmphOn); err != nil {
		log.Printf("Error al activar negrita: %v", err)
	}

	// Imprimir título
	if err = printer.TextLn("PRUEBA RAPIDA DE IMPRESION"); err != nil {
		log.Printf("Error al imprimir título: %v", err)
	}

	// Desactivar negrita
	if err = printer.SetEmphasis(types.EmphOff); err != nil {
		log.Printf("Error al desactivar negrita: %v", err)
	}

	// Línea separadora
	if err = printer.TextLn("================================"); err != nil {
		log.Printf("Error al imprimir línea: %v", err)
	}

	// Alinear a la izquierda
	if err = printer.SetJustification(types.AlignLeft); err != nil {
		log.Printf("Error al alinear izquierda: %v", err)
	}

	// Contenido
	if err := printer.TextLn("Esta es una prueba básica de la impresora completamente desacoplada."); err != nil {
		log.Printf("Error: %v", err)
	}

	// Mostrar las ventajas
	if err := printer.TextLn(""); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(types.EmphOn); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("Ventajas de la nueva arquitectura:"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.SetEmphasis(types.EmphOff); err != nil {
		log.Printf("Error: %v", err)
	}

	if err := printer.TextLn("- Protocolos intercambiables"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("- Conectores independientes"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("- Perfiles intercambiables"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("- Plantillas de ticket"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := printer.TextLn("- Procesamiento de imágenes mejorado"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Feed y corte
	if err := printer.Feed(3); err != nil {
		log.Printf("Error al alimentar papel: %v", err)
	}

	// Usar CutFull del paquete types (no del paquete escpos)
	if err := printer.Cut(types.CutFeed, 3); err != nil {
		log.Printf("Error al cortar: %v", err)
	}

	log.Println("Impresión completada!")
}
