package escpos

// Buffer limits
const (
	// MinBuf es el tamaño mínimo del buffer
	MinBuf int = 1
	// MaxBuf es el tamaño máximo del buffer
	MaxBuf int = 65535
)

// Reverse motion units and lines
const (
	// MaxReverseMotionUnits is the maximum number of motion units for reverse printing
	MaxReverseMotionUnits byte = 48
	// MaxReverseFeedLines is the maximum number of lines for reverse printing
	MaxReverseFeedLines byte = 2
)
