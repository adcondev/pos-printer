package character

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	ErrInvalidUnderlineMode = fmt.Errorf("invalid underline mode(try 0-2 or '0'..'2')")
	ErrInvalidRotationMode  = fmt.Errorf("invalid rotation mode(try 0-2 or '0'..'2')")
	ErrInvalidPrintColor    = fmt.Errorf("invalid print color(try 0-1 or '0'..'1')")
	ErrInvalidCharacterFont = fmt.Errorf("invalid character font(try 0-4 or '0'..'4')")
	ErrInvalidCodeTablePage = fmt.Errorf("invalid code table page(try check model support!)")
	ErrInvalidCharacterSet  = fmt.Errorf("invalid international character set(try check model support!)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance checks
var _ Capability = (*Commands)(nil)

// Capability defines the main interface for character-related commands
type Capability interface {
	SetRightSideCharacterSpacing(spacing byte) []byte
	SelectPrintModes(modeBits byte) []byte
	SetUnderlineMode(thickness byte) ([]byte, error)
	SetEmphasizedMode(mode byte) []byte
	SetDoubleStrikeMode(mode byte) []byte
	SelectCharacterFont(fontType byte) ([]byte, error)
	SelectInternationalCharacterSet(charset byte) ([]byte, error)
	Set90DegreeClockwiseRotationMode(rotationMode byte) ([]byte, error)
	SelectPrintColor(color byte) ([]byte, error)
	SelectCharacterCodeTable(page byte) ([]byte, error)
	SetUpsideDownMode(mode byte) []byte
	SelectCharacterSize(sizeConfig byte) []byte
	SetWhiteBlackReverseMode(mode byte) []byte
	SetSmoothingMode(mode byte) []byte
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
func (c *Commands) SetRightSideCharacterSpacing(n byte) []byte {
	return []byte{common.ESC, common.SP, n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	// TODO: placehold this constants for future validation maps
	// Font selection
	PrintModeFontA byte = 0x00 // Bit0 = 0
	PrintModeFontB byte = 0x01 // Bit0 = 1

	// Emphasis
	PrintModeEmphasizedOff byte = 0x00
	PrintModeEmphasizedOn  byte = 0x08 // Bit3

	// Character size
	PrintModeDoubleHeightOff byte = 0x00
	PrintModeDoubleHeightOn  byte = 0x10 // Bit4
	PrintModeDoubleWidthOff  byte = 0x00
	PrintModeDoubleWidthOn   byte = 0x20 // Bit5

	// Underline
	PrintModeUnderlineOff byte = 0x00
	PrintModeUnderlineOn  byte = 0x80 // Bit7
)

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
func (c *Commands) SelectPrintModes(n byte) []byte {
	return []byte{common.ESC, '!', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	// TODO: placehold this constants for future validation maps
	UnderlineModeOff  byte = 0x00 // n = 0
	UnderlineMode1Dot byte = 0x01 // n = 1 (1-dot)
	UnderlineMode2Dot byte = 0x02 // n = 2 (2-dots)

	// ASCII-code variants (some implementations send ASCII digit codes 48/49/50)
	UnderlineModeOffASCII  byte = '0'
	UnderlineMode1DotASCII byte = '1'
	UnderlineMode2DotASCII byte = '2'
)

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
func (c *Commands) SetUnderlineMode(n byte) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid values
	default:
		return nil, ErrInvalidUnderlineMode
	}
	return []byte{common.ESC, '-', n}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	// TODO: placehold this constants for future validation maps
	EOff byte = 0x00 // LSB = 0 -> emphasized OFF
	EOn  byte = 0x01 // LSB = 1 -> emphasized ON

	// ASCII-digit variants some implementations accept
	EOffASCII byte = '0'
	EOnASCII  byte = '1'
)

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
func (c *Commands) SetEmphasizedMode(n byte) []byte {
	return []byte{common.ESC, 'E', n}
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
func (c *Commands) SetDoubleStrikeMode(n byte) []byte {
	return []byte{common.ESC, 'G', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	FontA        byte = 0x00 // n = 0
	FontB        byte = 0x01 // n = 1
	FontC        byte = 0x02 // n = 2
	FontD        byte = 0x03 // n = 3
	FontE        byte = 0x04 // n = 4
	SpecialFontA byte = 97   // n = 97
	SpecialFontB byte = 98   // n = 98

	FontAASCII byte = '0'
	FontBASCII byte = '1'
	FontCASCII byte = '2'
	FontDASCII byte = '3'
	FontEASCII byte = '4'
)

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
func (c *Commands) SelectCharacterFont(n byte) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, 3, 4:
		// Numeric values
	case '0', '1', '2', '3', '4':
		// ASCII values
	case 97, 98:
		// Special fonts
	default:
		return nil, ErrInvalidCharacterFont
	}
	return []byte{common.ESC, 'M', n}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	CharsetUSA          byte = 0
	CharsetFrance       byte = 1
	CharsetGermany      byte = 2
	CharsetUK           byte = 3
	CharsetDenmarkI     byte = 4
	CharsetSweden       byte = 5
	CharsetItaly        byte = 6
	CharsetSpainI       byte = 7
	CharsetJapan        byte = 8
	CharsetNorway       byte = 9
	CharsetDenmarkII    byte = 10
	CharsetSpainII      byte = 11
	CharsetLatinAmerica byte = 12
	CharsetKorea        byte = 13
	CharsetSloveniaCro  byte = 14
	CharsetChina        byte = 15
	CharsetVietnam      byte = 16
	CharsetArabia       byte = 17

	// Extended India character sets (model-dependent)
	CharsetIndiaDevanagari byte = 66
	CharsetIndiaBengali    byte = 67
	CharsetIndiaTamil      byte = 68
	CharsetIndiaTelugu     byte = 69
	CharsetIndiaAssamese   byte = 70
	CharsetIndiaOriya      byte = 71
	CharsetIndiaKannada    byte = 72
	CharsetIndiaMalayalam  byte = 73
	CharsetIndiaGujarati   byte = 74
	CharsetIndiaPunjabi    byte = 75
	CharsetIndiaMarathi    byte = 82
)

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
func (c *Commands) SelectInternationalCharacterSet(n byte) ([]byte, error) {
	// Standard range
	if n <= 17 {
		return []byte{common.ESC, 'R', n}, nil
	}
	// India-specific range
	if n >= 66 && n <= 75 {
		return []byte{common.ESC, 'R', n}, nil
	}
	// Marathi
	if n == 82 {
		return []byte{common.ESC, 'R', n}, nil
	}
	return nil, ErrInvalidCharacterSet
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	Rotation90Off     byte = 0x00 // n = 0
	Rotation90On1Dot  byte = 0x01 // n = 1 (1-dot spacing)
	Rotation90On15Dot byte = 0x02 // n = 2 (1.5-dot spacing)

	// ASCII-digit variants
	Rotation90OffASCII     byte = '0'
	Rotation90On1DotASCII  byte = '1'
	Rotation90On15DotASCII byte = '2'
)

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
func (c *Commands) Set90DegreeClockwiseRotationMode(n byte) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, 2, '0', '1', '2':
		// Valid values
	default:
		return nil, ErrInvalidRotationMode
	}
	return []byte{common.ESC, 'V', n}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	PrintColorBlack      byte = 0x00 // n = 0
	PrintColorRed        byte = 0x01 // n = 1
	PrintColorBlackASCII byte = '0'
	PrintColorRedASCII   byte = '1'
)

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
func (c *Commands) SelectPrintColor(n byte) ([]byte, error) {
	// Validate allowed values
	switch n {
	case 0, 1, '0', '1':
		// Valid values
	default:
		return nil, ErrInvalidPrintColor
	}
	return []byte{common.ESC, 'r', n}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	CodeTablePage0  byte = 0  // PC437: USA, Standard Europe
	CodeTablePage1  byte = 1  // Katakana
	CodeTablePage2  byte = 2  // PC850: Multilingual
	CodeTablePage3  byte = 3  // PC860: Portuguese
	CodeTablePage4  byte = 4  // PC863: Canadian-French
	CodeTablePage5  byte = 5  // PC865: Nordic
	CodeTablePage6  byte = 6  // Hiragana
	CodeTablePage7  byte = 7  // One-pass Kanji
	CodeTablePage8  byte = 8  // One-pass Kanji
	CodeTablePage11 byte = 11 // PC851: Greek
	CodeTablePage12 byte = 12 // PC853: Turkish
	CodeTablePage13 byte = 13 // PC857: Turkish
	CodeTablePage14 byte = 14 // PC737: Greek
	CodeTablePage15 byte = 15 // ISO8859-7: Greek
	CodeTablePage16 byte = 16 // WPC1252
	CodeTablePage17 byte = 17 // PC866: Cyrillic #2
	CodeTablePage18 byte = 18 // PC852: Latin 2
	CodeTablePage19 byte = 19 // PC858: Euro
	CodeTablePage20 byte = 20 // Thai Character Code 42
	CodeTablePage21 byte = 21 // Thai Character Code 11
	CodeTablePage22 byte = 22 // Thai Character Code 13
	CodeTablePage23 byte = 23 // Thai Character Code 14
	CodeTablePage24 byte = 24 // Thai Character Code 16
	CodeTablePage25 byte = 25 // Thai Character Code 17
	CodeTablePage26 byte = 26 // Thai Character Code 18
	CodeTablePage30 byte = 30 // TCVN-3: Vietnamese
	CodeTablePage31 byte = 31 // TCVN-3: Vietnamese
	CodeTablePage32 byte = 32 // PC720: Arabic
	CodeTablePage33 byte = 33 // WPC775: Baltic Rim
	CodeTablePage34 byte = 34 // PC855: Cyrillic
	CodeTablePage35 byte = 35 // PC861: Icelandic
	CodeTablePage36 byte = 36 // PC862: Hebrew
	CodeTablePage37 byte = 37 // PC864: Arabic
	CodeTablePage38 byte = 38 // PC869: Greek
	CodeTablePage39 byte = 39 // ISO8859-2: Latin 2
	CodeTablePage40 byte = 40 // ISO8859-15: Latin 9
	CodeTablePage41 byte = 41 // PC1098: Farsi
	CodeTablePage42 byte = 42 // PC1118: Lithuanian
	CodeTablePage43 byte = 43 // PC1119: Lithuanian
	CodeTablePage44 byte = 44 // PC1125: Ukrainian
	CodeTablePage45 byte = 45 // WPC1250: Latin 2
	CodeTablePage46 byte = 46 // WPC1251: Cyrillic
	CodeTablePage47 byte = 47 // WPC1253: Greek
	CodeTablePage48 byte = 48 // WPC1254: Turkish
	CodeTablePage49 byte = 49 // WPC1255: Hebrew
	CodeTablePage50 byte = 50 // WPC1256: Arabic
	CodeTablePage51 byte = 51 // WPC1257: Baltic Rim
	CodeTablePage52 byte = 52 // WPC1258: Vietnamese
	CodeTablePage53 byte = 53 // KZ-1048: Kazakhstan

	// India-related pages (model-dependent)
	CodeTablePage66 byte = 66 // Devanagari
	CodeTablePage67 byte = 67 // Bengali
	CodeTablePage68 byte = 68 // Tamil
	CodeTablePage69 byte = 69 // Telugu
	CodeTablePage70 byte = 70 // Assamese
	CodeTablePage71 byte = 71 // Oriya
	CodeTablePage72 byte = 72 // Kannada
	CodeTablePage73 byte = 73 // Malayalam
	CodeTablePage74 byte = 74 // Gujarati
	CodeTablePage75 byte = 75 // Punjabi
	CodeTablePage82 byte = 82 // Marathi

	// Reserved / special
	CodeTablePage254 byte = 254
	CodeTablePage255 byte = 255
)

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
func (c *Commands) SelectCharacterCodeTable(n byte) ([]byte, error) {
	// Common pages
	validPages := map[byte]bool{
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
		return nil, ErrInvalidCodeTablePage
	}
	return []byte{common.ESC, 't', n}, nil
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	UpsideDownOff byte = 0x00 // LSB = 0 -> upside-down OFF
	UpsideDownOn  byte = 0x01 // LSB = 1 -> upside-down ON

	// ASCII-digit variants sometimes used by implementations
	UpsideDownOffASCII byte = '0'
	UpsideDownOnASCII  byte = '1'
)

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
func (c *Commands) SetUpsideDownMode(n byte) []byte {
	return []byte{common.ESC, '{', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	// CharSizeHeightMask is used to extract height bits from n
	CharSizeHeightMask byte = 0x07 // bits 0-2
	// CharSizeWidthMask is used to extract width bits from n
	CharSizeWidthMask byte = 0x70 // bits 4-6
	// CharSizeWidthShift is used to shift width bits to LSB position
	CharSizeWidthShift byte = 4
)

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	CharSize1x1 byte = 0x00 // width=1 height=1 (normal)
	CharSize2x1 byte = 0x10 // width=2 height=1
	CharSize3x1 byte = 0x20 // width=3 height=1
	CharSize4x1 byte = 0x30 // width=4 height=1
	CharSize1x2 byte = 0x01 // width=1 height=2
	CharSize2x2 byte = 0x11 // width=2 height=2 (double width & height)
)

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
func (c *Commands) SelectCharacterSize(n byte) []byte {
	return []byte{common.GS, '!', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	WhiteBlackReverseOff byte = 0x00 // LSB = 0 -> reverse OFF
	WhiteBlackReverseOn  byte = 0x01 // LSB = 1 -> reverse ON

	// ASCII-digit variants sometimes used by implementations
	WhiteBlackReverseOffASCII byte = '0'
	WhiteBlackReverseOnASCII  byte = '1'
)

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
func (c *Commands) SetWhiteBlackReverseMode(n byte) []byte {
	return []byte{common.GS, 'B', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go
const (
	SmoothingOff byte = 0x00 // LSB = 0 -> smoothing OFF
	SmoothingOn  byte = 0x01 // LSB = 1 -> smoothing ON

	// ASCII-digit variants sometimes accepted by implementations
	SmoothingOffASCII byte = '0'
	SmoothingOnASCII  byte = '1'
)

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
func (c *Commands) SetSmoothingMode(n byte) []byte {
	return []byte{common.GS, 'b', n}
}

// TODO: Mover esto a Wrapper de ESCPOS en pos.go

// BuildCharacterSize Helper functions for building character sizes
func BuildCharacterSize(width, height int) (byte, error) {
	if width < 1 || width > 8 {
		return 0, fmt.Errorf("width must be 1..8, got %d", width)
	}
	if height < 1 || height > 8 {
		return 0, fmt.Errorf("height must be 1..8, got %d", height)
	}
	w := byte(width-1) << 4
	h := byte(height - 1)
	return w | h, nil
}
