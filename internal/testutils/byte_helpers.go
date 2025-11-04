package testutils

import (
	"bytes"

	"github.com/adcondev/pos-printer/pkg/controllers/escpos/shared"
)

// Create a test helper for buffer management

// BufferBuilder helps in building byte buffers for tests
type BufferBuilder struct {
	buffer []byte
}

// NewBufferBuilder initializes a new BufferBuilder
func NewBufferBuilder() *BufferBuilder {
	return &BufferBuilder{buffer: make([]byte, 0)}
}

// Append adds bytes to the buffer
func (b *BufferBuilder) Append(cmd []byte) *BufferBuilder {
	b.buffer = append(b.buffer, cmd...)
	return b
}

// GetBuffer returns the constructed byte buffer
func (b *BufferBuilder) GetBuffer() []byte {
	return b.buffer
}

// ============================================================================
// Byte Slice Generators and Manipulators
// ============================================================================

// RepeatByte creates a byte slice of specified length filled with value
func RepeatByte(length int, value byte) []byte {
	result := make([]byte, length)
	for i := range result {
		result[i] = value
	}
	return result
}

// RangeBytes creates a byte slice with sequential values starting from start
func RangeBytes(length int, start byte) []byte {
	result := make([]byte, length)
	for i := range result {
		result[i] = start + byte(i)
	}
	return result
}

// ConcatBytes concatenates multiple byte slices into one
func ConcatBytes(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]byte, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// SplitBytes splits a byte slice at delimiter occurrences
func SplitBytes(data []byte, delimiter byte) [][]byte {
	return bytes.Split(data, []byte{delimiter})
}

// TODO: Check linter

// ReplaceBytes replaces all occurrences of old with new in data
func ReplaceBytes(data, old, news []byte) []byte {
	return bytes.ReplaceAll(data, old, news)
}

// TrimBytes removes leading and trailing bytes matching cutset
func TrimBytes(data []byte, cutset byte) []byte {
	return bytes.Trim(data, string(cutset))
}

// PadBytes pads data to length with padding byte
func PadBytes(data []byte, length int, padding byte) []byte {
	if len(data) >= length {
		return data
	}
	result := make([]byte, length)
	copy(result, data)
	for i := len(data); i < length; i++ {
		result[i] = padding
	}
	return result
}

// PadLeft pads data on the left to reach target length
func PadLeft(data []byte, length int, padding byte) []byte {
	if len(data) >= length {
		return data
	}
	padLen := length - len(data)
	result := make([]byte, length)
	for i := 0; i < padLen; i++ {
		result[i] = padding
	}
	copy(result[padLen:], data)
	return result
}

// ChunkBytes splits data into chunks of specified size
func ChunkBytes(data []byte, chunkSize int) [][]byte {
	if chunkSize <= 0 {
		return nil
	}
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

// ReverseBytes returns a reversed copy of the byte slice
func ReverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, j := 0, len(data)-1; i < len(data); i, j = i+1, j-1 {
		result[i] = data[j]
	}
	return result
}

// UniqueBytes returns unique bytes from data preserving order
func UniqueBytes(data []byte) []byte {
	seen := make(map[byte]bool)
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if !seen[b] {
			seen[b] = true
			result = append(result, b)
		}
	}
	return result
}

// CountByte counts occurrences of a byte in data
func CountByte(data []byte, target byte) int {
	count := 0
	for _, b := range data {
		if b == target {
			count++
		}
	}
	return count
}

// IndexOfByte returns the index of first occurrence of target, or -1
func IndexOfByte(data []byte, target byte) int {
	for i, b := range data {
		if b == target {
			return i
		}
	}
	return -1
}

// LastIndexOfByte returns the index of last occurrence of target, or -1
func LastIndexOfByte(data []byte, target byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == target {
			return i
		}
	}
	return -1
}

// FilterBytes returns only bytes matching the predicate
func FilterBytes(data []byte, predicate func(byte) bool) []byte {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if predicate(b) {
			result = append(result, b)
		}
	}
	return result
}

// MapBytes applies a transformation function to each byte
func MapBytes(data []byte, transform func(byte) byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = transform(b)
	}
	return result
}

// ============================================================================
// Command Builders
// ============================================================================

// BuildLittleEndianCommand creates a command with little-endian encoded value
func BuildLittleEndianCommand(prefix []byte, value uint16, suffix ...byte) []byte {
	nL, nH := shared.ToLittleEndian(value)
	return append(append(prefix, nL, nH), suffix...)
}

// BuildTabPositions creates a sequence of tab positions for testing
func BuildTabPositions(count int, spacing int) []byte {
	tabs := make([]byte, count)
	for i := range tabs {
		tabs[i] = byte((i + 1) * spacing)
	}
	return tabs
}

// BuildCommand constructs a command with variable parameters
func BuildCommand(cmd byte, subcmd byte, params ...byte) []byte {
	result := []byte{cmd, subcmd}
	return append(result, params...)
}

// BuildAlphanumeric creates test data with alphanumeric characters
func BuildAlphanumeric(size int) []byte {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, size)
	for i := range result {
		result[i] = chars[i%len(chars)]
	}
	return result
}

// BuildNumeric creates test data with numeric characters only
func BuildNumeric(size int) []byte {
	result := make([]byte, size)
	for i := range result {
		result[i] = '0' + byte(i%10)
	}
	return result
}

// BuildWithTerminator appends a terminator byte to data
func BuildWithTerminator(data []byte, terminator byte) []byte {
	result := make([]byte, len(data)+1)
	copy(result, data)
	result[len(data)] = terminator
	return result
}

// BuildWithLength prepends length byte(s) to data
func BuildWithLength(data []byte) []byte {
	result := make([]byte, 1+len(data))
	result[0] = byte(len(data))
	copy(result[1:], data)
	return result
}

// BuildWithLittleEndianLength prepends little-endian length to data
func BuildWithLittleEndianLength(data []byte) []byte {
	length := uint16(len(data)) //nolint:gosec
	nL, nH := shared.ToLittleEndian(length)
	result := make([]byte, 2+len(data))
	result[0] = nL
	result[1] = nH
	copy(result[2:], data)
	return result
}
