package print_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/adcondev/pos-printer/escpos/common"
	"github.com/adcondev/pos-printer/escpos/print"
)

// Naming Convention: Test{Struct}_{Method}_{Scenario}

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
			want:    []byte{'A', common.LF, 'B', common.HT, 'C', common.CR, 'D'},
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
	want := []byte{common.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.FormFeed() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndCarriageReturn_SingleByte(t *testing.T) {
	pc := &print.Commands{}

	got := pc.PrintAndCarriageReturn()
	want := []byte{common.CR}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.PrintAndCarriageReturn() = %#v, want %#v", got, want)
	}
}

func TestPrintCommands_PrintAndLineFeed_SingleByte(t *testing.T) {
	pc := &print.Commands{}

	got := pc.PrintAndLineFeed()
	want := []byte{common.LF}

	if !bytes.Equal(got, want) {
		t.Errorf("Commands.PrintAndLineFeed() = %#v, want %#v", got, want)
	}
}

// ============================================================================
// PagePrint Tests
// ============================================================================

func TestPagePrint_PrintDataInPageMode_ByteSequence(t *testing.T) {
	pp := &print.PagePrint{}

	got := pp.PrintDataInPageMode()
	want := []byte{common.ESC, common.FF}

	if !bytes.Equal(got, want) {
		t.Errorf("PagePrint.PrintDataInPageMode() = %#v, want %#v", got, want)
	}
}

func TestPagePrint_PrintAndReverseFeed_Validation(t *testing.T) {
	pp := &print.PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum reverse feed (0 units)",
				n:    0,
				want: []byte{common.ESC, 'K', 0},
			},
			{
				name: "typical reverse feed (10 units)",
				n:    10,
				want: []byte{common.ESC, 'K', 10},
			},
			{
				name: "maximum allowed reverse feed",
				n:    common.MaxReverseMotionUnits,
				want: []byte{common.ESC, 'K', common.MaxReverseMotionUnits},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := pp.PrintAndReverseFeed(tt.n)

				if err != nil {
					t.Errorf("PagePrint.PrintAndReverseFeed(%d) unexpected error: %v", tt.n, err)
				}
				if !bytes.Equal(got, tt.want) {
					t.Errorf("PagePrint.PrintAndReverseFeed(%d) = %#v, want %#v", tt.n, got, tt.want)
				}
			})
		}
	})

	t.Run("overflow error", func(t *testing.T) {
		n := common.MaxReverseMotionUnits + 1
		_, err := pp.PrintAndReverseFeed(n)

		if !errors.Is(err, common.ErrPrintReverseFeed) {
			t.Errorf("PagePrint.PrintAndReverseFeed(%d) error = %v, want %v", n, err, common.ErrPrintReverseFeed)
		}
	})
}

func TestPagePrint_PrintAndReverseFeedLines_Validation(t *testing.T) {
	pp := &print.PagePrint{}

	t.Run("valid range", func(t *testing.T) {
		tests := []struct {
			name string
			n    byte
			want []byte
		}{
			{
				name: "minimum lines (0)",
				n:    0,
				want: []byte{common.ESC, 'e', 0},
			},
			{
				name: "single line reverse",
				n:    1,
				want: []byte{common.ESC, 'e', 1},
			},
			{
				name: "maximum allowed lines",
				n:    common.MaxReverseFeedLines,
				want: []byte{common.ESC, 'e', common.MaxReverseFeedLines},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := pp.PrintAndReverseFeedLines(tt.n)

				if err != nil {
					t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) unexpected error: %v", tt.n, err)
				}
				if !bytes.Equal(got, tt.want) {
					t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
				}
			})
		}
	})

	t.Run("overflow error", func(t *testing.T) {
		n := common.MaxReverseFeedLines + 1
		_, err := pp.PrintAndReverseFeedLines(n)

		if !errors.Is(err, common.ErrPrintReverseFeedLines) {
			t.Errorf("PagePrint.PrintAndReverseFeedLines(%d) error = %v, want %v", n, err, common.ErrPrintReverseFeedLines)
		}
	})
}

func TestPagePrint_PrintAndFeedLines_ByteSequence(t *testing.T) {
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
			name: "typical feed (7 lines)",
			n:    7,
			want: []byte{common.ESC, 'd', 7},
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
				t.Errorf("PagePrint.PrintAndFeedLines(%d) unexpected error: %v", tt.n, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("PagePrint.PrintAndFeedLines(%d) = %#v, want %#v", tt.n, got, tt.want)
			}
		})
	}
}
