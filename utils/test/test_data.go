package test

var (
	// Common test barcodes
	UPCA12       = []byte("123456789012")
	UPCA11       = []byte("12345678901")
	UPCE8        = []byte("12345678")
	UPCE7        = []byte("1234567")
	UPCE6        = []byte("123456")
	EAN13_13     = []byte("1234567890123")
	EAN13_12     = []byte("123456789012")
	EAN8_8       = []byte("12345678")
	EAN8_7       = []byte("1234567")
	CODE39_NOSS  = []byte("ABC-123")
	CODE39_SS    = []byte("*ABC-123*")
	ITFEven      = []byte("123456")
	ITFMin       = []byte("12")
	ITFOdd       = []byte("12345")
	CODABAR_SS_1 = []byte("A123456B")
	CODABAR_SS_2 = []byte("C123456D")
	CODABAR_LSS  = []byte("a123456d")

	// Common invalid patterns
	NonNumeric    = []byte("12345A")
	MixedCase     = []byte("AbC123")
	InvalidCode39 = []byte("abc@123")

	// CODE39 allowed characters
	Code39Charset = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./")

	// CODABAR start/stop
	CodabarStartStop = []byte{'A', 'B', 'C', 'D', 'a', 'b', 'c', 'd'}
)
