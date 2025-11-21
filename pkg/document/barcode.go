package document

// BarcodeCommand represents a barcode command
type BarcodeCommand struct {
	Symbology string `json:"symbology"`              // UPC-A, EAN13, CODE128, etc.
	Data      string `json:"data"`                   // Barcode data
	Width     int    `json:"width,omitempty"`        // Module width (2-6)
	Height    int    `json:"height,omitempty"`       // Height in dots (1-255)
	HRIPos    string `json:"hri_position,omitempty"` // none, above, below, both
	HRIFont   string `json:"hri_font,omitempty"`     // A, B
	Align     string `json:"align,omitempty"`        // left, center, right
}
