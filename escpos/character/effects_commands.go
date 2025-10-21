package character

// SelectCharacterColor selects the character color.
//
// Format:
//
//	ASCII:   GS ( N pL pH fn m
//	Hex:     0x1D 0x28 0x4E 0x02 0x00 0x30 m
//	Decimal: 29 40 78 2 0 48 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 48–51 (model-dependent)
//
// Default:
//
//	m = 49 (Color 1)
//
// Parameters:
//
//	m: Selects the character color:
//	   48 -> None (non-printing dots)
//	   49 -> Color 1 (default)
//	   50 -> Color 2
//	   51 -> Color 3
//
// Notes:
//   - Applies to alphanumeric, Katakana, multilingual, user-defined, and
//     user-defined Kanji characters; does not affect graphics, bit images,
//     barcodes, or 2D codes
//   - m = 48 treats glyph dots as non-printing (useful with background/shadow)
//   - Underline prints in the selected character color
//   - Setting persists until ESC @, printer reset, or power-off
//
// Errors:
//
//	Returns ErrInvalidCharacterColor if m is not a valid color value (48-51).
func (ef *EffectsCommands) SelectCharacterColor(m byte) ([]byte, error) {
	// Validate allowed values
	if err := ValidateCharacterColor(m); err != nil {
		return nil, err
	}
	return []byte{GS, '(', 'N', 0x02, 0x00, 0x30, m}, nil
}

// SelectBackgroundColor selects the background color.
//
// Format:
//
//	ASCII:   GS ( N pL pH fn m
//	Hex:     0x1D 0x28 0x4E 0x02 0x00 0x31 m
//	Decimal: 29 40 78 2 0 49 m
//
// Range:
//
//	pL = 0x02, pH = 0x00
//	m = 48–51 (model-dependent)
//
// Default:
//
//	m = 48 (None)
//
// Parameters:
//
//	m: Selects the background color:
//	   48 -> None (no background dots printed)
//	   49 -> Color 1
//	   50 -> Color 2
//	   51 -> Color 3
//
// Notes:
//   - Background color does not affect spaces skipped by HT, ESC $, ESC \,
//     line spacing, or reverse print background
//   - Inter-character spacing (ESC SP, FS S) prints in this background color
//   - Settings persist until ESC @, printer reset, or power-off
//
// Errors:
//
//	Returns ErrInvalidBackgroundColor if m is not a valid color value (48-51).
func (ef *EffectsCommands) SelectBackgroundColor(m byte) ([]byte, error) {
	// Validate allowed values
	if err := ValidateBackgroundColor(m); err != nil {
		return nil, err
	}
	return []byte{GS, '(', 'N', 0x02, 0x00, 0x31, m}, nil
}

// SetCharacterShadowMode turns shading (shadow) mode on or off.
//
// Format:
//
//	ASCII:   GS ( N pL pH fn m a
//	Hex:     0x1D 0x28 0x4E pL pH 0x32 m a
//	Decimal: 29 40 78 pL pH 50 m a
//
// Range:
//
//	pL = 0x03, pH = 0x00
//	m = 0, 1, 48, 49
//	a = 48–51
//
// Default:
//
//	m = 0, a = 48
//
// Parameters:
//
//	m: Shadow mode (0 or 48 = OFF, 1 or 49 = ON)
//	a: Shadow color:
//	   48 -> None (not printed)
//	   49 -> Color 1
//	   50 -> Color 2
//	   51 -> Color 3
//
// Notes:
//   - Parameter a (shadowColor) MUST be supplied even when shadow mode is OFF
//   - Shadow mode prints a shadow in the specified shadow color
//   - Underline shadow is not printed
//   - Reverse print does not alter shadow color
//   - Settings persist until ESC @, printer reset, or power-off
//
// Errors:
//
//	Returns ErrInvalidShadowMode if m is not a valid mode value (0, 1, 48, 49).
//	Returns ErrInvalidShadowColor if a is not a valid color value (48-51).
func (ef *EffectsCommands) SetCharacterShadowMode(m byte, a byte) ([]byte, error) {
	// Validate allowed values
	if err := ValidateShadowMode(m); err != nil {
		return nil, err
	}
	if err := ValidateShadowColor(a); err != nil {
		return nil, err
	}
	return []byte{GS, '(', 'N', 0x03, 0x00, 0x32, m, a}, nil
}
