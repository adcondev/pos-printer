package posprinter

import (
	"fmt"
	"image"
	"log"

	"github.com/AdConDev/pos-printer/encoding"
	"github.com/AdConDev/pos-printer/imaging"
	"github.com/skip2/go-qrcode"

	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/profile"
	"github.com/AdConDev/pos-printer/protocol"
	"github.com/AdConDev/pos-printer/types"
	"github.com/AdConDev/pos-printer/utils"
)

// Printer define la interfaz de alto nivel para cualquier impresora
type Printer interface {
	// === Comandos básicos ===
	Initialize() error
	Close() error

	// === Formato de texto ===
	SetJustification(alignment types.Alignment) error
	SetFont(font types.Font) error
	SetEmphasis(on bool) error
	SetDoubleStrike(on bool) error
	SetUnderline(underline types.UnderlineMode) error

	// === Impresión de texto ===
	Text(str string) error
	TextLn(str string) error

	// Code Page y Character Sets
	CancelKanjiMode() error
	SetCharacterSet(charsetCode int) error

	// === Control de papel ===
	Cut(mode types.CutMode, lines int) error
	Feed(lines int) error

	// === Impresión de imágenes ===
	PrintImage(img image.Image) error
	PrintImageFromFile(filename string) error

	// TODO: Agregar más métodos según necesites
}

// GenericPrinter implementa Printer usando Protocol y Connector
type GenericPrinter struct {
	Protocol  protocol.Protocol
	Connector connector.Connector
	Profile   *profile.Profile

	// TODO: Agregar más campos si necesitas:
	// - Estado actual (font, alignment, etc.)
	// - Buffer de comandos
	// - Configuración
}

// NewGenericPrinter crea una nueva impresora genérica
func NewGenericPrinter(proto protocol.Protocol, conn connector.Connector, prof *profile.Profile) (*GenericPrinter, error) {
	if proto == nil {
		return nil, fmt.Errorf("el protocolo no puede ser nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("el conector no puede ser nil")
	}
	if prof == nil {
		return nil, fmt.Errorf("el perfil no puede ser nil")
	}

	printer := &GenericPrinter{
		Protocol:  proto,
		Connector: conn,
		Profile:   prof,
	}

	// Inicializar la impresora
	if err := printer.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize printer: %w", err)
	}

	return printer, nil
}

// GetProfile devuelve el perfil de la impresora
func (p *GenericPrinter) GetProfile() *profile.Profile {
	return p.Profile
}

// SetProfile establece un nuevo perfil
func (p *GenericPrinter) SetProfile(newProfile *profile.Profile) {
	p.Profile = newProfile
}

// === Implementación de la interfaz Printer ===

// Initialize inicializa la impresora
func (p *GenericPrinter) Initialize() error {
	cmd := p.Protocol.InitializePrinter()
	_, err := p.Connector.Write(cmd)
	return err
}

// Close cierra la conexión con la impresora
func (p *GenericPrinter) Close() error {
	// Primero enviar comandos de cierre del protocolo
	if closeCmd := p.Protocol.Close(); len(closeCmd) > 0 {
		_, _ = p.Connector.Write(closeCmd) // Ignorar error, vamos a cerrar de todos modos
	}

	// Luego cerrar el conector
	return p.Connector.Close()
}

// SetJustification establece la alineación del texto
func (p *GenericPrinter) SetJustification(alignment types.Alignment) error {
	cmd := p.Protocol.SetJustification(alignment)
	_, err := p.Connector.Write(cmd)
	// TODO: Si no hay error, guardar el estado actual
	return err
}

