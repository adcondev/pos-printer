package document

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/pos-printer/pkg/controllers/escpos/character"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/service"
)

// handleText maneja comandos de texto
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Style.Align) {
	case "center":
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case "right":
		if err := printer.AlignRight(); err != nil {
			return err
		}
	default:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	}

	// Aplicar estilo bold
	if cmd.Style.Bold {
		if err := printer.Bold(); err != nil {
			return err
		}
	}

	// Aplicar tamaño
	switch strings.ToLower(cmd.Style.Size) {
	case "2x2":
		if err := printer.DoubleSize(); err != nil {
			return err
		}
	case "3x3":
		// Implementar si la impresora soporta
		size, _ := character.NewSize(3, 3)
		cmdBytes := printer.Protocol.Character.SelectCharacterSize(size)
		if err := printer.Write(cmdBytes); err != nil {
			return err
		}
	default:
		if err := printer.NormalSize(); err != nil {
			return err
		}
	}

	// Imprimir texto
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

	// Resetear estilos para el siguiente comando
	if cmd.Style.Bold {
		cmdBytes := printer.Protocol.Character.SetEmphasizedMode(character.OffEm)
		err := printer.Write(cmdBytes)
		if err != nil {
			return err
		}
	}

	if cmd.Style.Size != "" && cmd.Style.Size != "normal" {
		err := printer.NormalSize()
		if err != nil {
			return err
		}
	}

	return nil
}

// handleImage maneja comandos de imagen
func (e *Executor) handleImage(printer *service.Printer, data json.RawMessage) error {
	var cmd ImageCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse image command: %w", err)
	}

	// Decodificar imagen desde base64
	img, err := graphics.LoadFromBase64(cmd.Code)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	// Configurar opciones de procesamiento
	opts := &graphics.Options{
		Width:          cmd.Width,
		Threshold:      cmd.Threshold,
		PreserveAspect: true,
		AutoRotate:     false,
	}

	// Si no se especifica ancho, usar el ancho del perfil
	if opts.Width == 0 {
		opts.Width = e.profile.DotsPerLine
	}

	// Si no se especifica threshold, usar valor por defecto
	if opts.Threshold == 0 {
		opts.Threshold = 128
	}

	// Configurar dithering
	switch strings.ToLower(cmd.Dithering) {
	case "atkinson":
		opts.Mode = graphics.Atkinson
	default:
		opts.Mode = graphics.Threshold
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Align) {
	case "center":
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case "right":
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

// handleSeparator maneja comandos de separador
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
		cmd.Length = e.profile.DotsPerLine / 12 // Aproximación para Font A
	}

	// Construir línea separadora
	line := strings.Repeat(cmd.Char, cmd.Length)

	return printer.PrintLine(line)
}

// handleFeed maneja comandos de avance de papel
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

// handleCut maneja comandos de corte
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
		return printer.FullCut(0)
	default: // partial
		return printer.PartialFeedAndCut(0)
	}
}

// handleQRPlaceholder placeholder para comando QR
func (e *Executor) handleQRPlaceholder(printer *service.Printer, data json.RawMessage) error {
	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}

	// TODO: Implementar cuando se agregue soporte QR
	// Por ahora, imprimir texto indicando QR
	err := printer.AlignCenter()
	if err != nil {
		return err
	}
	err = printer.PrintLine("[QR Code]")
	if err != nil {
		return err
	}
	err = printer.PrintLine(cmd.Data)
	if err != nil {
		return err
	}
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
}

// handleTablePlaceholder placeholder para comando de tabla
func (e *Executor) handleTablePlaceholder(printer *service.Printer, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}

	// TODO: Implementar renderizado de tablas
	// Por ahora, imprimir filas como texto simple
	for _, row := range cmd.Rows {
		line := strings.Join(row, " | ")
		err := printer.PrintLine(line)
		if err != nil {
			return err
		}
	}

	return nil
}
