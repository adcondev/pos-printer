package escpos

// Alignment define las alineaciones de texto estándar
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
	AlignJustified // Algunos protocolos podrían soportar esto
)

// Font define los tipos de fuente estándar
type Font byte

const (
	FontA Font = iota
	FontB
	FontC
	FontD
	FontE
	SpecialA
	SpecialB
)

// UnderlineMode define los modos de subrayado estándar
type UnderlineMode byte

const (
	UnderNone UnderlineMode = iota
	UnderSingle
	UnderDouble
)

// BarcodeType define los tipos de código de barras estándar
type BarcodeType int

const (
	UPCA BarcodeType = iota
	UPCE
	EAN13
	EAN8
	Code39
	ITF
	Codebar
)

// TextPositionBarcode define posiciones estándar para texto en códigos de barras
type TextPositionBarcode byte

const (
	NonePosBarcode TextPositionBarcode = iota
	AbovePosBarcode
	BelowPosBarcode
	BothPosBarcode
)

type BarcodeHeight byte

type BarcodeWidth byte

const (
	ExtraSmallWidth BarcodeWidth = iota
	SmallWidth
	MediumWidth
	LargeWidth
	ExtraLargeWidth
)

// CutPaper define modos de corte estándar
type CutPaper byte

const (
	FullCut CutPaper = iota
	PartialCut
)

// Density define densidades de impresión estándar para imágenes
type Density int

const (
	DensitySingle Density = iota
	DensityDouble
	DensityTriple
	DensityQuadruple
)

// QRModel define los modelos de código QR estándar
type QRModel byte

const (
	Model1 QRModel = iota // Modelo 1 (estándar)
	Model2                // Modelo 2 (recomendado y estándar)
)

// QRErrorCorrection define los niveles de corrección de errores estándar para códigos QR
type QRErrorCorrection byte

const (
	ECLow     QRErrorCorrection = iota // 7% de corrección
	ECMedium                           // 15% de corrección
	ECHigh                             // 25% de corrección
	ECHighest                          // 30% de corrección
)

// FIXME: Corregir los types con Min y Max, dejarlo en ESCPOS

// QRModuleSize define los tamaños de módulo estándar para códigos QR
type QRModuleSize byte

const (
	MinType QRModuleSize = 1
	MaxType QRModuleSize = 16
)

type Lines byte

// RealTimeStatus define los estados de tiempo real de la impresora
type RealTimeStatus byte

const (
	PrinterStatus RealTimeStatus = iota
	OfflineStatus
	ErrorStatus
	PaperSensorStatus
)

// CashDrawerPin define los pines del cajón de dinero
type CashDrawerPin byte

const (
	Pin2 CashDrawerPin = iota // Pin 1
	Pin5                      // Pin 2
)

// CashDrawerTimePulse define los tiempos de pulso del cajón de dinero
type CashDrawerTimePulse byte

const (
	Pulse100ms CashDrawerTimePulse = iota // 1x (100 ms)
	Pulse200ms                            // 2x (200 ms)
	Pulse300ms                            // 3x (300 ms)
	Pulse400ms                            // 4x (400 ms)
	Pulse500ms                            // 5x (500 ms)
	Pulse600ms                            // 6x (600 ms)
	Pulse700ms                            // 7x (700 ms)
	Pulse800ms                            // 8x (800 ms)
)

// EmphasizedMode define los modos de énfasis de texto
type EmphasizedMode byte

const (
	EmphOff EmphasizedMode = iota // Modo normal (sin énfasis)
	EmphOn                        // Modo enfatizado (negrita)
)

type TabColumnNumber []byte

type TabTotalPosition byte

type UserDefinedChar byte

type PrinterEnabled byte

const (
	EnaOff PrinterEnabled = iota
	EnaOn
)

type LineSpace byte

type BitImageMode byte

const (
	Mode8DotSingleDen  BitImageMode = iota // Modo de 8 puntos, densidad simple
	Mode8DotDoubleDen                      // Modo de 8 puntos, doble densidad
	Mode24DotSingleDen                     // Modo de 24 puntos, densidad simple
	Mode24DotDoubleDen                     // Modo de 24 puntos, doble densidad
)

type PrinterInitiated bool

const (
	OnInit  PrinterInitiated = true  // Inicializar la impresora
	OffInit PrinterInitiated = false // No inicializar la impresora
)
