package character

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// TODO: Verify if all types have corresponding brief tests

// ============================================================================
// Constant and Var Definitions
// ============================================================================

type Spacing byte

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
	NoDotAscii  UnderlineMode = '0'
	OneDotAscii UnderlineMode = '1'
	TwoDotAscii UnderlineMode = '2'
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
	// Rotation modes
	NoRotation      RotationMode = 0x00 // n = 0
	On90Dot1        RotationMode = 0x01 // n = 1 (1-dot spacing)
	On90Dot15       RotationMode = 0x02 // n = 2 (1.5-dot spacing)
	NoRotationAscii RotationMode = '0'
	On90Dot1Ascii   RotationMode = '1'
	On90Dot15Ascii  RotationMode = '2'
)

type PrintColor byte

const (
	// Print Color Modes
	Black      PrintColor = 0x00 // n = 0
	Red        PrintColor = 0x01 // n = 1
	BlackASCII PrintColor = '0'
	RedASCII   PrintColor = '1'
)

type CodeTable byte

const (
	// Character Code Table Pages (common values)
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
)

type ReverseMode byte

const (
	// White/black reverse modes
	OffRm      ReverseMode = 0x00 // LSB = 0 -> reverse OFF
	OnRm       ReverseMode = 0x01 // LSB = 1 -> reverse ON
	OffRmAscii ReverseMode = '0'
	OnRmAscii  ReverseMode = '1'
)

type SmoothingMode byte

const (
	// Smoothing modes
	OffSm      SmoothingMode = 0x00 // LSB = 0 -> smoothing OFF
	OnSm       SmoothingMode = 0x01 // LSB = 1 -> smoothing ON
	OffSmAscii SmoothingMode = '0'
	OnSmAscii  SmoothingMode = '1'
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

// Interface compliance checks
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
	return &Commands{
		Effects:        &EffectsCommands{},
		CodeConversion: &CodeConversionCommands{},
		UserDefined:    &UserDefinedCommands{},
	}
}

// SetRightSideCharacterSpacing
//
// Format:
//
//	ASCII: ESC SP n
//	Hex:   0x1B 0x20 n
//	Decimal: 27 32 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Sets the right-side character spacing to n × (horizontal or vertical
//	motion unit).
//
// Notes:
//   - The character spacing set by this command is effective for alphanumeric,
//     Kana, and user-defined characters.
//   - When characters are enlarged, the character spacing becomes n times the
//     normal value.
//   - In Standard mode the horizontal motion unit is used.
//   - In Page mode the vertical or horizontal motion unit is used according
//     to the print direction set by ESC T.
//   - When the starting position is set to the upper-left or lower-right of
//     the print area using ESC T, the horizontal motion unit is used.
//   - When the starting position is set to the upper-right or lower-left of
//     the print area using ESC T, the vertical motion unit is used.
//   - Character spacing can be set independently in Standard mode and in
//     Page mode; this command affects the spacing for the currently selected
//     mode.
//   - If the horizontal or vertical motion unit is changed after this
//     command is executed, the numeric character spacing value does not
//     change.
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off.
//   - This command is used to change the spacing between characters.
//
// Byte sequence:
//
//	ESC SP n -> 0x1B, 0x20, n
func (c *Commands) SetRightSideCharacterSpacing(n Spacing) []byte {
	return []byte{common.ESC, common.SP, byte(n)}
}

