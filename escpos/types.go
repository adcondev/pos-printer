package escpos

// Alignment define las alineaciones de texto estándar
type Alignment int

const (
	// AlignLeft Left alignment
	AlignLeft Alignment = iota
	// AlignCenter Center alignment
	AlignCenter
	// AlignRight Right alignment
	AlignRight
	// AlignJustified Justified alignment
	AlignJustified // Algunos protocolos podrían soportar esto
)

// Font define los tipos de fuente estándar
type Font byte

const (
	// FontA A (normal)
	FontA Font = iota
	// FontB B (smaller)
	FontB
	// FontC C
	FontC
	// FontD D
	FontD
	// FontE E
	FontE
	// SpecialA Special font A
	SpecialA
	// SpecialB Special font B
	SpecialB
)

// UnderlineMode define los modos de subrayado estándar
type UnderlineMode byte

const (
	// UnderNone No underline
	UnderNone UnderlineMode = iota
	// UnderSingle Single underline
	UnderSingle
	// UnderDouble Double underline
	UnderDouble
)

// BarcodeType define los tipos de código de barras estándar
type BarcodeType int

const (
	// UPCA Universal Product Code - A
	UPCA BarcodeType = iota
	// UPCE Universal Product Code - Compressed
	UPCE
	// EAN13 European Article Number - 13
	EAN13
	// EAN8 European Article Number - 8
	EAN8
	// Code39 CODE 39
	Code39
	// ITF Interleaved 2 of 5
	ITF
	// Codabar Lineal Barcode
	Codabar
)

// TextPositionBarcode define posiciones estándar para texto en códigos de barras
type TextPositionBarcode byte

const (
	// NonePosBarcode No text
	NonePosBarcode TextPositionBarcode = iota
	// AbovePosBarcode Text above the barcode
	AbovePosBarcode
	// BelowPosBarcode Text below the barcode
	BelowPosBarcode
	// BothPosBarcode Text both above and below the barcode
	BothPosBarcode
)

// BarcodeHeight define heights for barcodes
type BarcodeHeight byte

// BarcodeWidth define widths for barcodes
type BarcodeWidth byte

const (
	// ExtraSmallWidth Tiny width
	ExtraSmallWidth BarcodeWidth = iota
	// SmallWidth Small width
	SmallWidth
	// MediumWidth Medium width
	MediumWidth
	// LargeWidth Large width
	LargeWidth
	// ExtraLargeWidth Extra large width
	ExtraLargeWidth
)

// CutPaper define modos de corte estándar
type CutPaper byte

const (
	// FullCut Full paper cut
	FullCut CutPaper = iota
	// PartialCut  Partial paper cut
	PartialCut
)

// Density define densidades de impresión estándar para imágenes
type Density int

const (
	// DensitySingle No density (default)
	DensitySingle Density = iota
	// DensityDouble Double density
	DensityDouble
	// DensityTriple Triple density
	DensityTriple
	// DensityQuadruple Quadruple density
	DensityQuadruple
)

// QRModel define los modelos de código QR estándar
type QRModel byte

const (
	// Model1 QR Type 1 (standard)
	Model1 QRModel = iota
	// Model2 QR Type 2 (recommended)
	Model2
)

// QRErrorCorrection define los niveles de corrección de errores estándar para códigos QR
type QRErrorCorrection byte

const (
	// ECLow 7% error correction
	ECLow QRErrorCorrection = iota
	// ECMedium 15% error correction
	ECMedium
	// ECHigh 25% error correction
	ECHigh
	// ECHighest 30% error correction
	ECHighest
)

// QRModuleSize defines the size of QR code modules
type QRModuleSize byte

const (
	// MinType Minimum module size
	MinType QRModuleSize = 1
	// MaxType Maximum module size
	MaxType QRModuleSize = 16
)

// Lines defines the number of lines to feed
type Lines byte

// RealTimeStatus defines the real-time status types
type RealTimeStatus byte

const (
	// PrinterStatus defines the state of the printer
	PrinterStatus RealTimeStatus = iota
	// OfflineStatus defines the state of the printer's offline status
	OfflineStatus
	// ErrorStatus defines the state of the printer's error status
	ErrorStatus
	// PaperSensorStatus defines the state of the paper sensor
	PaperSensorStatus
)

// CashDrawerPin defines the pins for the cash drawer
type CashDrawerPin byte

const (
	// Pin2 Cash Drawer Pin 2
	Pin2 CashDrawerPin = iota
	// Pin5 Cash Drawer Pin 5
	Pin5
)

// CashDrawerTimePulse defines the pulse duration for the cash drawer
type CashDrawerTimePulse byte

const (
	// Pulse100ms 1x (100 ms)
	Pulse100ms CashDrawerTimePulse = iota
	// Pulse200ms 2x (200 ms)
	Pulse200ms
	// Pulse300ms 3x (300 ms)
	Pulse300ms
	// Pulse400ms 4x (400 ms)
	Pulse400ms
	// Pulse500ms 5x (500 ms)
	Pulse500ms
	// Pulse600ms 6x (600 ms)
	Pulse600ms
	// Pulse700ms 7x (700 ms)
	Pulse700ms
	// Pulse800ms 8x (800 ms)
	Pulse800ms
)

// EmphasizedMode defines the emphasized text modes
type EmphasizedMode byte

const (
	// EmphasizedOff Normal mode (no emphasis)
	EmphasizedOff EmphasizedMode = iota
	// EmphasizedOn Emphasized mode (bold)
	EmphasizedOn
)

// TabColumnNumber defines the tab column numbers
type TabColumnNumber []byte

// TabTotalPosition defines the total tab positions
type TabTotalPosition byte

// UserDefinedChar defines user-defined character slots
type UserDefinedChar byte

// PrinterEnabled defines whether the printer is enabled or disabled
type PrinterEnabled byte

const (
	// EnaOff Printer disabled
	EnaOff PrinterEnabled = iota
	// EnaOn Printer enabled
	EnaOn
)

// LineSpace defines the line spacing in dots
type LineSpace byte

// BitImageMode defines the bit image modes
type BitImageMode byte

const (
	// Mode8DotSingleDen 8-dot mode, single density
	Mode8DotSingleDen BitImageMode = iota
	// Mode8DotDoubleDen  8-dot mode, double density
	Mode8DotDoubleDen
	// Mode24DotSingleDen  24-dot mode, single density
	Mode24DotSingleDen
	// Mode24DotDoubleDen  24-dot mode, double density
	Mode24DotDoubleDen
)

// PrinterInitiated indicates whether to initialize the printer
type PrinterInitiated bool

const (
	// OnInit Printer initialized
	OnInit PrinterInitiated = true // Inicializar la impresora
	// OffInit Printer not initialized
	OffInit PrinterInitiated = false // No inicializar la impresora
)

type PrintMode byte

const (
	// TODO: placehold this constants for future validation maps
	PMFontA PrintMode = iota
	PMFontB
	PMEmphasizedOff
	PMEmphasizedOn
	PMDoubleHeightOff
	PMDoubleHeightOn
	PMDoubleWidthOff
	PMDoubleWidthOn
	PMUnderlineOff
	PMUnderlineOn
)
