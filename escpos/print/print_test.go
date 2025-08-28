package print_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Naming Convention: Test{Struct}_{Method}_{Scenario}

func TestPrintCommands_Format_ByteSequence(t *testing.T) {
	input := []byte("a\n\t\rB")
	want := []byte{'a', print.LF, common.HT, print.CR, 'B'}
	// clone input to avoid modifying original
	data := append([]byte(nil), input...)
	got := print.Formatting(data)
	if !bytes.Equal(got, want) {
		t.Errorf("Formatting(%q) = %v; want %v", input, got, want)
	}
}

// ============================================================================
// Commands Tests
// ============================================================================

func TestPrintCommands_Text_ValidInput(t *testing.T) {
	pc := &print.Commands{}

	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:    "simple text",
			input:   "Hello",
			want:    []byte("Hello"),
			wantErr: false,
		},
		{
			name:    "text with control characters",
			input:   "A\nB\tC\rD",
			want:    []byte{'A', print.LF, 'B', common.HT, 'C', print.CR, 'D'},
			wantErr: false,
		},
		{
			name:    "empty text returns error",
			input:   "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pc.Text(tt.input)

			// Check error state
			if (err != nil) != tt.wantErr {
				t.Errorf("Commands.Text(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			// Check result only if no error expected
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("Commands.Text(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPrintCommands_PrintAndFeedPaper_ByteSequence(t *testing.T) {
	pc := &print.Commands{}

	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			name: "minimum feed (0 units)",
			n:    0,
			want: []byte{common.ESC, 'J', 0},
		},
		{
			name: "typical feed (5 units)",
			n:    5,
			want: []byte{common.ESC, 'J', 5},
		},
		{
			name: "maximum feed (255 units)",
			n:    255,
			want: []byte{common.ESC, 'J', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pc.PrintAndFeedPaper(tt.n)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Commands.PrintAndFeedPaper(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}

func TestPrintCommands_FormFeed_SingleByte(t *testing.T) {
	pc := &print.Commands{}

	got := pc.FormFeed()
	want := []byte{print.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.FormFeed() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndCarriageReturn_SingleByte(t *testing.T) {
	pc := &print.Commands{}

	got := pc.PrintAndCarriageReturn()
	want := []byte{print.CR}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.PrintAndCarriageReturn() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndLineFeed_SingleByte(t *testing.T) {
	pc := &print.Commands{}

	got := pc.PrintAndLineFeed()
	want := []byte{print.LF}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.PrintAndLineFeed() = %#v, want %#v", got, want)
	}
}
