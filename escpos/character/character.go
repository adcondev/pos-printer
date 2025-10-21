// Package character implements ESC/POS commands for character formatting and appearance.
//
// ESC/POS is the command system used by thermal receipt printers to control
// character fonts, sizes, styles, effects, international character sets,
// code pages, and user-defined characters.
package character

import (
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for character formatting and appearance.
// ESC/POS is the command system used by thermal receipt printers to control
// character fonts, sizes, styles, effects, international character sets,
// code pages, and user-defined characters.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Spacing represents character spacing in dots
type Spacing byte

// PrintMode represents the print mode bits for character formatting
type PrintMode byte

const (
	// Print Mode bits
	FontAPm           PrintMode = 0x00
	FontBPm           PrintMode = 0x01
	EmphasizedOffPm   PrintMode = 0x00
	EmphasizedOnPm    PrintMode = 0x08
	DoubleHeightOffPm PrintMode = 0x00
	DoubleHeightOnPm  PrintMode = 0x10
	DoubleWidthOffPm  PrintMode = 0x00
	DoubleWidthOnPm   PrintMode = 0x20
	UnderlineOffPm    PrintMode = 0x00
	UnderlineOnPm     PrintMode = 0x80
)

type UnderlineMode byte

const (
	// Underline modes
	NoDot       UnderlineMode = 0x00
	OneDot      UnderlineMode = 0x01
	TwoDot      UnderlineMode = 0x02
	NoDotASCII  UnderlineMode = '0'
	OneDotASCII UnderlineMode = '1'
	TwoDotASCII UnderlineMode = '2'
)

type EmphasizedMode byte

const (
	// Emphasized modes
	OffEm EmphasizedMode = 0x00 // LSB = 0 -> emphasized OFF
	OnEm  EmphasizedMode = 0x01 // LSB = 1 -> emphasized ON
)

type DoubleStrikeMode byte

const (
	// Double-strike modes
	OffDsm DoubleStrikeMode = 0x00 // LSB = 0 -> double-strike OFF
	OnDsm  DoubleStrikeMode = 0x01 // LSB = 1 -> double-strike ON
)

type FontType byte

const (
	// Font types
	FontA        FontType = 0x00
	FontB        FontType = 0x01
	FontC        FontType = 0x02
	FontD        FontType = 0x03
	FontE        FontType = 0x04
	FontAAscii   FontType = '0'
	FontBAscii   FontType = '1'
	FontCAscii   FontType = '2'
	FontDAscii   FontType = '3'
	FontEAscii   FontType = '4'
	SpecialFontA FontType = 97
	SpecialFontB FontType = 98
)

type InternationalSet byte

const (
	// International character sets
	USA          InternationalSet = 0
	France       InternationalSet = 1
	Germany      InternationalSet = 2
	UK           InternationalSet = 3
	DenmarkI     InternationalSet = 4
	Sweden       InternationalSet = 5
	Italy        InternationalSet = 6
	SpainI       InternationalSet = 7
	Japan        InternationalSet = 8
	Norway       InternationalSet = 9
	DenmarkII    InternationalSet = 10
	SpainII      InternationalSet = 11
	LatinAmerica InternationalSet = 12
	Korea        InternationalSet = 13
	SloveniaCro  InternationalSet = 14
	China        InternationalSet = 15
	Vietnam      InternationalSet = 16
	Arabia       InternationalSet = 17

	// Extended India character sets (model-dependent)
	IndiaDevanagari InternationalSet = 66
	IndiaBengali    InternationalSet = 67
	IndiaTamil      InternationalSet = 68
	IndiaTelugu     InternationalSet = 69
	IndiaAssamese   InternationalSet = 70
	IndiaOriya      InternationalSet = 71
	IndiaKannada    InternationalSet = 72
	IndiaMalayalam  InternationalSet = 73
	IndiaGujarati   InternationalSet = 74
	IndiaPunjabi    InternationalSet = 75
	IndiaMarathi    InternationalSet = 82
)

type RotationMode byte

const (
	// NoRotation represents no rotation mode
	NoRotation RotationMode = 0x00
	// On90Dot1 represents 90-degree rotation with 1-dot spacing
	On90Dot1 RotationMode = 0x01
	// On90Dot15 represents 90-degree rotation with 1.5-dot spacing
	On90Dot15 RotationMode = 0x02
	// NoRotationASCII represents no rotation mode (ASCII mode)
	NoRotationASCII RotationMode = '0'
	// On90Dot1Ascii represents 90-degree rotation with 1-dot spacing (ASCII mode)
	On90Dot1Ascii RotationMode = '1'
	// On90Dot15Ascii represents 90-degree rotation with 1.5-dot spacing (ASCII mode)
	On90Dot15Ascii RotationMode = '2'
)

type PrintColor byte

const (
	// Black represents black print color
	Black PrintColor = 0x00
	// Red represents red print color
	Red PrintColor = 0x01
	// BlackASCII represents black print color (ASCII mode)
	BlackASCII PrintColor = '0'
	// RedASCII represents red print color (ASCII mode)
	RedASCII PrintColor = '1'
)

type CodeTable byte

const (
	// Character Code Table Pages (sharedcommands values)
	PC437         CodeTable = 0  // PC437: USA, Standard Europe
	Katakana      CodeTable = 1  // Katakana
	PC850         CodeTable = 2  // PC850: Multilingual
	PC860         CodeTable = 3  // PC860: Portuguese
	PC863         CodeTable = 4  // PC863: Canadian-French
	PC865         CodeTable = 5  // PC865: Nordic
	Hiragana      CodeTable = 6  // Hiragana
	OnePassKanji1 CodeTable = 7  // OnePassKanji1 Type 1
	OnePassKanji2 CodeTable = 8  // OnePassKanji2 Type 2
	PC851         CodeTable = 11 // PC851: Greek
	PC853         CodeTable = 12 // PC853: Turkish
	PC857         CodeTable = 13 // PC857: Turkish
	PC737         CodeTable = 14 // PC737: Greek
	ISO88597      CodeTable = 15 // ISO88597: Greek
	WPC1252       CodeTable = 16 // WPC1252
	PC866         CodeTable = 17 // PC866: Cyrillic #2
	PC852         CodeTable = 18 // PC852: Latin 2
	PC858         CodeTable = 19 // PC858: Euro
	ThaiCode42    CodeTable = 20 // ThaiCode42
	ThaiCode11    CodeTable = 21 // ThaiCode11
	ThaiCode13    CodeTable = 22 // ThaiCode13
	ThaiCode14    CodeTable = 23 // ThaiCode14
	ThaiCode16    CodeTable = 24 // ThaiCode16
	ThaiCode17    CodeTable = 25 // ThaiCode17
	ThaiCode18    CodeTable = 26 // ThaiCode18
	TCVN31        CodeTable = 30 // TCVN31 Type 1: Vietnamese
	TCVN32        CodeTable = 31 // TCVN32 Type 2: Vietnamese
	PC720         CodeTable = 32 // PC720: Arabic
	WPC775        CodeTable = 33 // WPC775: Baltic Rim
	PC855         CodeTable = 34 // PC855: Cyrillic
	PC861         CodeTable = 35 // PC861: Icelandic
	PC862         CodeTable = 36 // PC862: Hebrew
	PC864         CodeTable = 37 // PC864: Arabic
	PC869         CodeTable = 38 // PC869: Greek
	ISO88592      CodeTable = 39 // ISO88592: Latin 2
	ISO885915     CodeTable = 40 // ISO885915: Latin 9
	PC1098        CodeTable = 41 // PC1098: Farsi
	PC1118        CodeTable = 42 // PC1118: Lithuanian
	PC1119        CodeTable = 43 // PC1119: Lithuanian
	PC1125        CodeTable = 44 // PC1125: Ukrainian
	WPC1250       CodeTable = 45 // WPC1250: Latin 2
	WPC1251       CodeTable = 46 // WPC1251: Cyrillic
	WPC1253       CodeTable = 47 // WPC1253: Greek
	WPC1254       CodeTable = 48 // WPC1254: Turkish
	WPC1255       CodeTable = 49 // WPC1255: Hebrew
	WPC1256       CodeTable = 50 // WPC1256: Arabic
	WPC1257       CodeTable = 51 // WPC1257: Baltic Rim
	WPC1258       CodeTable = 52 // WPC1258: Vietnamese
	KZ1048        CodeTable = 53 // KZ1048: Kazakhstan

	// India-related pages (model-dependent)
	Devanagari CodeTable = 66 // Devanagari
	Bengali    CodeTable = 67 // Bengali
	Tamil      CodeTable = 68 // Tamil
	Telugu     CodeTable = 69 // Telugu
	Assamese   CodeTable = 70 // Assamese
	Oriya      CodeTable = 71 // Oriya
	Kannada    CodeTable = 72 // Kannada
	Malayalam  CodeTable = 73 // Malayalam
	Gujarati   CodeTable = 74 // Gujarati
	Punjabi    CodeTable = 75 // Punjabi
	Marathi    CodeTable = 82 // Marathi

	// Reserved / special
	Special254 CodeTable = 254 // Special254 (model-dependent)
	Special255 CodeTable = 255 // Special255 (model-dependent)
)

type UpsideDownMode byte

const (
	// Upside-down modes
	OffUdm UpsideDownMode = 0x00 // LSB = 0 -> upside-down OFF
	OnUdm  UpsideDownMode = 0x01 // LSB = 1 -> upside-down ON
)

type Size byte

const (
	// Character size configurations
	Size1x1 Size = 0x00 // width=1 height=1 (normal)
	Size2x1 Size = 0x10 // width=2 height=1
	Size3x1 Size = 0x20 // width=3 height=1
	Size4x1 Size = 0x30 // width=4 height=1
	Size1x2 Size = 0x01 // width=1 height=2
	Size2x2 Size = 0x11 // width=2 height=2 (double width & height)

	// Masks and shifts

	// HeightMask is used to extract height bits from n
	HeightMask Size = 0x07 // bits 0-2
	// WidthMask is used to extract width bits from n
	WidthMask Size = 0x70 // bits 4-6
	// WidthShift is used to shift width bits to LSB position
	WidthShift Size = 4
)

type ReverseMode byte

const (
	// White/black reverse modes
	OffRm      ReverseMode = 0x00 // LSB = 0 -> reverse OFF
	OnRm       ReverseMode = 0x01 // LSB = 1 -> reverse ON
	OffRmASCII ReverseMode = '0'
	OnRmASCII  ReverseMode = '1'
)

type SmoothingMode byte

const (
	// Smoothing modes
	OffSm      SmoothingMode = 0x00 // LSB = 0 -> smoothing OFF
	OnSm       SmoothingMode = 0x01 // LSB = 1 -> smoothing ON
	OffSmASCII SmoothingMode = '0'
	OnSmASCII  SmoothingMode = '1'
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrUnderlineMode   = fmt.Errorf("invalid underline mode(try 0-2 or '0'..'2')")
	ErrRotationMode    = fmt.Errorf("invalid rotation mode(try 0-2 or '0'..'2')")
	ErrPrintColor      = fmt.Errorf("invalid print color(try 0-1 or '0'..'1')")
	ErrCharacterFont   = fmt.Errorf("invalid character font(try 0-4 or '0'..'4')")
	ErrCodeTablePage   = fmt.Errorf("invalid code table page(try check model support!)")
	ErrCharacterSet    = fmt.Errorf("invalid international character set(try check model support!)")
	ErrCharacterWidth  = fmt.Errorf("invalid character width(try 1-8)")
	ErrCharacterHeight = fmt.Errorf("invalid character height(try 1-8)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Compile-time check that Commands implements Capability
var _ Capability = (*Commands)(nil)

// Capability defines the main interface for character-related commands
type Capability interface {
	SetRightSideCharacterSpacing(spacing Spacing) []byte
	SelectPrintModes(modeBits PrintMode) []byte
	SetUnderlineMode(thickness UnderlineMode) ([]byte, error)
	SetEmphasizedMode(mode EmphasizedMode) []byte
	SetDoubleStrikeMode(mode DoubleStrikeMode) []byte
	SelectCharacterFont(fontType FontType) ([]byte, error)
	SelectInternationalCharacterSet(charset InternationalSet) ([]byte, error)
	Set90DegreeClockwiseRotationMode(rotationMode RotationMode) ([]byte, error)
	SelectPrintColor(color PrintColor) ([]byte, error)
	SelectCharacterCodeTable(page CodeTable) ([]byte, error)
	SetUpsideDownMode(mode UpsideDownMode) []byte
	SelectCharacterSize(sizeConfig Size) []byte
	SetWhiteBlackReverseMode(mode ReverseMode) []byte
	SetSmoothingMode(mode SmoothingMode) []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands groups all character-related capabilities
type Commands struct {
	Effects        EffectsCapability
	CodeConversion CodeConversionCapability
	UserDefined    UserDefinedCapability
}

func NewCommands() *Commands {
	// Constructor that initializes sub-commands; returns pointer to avoid
	// copies when passing the structure.
	return &Commands{
		Effects:        &EffectsCommands{},
		CodeConversion: &CodeConversionCommands{},
		UserDefined:    &UserDefinedCommands{},
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// NewSize creates a Size value for character width and height magnification.
//
// Format:
//
//	Not applicable (helper function)
//
// Range:
//
//	width = 1–8
//	height = 1–8
//
// Default:
//
//	Not applicable
//
// Parameters:
//
//	width: Character width magnification (1-8)
//	height: Character height magnification (1-8)
//
// Notes:
//   - This is a helper function to build Size values for SelectCharacterSize
//   - The returned Size encodes both width and height in a single byte
//
// Errors:
//
//	Returns ErrCharacterWidth if width is outside the range 1-8.
//	Returns ErrCharacterHeight if height is outside the range 1-8.
func NewSize(width, height byte) (Size, error) {
	if width < 1 || width > 8 {
		return 0, ErrCharacterWidth
	}
	if height < 1 || height > 8 {
		return 0, ErrCharacterHeight
	}
	w := (width - 1) << 4
	h := height - 1
	return Size(w | h), nil
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateUnderlineMode validates if underline mode is valid
func ValidateUnderlineMode(mode UnderlineMode) error {
	switch mode {
	case NoDot, OneDot, TwoDot, NoDotASCII, OneDotASCII, TwoDotASCII:
		return nil
	default:
		return ErrUnderlineMode
	}
}

// ValidateFontType validates if font type is valid
func ValidateFontType(font FontType) error {
	switch font {
	case FontA, FontB, FontC, FontD, FontE,
		FontAAscii, FontBAscii, FontCAscii, FontDAscii, FontEAscii,
		SpecialFontA, SpecialFontB:
		return nil
	default:
		return ErrCharacterFont
	}
}

// ValidateInternationalSet validates if international character set is valid
func ValidateInternationalSet(charset InternationalSet) error {
	// Standard range
	if charset <= Arabia {
		return nil
	}
	// India-specific range
	if charset >= IndiaDevanagari && charset <= IndiaPunjabi {
		return nil
	}
	// Marathi
	if charset == IndiaMarathi {
		return nil
	}
	return ErrCharacterSet
}

// ValidateRotationMode validates if rotation mode is valid
func ValidateRotationMode(mode RotationMode) error {
	switch mode {
	case NoRotation, On90Dot1, On90Dot15,
		NoRotationASCII, On90Dot1Ascii, On90Dot15Ascii:
		return nil
	default:
		return ErrRotationMode
	}
}

// ValidatePrintColor validates if print color is valid
func ValidatePrintColor(color PrintColor) error {
	switch color {
	case Black, Red, BlackASCII, RedASCII:
		return nil
	default:
		return ErrPrintColor
	}
}

// ValidateCodeTable validates if code table page is valid
func ValidateCodeTable(page CodeTable) error {
	// Common pages 0-8
	if page <= OnePassKanji2 {
		return nil
	}
	// Pages 11-19
	if page >= PC851 && page <= PC858 {
		return nil
	}
	// Thai codes 20-26
	if page >= ThaiCode42 && page <= ThaiCode18 {
		return nil
	}
	// Vietnamese 30-31
	if page >= TCVN31 && page <= TCVN32 {
		return nil
	}
	// Pages 32-53
	if page >= PC720 && page <= KZ1048 {
		return nil
	}
	// India-specific pages 66-75
	if page >= Devanagari && page <= Punjabi {
		return nil
	}
	// Marathi
	if page == Marathi {
		return nil
	}
	// Special pages
	if page == Special254 || page == Special255 {
		return nil
	}
	return ErrCodeTablePage
}
