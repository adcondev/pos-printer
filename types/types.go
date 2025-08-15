package types

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
	BarcodeUPCA BarcodeType = iota
	BarcodeUPCE
	BarcodeEAN13
	BarcodeEAN8
	BarcodeCode39
	BarcodeITF
	BarcodeCodebar
	BarcodeCode93
	BarcodeCode128
)

// BarcodeTextPosition define posiciones estándar para texto en códigos de barras
type BarcodeTextPosition int

const (
	BarcodeTextNone BarcodeTextPosition = iota
	BarcodeTextAbove
	BarcodeTextBelow
	BarcodeTextBoth
)

// CutMode define modos de corte estándar
type CutMode int

const (
	CutFeed CutMode = iota
	Cut
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

// CharacterSet define los conjuntos de caracteres estándar
type CharacterSet int

const (
	CP437      CharacterSet = iota // CP437 U.S.A. / Standard Europe
	Katakana                       // Katakana (JIS X 0201)
	CP850                          // CP850 Multilingual
	CP860                          // CP860 Portuguese
	CP863                          // CP863 Canadian French
	CP865                          // CP865 Nordic
	WestEurope                     // WestEurope (ISO-8859-1)
	Greek                          // Greek (ISO-8859-7)
	Hebrew                         // Hebrew (ISO-8859-8)
	CP755                          // CP755 East Europe (not directly supported)
	Iran                           // Iran (CP720 Arabic)
	WCP1252                        // WCP1252 Windows-1252
	CP866                          // CP866 Cyrillic #2
	CP852                          // CP852 Latin2
	CP858                          // CP858 Multilingual + Euro
	IranII                         // IranII (CP864)
	Latvian                        // Latvian (Windows-1257)
)

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

// CashDrawerTime define los tiempos de pulso del cajón de dinero
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

// TODO: Tener en cuenta que es un modo y que son cantidades, cantidades se validad en ESCPOS en el rango definido.
// Los modos tendrían su respectivo map en ESCPOS. Definir Min y Max vuelve rígido el uso de los tipos. Solo crear Types.

// TODO: Agregar más tipos genéricos según necesites
// Por ejemplo:
// - QRCodeSize
// - PrintSpeed
// - CharacterSet
// etc.