// SelectPrintModes selects character font and style bits (emphasized,
// double-height, double-width, underline) together.
//
// Format:
//
//	ASCII: ESC ! n
//	Hex:   0x1B 0x21 n
//	Decimal: 27 33 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Selects the character font and styles (emphasized, double-height,
//	double-width, and underline) together by setting bits in the parameter
//	byte n. The bits have the following meanings:
//
//	Bit 0 (0x01) - Font selection
//	  0: Selects Font 1
//	  1: Selects Font 2
//
//	Bit 3 (0x08) - Emphasized mode
//	  0: Emphasized OFF
//	  1: Emphasized ON
//
//	Bit 4 (0x10) - Double-height mode
//	  0: Double-height OFF
//	  1: Double-height ON
//
//	Bit 5 (0x20) - Double-width mode
//	  0: Double-width OFF
//	  1: Double-width ON
//
//	Bit 7 (0x80) - Underline mode
//	  0: Underline OFF
//	  1: Underline ON
//
// Notes:
//   - Configurations for Font 1 and Font 2 differ by model. If the desired
//     font type cannot be selected with this command, use ESC M.
//   - Bits 0, 4, 5 and 7 affect 1-byte code characters. On some models,
//     bits 4, 5 and 7 also affect Korean characters.
//   - Emphasized mode (bit 3) is effective for both 1-byte and multi-byte
//     characters.
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, power is turned off, or one of these commands is executed:
//   - Bit 0 (font): ESC M
//   - Bit 3 (emphasized): ESC E
//   - Bit 4,5 (size): GS !
//   - Bit 7 (underline): ESC -
//   - When some characters in a line are double-height, all characters on
//     the line are aligned at the baseline.
//   - Double-width enlarges characters to the right from the left side of
//     the character. When both double-height and double-width are on,
//     characters become quadruple size.
//   - In Standard mode double-height enlarges in the paper-feed direction
//     and double-width enlarges perpendicular to paper feed. Rotating
//     characters 90° clockwise swaps the relationship.
//   - In Page mode double-height and double-width follow the character
//     orientation.
//   - Underline thickness is determined by ESC -, regardless of character
//     size. Underline color matches the printed character color (GS ( N
//     <Function 48>).
//   - The following are not underlined:
//   - 90° clockwise-rotated characters
//   - white/black reverse characters
//   - spaces set by HT, ESC $, and ESC \
//   - On printers with Automatic font replacement (GS ( E <Function 5> with
//     a = 111,112,113), the replacement font is selected by this command.
//
// Byte sequence:
//
//	ESC ! n -> 0x1B, 0x21, n
func (c *Commands) SelectPrintModes(n PrintMode) []byte {
	return []byte{common.ESC, '!', byte(n)}
}

// SetUnderlineMode sets underline mode on or off and selects underline thickness.
//
// Format:
//
//	ASCII: ESC - n
//	Hex:   0x1B 0x2D n
//	Decimal: 27 45 n
//
// Range:
//
//	n = 0, 1, 2, 48, 49, 50
//
// Default:
//
//	n = 0
//
// Description:
//
//	Turns underline mode on or off using n as follows:
//
//	  n = 0 or 48 -> Turns off underline mode
//	  n = 1 or 49 -> Turns on underline mode (1-dot thick)
//	  n = 2 or 50 -> Turns on underline mode (2-dots thick)
//
// Notes:
//   - The underline mode is effective for alphanumeric, Kana, and user-
//     defined characters. On some models it is also effective for Korean
//     characters.
//   - The underline color matches the printed character color (see
//     GS ( N <Function 48>).
//   - Changing the character size does not affect the current underline
//     thickness.
//   - When underline mode is turned off the underline thickness value is
//     retained but no underline is produced.
//   - The printer does not underline 90° clockwise-rotated characters,
//     white/black reverse characters, or spaces produced by HT, ESC $, and
//     ESC \.
//   - The setting remains in effect until ESC ! is executed, ESC @ is
//     executed, the printer is reset, or power is turned off.
//   - Some printer models support the 2-dot thick underline (n = 2 or 50).
//
// Byte sequence:
//
//	ESC - n -> 0x1B, 0x2D, n
func (c *Commands) SetUnderlineMode(n UnderlineMode) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid values
	default:
		return nil, ErrUnderlineMode
	}
	return []byte{common.ESC, '-', byte(n)}, nil
}

// SetEmphasizedMode turns emphasized (bold) mode on or off.
//
// Format:
//
//	ASCII: ESC E n
//	Hex:   0x1B 0x45 n
//	Decimal: 27 69 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Turns emphasized mode on or off. When the least-significant bit (LSB)
//	of n is 0, emphasized mode is turned off. When the LSB of n is 1,
//	emphasized mode is turned on.
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and
//     user-defined characters.
//   - Settings of this command remain in effect until ESC ! is executed,
//     ESC @ is executed, the printer is reset, or power is turned off.
//
// Byte sequence:
//
//	ESC E n -> 0x1B, 0x45, n
func (c *Commands) SetEmphasizedMode(n EmphasizedMode) []byte {
	return []byte{common.ESC, 'E', byte(n)}
}

// SetDoubleStrikeMode turns double-strike mode on or off.
//
// Format:
//
//	ASCII: ESC G n
//	Hex:   0x1B 0x47 n
//	Decimal: 27 71 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Turns double-strike mode on or off. When the least-significant bit (LSB)
//	of n is 0, double-strike mode is turned off. When the LSB of n is 1,
//	double-strike mode is turned on.
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and
//     user-defined characters.
//   - Settings of this command remain in effect until ESC ! is executed,
//     ESC @ is executed, the printer is reset, or power is turned off.
//
// Byte sequence:
//
//	ESC G n -> 0x1B, 0x47, n
func (c *Commands) SetDoubleStrikeMode(n DoubleStrikeMode) []byte {
	return []byte{common.ESC, 'G', byte(n)}
}

