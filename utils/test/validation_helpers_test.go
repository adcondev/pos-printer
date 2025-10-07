package test

import "testing"

func TestIsUppercase(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"uppercase letters", []byte("HELLO"), true},
		{"mixed case", []byte("Hello"), false},
		{"lowercase letters", []byte("hello"), false},
		{"uppercase with numbers", []byte("HELLO123"), true},
		{"uppercase with spaces", []byte("HELLO WORLD"), true},
		{"empty slice", []byte{}, false},
		{"only numbers", []byte("12345"), false},
		{"only symbols", []byte("!@#$%"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUppercase(tt.data)
			if result != tt.expected {
				t.Errorf("IsUppercase(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsLowercase(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"lowercase letters", []byte("hello"), true},
		{"mixed case", []byte("Hello"), false},
		{"uppercase letters", []byte("HELLO"), false},
		{"lowercase with numbers", []byte("hello123"), true},
		{"lowercase with spaces", []byte("hello world"), true},
		{"empty slice", []byte{}, false},
		{"only numbers", []byte("12345"), false},
		{"only symbols", []byte("!@#$%"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLowercase(tt.data)
			if result != tt.expected {
				t.Errorf("IsLowercase(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsASCII(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"ascii text", []byte("Hello World!"), true},
		{"ascii with numbers", []byte("abc123"), true},
		{"ascii symbols", []byte("!@#$%^&*()"), true},
		{"empty slice", []byte{}, true},
		{"non-ascii", []byte{128, 200, 255}, false},
		{"mixed ascii and non-ascii", []byte{'a', 'b', 128}, false},
		{"null byte", []byte{0}, true},
		{"boundary 127", []byte{127}, true},
		{"boundary 128", []byte{128}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsASCII(tt.data)
			if result != tt.expected {
				t.Errorf("IsASCII(%v) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsPrintableASCII(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"printable text", []byte("Hello World!"), true},
		{"printable with numbers", []byte("abc123"), true},
		{"printable symbols", []byte("!@#$%^&*()"), true},
		{"empty slice", []byte{}, true},
		{"with newline", []byte("hello\n"), false},
		{"with tab", []byte("hello\t"), false},
		{"with null", []byte{0}, false},
		{"boundary space (32)", []byte{32}, true},
		{"boundary tilde (126)", []byte{126}, true},
		{"below range (31)", []byte{31}, false},
		{"above range (127)", []byte{127}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrintableASCII(tt.data)
			if result != tt.expected {
				t.Errorf("IsPrintableASCII(%v) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestHasNullTerminator(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"null terminated", []byte("hello\x00"), true},
		{"not null terminated", []byte("hello"), false},
		{"empty slice", []byte{}, false},
		{"only null", []byte{0}, true},
		{"null in middle", []byte("hel\x00lo"), false},
		{"multiple nulls at end", []byte("hello\x00\x00"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasNullTerminator(tt.data)
			if result != tt.expected {
				t.Errorf("HasNullTerminator(%v) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsEvenLength(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"even length 2", []byte("hi"), true},
		{"even length 4", []byte("test"), true},
		{"odd length 3", []byte("bye"), false},
		{"odd length 5", []byte("hello"), false},
		{"empty slice (0)", []byte{}, true},
		{"length 1", []byte("a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEvenLength(tt.data)
			if result != tt.expected {
				t.Errorf("IsEvenLength(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsInRange(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		min      byte
		max      byte
		expected bool
	}{
		{"all in range", []byte{5, 10, 15}, 0, 20, true},
		{"below range", []byte{5, 10, 15}, 10, 20, false},
		{"above range", []byte{5, 10, 15}, 0, 10, false},
		{"at boundaries", []byte{10, 15, 20}, 10, 20, true},
		{"empty slice", []byte{}, 0, 100, true},
		{"single byte in range", []byte{50}, 0, 100, true},
		{"single byte out of range", []byte{150}, 0, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInRange(tt.data, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("IsInRange(%v, %d, %d) = %v, expected %v", tt.data, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

func TestContainsOnly(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		allowed  []byte
		expected bool
	}{
		{"all allowed", []byte("abc"), []byte("abcdef"), true},
		{"contains disallowed", []byte("abcx"), []byte("abcdef"), false},
		{"empty data", []byte{}, []byte("abc"), true},
		{"empty allowed", []byte("a"), []byte{}, false},
		{"exact match", []byte("abc"), []byte("abc"), true},
		{"repeated allowed", []byte("aaa"), []byte("a"), true},
		{"digits only", []byte("123"), []byte("0123456789"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsOnly(tt.data, tt.allowed)
			if result != tt.expected {
				t.Errorf("ContainsOnly(%q, %q) = %v, expected %v", tt.data, tt.allowed, result, tt.expected)
			}
		})
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		targets  []byte
		expected bool
	}{
		{"contains target", []byte("hello"), []byte("aeiou"), true},
		{"no target", []byte("xyz"), []byte("abc"), false},
		{"empty data", []byte{}, []byte("abc"), false},
		{"empty targets", []byte("abc"), []byte{}, false},
		{"multiple matches", []byte("aeiou"), []byte("aeiou"), true},
		{"single match", []byte("xyzae"), []byte("a"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsAny(tt.data, tt.targets)
			if result != tt.expected {
				t.Errorf("ContainsAny(%q, %q) = %v, expected %v", tt.data, tt.targets, result, tt.expected)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		min      int
		max      int
		expected bool
	}{
		{"within range", []byte("hello"), 1, 10, true},
		{"at min boundary", []byte("h"), 1, 10, true},
		{"at max boundary", []byte("hellohello"), 1, 10, true},
		{"below min", []byte{}, 1, 10, false},
		{"above max", []byte("hello world!"), 1, 10, false},
		{"empty in range", []byte{}, 0, 10, true},
		{"zero min max", []byte{}, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateLength(tt.data, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("ValidateLength(%q, %d, %d) = %v, expected %v", tt.data, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"all digits", []byte("123456"), true},
		{"with letters", []byte("123abc"), false},
		{"with spaces", []byte("123 456"), false},
		{"with symbols", []byte("123!"), false},
		{"empty slice", []byte{}, true},
		{"single digit", []byte("5"), true},
		{"zero", []byte("0"), true},
		{"leading zeros", []byte("000123"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumeric(tt.data)
			if result != tt.expected {
				t.Errorf("IsNumeric(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"all letters", []byte("hello"), true},
		{"mixed case", []byte("HelloWorld"), true},
		{"with numbers", []byte("hello123"), false},
		{"with spaces", []byte("hello world"), false},
		{"with symbols", []byte("hello!"), false},
		{"empty slice", []byte{}, false},
		{"uppercase only", []byte("HELLO"), true},
		{"lowercase only", []byte("world"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAlpha(tt.data)
			if result != tt.expected {
				t.Errorf("IsAlpha(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"letters and numbers", []byte("hello123"), true},
		{"only letters", []byte("hello"), true},
		{"only numbers", []byte("123"), true},
		{"with spaces", []byte("hello 123"), false},
		{"with symbols", []byte("hello!123"), false},
		{"empty slice", []byte{}, true},
		{"mixed case", []byte("HelloWorld123"), true},
		{"underscore", []byte("hello_123"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAlphanumeric(tt.data)
			if result != tt.expected {
				t.Errorf("IsAlphanumeric(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsHexadecimal(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"valid hex lowercase", []byte("abc123"), true},
		{"valid hex uppercase", []byte("ABC123"), true},
		{"valid hex mixed", []byte("AbC123"), true},
		{"all digits", []byte("123456"), true},
		{"all letters a-f", []byte("abcdef"), true},
		{"all letters A-F", []byte("ABCDEF"), true},
		{"invalid letter g", []byte("abcg"), false},
		{"invalid letter G", []byte("ABCG"), false},
		{"with symbols", []byte("abc!"), false},
		{"empty slice", []byte{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHexadecimal(tt.data)
			if result != tt.expected {
				t.Errorf("IsHexadecimal(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"valid binary", []byte("101010"), true},
		{"all zeros", []byte("000000"), true},
		{"all ones", []byte("111111"), true},
		{"with digit 2", []byte("1012"), false},
		{"with letters", []byte("10a10"), false},
		{"with spaces", []byte("101 010"), false},
		{"empty slice", []byte{}, true},
		{"single zero", []byte("0"), true},
		{"single one", []byte("1"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBinary(tt.data)
			if result != tt.expected {
				t.Errorf("IsBinary(%q) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"only spaces", []byte("   "), true},
		{"only tabs", []byte("\t\t\t"), true},
		{"only newlines", []byte("\n\n\n"), true},
		{"only carriage returns", []byte("\r\r\r"), true},
		{"mixed whitespace", []byte(" \t\n\r"), true},
		{"with text", []byte(" hello "), false},
		{"empty slice", []byte{}, false},
		{"single space", []byte(" "), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWhitespace(tt.data)
			if result != tt.expected {
				t.Errorf("IsWhitespace(%v) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		prefix   []byte
		expected bool
	}{
		{"has prefix", []byte("hello world"), []byte("hello"), true},
		{"no prefix", []byte("hello world"), []byte("world"), false},
		{"exact match", []byte("hello"), []byte("hello"), true},
		{"empty prefix", []byte("hello"), []byte{}, true},
		{"empty data", []byte{}, []byte("hello"), false},
		{"both empty", []byte{}, []byte{}, true},
		{"prefix longer", []byte("hi"), []byte("hello"), false},
		{"single byte match", []byte("hello"), []byte("h"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StartsWith(tt.data, tt.prefix)
			if result != tt.expected {
				t.Errorf("StartsWith(%q, %q) = %v, expected %v", tt.data, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		suffix   []byte
		expected bool
	}{
		{"has suffix", []byte("hello world"), []byte("world"), true},
		{"no suffix", []byte("hello world"), []byte("hello"), false},
		{"exact match", []byte("hello"), []byte("hello"), true},
		{"empty suffix", []byte("hello"), []byte{}, true},
		{"empty data", []byte{}, []byte("hello"), false},
		{"both empty", []byte{}, []byte{}, true},
		{"suffix longer", []byte("hi"), []byte("hello"), false},
		{"single byte match", []byte("hello"), []byte("o"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EndsWith(tt.data, tt.suffix)
			if result != tt.expected {
				t.Errorf("EndsWith(%q, %q) = %v, expected %v", tt.data, tt.suffix, result, tt.expected)
			}
		})
	}
}
