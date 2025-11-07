// Package main implements an example of printing a document using the pos-printer library in JSON print job.
package main

import (
	"log"
	"os"

	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

func main() {
	// 1. Crear perfil de impresora
	prof := profile.CreateECPM80250()

	// 2. Crear conexión
	conn, err := connection.NewWindowsPrintConnector(prof.Model)
	if err != nil {
		log.Fatalf("Failed to create connector: %v", err)
	}
	defer func(conn *connection.WindowsPrintConnector) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connector: %v", err)
		}
	}(conn)

	// 3. Crear protocolo
	proto := composer.NewEscpos()

	// 4. Crear servicio de impresora
	printer, err := service.NewPrinter(proto, prof, conn)
	if err != nil {
		log.Panicf("Failed to create printer: %v", err)
	}
	defer func(printer *service.Printer) {
		err := printer.Close()
		if err != nil {
			log.Fatalf("Failed to close printer: %v", err)
		}
	}(printer)

	// 5. Crear ejecutor de documentos
	executor := document.NewExecutor(printer)

	// Opción A: Cargar documento JSON desde archivo
	jsonData, err := os.ReadFile("./examples/document/ticket.json")
	if err != nil {
		log.Panicf("Failed to read JSON file: %v", err)
	}

	if err := executor.ExecuteJSON(jsonData); err != nil {
		log.Panicf("Failed to execute document: %v", err)
	}

	// Opción B: Construir documento programáticamente
	builder := document.NewBuilder()
	doc := builder.
		SetProfile("80mm EC-PM-80250", 576, "PC850").
		AddText("MI TIENDA", &document.TextStyle{
			Align: "center",
			Bold:  true,
			Size:  "2x2",
		}).
		AddSeparator("=", 48).
		AddText("Ticket de Venta", nil).
		AddFeed(1).
		AddText("Total: $100.00", &document.TextStyle{
			Align: "right",
			Bold:  true,
			Size:  "2x2",
		}).
		AddFeed(3).
		AddCut("partial", 0).
		Build()

	if err := executor.Execute(doc); err != nil {
		log.Panicf("Failed to execute document: %v", err)
	}

	log.Println("Document printed successfully!")
}
