package common

import "errors"

// Related errors for printer operations
var (
	// ErrEmptyBuffer buffer is empty
	ErrEmptyBuffer = errors.New("can't print an empty buffer")
	// ErrBufferOverflow buffer is too large
	ErrBufferOverflow = errors.New("can't print overflowed buffer (protocol max 64KB; model may be lower)")
	// ErrPrintReverseFeed invalid number of motion units for reverse print
	ErrPrintReverseFeed = errors.New("invalid number of motion units for reverse print (protocol max 48; model may be lower)")
	// ErrPrintReverseFeedLines invalid number of lines for reverse print
	ErrPrintReverseFeedLines = errors.New("invalid number of lines for reverse print (protocol max 2; model may be lower)")
	// ErrLengthOutOfRange length is out of range (0-65535)
	ErrLengthOutOfRange = errors.New("length is out of range (0-65535)")
)