// SetFont establece la fuente
func (p *GenericPrinter) SetFont(font types.Font) error {
	cmd, err := p.Protocol.SelectCharacterFont(font)
	if err != nil {
		return fmt.Errorf("character selection error: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	return err
}

// SetEmphasis activa/desactiva negrita
func (p *GenericPrinter) SetEmphasis(on types.EmphasizedMode) error {
	cmd, err := p.Protocol.TurnEmphasizedMode(on)
	if err != nil {
		return err
	}
	_, err = p.Connector.Write(cmd)
	return err
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *GenericPrinter) SetDoubleStrike(on bool) error {
	cmd := p.Protocol.SetDoubleStrike(on)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetUnderline configura el subrayado
func (p *GenericPrinter) SetUnderline(underline types.UnderlineMode) error {
	cmd, err := p.Protocol.TurnUnderlineMode(underline)
	if err != nil {
		return fmt.Errorf("underline mode error: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	return err
}

// SetCharacterSet cambia el juego de caracteres activo
func (p *GenericPrinter) SetCharacterSet(charsetCode types.CharacterSet) error {
	// Verificar que el charset esté soportado por el perfil
	supported := false
	for _, cs := range p.GetSupportedCharsets() {
		if cs == charsetCode {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("el cáracter set %d no está soportado en el perfil de la impresora", charsetCode)
	}

	// Enviar comando al protocolo
	cmd := p.Protocol.SelectCharacterTable(charsetCode)
	if _, err := p.Connector.Write(cmd); err != nil {
		return err
	}

	// Actualizar el charset activo en el perfil
	p.Profile.ActiveCharSet = charsetCode

	return nil
}

// Text imprime texto con codificación apropiada
func (p *GenericPrinter) Text(str string) error {
	// Codificar usando el charset activo del perfil
	encoded, err := encoding.EncodeString(str, p.Profile.ActiveCharSet)
	if err != nil {
		// Fallback: intentar con el charset por defecto
		log.Printf("Error codificando con charset %d, volviendo a default: %v",
			p.Profile.ActiveCharSet, err)
		encoded, err = encoding.EncodeString(str, p.Profile.DefaultCharSet)
		if err != nil {
			return fmt.Errorf("fallo al codificar texto: %w", err)
		}
	}

	// Enviar al protocolo como bytes raw
	cmd := p.Protocol.Text(string(encoded))
	_, err = p.Connector.Write(cmd)
	return err
}

// TextLn imprime texto con salto de línea
func (p *GenericPrinter) TextLn(str string) error {
	// Similar a Text pero agregando LF
	encoded, err := encoding.EncodeString(str, p.Profile.ActiveCharSet)
	if err != nil {
		encoded, err = encoding.EncodeString(str, p.Profile.DefaultCharSet)
		if err != nil {
			log.Printf("failed to encode text: %v", err)
			cmd := p.Protocol.TextLn(" err ")
			_, err = p.Connector.Write(cmd)
			return err
		}
	}

	// El protocolo agrega el LF
	cmd := p.Protocol.TextLn(string(encoded))
	_, err = p.Connector.Write(cmd)
	return err
}

// Cut corta el papel
func (p *GenericPrinter) Cut(mode types.CutMode, lines int) error {
	// TODO: Verificar si la impresora tiene cutter con HasCapability
	cmd := p.Protocol.Cut(mode, lines) // 0 lines feed antes del corte
	_, err := p.Connector.Write(cmd)
	return err
}

// Feed alimenta papel
func (p *GenericPrinter) Feed(lines int) error {
	cmd := p.Protocol.Feed(lines)
	_, err := p.Connector.Write(cmd)
	return err
}

// PrintImageOptions contiene opciones para imprimir imágenes
type PrintImageOptions struct {
	Density    types.Density
	DitherMode imaging.DitherMode
	Threshold  uint8
	Width      int
}

// DefaultPrintImageOptions devuelve opciones por defecto
func DefaultPrintImageOptions() PrintImageOptions {
	return PrintImageOptions{
		Density:    types.DensitySingle,
		DitherMode: imaging.DitherNone,
		Threshold:  128,
		Width:      256,
	}
}

// PrintImage imprime una imagen con opciones por defecto
func (p *GenericPrinter) PrintImage(img image.Image) error {
	opts := DefaultPrintImageOptions()
	return p.PrintImageWithOptions(img, opts)
}

// PrintImageWithOptions imprime una imagen con opciones específicas
func (p *GenericPrinter) PrintImageWithOptions(img image.Image, opts PrintImageOptions) error {
	// Verificar soporte de imágenes
	if !p.Profile.HasImageSupport() {
		return fmt.Errorf("protocol %s does not support images", p.Protocol.Name())
	}

	// Crear PrintRasterBitImage
	resizedImg := imaging.ResizeToWidth(img, opts.Width, p.Profile.DotsPerLine)
	printImg := imaging.NewPrintImage(resizedImg, opts.DitherMode)
	printImg.Threshold = opts.Threshold

	// Aplicar dithering si se especificó
	if opts.DitherMode != imaging.DitherNone {
		if err := printImg.ApplyDithering(opts.DitherMode); err != nil {
			return fmt.Errorf("failed to apply dithering: %w", err)
		}
	}

	// Generar comandos
	cmd, err := p.Protocol.PrintRasterBitImage(printImg, opts.Density)
	if err != nil {
		return fmt.Errorf("failed to generate imaging commands: %w", err)
	}

	// Enviar a la impresora
	_, err = p.Connector.Write(cmd)
	return err
}

// TODO: Implementar PrintImageFromFile en GenericPrinter

func (p *GenericPrinter) PrintImageFromFile(filename string) error {
	file, err := utils.SafeOpen(filename)
	if err != nil {
		return fmt.Errorf("failed to open imaging file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("printer: error al cerrar imagen: %s", cerr)
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode imaging: %w", err)
	}
	return p.PrintImage(img)
}

func (p *GenericPrinter) GetSupportedCharsets() []types.CharacterSet {
	// Retorna los juegos de caracteres soportados por el perfil
	return p.Profile.CharacterSets
}

func (p *GenericPrinter) CancelKanjiMode() error {
	// Enviar comando para cancelar modo Kanji
	cmd := p.Protocol.CancelKanjiMode()
	_, err := p.Connector.Write(cmd)
	return err
}

// FIXME: Revisar los hardcodeos de ecLevel, tamaño y componentType. Usar constantes definidas de interfaz y traducción a protocolo.
// TODO: Que tanto revisar que tanto varia el comando entre protocolos. Asegurar compatibilidad.

// PrintQR imprime un código QR con datos, versión y nivel de corrección de errores
func (p *GenericPrinter) PrintQR(
	data string,
	model types.QRModel,
	ecLevel types.QRErrorCorrection,
	moduleSize types.QRModuleSize,
	size int,
) error {
	// Verificar soporte
	if !p.Profile.SupportsQR && p.Profile.HasImageSupport() {
		log.Printf("El perfil no soporta QR nativo, usando imagen como fallback")
		// Fallback a imagen
		qrCode, err := qrcode.New(data, qrcode.RecoveryLevel(ecLevel))
		if err != nil {
			return fmt.Errorf("error al generar QR desde imagen: %w", err)
		}
		if size <= 0 {
			return fmt.Errorf("tamaño inválido para QR: %d", size)
		}
		qrImage := qrCode.Image(size)
		return p.PrintImage(qrImage)
	}

	// Usar el comando nativo
	cmdLines, err := p.Protocol.PrintQR(data, model, moduleSize, ecLevel)
	if err != nil {
		return fmt.Errorf("error al generar QR: %w", err)
	}

	for _, cmd := range cmdLines {
		_, err = p.Connector.Write(cmd)
		if err != nil {
			return fmt.Errorf("error al enviar QR a la impresora: %w", err)
		}
	}

	return err
}

// Implementar el resto de métodos que necesites
