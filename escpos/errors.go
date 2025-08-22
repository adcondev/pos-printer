package escpos

import "errors"

// Related errors for printer operations
var (
	// errEmptyBuffer buffer is empty
	errEmptyBuffer = errors.New("can't print an empty buffer")
	// errBufferOverflow buffer is too large
	errBufferOverflow = errors.New("can't print overflowed buffer (protocol max 64KB; model may be lower)")
	// errPrintReverseFeed invalid number of motion units for reverse print
	errPrintReverseFeed = errors.New("invalid number of motion units for reverse print (protocol max 48; model may be lower)")
	// errPrintReverseFeedLines invalid number of lines for reverse print
	errPrintReverseFeedLines = errors.New("invalid number of lines for reverse print (protocol max 2; model may be lower)")
	// errNegativeInt negative integer where positive was expected
	errNegativeInt = errors.New("invalid negative integer where positive was expected")
)
