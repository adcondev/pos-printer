package test

import "github.com/adcondev/pos-printer/escpos/common"

// BuildLittleEndianCommand creates a command with little-endian encoded value
func BuildLittleEndianCommand(prefix []byte, value uint16, suffix ...byte) []byte {
	nL, nH := common.ToLittleEndian(value)
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

// BuildTestData creates a byte array of specified size with pattern
func BuildTestData(size int, pattern byte) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = pattern
	}
	return data
}

// BuildCommand constructs a command with variable parameters
func BuildCommand(cmd byte, subcmd byte, params ...byte) []byte {
	result := []byte{cmd, subcmd}
	return append(result, params...)
}

// BuildSequence creates a sequence of bytes from start to end (inclusive)
func BuildSequence(start, end byte) []byte {
	if start > end {
		return []byte{}
	}
	result := make([]byte, end-start+1)
	for i := range result {
		result[i] = start + byte(i)
	}
	return result
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

// BuildPattern creates a repeating pattern of bytes
func BuildPattern(size int, pattern []byte) []byte {
	if len(pattern) == 0 {
		return make([]byte, size)
	}
	result := make([]byte, size)
	for i := range result {
		result[i] = pattern[i%len(pattern)]
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
	length := uint16(len(data))
	nL, nH := common.ToLittleEndian(length)
	result := make([]byte, 2+len(data))
	result[0] = nL
	result[1] = nH
	copy(result[2:], data)
	return result
}