// SelectCharacterFont selects a character font.
//
// Format:
//
//	ASCII: ESC M n
//	Hex:   0x1B 0x4D n
//	Decimal: 27 77 n
//
// Range:
//
//	Depending on the model: 0–4, 48–52, 97, 98
//
// Default:
//
//	Depending on the model: n = 0 or n = 1
//
// Description:
//
//	Selects a character font using n as follows:
//
//	  n = 0 or 48 -> Font A
//	  n = 1 or 49 -> Font B
//	  n = 2 or 50 -> Font C
//	  n = 3 or 51 -> Font D
//	  n = 4 or 52 -> Font E
//	  n = 97      -> Special font A
//	  n = 98      -> Special font B
//
// Notes:
//   - The selected character font is effective for alphanumeric, Kana, and
//     user-defined characters.
//   - Configurations of Font A and Font B depend on the printer model.
//   - Settings remain in effect until ESC ! is executed, ESC @ is executed,
//     the printer is reset, or the power is turned off.
//   - On printers with the Automatic font replacement function, the
//     replaced font selected by GS ( E <Function 5> (a = 111, 112, 113)
//     is selected by this command.
//
// Byte sequence:
//
//	ESC M n -> 0x1B, 0x4D, n
func (c *Commands) SelectCharacterFont(n FontType) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, 3, 4:
		// Numeric values
	case '0', '1', '2', '3', '4':
		// ASCII values
	case 97, 98:
		// Special fonts
	default:
		return nil, ErrCharacterFont
	}
	return []byte{common.ESC, 'M', byte(n)}, nil
}

// SelectInternationalCharacterSet selects an international character set.
//
// Format:
//
//	ASCII: ESC R n
//	Hex:   0x1B 0x52 n
//	Decimal: 27 82 n
//
// Range:
//
//	Different depending on the printer model (common: 0–17; some models
//	support extended India codes such as 66–75, 82).
//
// Default:
//
//	Depends on the printer model.
//	  - Other models: n = 0
//	  - Japanese models: n = 8
//	  - Korean models: n = 13
//	  - Simplified Chinese models: n = 15
//
// Description:
//
//	Selects an international character set using the parameter n. Typical
//	values map to countries/regions as follows:
//
//	  0   U.S.A.
//	  1   France
//	  2   Germany
//	  3   U.K.
//	  4   Denmark I
//	  5   Sweden
//	  6   Italy
//	  7   Spain I
//	  8   Japan
//	  9   Norway
//	  10  Denmark II
//	  11  Spain II
//	  12  Latin America
//	  13  Korea
//	  14  Slovenia / Croatia
//	  15  China
//	  16  Vietnam
//	  17  Arabia
//
//	Some models support additional India-specific character sets:
//
//	  66  India (Devanagari)
//	  67  India (Bengali)
//	  68  India (Tamil)
//	  69  India (Telugu)
//	  70  India (Assamese)
//	  71  India (Oriya)
//	  72  India (Kannada)
//	  73  India (Malayalam)
//	  74  India (Gujarati)
//	  75  India (Punjabi)
//	  82  India (Marathi)
//
// Notes:
//   - The selected international character set remains in effect until ESC @
//     is executed, the printer is reset, or power is turned off.
//   - Refer to the printer's Character Code Tables for model-specific
//     mappings and supported characters.
//
// Byte sequence:
//
//	ESC R n -> 0x1B, 0x52, n
func (c *Commands) SelectInternationalCharacterSet(n InternationalSet) ([]byte, error) {
	// Standard range
	if n <= 17 {
		return []byte{common.ESC, 'R', byte(n)}, nil
	}
	// India-specific range
	if n >= 66 && n <= 75 {
		return []byte{common.ESC, 'R', byte(n)}, nil
	}
	// Marathi
	if n == 82 {
		return []byte{common.ESC, 'R', byte(n)}, nil
	}
	return nil, ErrCharacterSet
}

