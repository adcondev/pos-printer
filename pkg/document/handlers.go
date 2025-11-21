package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/internal/load"
	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/printer"
	"github.com/adcondev/pos-printer/pkg/tables"
)

const (
	center = "center"
	right  = "right"
	left   = "left"
)

// handleText manages text commands
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Si hay label con diferente alineación que el content
	if cmd.Label.Text != "" {
		// Aplicar estilo del label
		if err := e.applyTextStyle(printer, cmd.Label.Style); err != nil {
			return err
		}

		// Construir texto del label
		labelText := cmd.Label.Text
		if cmd.Label.Separator != "" {
			labelText += cmd.Label.Separator
		}

		// Si las alineaciones son diferentes, imprimir label y resetear
		if cmd.Label.Style.Align != cmd.Content.Style.Align {
			if err := printer.PrintLine(labelText); err != nil {
				return err
			}
			// Reset completo antes del content
			if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
				return err
			}
		} else {
			// Misma alineación, imprimir inline con espacio
			if !strings.HasSuffix(labelText, " ") {
				labelText += " "
			}
			if err := printer.Print(labelText); err != nil {
				return err
			}
			// Reset solo los estilos diferentes
			if err := e.resetDifferingStyles(printer, cmd.Label.Style, cmd.Content.Style); err != nil {
				return err
			}
		}
	}

	// Aplicar estilo del contenido
	if err := e.applyTextStyle(printer, cmd.Content.Style); err != nil {
		return fmt.Errorf("failed to apply content style: %w", err)
	}

	// Imprimir contenido
	if cmd.Content.Text != "" {
		if cmd.NewLine {
			if err := printer.PrintLine(cmd.Content.Text); err != nil {
				return err
			}
		} else {
			if err := printer.Print(cmd.Content.Text); err != nil {
				return err
			}
		}
	}

	// Resetear estilos del contenido
	if err := e.resetTextStyle(printer, cmd.Content.Style); err != nil {
		return fmt.Errorf("failed to reset content style: %w", err)
	}

	return nil
}

func (e *Executor) resetDifferingStyles(printer *service.Printer, labelStyle, contentStyle TextStyle) error {
	// Reset bold si difiere
	if labelStyle.Bold != contentStyle.Bold {
		if labelStyle.Bold {
			if err := printer.DisableBold(); err != nil {
				return err
			}
		} else {
			if err := printer.EnableBold(); err != nil {
				return err
			}
		}
	}

	// Reset size si difiere
	if labelStyle.Size != contentStyle.Size {
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	// Reset underline si difiere
	if labelStyle.Underline != contentStyle.Underline {
		if err := printer.NoDot(); err != nil {
			return err
		}
	}

	// Reset inverse si difiere
	if labelStyle.Inverse != contentStyle.Inverse {
		if labelStyle.Inverse {
			if err := printer.InverseOff(); err != nil {
				return err
			}
		} else {
			if err := printer.InverseOn(); err != nil {
				return err
			}
		}
	}

	// Reset font si difiere
	if labelStyle.Font != contentStyle.Font {
		if err := printer.FontA(); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) applyAlign(printer *service.Printer, align string) error {
	// Aplicar alineación
	switch strings.ToLower(align) {
	case center:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case right:
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case left:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		log.Printf("Unknown alignment: %s, using left", align)
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) applySize(printer *service.Printer, size string) error {
	if size != "" {
		switch ss := strings.ToLower(size); ss {
		case "1x1", "1":
			if err := printer.SingleSize(); err != nil {
				return err
			}
		case "2x2", "2":
			if err := printer.DoubleSize(); err != nil {
				return err
			}
		case "3x3", "3":
			if err := printer.TripleSize(); err != nil {
				return err
			}
		case "4x4", "4":
			if err := printer.QuadraSize(); err != nil {
				return err
			}
		case "5x5", "5":
			if err := printer.PentaSize(); err != nil {
				return err
			}
		case "6x6", "6":
			if err := printer.HexaSize(); err != nil {
				return err
			}
		case "7x7", "7":
			if err := printer.HeptaSize(); err != nil {
				return err
			}
		case "8x8", "8":
			if err := printer.OctaSize(); err != nil {
				return err
			}
		default:
			// Intentar parsear tamaño personalizado WxH
			if len(ss) == 3 && ss[1] == 'x' {
				parts := strings.Split(ss, "x")
				widthMultiplier := parts[0][0] - '0'
				heightMultiplier := parts[1][0] - '0'
				if err := printer.CustomSize(widthMultiplier, heightMultiplier); err != nil {
					return err
				}
				log.Printf("Applied custom text size: %s", size)
			} else {
				if err := printer.SingleSize(); err != nil {
					return err
				}
				log.Printf("Unknown text size: %s, using single size", size)
			}
		}
	} else {
		// Default size
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) applyUnderline(printer *service.Printer, underline string) error {
	switch strings.ToLower(underline) {
	case "0", "0pt":
		// No underline
		err := printer.NoDot()
		if err != nil {
			return err
		}
	case "1", "1pt":
		err := printer.OneDot()
		if err != nil {
			return err
		}
	case "2", "2pt":
		err := printer.TwoDot()
		if err != nil {
			return err
		}
	default:
		err := printer.NoDot()
		if err != nil {
			return err
		}
		log.Printf("Unknown underline style: %s, using none", underline)
	}
	return nil
}

func (e *Executor) applyFont(printer *service.Printer, font string) error {
	switch strings.ToLower(font) {
	case "a":
		if err := printer.FontA(); err != nil {
			return err
		}
	case "b":
		if err := printer.FontB(); err != nil {
			return err
		}
	default:
		log.Printf("Unknown font: %s, using Font A", font)
		if err := printer.FontA(); err != nil {
			return err
		}
	}
	return nil
}

// applyTextStyle aplica los estilos de texto especificados
func (e *Executor) applyTextStyle(printer *service.Printer, style TextStyle) error {
	// Aplicar alineación
	err := e.applyAlign(printer, style.Align)
	if err != nil {
		return fmt.Errorf("failed to apply alignment: %w", err)
	}

	// Aplicar bold
	if style.Bold {
		if err := printer.EnableBold(); err != nil {
			return err
		}
	}

	// Aplicar tamaño
	err = e.applySize(printer, style.Size)
	if err != nil {
		return fmt.Errorf("failed to apply size: %w", err)
	}

	// Aplicar underline
	err = e.applyUnderline(printer, style.Underline)
	if err != nil {
		return fmt.Errorf("failed to apply underline: %w", err)
	}

	// Aplicar inverse
	if style.Inverse {
		err := printer.InverseOn()
		if err != nil {
			return err
		}
	}

	// Apply font
	err = e.applyFont(printer, style.Font)
	if err != nil {
		return fmt.Errorf("failed to apply font: %w", err)
	}

	return nil
}

// resetTextStyle resetea los estilos aplicados
func (e *Executor) resetTextStyle(printer *service.Printer, style TextStyle) error {
	// Reset bold
	if style.Bold {
		if err := printer.DisableBold(); err != nil {
			return err
		}
	}

	// Reset size
	if style.Size != "" {
		if err := printer.SingleSize(); err != nil {
			return err
		}
	}

	// Reset underline
	if style.Underline != "" {
		if err := printer.NoDot(); err != nil {
			return err
		}
	}

	// Reset inverse
	if style.Inverse {
		if err := printer.InverseOff(); err != nil {
			return err
		}
	}

	// Reset font
	if style.Font != "" {
		if err := printer.FontA(); err != nil {
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

// TODO: Add multiple char support for separator (_-=-, etc.)

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
		// TODO: Verify the following line for different fonts
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
