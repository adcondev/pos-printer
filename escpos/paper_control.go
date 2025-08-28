package escpos

import (
	"fmt"

	"github.com/adcondev/pos-printer/escpos/common"
)

var cutMap = map[CutPaper]byte{
	FullCut:    '0',
	PartialCut: '1',
}

// Cut genera comando de corte
func (c *Commands) Cut(mode CutPaper) ([]byte, error) {
	cut, ok := cutMap[mode]
	if !ok {
		return nil, fmt.Errorf("invalid cut mode: %v", mode)
	}

	return []byte{common.GS, 'V', cut}, nil
}