// Set90DegreeClockwiseRotationMode turns 90° clockwise rotation mode on or off.
//
// Format:
//
//	ASCII: ESC V n
//	Hex:   0x1B 0x56 n
//	Decimal: 27 86 n
//
// Range:
//
//	n = 0–2, 48–50
//
// Default:
//
//	n = 0
//
// Description:
//
//	In Standard mode, turns 90° clockwise rotation mode on or off for
//	characters according to n:
//
//	  n = 0 or 48 -> Turns off 90° clockwise rotation mode
//	  n = 1 or 49 -> Turns on 90° clockwise rotation mode (1-dot character spacing)
//	  n = 2 or 50 -> Turns on 90° clockwise rotation mode (1.5-dot character spacing)
//
// Notes:
//   - This mode is effective for alphanumeric, Kana, multilingual, and
//     user-defined characters.
//   - When underline mode is turned on, the printer does not underline
//     90° clockwise-rotated characters.
//   - When character orientation changes in 90° clockwise rotation mode,
//     the relationship between vertical and horizontal directions is
//     reversed.
//   - The 90° clockwise rotation mode has no effect in Page mode.
//   - Some printer models support n = 2 (1.5-dot spacing); some models have
//     fonts for which 90° rotation is not effective.
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, or power is turned off.
//
// Byte sequence:
//
//	ESC V n -> 0x1B, 0x56, n
func (c *Commands) Set90DegreeClockwiseRotationMode(n RotationMode) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid values
	default:
		return nil, ErrRotationMode
	}
	return []byte{common.ESC, 'V', byte(n)}, nil
}

// SelectPrintColor selects the print color.
//
// Format:
//
//	ASCII: ESC r n
//	Hex:   0x1B 0x72 n
//	Decimal: 27 114 n
//
// Range:
//
//	n = 0, 1, 48, 49
//
// Default:
//
//	n = 0
//
// Description:
//
//	Selects a print color using n as follows:
//
//	  n = 0 or 48 -> Black
//	  n = 1 or 49 -> Red
//
// Notes:
//   - In Standard mode this command is enabled only when processed at the
//     Beginning of the line.
//   - In Page mode the color setting is applied to all data collectively
//     printed by FF (in Page mode) or ESC FF.
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or power is turned off.
//   - For printers that support two-color printing, GS ( N and GS ( L / GS 8 L
//     are available to define and control character/background/graphics
//     color layers. Use model-specific GS ( N / GS ( L / GS 8 L commands when
//     available for more advanced two-color workflows.
//
// Byte sequence:
//
//	ESC r n -> 0x1B, 0x72, n
func (c *Commands) SelectPrintColor(n PrintColor) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, '0', '1':
		// Valid values
	default:
		return nil, ErrPrintColor
	}
	return []byte{common.ESC, 'r', byte(n)}, nil
}

