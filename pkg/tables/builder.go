// Package tables provides table generation and rendering for ESC/POS printers
package tables

import (
	"encoding/json"
	"io"
)

// Builder helps construct tables programmatically
type Builder struct {
	definition Definition
	rows       []Row
	options    *Options
}

// NewBuilder creates a new table builder
func NewBuilder() *Builder {
	return &Builder{
		definition: Definition{
			Columns: []Column{},
		},
		rows:    []Row{},
		options: DefaultOptions(),
	}
}

// AddColumn adds a column to the table definition
func (b *Builder) AddColumn(header string, width int, align Alignment) *Builder {
	b.definition.Columns = append(b.definition.Columns, Column{
		Header: header,
		Width:  width,
		Align:  align,
	})
	return b
}

// SetPaperWidth sets the paper width in characters
func (b *Builder) SetPaperWidth(width int) *Builder {
	b.options.PaperWidth = width
	b.definition.PaperWidth = width
	return b
}

// SetOptions configures rendering options
func (b *Builder) SetOptions(opts *Options) *Builder {
	if opts != nil {
		b.options = opts
	}
	return b
}

// AddRow adds a data row
func (b *Builder) AddRow(cells ...string) *Builder {
	b.rows = append(b.rows, Row(cells))
	return b
}

// Build creates the final Data structure
func (b *Builder) Build() *Data {
	return &Data{
		Definition:  b.definition,
		ShowHeaders: b.options.ShowHeaders,
		Rows:        b.rows,
	}
}

// ToJSON converts the table data to JSON
func (b *Builder) ToJSON() ([]byte, error) {
	data := b.Build()
	return json.MarshalIndent(data, "", "  ")
}

// RenderTo renders the table to the specified writer
func (b *Builder) RenderTo(w io.Writer) error {
	data := b.Build()
	engine := NewEngine(&b.definition, b.options)
	return engine.Render(w, data)
}
