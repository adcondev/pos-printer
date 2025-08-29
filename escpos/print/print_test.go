package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// ============================================================================
// Utility Functions Tests
// ============================================================================

// Naming Convention: TestUtility_{Function}_{Optional Scenario}

func TestUtility_Formatting_CharacterReplacement(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "replaces newline with LF",
			input: []byte("Hello\nWorld"),
			want:  []byte{'H', 'e', 'l', 'l', 'o', print.LF, 'W', 'o', 'r', 'l', 'd'},
		},
		{
			name:  "replaces tab with HT",
			input: []byte("Col1\tCol2"),
			want:  []byte{'C', 'o', 'l', '1', common.HT, 'C', 'o', 'l', '2'},
		},
		{
			name:  "replaces carriage return with CR",
			input: []byte("Line1\rLine2"),
			want:  []byte{'L', 'i', 'n', 'e', '1', print.CR, 'L', 'i', 'n', 'e', '2'},
		},
		{
			name:  "handles multiple replacements",
			input: []byte("A\nB\tC\rD"),
			want:  []byte{'A', print.LF, 'B', common.HT, 'C', print.CR, 'D'},
		},
		{
			name:  "preserves regular characters",
			input: []byte("NoSpecialChars123!@#"),
			want:  []byte("NoSpecialChars123!@#"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone input to avoid modifying original
			data := append([]byte(nil), tt.input...)
			got := print.Formatting(data)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Formatting(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Commands Tests
// ============================================================================

// Naming Convention: Test{Struct}_{Method}_{Optional Scenario}

func TestCommands_Text(t *testing.T) {
	cmd := print.NewCommands()

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
			name:    "text with newline",
			input:   "Line1\nLine2",
			want:    []byte{'L', 'i', 'n', 'e', '1', print.LF, 'L', 'i', 'n', 'e', '2'},
			wantErr: false,
		},
		{
			name:    "text with tab",
			input:   "A\tB",
			want:    []byte{'A', common.HT, 'B'},
			wantErr: false,
		},
		{
			name:    "text with carriage return",
			input:   "A\rB",
			want:    []byte{'A', print.CR, 'B'},
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
			got, err := cmd.Text(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Text(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("Text(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCommands_PrintAndFeedPaper(t *testing.T) {
	cmd := print.NewCommands()

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
			name: "typical feed (30 units)",
			n:    30,
			want: []byte{common.ESC, 'J', 30},
		},
		{
			name: "maximum feed (255 units)",
			n:    255,
			want: []byte{common.ESC, 'J', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.PrintAndFeedPaper(tt.n)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndFeedPaper(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}

func TestCommands_FormFeed(t *testing.T) {
	cmd := print.NewCommands()
	got := cmd.FormFeed()
	want := []byte{print.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("FormFeed() = %#v, want %#v", got, want)
	}
}

func TestCommands_PrintAndCarriageReturn(t *testing.T) {
	cmd := print.NewCommands()
	got := cmd.PrintAndCarriageReturn()
	want := []byte{print.CR}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintAndCarriageReturn() = %#v, want %#v", got, want)
	}
}

func TestCommands_PrintAndLineFeed(t *testing.T) {
	cmd := print.NewCommands()
	got := cmd.PrintAndLineFeed()
	want := []byte{print.LF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintAndLineFeed() = %#v, want %#v", got, want)
	}
}

// ============================================================================
// PagePrint Tests
// ============================================================================

func TestPagePrint_PrintDataInPageMode(t *testing.T) {
	pp := &print.PagePrint{}
	got := pp.PrintDataInPageMode()
	want := []byte{common.ESC, print.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_CancelData(t *testing.T) {
	pp := &print.PagePrint{}
	got := pp.CancelData()
	want := []byte{print.CAN}

	if !bytes.Equal(got, want) {
		t.Errorf("CancelData() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_PrintAndReverseFeed(t *testing.T) {
	pp := &print.PagePrint{}

	tests := []struct {
		name    string
		n       byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "minimum reverse feed (0 units)",
			n:       0,
			want:    []byte{common.ESC, 'K', 0},
			wantErr: false,
		},
		{
			name:    "typical reverse feed (10 units)",
			n:       10,
			want:    []byte{common.ESC, 'K', 10},
			wantErr: false,
		},
		{
			name:    "maximum allowed reverse feed",
			n:       print.MaxReverseMotionUnits,
			want:    []byte{common.ESC, 'K', print.MaxReverseMotionUnits},
			wantErr: false,
		},
		{
			name:    "exceeds maximum returns error",
			n:       print.MaxReverseMotionUnits + 1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndReverseFeed(tt.n)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintAndReverseFeed(%d) error = %v, wantErr %v", tt.n, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndReverseFeed(%d) = %#v, want %#v", tt.n, got, tt.want)
			}

			if tt.wantErr && !errors.Is(err, print.ErrPrintReverseFeed) {
				t.Errorf("PrintAndReverseFeed(%d) error = %v, want %v", tt.n, err, print.ErrPrintReverseFeed)
			}
		})
	}
}

func TestPagePrint_PrintAndReverseFeedLines(t *testing.T) {
	pp := &print.PagePrint{}

	tests := []struct {
		name    string
		n       byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "minimum lines (0)",
			n:       0,
			want:    []byte{common.ESC, 'e', 0},
			wantErr: false,
		},
		{
			name:    "single line reverse",
			n:       1,
			want:    []byte{common.ESC, 'e', 1},
			wantErr: false,
		},
		{
			name:    "maximum allowed lines",
			n:       print.MaxReverseFeedLines,
			want:    []byte{common.ESC, 'e', print.MaxReverseFeedLines},
			wantErr: false,
		},
		{
			name:    "exceeds maximum returns error",
			n:       print.MaxReverseFeedLines + 1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndReverseFeedLines(tt.n)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintAndReverseFeedLines(%d) error = %v, wantErr %v", tt.n, err, tt.wantErr)
				return
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndReverseFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
			}

			if tt.wantErr && !errors.Is(err, print.ErrPrintReverseFeedLines) {
				t.Errorf("PrintAndReverseFeedLines(%d) error = %v, want %v", tt.n, err, print.ErrPrintReverseFeedLines)
			}
		})
	}
}

func TestPagePrint_PrintAndFeedLines(t *testing.T) {
	pp := &print.PagePrint{}

	tests := []struct {
		name string
		n    byte
		want []byte
	}{
		{
			name: "no feed (0 lines)",
			n:    0,
			want: []byte{common.ESC, 'd', 0},
		},
		{
			name: "typical feed (5 lines)",
			n:    5,
			want: []byte{common.ESC, 'd', 5},
		},
		{
			name: "maximum feed (255 lines)",
			n:    255,
			want: []byte{common.ESC, 'd', 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pp.PrintAndFeedLines(tt.n)

			if err != nil {
				t.Errorf("PrintAndFeedLines(%d) unexpected error: %v", tt.n, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PrintAndFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}
