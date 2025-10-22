package pos

import (
	"fmt"
	"image"
	"log"

	"github.com/skip2/go-qrcode"

	"github.com/adcondev/pos-printer/encoding"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/barcode"
	"github.com/adcondev/pos-printer/imaging"

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
		printer.Escpos = escpos.NewEscposProtocol()
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
func (ep *EscposPrinter) GetProfile() *profile.Escpos {
	return ep.Profile
}

// SetProfile establece un nuevo perfil
func (ep *EscposPrinter) SetProfile(newProfile *profile.Escpos) {
	ep.Profile = newProfile
}

// === Aux ===

// hasProtocol verifica si algún protocolo está configurado
func (ep *EscposPrinter) hasProtocol() bool {
	_, ok := protoMap[ep.protocolType]
	return ok
}

// GetProtocolName devuelve el nombre del protocolo activo
func (ep *EscposPrinter) GetProtocolName() string {
	if !ep.hasProtocol() {
		return "Unknown Protocol"
	}
	return protoMap[ep.protocolType]
}

// === Implementación de la interfaz Printer ===

// Initialize inicializa la impresora
func (ep *EscposPrinter) Initialize() error {
	var cmd []byte
	switch protoMap[ep.protocolType] {
	case "Escpos":
		cmd = ep.Escpos.InitializePrinter()
	case "ZPL":
		// cmd := p.ZPL.InitializePrinter()
	case "PDF":
		// cmd := p.PDF.InitializePrinter()
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[ep.protocolType])
	}

	_, err := ep.Connector.Write(cmd)
	if err == nil {
		ep.initialized = true
	}
	return err
}

// Close cierra la conexión con la impresora
func (ep *EscposPrinter) Close() error {
	// Primero enviar comandos de cierre del protocolo
	// cerrar el conector
	return ep.Connector.Close()
}

// SetJustification establece la alineación del texto
func (ep *EscposPrinter) SetJustification(alignment escpos.Alignment) error {
	var cmd []byte
	var err error
	switch protoMap[ep.protocolType] {
	case "Escpos":
		cmd, err = ep.Escpos.SetJustification(alignment)
		if err != nil {
			return fmt.Errorf("error al establecer alineación: %w", err)
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[ep.protocolType])
	}

	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer alineación: %w", err)
	}
	ep.activeAlignment = alignment

	return err
}

// SetFont establece la fuente
func (ep *EscposPrinter) SetFont(font escpos.Font) error {
	var cmd []byte
	var err error
	switch protoMap[ep.protocolType] {
	case "Escpos":
		cmd, err = ep.Escpos.SelectCharacterFont(font)
		if err != nil {
			return err
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[ep.protocolType])
	}

	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer fuente: %w", err)
	}
	ep.activeFont = font

	return err
}

// SetEmphasis activa/desactiva negrita
func (ep *EscposPrinter) SetEmphasis(on escpos.EmphasizedMode) error {
	var cmd []byte
	var err error
	switch protoMap[ep.protocolType] {
	case "Escpos":
		cmd, err = ep.Escpos.TurnEmphasizedMode(on)
		if err != nil {
			return fmt.Errorf("escpos: error al establecer negrita: %w", err)
		}
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[ep.protocolType])
	}

	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al establecer negrita: %w", err)
	}
	ep.activeEmphasis = on
	return err
}

// SetDoubleStrike activa/desactiva doble golpe
func (ep *EscposPrinter) SetDoubleStrike(on bool) error {

	cmd := ep.Escpos.SetDoubleStrike(on)
	_, err := ep.Connector.Write(cmd)
	return err
}

