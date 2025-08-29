package pos

import (
	"fmt"
	"image"
	"log"

	"github.com/adcondev/pos-printer/encoding"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/imaging"
	"github.com/skip2/go-qrcode"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/internal"
	"github.com/adcondev/pos-printer/profile"
)

// Protocol defines the printing protocol (Escpos, ZPL, PDF, etc.
type Protocol byte

const (
	// EscposProto defines the Escpos protocol
	EscposProto Protocol = iota
	// ZplProto defines the ZPL protocol
	ZplProto
	// PdfProto defines the PDF generation protocol
	PdfProto
)

// EscposPrinter implementa Printer usando Profile y Connector
type EscposPrinter struct {
	// Componentes obligatorios
	Connector connector.Connector
	Profile   *profile.Escpos

	// Commands
	Escpos *escpos.Protocol

	// Command Types
	// PrintDataInPageMode *escpos.TextCommands

	// Protocolo activo
	protocolType Protocol

	// Opcional: Estado interno
	initialized     escpos.PrinterInitiated
	activeFont      escpos.Font
	activeAlignment escpos.Alignment
	activeUnderline escpos.UnderlineMode
	activeCharset   encoding.CharacterSet
	activeEmphasis  escpos.EmphasizedMode
}

var protoMap = map[Protocol]string{
	EscposProto: "Escpos",
	// types.ZplProto:    "ZPL",
	// types.PdfProto:    "PDF",
}

// NewPrinter crea una nueva impresora genérica
func NewPrinter(proto Protocol, conn connector.Connector, prof *profile.Escpos) (*EscposPrinter, error) {
	protoType, ok := protoMap[proto]
	if !ok {
		return nil, fmt.Errorf("not defined protocol: %d", proto)
	}
	if conn == nil {
		return nil, fmt.Errorf("el conector no puede ser nil")
	}
	if prof == nil {
		return nil, fmt.Errorf("el perfil no puede ser nil")
	}

	printer := &EscposPrinter{
		Connector:       conn,
		Profile:         prof,
		initialized:     false,
		activeFont:      escpos.FontA,
		activeAlignment: escpos.AlignLeft,
		activeUnderline: escpos.UnderNone,
		activeCharset:   prof.DefaultCharSet,
		activeEmphasis:  escpos.EmphasizedOff,
		protocolType:    proto,
	}

	switch protoMap[proto] {
	case "Escpos":
		printer.Escpos = escpos.NewEscposCommands()
	case "ZPL":
		// printer.zpl = zpl.NewZPLProtocol()
		return nil, fmt.Errorf("protocol %s not released yet", protoType)
	case "PDF":
		// printer.pdf = pdf.NewPDFProtocol()
	default:
		return nil, fmt.Errorf("protocol %s not released yet", protoType)
	}

	return printer, nil
}

// GetProfile devuelve el perfil de la impresora
func (p *EscposPrinter) GetProfile() *profile.Escpos {
	return p.Profile
}

// SetProfile establece un nuevo perfil
func (p *EscposPrinter) SetProfile(newProfile *profile.Escpos) {
	p.Profile = newProfile
}

// === Aux ===

// hasProtocol verifica si algún protocolo está configurado
func (p *EscposPrinter) hasProtocol() bool {
	_, ok := protoMap[p.protocolType]
	return ok
}

// GetProtocolName devuelve el nombre del protocolo activo
func (p *EscposPrinter) GetProtocolName() string {
	if !p.hasProtocol() {
		return "Unknown Protocol"
	}
	return protoMap[p.protocolType]
}

// === Implementación de la interfaz Printer ===

// Initialize inicializa la impresora
func (p *EscposPrinter) Initialize() error {
	var cmd []byte
	switch protoMap[p.protocolType] {
	case "Escpos":
		cmd = p.Escpos.InitializePrinter()
	case "ZPL":
		// cmd := p.ZPL.InitializePrinter()
	case "PDF":
		// cmd := p.PDF.InitializePrinter()
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[p.protocolType])
	}

	_, err := p.Connector.Write(cmd)
	if err == nil {
		p.initialized = true
	}
	return err
}

// Close cierra la conexión con la impresora
func (p *EscposPrinter) Close() error {
	// Primero enviar comandos de cierre del protocolo
	// cerrar el conector
	return p.Connector.Close()
}

// SetJustification establece la alineación del texto
func (p *EscposPrinter) SetJustification(alignment escpos.Alignment) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "Escpos":
		cmd, err = p.Escpos.SetJustification(alignment)
		if err != nil {
			return fmt.Errorf("error al establecer alineación: %w", err)
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[p.protocolType])
	}

	_, err = p.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer alineación: %w", err)
	}
	p.activeAlignment = alignment

	return err
}

// SetFont establece la fuente
func (p *EscposPrinter) SetFont(font escpos.Font) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "Escpos":
		cmd, err = p.Escpos.SelectCharacterFont(font)
		if err != nil {
			return err
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[p.protocolType])
	}

	_, err = p.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer fuente: %w", err)
	}
	p.activeFont = font

	return err
}

