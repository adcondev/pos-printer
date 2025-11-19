package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	"github.com/adcondev/pos-printer/pkg/commands/character"
	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/printer"
	"github.com/adcondev/pos-printer/pkg/tables"
)

const (
	center = "center"
	right  = "right"
)

// handleText manages text commands (actualización)
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Si hay un label, imprimirlo primero con su estilo
	if cmd.Label != "" {
		// Aplicar estilo del label
		if err := e.applyTextStyle(printer, cmd.LabelStyle); err != nil {
			return fmt.Errorf("failed to apply label style: %w", err)
		}

		// TODO: Define custom separator between label and content (":\n", " - ", etc.)

		// Imprimir label sin salto de línea
		labelText := cmd.Label
		if !strings.HasSuffix(labelText, ":") {
			labelText += ": "
		} else {
			labelText += " "
		}

		if err := printer.Print(labelText); err != nil {
			return fmt.Errorf("failed to print label: %w", err)
		}

		// Resetear estilos del label antes de aplicar los del contenido
		if err := e.resetTextStyle(printer, cmd.LabelStyle); err != nil {
			return fmt.Errorf("failed to reset label style: %w", err)
		}
	}

	// Aplicar estilo del contenido
	if err := e.applyTextStyle(printer, cmd.Style); err != nil {
		return fmt.Errorf("failed to apply content style: %w", err)
	}

	// Imprimir contenido
	if cmd.Content != "" {
		if cmd.NewLine {
			if err := printer.PrintLine(cmd.Content); err != nil {
				return err
			}
		} else {
			if err := printer.Print(cmd.Content); err != nil {
				return err
			}
		}
	}

	// Resetear estilos del contenido
	if err := e.resetTextStyle(printer, cmd.Style); err != nil {
		return fmt.Errorf("failed to reset content style: %w", err)
	}

	return nil
}

// applyTextStyle aplica los estilos de texto especificados
func (e *Executor) applyTextStyle(printer *service.Printer, style TextStyle) error {
	// Aplicar alineación
	switch strings.ToLower(style.Align) {
	case center:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case right:
		if err := printer.AlignRight(); err != nil {
			return err
		}
	default:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	}

	// Aplicar bold
	if style.Bold {
		if err := printer.Bold(); err != nil {
			return err
		}
	}

	// Aplicar tamaño
	switch strings.ToLower(style.Size) {
	case "2x2":
		if err := printer.DoubleSize(); err != nil {
			return err
		}
	case "3x3":
		size, _ := character.NewSize(3, 3)
		cmdBytes := printer.Protocol.Character.SelectCharacterSize(size)
		if err := printer.Write(cmdBytes); err != nil {
			return err
		}
	default:
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	// Aplicar underline
	if style.Underline {
		cmd, _ := printer.Protocol.Character.SetUnderlineMode(character.OneDot)
		if err := printer.Write(cmd); err != nil {
			return err
		}
	}

	// Aplicar inverse
	if style.Inverse {
		cmd := printer.Protocol.Character.SetWhiteBlackReverseMode(character.OnRm)
		if err := printer.Write(cmd); err != nil {
			return err
		}
	}

	return nil
}

// resetTextStyle resetea los estilos aplicados
func (e *Executor) resetTextStyle(printer *service.Printer, style TextStyle) error {
	// Reset bold
	if style.Bold {
		cmdBytes := printer.Protocol.Character.SetEmphasizedMode(character.OffEm)
		if err := printer.Write(cmdBytes); err != nil {
			return err
		}
	}

	// Reset size
	if style.Size != "" && style.Size != "normal" {
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	// Reset underline
	if style.Underline {
		cmd, _ := printer.Protocol.Character.SetUnderlineMode(character.NoDot)
		if err := printer.Write(cmd); err != nil {
			return err
		}
	}

	// Reset inverse
	if style.Inverse {
		cmd := printer.Protocol.Character.SetWhiteBlackReverseMode(character.OffRm)
		if err := printer.Write(cmd); err != nil {
			return err
		}
	}

	return nil
}

// handleImage manages image commands
func (e *Executor) handleImage(printer *service.Printer, data json.RawMessage) error {
	var cmd ImageCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse image command: %w", err)
	}

	// Decodificar imagen desde base64
	img, format, err := load.ImgFromBase64(cmd.Code)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	log.Printf("Loaded image with format: %s", format)

	// Configurar opciones de procesamiento
	opts := &graphics.ImgOptions{
		PixelWidth:     cmd.PixelWidth,
		Threshold:      cmd.Threshold,
		PreserveAspect: true,
		AutoRotate:     false,
	}

	// Si no se especifica ancho, usar el ancho del perfil
	if opts.PixelWidth == 0 {
		opts.PixelWidth = 256
	}

	// Si no se especifica threshold, usar valor por defecto
	if opts.Threshold == 0 {
		opts.Threshold = 128
	}

	// Configurar dithering
	switch strings.ToLower(cmd.Dithering) {
	case "atkinson":
		opts.Dithering = graphics.Atkinson
	default:
		opts.Dithering = graphics.Threshold
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	// Procesar imagen
	pipeline := graphics.NewPipeline(opts)
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	// Imprimir bitmap
	if err := printer.PrintBitmap(bitmap); err != nil {
		return fmt.Errorf("failed to print bitmap: %w", err)
	}

	// Resetear alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
}

// handleSeparator manages separator commands
func (e *Executor) handleSeparator(printer *service.Printer, data json.RawMessage) error {
	var cmd SeparatorCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse separator command: %w", err)
	}

	// Valores por defecto
	if cmd.Char == "" {
		cmd.Char = "-"
	}
	if cmd.Length == 0 {
		// Usar ancho del papel en caracteres (aproximado)
		cmd.Length = e.printer.Profile.DotsPerLine / 12 // Aproximación para Font A
	}

	// Construir línea separadora
	line := strings.Repeat(cmd.Char, cmd.Length)

	return printer.PrintLine(line)
}

// handleFeed manages feed commands
func (e *Executor) handleFeed(printer *service.Printer, data json.RawMessage) error {
	var cmd FeedCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse feed command: %w", err)
	}

	if cmd.Lines <= 0 {
		cmd.Lines = 1
	}

	return printer.FeedLines(byte(cmd.Lines))
}