// SelectCharacterCodeTable selects a character code table page.
//
// Format:
//
//	ASCII: ESC t n
//	Hex:   0x1B 0x74 n
//	Decimal: 27 116 n
//
// Range:
//
//	Different depending on the printer model (see constants below for common pages).
//
// Default:
//
//	n = 0
//
// Description:
//
//	Selects a page n from the character code table. Typical page mappings:
//
//	  0   Page 0  [PC437: USA, Standard Europe]
//	  1   Page 1  [Katakana]
//	  2   Page 2  [PC850: Multilingual]
//	  3   Page 3  [PC860: Portuguese]
//	  4   Page 4  [PC863: Canadian-French]
//	  5   Page 5  [PC865: Nordic]
//	  6   Page 6  [Hiragana]
//	  7   Page 7  [One-pass printing Kanji characters]
//	  8   Page 8  [One-pass printing Kanji characters]
//	  11  Page 11 [PC851: Greek]
//	  12  Page 12 [PC853: Turkish]
//	  13  Page 13 [PC857: Turkish]
//	  14  Page 14 [PC737: Greek]
//	  15  Page 15 [ISO8859-7: Greek]
//	  16  Page 16 [WPC1252]
//	  17  Page 17 [PC866: Cyrillic #2]
//	  18  Page 18 [PC852: Latin 2]
//	  19  Page 19 [PC858: Euro]
//	  20  Page 20 [Thai Character Code 42]
//	  21  Page 21 [Thai Character Code 11]
//	  22  Page 22 [Thai Character Code 13]
//	  23  Page 23 [Thai Character Code 14]
//	  24  Page 24 [Thai Character Code 16]
//	  25  Page 25 [Thai Character Code 17]
//	  26  Page 26 [Thai Character Code 18]
//	  30  Page 30 [TCVN-3: Vietnamese]
//	  31  Page 31 [TCVN-3: Vietnamese]
//	  32  Page 32 [PC720: Arabic]
//	  33  Page 33 [WPC775: Baltic Rim]
//	  34  Page 34 [PC855: Cyrillic]
//	  35  Page 35 [PC861: Icelandic]
//	  36  Page 36 [PC862: Hebrew]
//	  37  Page 37 [PC864: Arabic]
//	  38  Page 38 [PC869: Greek]
//	  39  Page 39 [ISO8859-2: Latin 2]
//	  40  Page 40 [ISO8859-15: Latin 9]
//	  41  Page 41 [PC1098: Farsi]
//	  42  Page 42 [PC1118: Lithuanian]
//	  43  Page 43 [PC1119: Lithuanian]
//	  44  Page 44 [PC1125: Ukrainian]
//	  45  Page 45 [WPC1250: Latin 2]
//	  46  Page 46 [WPC1251: Cyrillic]
//	  47  Page 47 [WPC1253: Greek]
//	  48  Page 48 [WPC1254: Turkish]
//	  49  Page 49 [WPC1255: Hebrew]
//	  50  Page 50 [WPC1256: Arabic]
//	  51  Page 51 [WPC1257: Baltic Rim]
//	  52  Page 52 [WPC1258: Vietnamese]
//	  53  Page 53 [KZ-1048: Kazakhstan]
//	  66  Page 66 [Devanagari]
//	  67  Page 67 [Bengali]
//	  68  Page 68 [Tamil]
//	  69  Page 69 [Telugu]
//	  70  Page 70 [Assamese]
//	  71  Page 71 [Oriya]
//	  72  Page 72 [Kannada]
//	  73  Page 73 [Malayalam]
//	  74  Page 74 [Gujarati]
//	  75  Page 75 [Punjabi]
//	  82  Page 82 [Marathi]
//	  254 Page 254
//	  255 Page 255
//
// Notes:
//   - The alphanumeric range (ASCII 0x20–0x7F / decimal 32–127) is the same
//     across pages; differences appear in the extended range (0x80–0xFF).
//   - The selected code table remains in effect until ESC @ is executed, the
//     printer is reset, or power is turned off.
//   - Consult your printer's Character Code Tables for exact glyph mappings
//     per page and model-specific supported pages.
//
// Byte sequence:
//
//	ESC t n -> 0x1B, 0x74, n
func (c *Commands) SelectCharacterCodeTable(n CodeTable) ([]byte, error) {
	// Common pages
	validPages := map[CodeTable]bool{
		0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true,
		11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
		20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true,
		30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
		40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
		50: true, 51: true, 52: true, 53: true,
		66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true, 74: true, 75: true,
		82:  true,
		254: true, 255: true,
	}

	if !validPages[n] {
		return nil, ErrCodeTablePage
	}
	return []byte{common.ESC, 't', byte(n)}, nil
}

// SetUpsideDownMode turns upside-down (180° rotated) print mode on or off.
//
// Format:
//
//	ASCII: ESC { n
//	Hex:   0x1B 0x7B n
//	Decimal: 27 123 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	In Standard mode, turns upside-down print mode on or off. When the
//	least-significant bit (LSB) of n is 0, upside-down print mode is turned
//	off. When the LSB of n is 1, upside-down print mode is turned on.
//
// Notes:
//   - In Standard mode this command is only valid when processed at the
//     beginning of a line.
//   - Upside-down mode is effective for all Standard-mode data except certain
//     graphics and obsolete raster/variable-size image commands (see model
//     documentation).
//   - The mode has no effect in Page mode.
//   - When turned on, characters are printed rotated 180° from right to
//     left. The line printing order is not reversed, so take care with the
//     order of transmitted data.
//   - Settings remain in effect until ESC @ is executed, the printer is
//     reset, or power is turned off.
//
// Byte sequence:
//
//	ESC { n -> 0x1B, 0x7B, n
func (c *Commands) SetUpsideDownMode(n UpsideDownMode) []byte {
	return []byte{common.ESC, '{', byte(n)}
}

// TODO: Check if SelectCharacterSize need extra work with the masks and values.
// TODO: Check if conditionals are needed too. According to:
// [Range]
// n = 0xxx0xxxb (n = 0 – 7, 16 – 23, 32 – 39, 48 – 55, 64 – 71, 80 – 87, 96 – 103, 112 – 119)
// (Enlargement in vertical direction: 1–8, Enlargement in horizontal direction: 1–8)

