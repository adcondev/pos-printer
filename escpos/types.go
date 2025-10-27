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

// EmphasizedMode defines the emphasized text modes
type EmphasizedMode byte

const (
	// EmphasizedOff Normal mode (no emphasis)
	EmphasizedOff EmphasizedMode = iota
	// EmphasizedOn Emphasized mode (bold)
	EmphasizedOn
)

// PrinterInitiated indicates whether to initialize the printer
type PrinterInitiated bool
