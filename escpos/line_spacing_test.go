package escpos

import (
	"bytes"
	"testing"
)

// Naming Convention: Test{Struct}_{Method}_{Scenario}

func TestLineSpacingCommands_SetLineSpacing(t *testing.T) {
	lsc := &LineSpacingCommands{}
	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			"minimum spacing (0 dots)",
			0,
			[]byte{ESC, '3', 0},
		},
		{
			"typical spacing (30 dots)",
			30,
			[]byte{ESC, '3', 30},
		},
		{
			"maximum spacing (255 dots)",
			255,
			[]byte{ESC, '3', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lsc.SetLineSpacing(tt.n)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SetLineSpacing(%d) = %#v; want %#v", tt.n, got, tt.want)
			}
		})
	}
}

func TestLineSpacingCommands_SelectDefaultLineSpacing(t *testing.T) {
	lsc := &LineSpacingCommands{}
	got := lsc.SelectDefaultLineSpacing()
	want := []byte{ESC, '2'}
	if !bytes.Equal(got, want) {
		t.Errorf("SelectDefaultLineSpacing() = %#v; want %#v", got, want)
	}
}
