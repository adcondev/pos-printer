// Package models defines the data structures for ESC/POS document representation.
package models

import (
	"encoding/json"
)

// Document represents a complete ESC/POS document structure
type Document struct {
	Config *Config `json:"config,omitempty"`
	Header []Block `json:"header,omitempty"`
	Body   []Block `json:"body,omitempty"`
	Footer []Block `json:"footer,omitempty"`
}

// Config holds configuration settings for the ESC/POS document
type Config struct {
	Printer      string `json:"printer,omitempty"`
	PaperWidth   int    `json:"paper_width,omitempty"`
	CharsPerLine int    `json:"chars_per_line,omitempty"`
	Encoding     string `json:"encoding,omitempty"`
}

// Block represents a single content block in the ESC/POS document
type Block struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TextBlock represents a text content block
type TextBlock struct {
	Content string   `json:"content"`
	Style   []string `json:"style,omitempty"`
	Size    string   `json:"size,omitempty"`
	Align   string   `json:"align,omitempty"`
}

// ImageBlock represents an image content block
type ImageBlock struct {
	Source    string `json:"source"`
	Data      string `json:"data"`
	Algorithm string `json:"algorithm"`
	Align     string `json:"align,omitempty"`
}

// TODO: Table data input and parsing will be pending.
