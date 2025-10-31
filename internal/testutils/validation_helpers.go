package testutils

// IsUppercase checks if all alphabetic bytes are uppercase
func IsUppercase(data []byte) bool {
	hasAlpha := false
	for _, b := range data {
		if b >= 'a' && b <= 'z' {
			return false
		}
		if b >= 'A' && b <= 'Z' {
			hasAlpha = true
		}
	}
	return hasAlpha
}

// IsLowercase checks if all alphabetic bytes are lowercase
func IsLowercase(data []byte) bool {
	hasAlpha := false
	for _, b := range data {
		if b >= 'A' && b <= 'Z' {
			return false
		}
		if b >= 'a' && b <= 'z' {
			hasAlpha = true
		}
	}
	return hasAlpha
}

// IsASCII checks if all bytes are valid ASCII (0-127)
func IsASCII(data []byte) bool {
	for _, b := range data {
		if b > 127 {
			return false
		}
	}
	return true
}

// IsPrintableASCII checks if all bytes are printable ASCII (32-126)
func IsPrintableASCII(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
}

// HasNullTerminator checks if data ends with null byte
func HasNullTerminator(data []byte) bool {
	return len(data) > 0 && data[len(data)-1] == 0
}

// IsEvenLength checks if byte slice has even length
func IsEvenLength(data []byte) bool {
	return len(data)%2 == 0
}

// IsInRange checks if all bytes are within specified range
func IsInRange(data []byte, min, max byte) bool {
	for _, b := range data {
		if b < min || b > max {
			return false
		}
	}
	return true
}

// ContainsOnly checks if data contains only bytes from allowed set
func ContainsOnly(data []byte, allowed []byte) bool {
	allowedMap := make(map[byte]bool)
	for _, b := range allowed {
		allowedMap[b] = true
	}
	for _, b := range data {
		if !allowedMap[b] {
			return false
		}
	}
	return true
}

// ContainsAny checks if data contains any byte from targets
func ContainsAny(data []byte, targets []byte) bool {
	targetMap := make(map[byte]bool)
	for _, b := range targets {
		targetMap[b] = true
	}
	for _, b := range data {
		if targetMap[b] {
			return true
		}
	}
	return false
}

// ValidateLength checks if data length is within bounds
func ValidateLength(data []byte, min, max int) bool {
	length := len(data)
	return length >= min && length <= max
}

// IsNumeric checks if all bytes are numeric ASCII characters
func IsNumeric(data []byte) bool {
	if len(data) == 0 {
		return true // empty is considered valid
	}
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// IsAlpha checks if all bytes are alphabetic ASCII characters
func IsAlpha(data []byte) bool {
	if len(data) == 0 {
		return false // empty is not alpha
	}
	for _, b := range data {
		if (b < 'A' || b > 'Z') && (b < 'a' || b > 'z') {
			return false
		}
	}
	return true
}

// IsAlphanumeric checks if all bytes are alphanumeric ASCII characters
func IsAlphanumeric(data []byte) bool {
	if len(data) == 0 {
		return true // empty is considered valid
	}
	for _, b := range data {
		if !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')) {
			return false
		}
	}
	return true
}

// IsHexadecimal checks if all bytes are valid hexadecimal characters
func IsHexadecimal(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	for _, b := range data {
		if !((b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')) {
			return false
		}
	}
	return true
}

// IsBinary checks if all bytes are '0' or '1'
func IsBinary(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	for _, b := range data {
		if b != '0' && b != '1' {
			return false
		}
	}
	return true
}

// IsWhitespace checks if all bytes are whitespace characters
func IsWhitespace(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	for _, b := range data {
		if b != ' ' && b != '\t' && b != '\n' && b != '\r' {
			return false
		}
	}
	return true
}

// StartsWith checks if data starts with prefix
func StartsWith(data, prefix []byte) bool {
	if len(data) < len(prefix) {
		return false
	}
	for i, b := range prefix {
		if data[i] != b {
			return false
		}
	}
	return true
}

// EndsWith checks if data ends with suffix
func EndsWith(data, suffix []byte) bool {
	if len(data) < len(suffix) {
		return false
	}
	offset := len(data) - len(suffix)
	for i, b := range suffix {
		if data[offset+i] != b {
			return false
		}
	}
	return true
}
