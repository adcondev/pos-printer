// Package document provides structures and functions to build print documents.
package document

import (
	"encoding/json"
	"fmt"
)

// PrintJob representa los datos de un documento de impresión
type PrintJob struct {
	Data Document `json:"data"`
}

// Document representa un documento de impresión completo
type Document struct {
	Version  string        `json:"version"`
	Profile  ProfileConfig `json:"profile"`
	Commands []Command     `json:"commands"`
}

// ProfileConfig configuración del perfil de impresora
type ProfileConfig struct {
	Model      string `json:"model"`
	PaperWidth int    `json:"paper_width"`
	CodeTable  string `json:"code_table"`
	DPI        int    `json:"dpi,omitempty"`
	HasQR      bool   `json:"has_qr"` // Indica si soporta QR nativo
}

// Command represents a single command in the document
type Command struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TextCommand represents a text command
type TextCommand struct {
	Content string    `json:"content"`
	Style   TextStyle `json:"style,omitempty"`
	NewLine bool      `json:"newline,omitempty"`
}

// TextStyle estilo de texto
type TextStyle struct {
	Align     string `json:"align,omitempty"` // left, center, right
	Bold      bool   `json:"bold,omitempty"`
	Size      string `json:"size,omitempty"` // normal, 2x2, 3x3
	Underline bool   `json:"underline,omitempty"`
	Inverse   bool   `json:"inverse,omitempty"`
}

// ImageCommand represents an image command
type ImageCommand struct {
	Code       string `json:"code"`                  // Base64
	Format     string `json:"format,omitempty"`      // png, jpg, bmp
	PixelWidth int    `json:"pixel_width,omitempty"` // Ancho deseado en píxeles // FIXME: Considerar JSON
	Align      string `json:"align,omitempty"`       // Alineación
	Threshold  byte   `json:"threshold,omitempty"`   // Umbral B/N (0-255)
	Dithering  string `json:"dithering,omitempty"`   // threshold, atkinson
	Scaling    string `json:"scaling,omitempty"`     // bilinear, nns // FIXME: Considerar JSON
}

// SeparatorCommand represents a separator command
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`   // Carácter a usar
	Length int    `json:"length,omitempty"` // Longitud en caracteres
}

// FeedCommand represents a feed command
type FeedCommand struct {
	Lines int `json:"lines"` // Líneas a avanzar
}

// CutCommand represents a cut command
type CutCommand struct {
	Mode string `json:"mode,omitempty"` // full, partial
	Feed int    `json:"feed,omitempty"` // Líneas antes del corte
}

// TODO: Review for more useful fields

// QRCommand actualizado para soportar todas las opciones
type QRCommand struct {
	Data      string `json:"data"`                 // Datos del QR (URL, texto, etc.)
	HumanText string `json:"human_text,omitempty"` // Texto a mostrar debajo del QR

	// Opciones básicas
	PixelWidth int    `json:"pixel_width,omitempty"` // Pixel size
	Correction string `json:"correction,omitempty"`  // L, M, Q, H
	Align      string `json:"align,omitempty"`       // left, center, right

	// Opciones avanzadas (solo imagen)
	Logo        string `json:"logo_path,omitempty"`    // Ruta relativa al logo
	CircleShape bool   `json:"circle_shape,omitempty"` // Usar bloques circulares
}

// TableCommand represents a table command (WIP)
type TableCommand struct {
	Columns []TableColumn `json:"columns"`
	Rows    [][]string    `json:"rows"`
	// TODO: Adequate fields still not defined
}

// TableColumn define una columna de tabla
type TableColumn struct {
	Width int    `json:"width"`
	Align string `json:"align,omitempty"`
}

// ParseDocument parsea un documento JSON
func ParseDocument(data []byte) (*Document, error) {
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Validación básica
	if doc.Version == "" {
		// TODO: Review an smart way to handle versions
		doc.Version = "1.0"
	}

	if len(doc.Commands) == 0 {
		return nil, fmt.Errorf("document must contain at least one command")
	}

	return &doc, nil
}

// ToJSON convierte el documento a JSON
func (d *Document) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
