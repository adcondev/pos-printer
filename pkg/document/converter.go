// Package document proporciona estructuras y funciones para construir documentos de impresión.
package document

import (
	"encoding/json"
	"log"
)

// Builder ayuda a construir documentos programáticamente
type Builder struct {
	doc *Document
}

// NewBuilder crea un nuevo constructor de documentos
func NewBuilder() *Builder {
	return &Builder{
		doc: &Document{
			Version:  "1.0",
			Commands: []Command{},
		},
	}
}

// SetProfile configura el perfil de impresora
func (b *Builder) SetProfile(model string, width int, codeTable string) *Builder {
	b.doc.Profile = ProfileConfig{
		Model:     model,
		Width:     width,
		CodeTable: codeTable,
	}
	return b
}

// AddText generates a text command
func (b *Builder) AddText(content string, style *TextStyle) *Builder {
	cmd := TextCommand{
		Content: content,
		NewLine: true,
	}
	if style != nil {
		cmd.Style = *style
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling text command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "text",
		Data: data,
	})
	return b
}

// AddImage creates an image command
func (b *Builder) AddImage(base64Data string, width int, align string) *Builder {
	cmd := ImageCommand{
		Code:      base64Data,
		Width:     width,
		Align:     align,
		Dithering: "threshold",
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling image command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "image",
		Data: data,
	})
	return b
}

// AddSeparator agrega un separador
func (b *Builder) AddSeparator(char string, length int) *Builder {
	cmd := SeparatorCommand{
		Char:   char,
		Length: length,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling separator command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "separator",
		Data: data,
	})
	return b
}

// AddFeed agrega avance de papel
func (b *Builder) AddFeed(lines int) *Builder {
	cmd := FeedCommand{Lines: lines}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling feed command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "feed",
		Data: data,
	})
	return b
}

// AddCut agrega corte de papel
func (b *Builder) AddCut(mode string, feed int) *Builder {
	cmd := CutCommand{
		Mode: mode,
		Feed: feed,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling cut command: %v", err)
	}
	b.doc.Commands = append(b.doc.Commands, Command{
		Type: "cut",
		Data: data,
	})
	return b
}

// Build construye el documento final
func (b *Builder) Build() *Document {
	return b.doc
}
