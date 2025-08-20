package escpos

import (
	"fmt"
)

var cutMap = map[CutPaper]byte{
	FullCut:    '0',
	PartialCut: '1',
}

// Cut genera comando de corte
func (p *Commands) Cut(mode CutPaper) ([]byte, error) {
	cut, ok := cutMap[mode]
	if !ok {
		return nil, fmt.Errorf("invalid cut mode: %v", mode)
	}

	return []byte{GS, 'V', cut}, nil
}

// Feed genera comando de alimentación de papel
func (p *Commands) Feed(lines byte) []byte {

	// ESC d n - Print and feed n lines
	return []byte{ESC, 'd', lines}
}

func (p *Commands) FeedReverse(lines byte) []byte {
	return []byte{ESC, 'e', lines}
}
