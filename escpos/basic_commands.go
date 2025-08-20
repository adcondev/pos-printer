package escpos

import "strings"

const (
	LF byte = 0x0A

	ESC byte = 0x1B
)

// Text convierte texto a bytes con encoding apropiado
func (p *Commands) Text(str string) []byte {
	cmd := strings.ReplaceAll(str, "\n", string(LF))
	return []byte(cmd)
}

// Ln agrega un salto de línea al final
func (p *Commands) Ln(str string) []byte {
	text := p.Text(str)
	// Agregar LF al final
	return append(text, LF)
}

// Raw envía bytes sin procesar
func (p *Commands) Raw(str string) []byte {
	return []byte(str)
}

func PrintAndFeedPaper(n byte) []byte {
	return []byte{ESC, 'J', n}
}