// SetEmphasis activa/desactiva negrita
func (p *EscposPrinter) SetEmphasis(on escpos.EmphasizedMode) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "Escpos":
		cmd, err = p.Escpos.TurnEmphasizedMode(on)
		if err != nil {
			return fmt.Errorf("escpos: error al establecer negrita: %w", err)
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[p.protocolType])
	}

	_, err = p.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer negrita: %w", err)
	}
	p.activeEmphasis = on
	return err
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *EscposPrinter) SetDoubleStrike(on bool) error {

	cmd := p.Escpos.SetDoubleStrike(on)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetUnderline configura el subrayado
func (p *EscposPrinter) SetUnderline(underline escpos.UnderlineMode) error {
	cmd, err := p.Escpos.TurnUnderlineMode(underline)
	if err != nil {
		return fmt.Errorf("underline mode error: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	return err
}

// SetCharacterSet cambia el juego de caracteres activo
func (p *EscposPrinter) SetCharacterSet(charsetCode encoding.CharacterSet) error {
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
	cmd, err := p.Escpos.SelectCharacterTable(charsetCode)
	if err != nil {
		return fmt.Errorf("error al generar comando de cambio de charset: %w", err)
	}
	if _, err = p.Connector.Write(cmd); err != nil {
		return err
	}

	// Actualizar el charset activo en el perfil
	p.activeCharset = charsetCode

	return nil
}

// Print imprime texto con codificación apropiada
func (p *EscposPrinter) Print(str string) error {
	// Codificar usando el charset activo del perfil
	encoded, err := encoding.EncodeString(str, p.activeCharset)
	if err != nil {
		// Fallback: intentar con el charset por defecto
		log.Printf("Error codificando con charset %d, volviendo a default: %v",
			p.activeCharset, err)
		encoded, err = encoding.EncodeString(str, p.Profile.DefaultCharSet)
		if err != nil {
			return fmt.Errorf("fallo al codificar texto %s: %w", str, err)
		}
	}

	// Enviar al protocolo como bytes raw
	cmd, err := p.Escpos.Print.Text(string(encoded))
	if err != nil {
		return fmt.Errorf("error al generar comando de impresión: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar comando de impresión: %v", err)
	}
	return nil
}

// TextLn imprime texto con salto de línea
func (p *EscposPrinter) TextLn(str string) error {
	// Similar a PrintDataInPageMode pero agregando LF
	encoded, err := encoding.EncodeString(str+"\n", p.activeCharset)
	if err != nil {
		log.Printf("trying default... failed to encode text \"%s\": %v", str, encoding.Registry[p.activeCharset].Name)
		encoded, err = encoding.EncodeString(str+"\n", p.Profile.DefaultCharSet)
		if err != nil {
			log.Printf("failed to encode text \"%s\": %v", str, encoding.Registry[p.Profile.DefaultCharSet].Name)
			return err
		}
	}

	// El protocolo agrega el LF
	cmd, err := p.Escpos.Print.Text(string(encoded))
	if err != nil {
		return fmt.Errorf("failed to generate print command: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("failed to send print command: %v", err)
	}
	return err
}

// Cut corta el papel
func (p *EscposPrinter) Cut(mode escpos.CutPaper) error {
	cmd, err := p.Escpos.Cut(mode) // 0 lines feed antes del corte
	if err != nil {
		return fmt.Errorf("error al generar comando de corte: %w", err)
	}
	_, err = p.Connector.Write(cmd)
	return err
}

// Feed alimenta papel
func (p *EscposPrinter) Feed(lines byte) error {
	cmd := p.Escpos.Print.PrintAndFeedPaper(lines)
	_, err := p.Connector.Write(cmd)
	return err
}

// PrintImageOptions contiene opciones para imprimir imágenes
type PrintImageOptions struct {
	Density    escpos.Density
	DitherMode imaging.DitherMode
	Threshold  uint8
	Width      int
}

// DefaultPrintImageOptions devuelve opciones por defecto
func DefaultPrintImageOptions() PrintImageOptions {
	return PrintImageOptions{
		Density:    escpos.DensitySingle,
		DitherMode: imaging.DitherNone,
		Threshold:  128,
		Width:      256,
	}
}

// PrintImage imprime una imagen con opciones por defecto
func (p *EscposPrinter) PrintImage(img image.Image) error {
	opts := DefaultPrintImageOptions()
	return p.PrintImageWithOptions(img, opts)
}

// PrintImageWithOptions imprime una imagen con opciones específicas
func (p *EscposPrinter) PrintImageWithOptions(img image.Image, opts PrintImageOptions) error {
	// Verificar soporte de imágenes
	if !p.Profile.HasImageSupport() {
		return fmt.Errorf("printer %s does not support images", p.Profile.ModelInfo())
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
	cmd, err := p.Escpos.PrintRasterBitImage(printImg, opts.Density)
	if err != nil {
		return fmt.Errorf("failed to generate imaging commands: %w", err)
	}

	// Enviar a la impresora
	_, err = p.Connector.Write(cmd)
	return err
}

// PrintImageFromFile opens and prints an image from a file
func (p *EscposPrinter) PrintImageFromFile(filename string) error {
	file, err := internal.SafeOpen(filename)
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

// GetSupportedCharsets returns the character sets supported by the printer profile
func (p *EscposPrinter) GetSupportedCharsets() []encoding.CharacterSet {
	// Retorna los juegos de caracteres soportados por el perfil
	return p.Profile.CharacterSets
}

// CancelKanjiMode deactivates Kanji mode
func (p *EscposPrinter) CancelKanjiMode() error {
	// Enviar comando para cancelar modo Kanji
	cmd := p.Escpos.CancelKanjiMode()
	_, err := p.Connector.Write(cmd)
	return err
}

// PrintQR imprime un código QR con datos, versión y nivel de corrección de errores
func (p *EscposPrinter) PrintQR(
	data string,
	model escpos.QRModel,
	ecLevel escpos.QRErrorCorrection,
	moduleSize escpos.QRModuleSize,
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
	cmdLines, err := p.Escpos.PrintQR(data, model, moduleSize, ecLevel)
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
