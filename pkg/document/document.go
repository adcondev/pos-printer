package document

import (
	"encoding/json"
	"fmt"
)

// PrintData representa los datos de un documento de impresión
type PrintData struct {
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
	Model     string `json:"model"`      // Modelo de impresora
	Width     int    `json:"width"`      // Ancho en dots
	CodeTable string `json:"code_table"` // Tabla de caracteres
}

// Command representa un comando individual
type Command struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TextCommand comando de texto
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

// ImageCommand comando de imagen
type ImageCommand struct {
	Code      string `json:"code"`                // Base64
	Format    string `json:"format,omitempty"`    // png, jpg, bmp
	Width     int    `json:"width,omitempty"`     // Ancho deseado
	Align     string `json:"align,omitempty"`     // Alineación
	Threshold byte   `json:"threshold,omitempty"` // Umbral B/N (0-255)
	Dithering string `json:"dithering,omitempty"` // threshold, atkinson
}

// SeparatorCommand comando de separador
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`   // Carácter a usar
	Length int    `json:"length,omitempty"` // Longitud en caracteres
}

// FeedCommand comando de avance
type FeedCommand struct {
	Lines int `json:"lines"` // Líneas a avanzar
}

// CutCommand comando de corte
type CutCommand struct {
	Mode string `json:"mode,omitempty"` // full, partial
	Feed int    `json:"feed,omitempty"` // Líneas antes del corte
}

// QRCommand comando QR (preparado para futura implementación)
type QRCommand struct {
	Data  string `json:"data"`
	Size  string `json:"size,omitempty"`  // small, medium, large
	Align string `json:"align,omitempty"` // Alineación
	// TODO: Implementar cuando se agregue soporte QR
}

// TableCommand comando de tabla (preparado para futura implementación)
type TableCommand struct {
	Columns []TableColumn `json:"columns"`
	Rows    [][]string    `json:"rows"`
	// TODO: Implementar cuando se agregue soporte de tablas
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