const (
	// HeightMask is used to extract height bits from n
	HeightMask Size = 0x07 // bits 0-2
	// WidthMask is used to extract width bits from n
	WidthMask Size = 0x70 // bits 4-6
	// WidthShift is used to shift width bits to LSB position
	WidthShift Size = 4
)

// BuildCharacterSize Helper functions for building character sizes
func BuildCharacterSize(width, height byte) (Size, error) {
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

// SelectCharacterSize selects character size (width and height magnification).
//
// Format:
//
//	ASCII: GS ! n
//	Hex:   0x1D 0x21 n
//	Decimal: 29 33 n
//
// Range:
//
//	n = 0xxx0xxxb (width and height encoded in a single byte)
//	Width magnification: 1–8
//	Height magnification: 1–8
//	Valid n examples: 0–7, 16–23, 32–39, 48–55, 64–71, 80–87, 96–103, 112–119
//
// Default:
//
//	n = 0 (normal size)
//
// Description:
//
//	Selects character size (height and width magnification). Bits in n are
//	used as follows:
//
//	  Bits 0–2: Height magnification (value 0..7 -> x1..x8 where stored value = height-1)
//	  Bits 4–6: Width magnification  (value 0..7 -> x1..x8 where stored value = width-1)
//
//	In other words:
//	  n = ((width-1) << 4) | (height-1)
//
// Notes:
//   - The character size set by this command is effective for alphanumeric,
//     Kana, multilingual, and user-defined characters.
//   - When characters on a line have different heights, they are aligned at
//     the baseline.
//   - Width enlargement extends characters to the right from the left side.
//   - ESC ! can also toggle double-width and double-height modes.
//   - In Standard mode double-height enlarges in the paper-feed direction
//     and double-width enlarges perpendicular to the paper feed. In 90°
//     rotated mode the relationship is reversed. In Page mode the size
//     follows the character orientation.
//   - The setting for alphanumeric and Katakana remains until ESC !,
//     ESC @, reset, or power-off. For Kanji/multilingual chars the setting
//     remains until FS !, FS W, ESC @, reset, or power-off.
//
// Byte sequence:
//
//	GS ! n -> 0x1D, 0x21, n
func (c *Commands) SelectCharacterSize(n Size) []byte {
	return []byte{common.GS, '!', byte(n)}
}

// SetWhiteBlackReverseMode turns white/black reverse print mode on or off.
//
// Format:
//
//	ASCII: GS B n
//	Hex:   0x1D 0x42 n
//	Decimal: 29 66 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Turns white/black reverse print mode on or off. When the least-significant
//	bit (LSB) of n is 0, reverse mode is turned off. When the LSB of n is 1,
//	reverse mode is turned on.
//
// Notes:
//   - The white/black reverse print is effective for both single-byte and
//     multi-byte code characters.
//   - When reverse mode is turned on, characters are printed in white on a
//     black background.
//   - Reverse mode affects right-side character spacing set by ESC SP and
//     left/right spacing of multi-byte characters set by FS S.
//   - Reverse mode does not affect line spacing or spaces skipped by HT,
//     ESC $, or ESC \.
//   - When underline mode is turned on, the printer does not underline
//     white/black reversed characters.
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or the power is turned off.
//
// Byte sequence:
//
//	GS B n -> 0x1D, 0x42, n
func (c *Commands) SetWhiteBlackReverseMode(n ReverseMode) []byte {
	return []byte{common.GS, 'B', byte(n)}
}

// SetSmoothingMode turns smoothing mode on or off.
//
// Format:
//
//	ASCII: GS b n
//	Hex:   0x1D 0x62 n
//	Decimal: 29 98 n
//
// Range:
//
//	n = 0–255
//
// Default:
//
//	n = 0
//
// Description:
//
//	Turns smoothing mode on or off. When the least-significant bit (LSB) of
//	n is 0, smoothing mode is turned off. When the LSB of n is 1, smoothing
//	mode is turned on.
//
// Notes:
//   - The smoothing mode is effective for quadruple-size or larger characters
//     (alphanumeric, Kana, multilingual, and user-defined characters).
//   - The setting remains in effect until ESC @ is executed, the printer is
//     reset, or the power is turned off.
//
// Byte sequence:
//
//	GS b n -> 0x1D, 0x62, n
func (c *Commands) SetSmoothingMode(n SmoothingMode) []byte {
	return []byte{common.GS, 'b', byte(n)}
}
