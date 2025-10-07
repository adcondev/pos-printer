package test

import "bytes"

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

// ReplaceBytes replaces all occurrences of old with new in data
func ReplaceBytes(data, old, new []byte) []byte {
	return bytes.ReplaceAll(data, old, new)
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
