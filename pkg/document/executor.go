package document

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/pkg/controllers/escpos/character"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/service"
)

// Executor ejecuta documentos de impresión
type Executor struct {
	printer  *service.Printer
	handlers map[string]CommandHandler
	profile  *profile.Escpos
}

// CommandHandler función que maneja un tipo de comando
type CommandHandler func(printer *service.Printer, data json.RawMessage) error

// NewExecutor crea un nuevo ejecutor
func NewExecutor(printer *service.Printer) *Executor {
	e := &Executor{
		printer:  printer,
		handlers: make(map[string]CommandHandler),
		profile:  &printer.Profile,
	}

	// Registrar handlers básicos
	e.RegisterHandler("text", e.handleText)
	e.RegisterHandler("image", e.handleImage)
	e.RegisterHandler("separator", e.handleSeparator)
	e.RegisterHandler("feed", e.handleFeed)
	e.RegisterHandler("cut", e.handleCut)

	// Handlers preparados para futura implementación
	e.RegisterHandler("qr", e.handleQRPlaceholder)
	e.RegisterHandler("table", e.handleTablePlaceholder)

	return e
}

// RegisterHandler registra un nuevo manejador de comando
func (e *Executor) RegisterHandler(cmdType string, handler CommandHandler) {
	e.handlers[cmdType] = handler
}

// Execute ejecuta un documento completo
func (e *Executor) Execute(doc *Document) error {
	// Inicializar impresora
	if err := e.printer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize printer: %w", err)
	}

	// Configurar code table si se especifica
	if doc.Profile.CodeTable != "" {
		if err := e.setCodeTable(doc.Profile.CodeTable); err != nil {
			log.Printf("Warning: failed to set code table %s: %v", doc.Profile.CodeTable, err)
		}
	}

	// Ejecutar cada comando
	for i, cmd := range doc.Commands {
		handler, exists := e.handlers[cmd.Type]
		if !exists {
			return fmt.Errorf("unknown command type at position %d: %s", i, cmd.Type)
		}

		if err := handler(e.printer, cmd.Data); err != nil {
			return fmt.Errorf("command %d (%s) failed: %w", i, cmd.Type, err)
		}
	}

	return nil
}

// setCodeTable configura la tabla de caracteres
func (e *Executor) setCodeTable(tableName string) error {
	// Mapa de nombres a constantes
	tables := map[string]character.CodeTable{
		"PC437":   character.PC437,
		"PC850":   character.PC850,
		"PC852":   character.PC852,
		"WPC1252": character.WPC1252,
		// Agregar más según necesidad
	}

	table, ok := tables[tableName]
	if !ok {
		return fmt.Errorf("unknown code table: %s", tableName)
	}

	return e.printer.SetCodeTable(table)
}

// ExecuteJSON ejecuta un documento desde JSON
func (e *Executor) ExecuteJSON(data []byte) error {
	doc, err := ParseDocument(data)
	if err != nil {
		return err
	}
	return e.Execute(doc)
}