// SetUnderline configura el subrayado
func (ep *EscposPrinter) SetUnderline(underline escpos.UnderlineMode) error {
	cmd, err := ep.Escpos.TurnUnderlineMode(underline)
	if err != nil {
		return fmt.Errorf("underline mode error: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	return err
}

// SetCharacterSet cambia el juego de caracteres activo
func (ep *EscposPrinter) SetCharacterSet(charsetCode encoding.CharacterSet) error {
	// Verificar que el charset esté soportado por el perfil
	supported := false
	for _, cs := range ep.GetSupportedCharsets() {
		if cs == charsetCode {
			supported = true
			break
		}
	}

	if !supported {
		return fmt.Errorf("el cáracter set %d no está soportado en el perfil de la impresora", charsetCode)
	}

	// Enviar comando al protocolo
	cmd, err := ep.Escpos.SelectCharacterTable(charsetCode)
	if err != nil {
		return fmt.Errorf("error al generar comando de cambio de charset: %w", err)
	}
	if _, err = ep.Connector.Write(cmd); err != nil {
		return err
	}

	// Actualizar el charset activo en el perfil
	ep.activeCharset = charsetCode

	return nil
}

// Print imprime texto con codificación apropiada
func (ep *EscposPrinter) Print(str string) error {
	// Codificar usando el charset activo del perfil
	encoded, err := encoding.EncodeString(str, ep.activeCharset)
	if err != nil {
		// Fallback: intentar con el charset por defecto
		log.Printf("Error codificando con charset %d, volviendo a default: %v",
			ep.activeCharset, err)
		encoded, err = encoding.EncodeString(str, ep.Profile.DefaultCharSet)
		if err != nil {
			return fmt.Errorf("fallo al codificar texto %s: %w", str, err)
		}
	}

	// Enviar al protocolo como bytes raw
	cmd, err := ep.Escpos.Print.Text(string(encoded))
	if err != nil {
		return fmt.Errorf("error al generar comando de impresión: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar comando de impresión: %v", err)
	}
	return nil
}

// TextLn imprime texto con salto de línea
func (ep *EscposPrinter) TextLn(str string) error {
	// Similar a PrintDataInPageMode pero agregando LF
	encoded, err := encoding.EncodeString(str+"\n", ep.activeCharset)
	if err != nil {
		log.Printf("trying default... failed to encode text \"%s\": %v", str, encoding.Registry[ep.activeCharset].Name)
		encoded, err = encoding.EncodeString(str+"\n", ep.Profile.DefaultCharSet)
		if err != nil {
			log.Printf("failed to encode text \"%s\": %v", str, encoding.Registry[ep.Profile.DefaultCharSet].Name)
			return err
		}
	}

	// El protocolo agrega el LF
	cmd, err := ep.Escpos.Print.Text(string(encoded))
	if err != nil {
		return fmt.Errorf("failed to generate print command: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("failed to send print command: %v", err)
	}
	return err
}

// Cut corta el papel
func (ep *EscposPrinter) Cut(mode escpos.CutPaper) error {
	cmd, err := ep.Escpos.Cut(mode) // 0 lines feed antes del corte
	if err != nil {
		return fmt.Errorf("error al generar comando de corte: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	return err
}

// Feed alimenta papel
func (ep *EscposPrinter) Feed(lines byte) error {
	cmd := ep.Escpos.Print.PrintAndFeedPaper(lines)
	_, err := ep.Connector.Write(cmd)
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
func (ep *EscposPrinter) PrintImage(img image.Image) error {
	opts := DefaultPrintImageOptions()
	return ep.PrintImageWithOptions(img, opts)
}

// PrintImageWithOptions imprime una imagen con opciones específicas
func (ep *EscposPrinter) PrintImageWithOptions(img image.Image, opts PrintImageOptions) error {
	// Verificar soporte de imágenes
	if !ep.Profile.HasImageSupport() {
		return fmt.Errorf("printer %s does not support images", ep.Profile.ModelInfo())
	}

	// Crear PrintRasterBitImage
	resizedImg := imaging.ResizeToWidth(img, opts.Width, ep.Profile.DotsPerLine)
	printImg := imaging.NewPrintImage(resizedImg, opts.DitherMode)
	printImg.Threshold = opts.Threshold

	// Aplicar dithering si se especificó
	if opts.DitherMode != imaging.DitherNone {
		if err := printImg.ApplyDithering(opts.DitherMode); err != nil {
			return fmt.Errorf("failed to apply dithering: %w", err)
		}
	}

	// Generar comandos
	cmd, err := ep.Escpos.PrintRasterBitImage(printImg, opts.Density)
	if err != nil {
		return fmt.Errorf("failed to generate imaging commands: %w", err)
	}

	// Enviar a la impresora
	_, err = ep.Connector.Write(cmd)
	return err
}

// PrintImageFromFile opens and prints an image from a file
func (ep *EscposPrinter) PrintImageFromFile(filename string) error {
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
	return ep.PrintImage(img)
}

// GetSupportedCharsets returns the character sets supported by the printer profile
func (ep *EscposPrinter) GetSupportedCharsets() []encoding.CharacterSet {
	// Retorna los juegos de caracteres soportados por el perfil
	return ep.Profile.CharacterSets
}

// CancelKanjiMode deactivates Kanji mode
func (ep *EscposPrinter) CancelKanjiMode() error {
	// Enviar comando para cancelar modo Kanji
	cmd := ep.Escpos.CancelKanjiMode()
	_, err := ep.Connector.Write(cmd)
	return err
}

// PrintQR imprime un código QR con datos, versión y nivel de corrección de errores
func (ep *EscposPrinter) PrintQR(
	data string,
	model escpos.QRModel,
	ecLevel escpos.QRErrorCorrection,
	moduleSize escpos.QRModuleSize,
	size int,
) error {
	// Verificar soporte
	if !ep.Profile.SupportsQR && ep.Profile.HasImageSupport() {
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
		return ep.PrintImage(qrImage)
	}

	// Usar el comando nativo
	cmdLines, err := ep.Escpos.PrintQR(data, model, moduleSize, ecLevel)
	if err != nil {
		return fmt.Errorf("error al generar QR: %w", err)
	}

	for _, cmd := range cmdLines {
		_, err = ep.Connector.Write(cmd)
		if err != nil {
			return fmt.Errorf("error al enviar QR a la impresora: %w", err)
		}
	}

	return err
}

func (ep *EscposPrinter) PrintBarcode(symbol barcode.Symbology, data string) error {
	cmd, err := ep.Escpos.Barcode.PrintBarcode(symbol, []byte(data))
	if err != nil {
		return fmt.Errorf("error al generar código de barras: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar código de barras a la impresora: %w", err)
	}
	return nil
}

func (ep *EscposPrinter) SetBarcodeHRIPosition(position barcode.HRIPosition) error {
	cmd, err := ep.Escpos.Barcode.SelectHRICharacterPosition(position)
	if err != nil {
		return fmt.Errorf("error al generar comando de posición HRI: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar comando de posición HRI a la impresora: %w", err)
	}
	return nil
}

func (ep *EscposPrinter) SetBarcodeHeight(height barcode.Height) error {
	cmd, err := ep.Escpos.Barcode.SetBarcodeHeight(height)
	if err != nil {
		return fmt.Errorf("error al generar comando de altura de código de barras: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar comando de altura de código de barras a la impresora: %w", err)
	}
	return nil
}

func (ep *EscposPrinter) SetBarcodeWidth(width barcode.Width) error {
	cmd, err := ep.Escpos.Barcode.SetBarcodeWidth(width)
	if err != nil {
		return fmt.Errorf("error al generar comando de ancho de código de barras: %w", err)
	}
	_, err = ep.Connector.Write(cmd)
	if err != nil {
		return fmt.Errorf("error al enviar comando de ancho de código de barras a la impresora: %w", err)
	}
	return nil
}

// SetTextSize sets the text size, width and height multipliers.
// 0 = 1x (Normal size),
// 1 = 2x (Double size),
// 2 = 3x (Triple size),
// 7 = 8x (Maximum size)
func (p *EscposPrinter) SetTextSize(widthMultiplier, heightMultiplier int) error {
	cmd := p.Escpos.SetTextSize(widthMultiplier, heightMultiplier)
	_, err := p.Connector.Write(cmd)
	return err
}

// Implementar el resto de métodos que necesites
