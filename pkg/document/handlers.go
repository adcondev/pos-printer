package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/commands/character"
	posqr "github.com/adcondev/pos-printer/pkg/commands/qrcode"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/printer"
)

const (
	center = "center"
	right  = "right"
)

// handleText manages text commands
func (e *Executor) handleText(printer *service.Printer, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	// Aplicar alineación
	switch strings.ToLower(cmd.Style.Align) {
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

	// Reset styles for subsequent command
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

// handleImage manages image commands
func (e *Executor) handleImage(printer *service.Printer, data json.RawMessage) error {
	var cmd ImageCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse image command: %w", err)
	}

	// Decodificar imagen desde base64
	img, err := graphics.ImgFromBase64(cmd.Code)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

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
		return fmt.Errorf("QR data too long: %d bytes (maximum %d)", len(data), posqr.MaxDataLength)
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

	if cmd.LogoPath != "" {
		opts.LogoPath = cmd.LogoPath
		if cmd.LogoSizeMulti > 0 {
			opts.LogoSizeMulti = cmd.LogoSizeMulti
		}
	}

	if cmd.HalftonePath != "" {
		opts.HalftonePath = cmd.HalftonePath
		// Si hay halftone, desactivar circle shape
		if cmd.CircleShape {
			log.Printf("warning: halftone and circle_shape cannot be used together, prioritizing halftone")
		}
		opts.CircleShape = false
	} else {
		// Solo aplicar circle shape si no hay halftone
		opts.CircleShape = cmd.CircleShape
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

	return err
}

// handleTablePlaceholder manages table commands (WIP)
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
