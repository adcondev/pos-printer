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

// TextLn agrega un salto de línea al final
func (p *Commands) TextLn(str string) []byte {
	text := p.Text(str)
	// Agregar LF al final
	return append(text, LF)
}

// TextRaw envía bytes sin procesar
func (p *Commands) TextRaw(str string) []byte {
	return []byte(str)
}

func PrintAndFeedPaper(n byte) []byte {
	return []byte{ESC, 'J', n}
}