// handleCut manages cut commands
func (e *Executor) handleCut(printer *service.Printer, data json.RawMessage) error {
	var cmd CutCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse cut command: %w", err)
	}

	// Avance antes del corte si se especifica
	if cmd.Feed > 0 {
		err := printer.FeedLines(byte(cmd.Feed))
		if err != nil {
			return err
		}
	}

	// Ejecutar corte
	switch strings.ToLower(cmd.Mode) {
	case "full":
		return printer.FullFeedAndCut(0)
	default: // partial
		return printer.PartialFeedAndCut(0)
	}
}

// TODO: Manage text_under and text_above options instead of human_text

// handleQR manges QR code commands
func (e *Executor) handleQR(printer *service.Printer, data json.RawMessage) error {
	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse QR command: %w", err)
	}

	// Validación de datos
	if cmd.Data == "" {
		return fmt.Errorf("QR data cannot be empty")
	}
	if len(cmd.Data) > posqr.MaxDataLength {
		return fmt.Errorf("QR data too long: %d bytes (maximum %d)", len(cmd.Data), posqr.MaxDataLength)
	}

	// Construir opciones
	opts := graphics.DefaultQROptions()

	if cmd.PixelWidth > 0 {
		opts.PixelWidth = cmd.PixelWidth
	} else {
		// Usar 50% del ancho del papel por defecto
		opts.PixelWidth = e.printer.Profile.DotsPerLine / 2
	}

	// Mapear corrección de errores
	switch strings.ToUpper(cmd.Correction) {
	case "L":
		opts.ErrorCorrection = posqr.LevelL
	case "Q":
		opts.ErrorCorrection = posqr.LevelQ
	case "H":
		opts.ErrorCorrection = posqr.LevelH
	default:
		opts.ErrorCorrection = posqr.LevelM
	}

	if cmd.Logo != "" {
		opts.LogoData = cmd.Logo
	}

	// Solo aplicar circle shape si no hay halftone
	opts.CircleShape = cmd.CircleShape

	// Aplicar alineación
	switch strings.ToLower(cmd.Align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	// Imprimir QR
	err := printer.PrintQR(cmd.Data, opts)
	if err != nil {
		return err
	}

	// Imprimir texto humano si existe
	if cmd.HumanText != "" {
		// Centrar el texto debajo del QR si el QR está centrado
		if strings.ToLower(cmd.Align) == "center" {
			if err := printer.AlignCenter(); err != nil {
				return err
			}
		}
		if err := printer.PrintLine(cmd.HumanText); err != nil {
			return err
		}
	}

	// Restaurar alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
}

// TODO: Consider a title fields for tables

// handleTable manages table commands
func (e *Executor) handleTable(printer *service.Printer, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse table command: %w", err)
	}

	if cmd.Options != nil {
		// Validar que ColumnSpacing no sea negativo
		if cmd.Options.ColumnSpacing < 0 {
			return fmt.Errorf("column_spacing cannot be negative")
		}
	}
	// Validate table command
	if len(cmd.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column defined")
	}

	// Create table options with defaults
	opts := &tables.Options{
		ShowHeaders:   cmd.ShowHeaders,
		WordWrap:      true,
		ColumnSpacing: 1,
		HeaderStyle:   tables.Style{Bold: true},
	}

	// Apply custom options if provided
	if cmd.Options != nil {
		opts.WordWrap = cmd.Options.WordWrap
		if cmd.Options.HeaderBold {
			opts.HeaderStyle.Bold = true
		}
		if cmd.Options.ColumnSpacing > 0 {
			opts.ColumnSpacing = cmd.Options.ColumnSpacing
		}
	}

	// Set paper width
	switch {
	case cmd.Definition.PaperWidth > 0:
		opts.PaperWidth = cmd.Definition.PaperWidth
	case printer.Profile.PrintWidth > 0:
		opts.PaperWidth = printer.Profile.PrintWidth
	default:
		if printer.Profile.PaperWidth >= 80 {
			opts.PaperWidth = tables.PaperWidth80mm
		} else {
			opts.PaperWidth = tables.PaperWidth58mm
		}
	}

	// Create table engine
	engine := tables.NewEngine(&cmd.Definition, opts)

	// Prepare table data
	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: cmd.ShowHeaders,
		Rows:        make([]tables.Row, len(cmd.Rows)),
	}

	// Convert rows
	for i, row := range cmd.Rows {
		tableData.Rows[i] = row
	}

	// Render table to string
	var buf strings.Builder
	if err := engine.Render(&buf, tableData); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Options.Align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	err := printer.Print(buf.String())
	if err != nil {
		return err
	}

	// Restaurar alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	// Send the raw output (includes ESC/POS commands for bold)
	return nil
}
