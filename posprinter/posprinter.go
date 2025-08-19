package posprinter

import (
	"fmt"
	"image"
	"log"

	"github.com/AdConDev/pos-printer/encoding"
	"github.com/AdConDev/pos-printer/imaging"
	"github.com/AdConDev/pos-printer/protocol/escpos"
	"github.com/AdConDev/pos-printer/types"
	"github.com/skip2/go-qrcode"

	"github.com/AdConDev/pos-printer/connector"
	"github.com/AdConDev/pos-printer/profile"
	"github.com/AdConDev/pos-printer/utils"
)

// GenericPrinter implementa Printer usando Protocol y Connector
type GenericPrinter struct {
	// Componentes obligatorios
	Connector connector.Connector
	Profile   *profile.Profile

	// Protocolos - solo uno debe estar activo (no-nil)
	ESCPOS *escpos.Commands
	// ZPL *zpl.Commands
	// PDF *pdf.Commands

	// Protocolo activo
	protocolType types.Protocol

	// Opcional: Estado interno
	initialized     types.PrinterInitiated
	activeFont      types.Font
	activeAlignment types.Alignment
	activeUnderline types.UnderlineMode
	activeCharset   types.CharacterSet
	activeEmphasis  types.EmphasizedMode
}

var protoMap = map[types.Protocol]string{
	types.EscposProto: "ESCPOS",
	types.ZplProto:    "ZPL",
	types.PdfProto:    "PDF",
}

// NewGenericPrinter crea una nueva impresora genérica
func NewGenericPrinter(proto types.Protocol, conn connector.Connector, prof *profile.Profile) (*GenericPrinter, error) {
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

	printer := &GenericPrinter{
		Connector:       conn,
		Profile:         prof,
		initialized:     false,
		activeFont:      types.FontA,
		activeAlignment: types.AlignLeft,
		activeUnderline: types.UnderNone,
		activeCharset:   prof.DefaultCharSet,
		activeEmphasis:  types.EmphOff,
		protocolType:    proto,
	}

	switch protoMap[proto] {
	case "ESCPOS":
		printer.ESCPOS = escpos.NewESCPOSProtocol()
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
func (p *GenericPrinter) GetProfile() *profile.Profile {
	return p.Profile
}

// SetProfile establece un nuevo perfil
func (p *GenericPrinter) SetProfile(newProfile *profile.Profile) {
	p.Profile = newProfile
}

// === Aux ===

// hasProtocol verifica si algún protocolo está configurado
func (p *GenericPrinter) hasProtocol() bool {
	_, ok := protoMap[p.protocolType]
	return ok
}

// GetProtocolName devuelve el nombre del protocolo activo
func (p *GenericPrinter) GetProtocolName() string {
	if !p.hasProtocol() {
		return "Unknown Protocol"
	}
	return protoMap[p.protocolType]
}

// === Implementación de la interfaz Printer ===

// Initialize inicializa la impresora
func (p *GenericPrinter) Initialize() error {
	var cmd []byte
	switch protoMap[p.protocolType] {
	case "ESCPOS":
		cmd = p.ESCPOS.InitializePrinter()
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
func (p *GenericPrinter) Close() error {
	// Primero enviar comandos de cierre del protocolo
	var cmd []byte
	switch protoMap[p.protocolType] {
	case "ESCPOS":
		cmd = p.ESCPOS.Close()
	case "ZPL":
		// ZPL
	case "PDF":
		// PDF
	default:
		return fmt.Errorf("protocol %s not implemented", protoMap[p.protocolType])
	}
	// Luego cerrar el conector
	_, err := p.Connector.Write(cmd)
	if err == nil {
		p.initialized = false
	}
	return p.Connector.Close()
}

// SetJustification establece la alineación del texto
func (p *GenericPrinter) SetJustification(alignment types.Alignment) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "ESCPOS":
		cmd, err = p.ESCPOS.SetJustification(alignment)
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
	} else {
		p.activeAlignment = alignment
	}

	return err
}

// SetFont establece la fuente
func (p *GenericPrinter) SetFont(font types.Font) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "ESCPOS":
		cmd, err = p.ESCPOS.SelectCharacterFont(font)
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
	} else {
		p.activeFont = font
	}

	return err
}

// SetEmphasis activa/desactiva negrita
func (p *GenericPrinter) SetEmphasis(on types.EmphasizedMode) error {
	var cmd []byte
	var err error
	switch protoMap[p.protocolType] {
	case "ESCPOS":
		cmd, err = p.ESCPOS.TurnEmphasizedMode(on)
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
	} else {
		p.activeEmphasis = on
	}
	return err
}

// SetDoubleStrike activa/desactiva doble golpe
func (p *GenericPrinter) SetDoubleStrike(on bool) error {
	
	cmd := p.ESCPOS.SetDoubleStrike(on)
	_, err := p.Connector.Write(cmd)
	return err
}

// SetUnderline configura el subrayado
func (p *GenericPrinter) SetUnderline(underline types.UnderlineMode) error {
	cmd, err := p.ESCPOS.TurnUnderlineMode(underline)
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
	cmd := p.ESCPOS.SelectCharacterTable(charsetCode)
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
	cmd := p.ESCPOS.Text(string(encoded))
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
			cmd := p.ESCPOS.TextLn(" err ")
			_, err = p.Connector.Write(cmd)
			return err
		}
	}

	// El protocolo agrega el LF
	cmd := p.ESCPOS.TextLn(string(encoded))
	_, err = p.Connector.Write(cmd)
	return err
}

// Cut corta el papel
func (p *GenericPrinter) Cut(mode types.CutMode, lines int) error {
	// TODO: Verificar si la impresora tiene cutter con HasCapability
	cmd := p.ESCPOS.Cut(mode, lines) // 0 lines feed antes del corte
	_, err := p.Connector.Write(cmd)
	return err
}

// Feed alimenta papel
func (p *GenericPrinter) Feed(lines int) error {
	cmd := p.ESCPOS.Feed(lines)
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
	cmd, err := p.ESCPOS.PrintRasterBitImage(printImg, opts.Density)
	if err != nil {
		return fmt.Errorf("failed to generate imaging commands: %w", err)
	}

	// Enviar a la impresora
	_, err = p.Connector.Write(cmd)
	return err
}

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
	cmd := p.ESCPOS.CancelKanjiMode()
	_, err := p.Connector.Write(cmd)
	return err
}

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
	cmdLines, err := p.ESCPOS.PrintQR(data, model, moduleSize, ecLevel)
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
